package json

import (
	"encoding"
	"fmt"
	"math"
	"reflect"
	"strings"
	"unsafe"
)

type compiledCode struct {
	code    *opcode
	linked  bool // whether recursive code already have linked
	curLen  uintptr
	nextLen uintptr
}

type opcodeSet struct {
	code       *opcode
	codeLength int
}

var (
	marshalJSONType = reflect.TypeOf((*Marshaler)(nil)).Elem()
	marshalTextType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func encodeCompileToGetCodeSetSlowPath(typeptr uintptr) (*opcodeSet, error) {
	opcodeMap := loadOpcodeMap()
	if codeSet, exists := opcodeMap[typeptr]; exists {
		return codeSet, nil
	}

	// noescape trick for header.typ ( reflect.*rtype )
	copiedType := *(**rtype)(unsafe.Pointer(&typeptr))

	code, err := encodeCompileHead(&encodeCompileContext{
		typ:                      copiedType,
		root:                     true,
		structTypeToCompiledCode: map[uintptr]*compiledCode{},
	})
	if err != nil {
		return nil, err
	}
	code = copyOpcode(code)
	codeLength := code.totalLength()
	codeSet := &opcodeSet{
		code:       code,
		codeLength: codeLength,
	}
	storeOpcodeSet(typeptr, codeSet, opcodeMap)
	return codeSet, nil
}

func encodeCompileHead(ctx *encodeCompileContext) (*opcode, error) {
	typ := ctx.typ
	switch {
	case typ.Implements(marshalJSONType):
		return encodeCompileMarshalJSON(ctx)
	case rtype_ptrTo(typ).Implements(marshalJSONType):
		return encodeCompileMarshalJSONPtr(ctx)
	case typ.Implements(marshalTextType):
		return encodeCompileMarshalText(ctx)
	case rtype_ptrTo(typ).Implements(marshalTextType):
		return encodeCompileMarshalTextPtr(ctx)
	}
	isPtr := false
	orgType := typ
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		isPtr = true
	}
	if typ.Kind() == reflect.Map {
		return encodeCompileMap(ctx.withType(typ), isPtr)
	} else if typ.Kind() == reflect.Struct {
		code, err := encodeCompileStruct(ctx.withType(typ), isPtr)
		if err != nil {
			return nil, err
		}
		encodeConvertHeadOnlyCode(code, isPtr)
		encodeOptimizeStructEnd(code)
		encodeLinkRecursiveCode(code)
		return code, nil
	} else if isPtr && typ.Implements(marshalTextType) {
		typ = orgType
	} else if isPtr && typ.Implements(marshalJSONType) {
		typ = orgType
	}
	code, err := encodeCompile(ctx.withType(typ))
	if err != nil {
		return nil, err
	}
	encodeConvertHeadOnlyCode(code, isPtr)
	encodeOptimizeStructEnd(code)
	encodeLinkRecursiveCode(code)
	return code, nil
}

func encodeLinkRecursiveCode(c *opcode) {
	for code := c; code.op != opEnd && code.op != opStructFieldRecursiveEnd; {
		switch code.op {
		case opStructFieldRecursive,
			opStructFieldPtrAnonymousHeadRecursive,
			opStructFieldAnonymousHeadRecursive:
			if code.jmp.linked {
				code = code.next
				continue
			}
			code.jmp.code = copyOpcode(code.jmp.code)
			c := code.jmp.code
			c.end.next = newEndOp(&encodeCompileContext{})
			c.op = c.op.ptrHeadToHead()

			beforeLastCode := c.end
			lastCode := beforeLastCode.next

			lastCode.idx = beforeLastCode.idx + uintptrSize
			lastCode.elemIdx = lastCode.idx + uintptrSize

			// extend length to alloc slot for elemIdx
			totalLength := uintptr(code.totalLength() + 1)
			nextTotalLength := uintptr(c.totalLength() + 1)

			c.end.next.op = opStructFieldRecursiveEnd

			code.jmp.curLen = totalLength
			code.jmp.nextLen = nextTotalLength
			code.jmp.linked = true

			encodeLinkRecursiveCode(code.jmp.code)
			code = code.next
			continue
		}
		switch code.op.codeType() {
		case codeArrayElem, codeSliceElem, codeMapKey:
			code = code.end
		default:
			code = code.next
		}
	}
}

func encodeOptimizeStructEnd(c *opcode) {
	for code := c; code.op != opEnd; {
		if code.op == opStructFieldRecursive {
			// ignore if exists recursive operation
			return
		}
		switch code.op.codeType() {
		case codeArrayElem, codeSliceElem, codeMapKey:
			code = code.end
		default:
			code = code.next
		}
	}

	for code := c; code.op != opEnd; {
		switch code.op.codeType() {
		case codeArrayElem, codeSliceElem, codeMapKey:
			code = code.end
		case codeStructEnd:
			switch code.op {
			case opStructEnd:
				prev := code.prevField
				if strings.Contains(prev.op.String(), "Head") {
					// not exists field
					code = code.next
					break
				}
				if prev.op != prev.op.fieldToEnd() {
					prev.op = prev.op.fieldToEnd()
					prev.next = code.next
				}
				code = code.next
			default:
				code = code.next
			}
		default:
			code = code.next
		}
	}
}

func encodeConvertHeadOnlyCode(c *opcode, isPtrHead bool) {
	if c.nextField == nil {
		return
	}
	if c.nextField.op.codeType() != codeStructEnd {
		return
	}
	switch c.op {
	case opStructFieldHead:
		encodeConvertHeadOnlyCode(c.next, false)
		if !strings.Contains(c.next.op.String(), "Only") {
			return
		}
		c.op = opStructFieldHeadOnly
	case opStructFieldHeadOmitEmpty:
		encodeConvertHeadOnlyCode(c.next, false)
		if !strings.Contains(c.next.op.String(), "Only") {
			return
		}
		c.op = opStructFieldHeadOmitEmptyOnly
	case opStructFieldHeadStringTag:
		encodeConvertHeadOnlyCode(c.next, false)
		if !strings.Contains(c.next.op.String(), "Only") {
			return
		}
		c.op = opStructFieldHeadStringTagOnly
	case opStructFieldPtrHead:
	}

	if strings.Contains(c.op.String(), "Marshal") {
		return
	}
	if strings.Contains(c.op.String(), "Slice") {
		return
	}
	if strings.Contains(c.op.String(), "Map") {
		return
	}

	isPtrOp := strings.Contains(c.op.String(), "Ptr")
	if isPtrOp && !isPtrHead {
		c.op = c.op.headToOnlyHead()
	} else if !isPtrOp && isPtrHead {
		c.op = c.op.headToPtrHead().headToOnlyHead()
	} else if isPtrOp && isPtrHead {
		c.op = c.op.headToPtrHead().headToOnlyHead()
	}
}

func encodeImplementsMarshaler(typ *rtype) bool {
	switch {
	case typ.Implements(marshalJSONType):
		return true
	case rtype_ptrTo(typ).Implements(marshalJSONType):
		return true
	case typ.Implements(marshalTextType):
		return true
	case rtype_ptrTo(typ).Implements(marshalTextType):
		return true
	}
	return false
}

func encodeCompile(ctx *encodeCompileContext) (*opcode, error) {
	typ := ctx.typ
	switch {
	case typ.Implements(marshalJSONType):
		return encodeCompileMarshalJSON(ctx)
	case rtype_ptrTo(typ).Implements(marshalJSONType):
		return encodeCompileMarshalJSONPtr(ctx)
	case typ.Implements(marshalTextType):
		return encodeCompileMarshalText(ctx)
	case rtype_ptrTo(typ).Implements(marshalTextType):
		return encodeCompileMarshalTextPtr(ctx)
	}
	switch typ.Kind() {
	case reflect.Ptr:
		return encodeCompilePtr(ctx)
	case reflect.Slice:
		elem := typ.Elem()
		if !encodeImplementsMarshaler(elem) && elem.Kind() == reflect.Uint8 {
			return encodeCompileBytes(ctx)
		}
		return encodeCompileSlice(ctx)
	case reflect.Array:
		return encodeCompileArray(ctx)
	case reflect.Map:
		return encodeCompileMap(ctx, true)
	case reflect.Struct:
		return encodeCompileStruct(ctx, false)
	case reflect.Interface:
		return encodeCompileInterface(ctx)
	case reflect.Int:
		return encodeCompileInt(ctx)
	case reflect.Int8:
		return encodeCompileInt8(ctx)
	case reflect.Int16:
		return encodeCompileInt16(ctx)
	case reflect.Int32:
		return encodeCompileInt32(ctx)
	case reflect.Int64:
		return encodeCompileInt64(ctx)
	case reflect.Uint:
		return encodeCompileUint(ctx)
	case reflect.Uint8:
		return encodeCompileUint8(ctx)
	case reflect.Uint16:
		return encodeCompileUint16(ctx)
	case reflect.Uint32:
		return encodeCompileUint32(ctx)
	case reflect.Uint64:
		return encodeCompileUint64(ctx)
	case reflect.Uintptr:
		return encodeCompileUint(ctx)
	case reflect.Float32:
		return encodeCompileFloat32(ctx)
	case reflect.Float64:
		return encodeCompileFloat64(ctx)
	case reflect.String:
		return encodeCompileString(ctx)
	case reflect.Bool:
		return encodeCompileBool(ctx)
	}
	return nil, &UnsupportedTypeError{Type: rtype2type(typ)}
}

func encodeCompileKey(ctx *encodeCompileContext) (*opcode, error) {
	typ := ctx.typ
	switch {
	case rtype_ptrTo(typ).Implements(marshalJSONType):
		return encodeCompileMarshalJSONPtr(ctx)
	case rtype_ptrTo(typ).Implements(marshalTextType):
		return encodeCompileMarshalTextPtr(ctx)
	}
	switch typ.Kind() {
	case reflect.Ptr:
		return encodeCompilePtr(ctx)
	case reflect.Interface:
		return encodeCompileInterface(ctx)
	case reflect.String:
		return encodeCompileString(ctx)
	case reflect.Int:
		return encodeCompileIntString(ctx)
	case reflect.Int8:
		return encodeCompileInt8String(ctx)
	case reflect.Int16:
		return encodeCompileInt16String(ctx)
	case reflect.Int32:
		return encodeCompileInt32String(ctx)
	case reflect.Int64:
		return encodeCompileInt64String(ctx)
	case reflect.Uint:
		return encodeCompileUintString(ctx)
	case reflect.Uint8:
		return encodeCompileUint8String(ctx)
	case reflect.Uint16:
		return encodeCompileUint16String(ctx)
	case reflect.Uint32:
		return encodeCompileUint32String(ctx)
	case reflect.Uint64:
		return encodeCompileUint64String(ctx)
	case reflect.Uintptr:
		return encodeCompileUintString(ctx)
	}
	return nil, &UnsupportedTypeError{Type: rtype2type(typ)}
}

func encodeCompilePtr(ctx *encodeCompileContext) (*opcode, error) {
	ptrOpcodeIndex := ctx.opcodeIndex
	ptrIndex := ctx.ptrIndex
	ctx.incIndex()
	code, err := encodeCompile(ctx.withType(ctx.typ.Elem()))
	if err != nil {
		return nil, err
	}
	ptrHeadOp := code.op.headToPtrHead()
	if code.op != ptrHeadOp {
		code.op = ptrHeadOp
		code.decOpcodeIndex()
		ctx.decIndex()
		return code, nil
	}
	c := ctx.context()
	c.opcodeIndex = ptrOpcodeIndex
	c.ptrIndex = ptrIndex
	return newOpCodeWithNext(c, opPtr, code), nil
}

func encodeCompileMarshalJSON(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opMarshalJSON)
	ctx.incIndex()
	return code, nil
}

func encodeCompileMarshalJSONPtr(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx.withType(rtype_ptrTo(ctx.typ)), opMarshalJSON)
	ctx.incIndex()
	return code, nil
}

func encodeCompileMarshalText(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opMarshalText)
	ctx.incIndex()
	return code, nil
}

func encodeCompileMarshalTextPtr(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx.withType(rtype_ptrTo(ctx.typ)), opMarshalText)
	ctx.incIndex()
	return code, nil
}

const intSize = 32 << (^uint(0) >> 63)

func encodeCompileInt(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opInt)
	switch intSize {
	case 32:
		code.mask = math.MaxUint32
		code.rshiftNum = 31
	default:
		code.mask = math.MaxUint64
		code.rshiftNum = 63
	}
	ctx.incIndex()
	return code, nil
}

func encodeCompileInt8(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opInt)
	code.mask = math.MaxUint8
	code.rshiftNum = 7
	ctx.incIndex()
	return code, nil
}

func encodeCompileInt16(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opInt)
	code.mask = math.MaxUint16
	code.rshiftNum = 15
	ctx.incIndex()
	return code, nil
}

func encodeCompileInt32(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opInt)
	code.mask = math.MaxUint32
	code.rshiftNum = 31
	ctx.incIndex()
	return code, nil
}

func encodeCompileInt64(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opInt)
	code.mask = math.MaxUint64
	code.rshiftNum = 63
	ctx.incIndex()
	return code, nil
}

func encodeCompileUint(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opUint)
	switch intSize {
	case 32:
		code.mask = math.MaxUint32
		code.rshiftNum = 31
	default:
		code.mask = math.MaxUint64
		code.rshiftNum = 63
	}
	ctx.incIndex()
	return code, nil
}

func encodeCompileUint8(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opUint)
	code.mask = math.MaxUint8
	code.rshiftNum = 7
	ctx.incIndex()
	return code, nil
}

func encodeCompileUint16(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opUint)
	code.mask = math.MaxUint16
	code.rshiftNum = 15
	ctx.incIndex()
	return code, nil
}

func encodeCompileUint32(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opUint)
	code.mask = math.MaxUint32
	code.rshiftNum = 31
	ctx.incIndex()
	return code, nil
}

func encodeCompileUint64(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opUint)
	code.mask = math.MaxUint64
	code.rshiftNum = 63
	ctx.incIndex()
	return code, nil
}

func encodeCompileIntString(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opIntString)
	switch intSize {
	case 32:
		code.mask = math.MaxUint32
		code.rshiftNum = 31
	default:
		code.mask = math.MaxUint64
		code.rshiftNum = 63
	}
	ctx.incIndex()
	return code, nil
}

func encodeCompileInt8String(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opIntString)
	code.mask = math.MaxUint8
	code.rshiftNum = 7
	ctx.incIndex()
	return code, nil
}

func encodeCompileInt16String(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opIntString)
	code.mask = math.MaxUint16
	code.rshiftNum = 15
	ctx.incIndex()
	return code, nil
}

func encodeCompileInt32String(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opIntString)
	code.mask = math.MaxUint32
	code.rshiftNum = 31
	ctx.incIndex()
	return code, nil
}

func encodeCompileInt64String(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opIntString)
	code.mask = math.MaxUint64
	code.rshiftNum = 63
	ctx.incIndex()
	return code, nil
}

func encodeCompileUintString(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opUintString)
	switch intSize {
	case 32:
		code.mask = math.MaxUint32
		code.rshiftNum = 31
	default:
		code.mask = math.MaxUint64
		code.rshiftNum = 63
	}
	ctx.incIndex()
	return code, nil
}

func encodeCompileUint8String(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opUintString)
	code.mask = math.MaxUint8
	code.rshiftNum = 7
	ctx.incIndex()
	return code, nil
}

func encodeCompileUint16String(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opUintString)
	code.mask = math.MaxUint16
	code.rshiftNum = 15
	ctx.incIndex()
	return code, nil
}

func encodeCompileUint32String(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opUintString)
	code.mask = math.MaxUint32
	code.rshiftNum = 31
	ctx.incIndex()
	return code, nil
}

func encodeCompileUint64String(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opUintString)
	code.mask = math.MaxUint64
	code.rshiftNum = 63
	ctx.incIndex()
	return code, nil
}

func encodeCompileFloat32(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opFloat32)
	ctx.incIndex()
	return code, nil
}

func encodeCompileFloat64(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opFloat64)
	ctx.incIndex()
	return code, nil
}

func encodeCompileString(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opString)
	ctx.incIndex()
	return code, nil
}

func encodeCompileBool(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opBool)
	ctx.incIndex()
	return code, nil
}

func encodeCompileBytes(ctx *encodeCompileContext) (*opcode, error) {
	code := newOpCode(ctx, opBytes)
	ctx.incIndex()
	return code, nil
}

func encodeCompileInterface(ctx *encodeCompileContext) (*opcode, error) {
	code := newInterfaceCode(ctx)
	ctx.incIndex()
	return code, nil
}

func encodeCompileSlice(ctx *encodeCompileContext) (*opcode, error) {
	ctx.root = false
	elem := ctx.typ.Elem()
	size := elem.Size()

	header := newSliceHeaderCode(ctx)
	ctx.incIndex()

	code, err := encodeCompile(ctx.withType(ctx.typ.Elem()).incIndent())
	if err != nil {
		return nil, err
	}

	// header => opcode => elem => end
	//             ^        |
	//             |________|

	elemCode := newSliceElemCode(ctx, header, size)
	ctx.incIndex()

	end := newOpCode(ctx, opSliceEnd)
	ctx.incIndex()

	header.elem = elemCode
	header.end = end
	header.next = code
	code.beforeLastCode().next = (*opcode)(unsafe.Pointer(elemCode))
	elemCode.next = code
	elemCode.end = end
	return (*opcode)(unsafe.Pointer(header)), nil
}

func encodeCompileArray(ctx *encodeCompileContext) (*opcode, error) {
	ctx.root = false
	typ := ctx.typ
	elem := typ.Elem()
	alen := typ.Len()
	size := elem.Size()

	header := newArrayHeaderCode(ctx, alen)
	ctx.incIndex()

	code, err := encodeCompile(ctx.withType(elem).incIndent())
	if err != nil {
		return nil, err
	}
	// header => opcode => elem => end
	//             ^        |
	//             |________|

	elemCode := newArrayElemCode(ctx, header, alen, size)
	ctx.incIndex()

	end := newOpCode(ctx, opArrayEnd)
	ctx.incIndex()

	header.elem = elemCode
	header.end = end
	header.next = code
	code.beforeLastCode().next = (*opcode)(unsafe.Pointer(elemCode))
	elemCode.next = code
	elemCode.end = end
	return (*opcode)(unsafe.Pointer(header)), nil
}

//go:linkname mapiterinit reflect.mapiterinit
//go:noescape
func mapiterinit(mapType *rtype, m unsafe.Pointer) unsafe.Pointer

//go:linkname mapiterkey reflect.mapiterkey
//go:noescape
func mapiterkey(it unsafe.Pointer) unsafe.Pointer

//go:linkname mapiternext reflect.mapiternext
//go:noescape
func mapiternext(it unsafe.Pointer)

//go:linkname maplen reflect.maplen
//go:noescape
func maplen(m unsafe.Pointer) int

func encodeCompileMap(ctx *encodeCompileContext, withLoad bool) (*opcode, error) {
	// header => code => value => code => key => code => value => code => end
	//                                     ^                       |
	//                                     |_______________________|
	ctx = ctx.incIndent()
	header := newMapHeaderCode(ctx, withLoad)
	ctx.incIndex()

	typ := ctx.typ
	keyType := ctx.typ.Key()
	keyCode, err := encodeCompileKey(ctx.withType(keyType))
	if err != nil {
		return nil, err
	}

	value := newMapValueCode(ctx, header)
	ctx.incIndex()

	valueType := typ.Elem()
	valueCode, err := encodeCompile(ctx.withType(valueType))
	if err != nil {
		return nil, err
	}

	key := newMapKeyCode(ctx, header)
	ctx.incIndex()

	ctx = ctx.decIndent()

	header.mapKey = key
	header.mapValue = value

	end := newMapEndCode(ctx, header)
	ctx.incIndex()

	header.next = keyCode
	keyCode.beforeLastCode().next = (*opcode)(unsafe.Pointer(value))
	value.next = valueCode
	valueCode.beforeLastCode().next = (*opcode)(unsafe.Pointer(key))
	key.next = keyCode

	header.end = end
	key.end = end
	value.end = end

	return (*opcode)(unsafe.Pointer(header)), nil
}

func encodeTypeToHeaderType(ctx *encodeCompileContext, code *opcode) opType {
	switch code.op {
	case opPtr:
		ptrNum := 1
		c := code
		ctx.decIndex()
		for {
			if code.next.op == opPtr {
				ptrNum++
				code = code.next
				ctx.decIndex()
				continue
			}
			break
		}
		c.ptrNum = ptrNum
		if ptrNum > 1 {
			switch code.next.op {
			case opInt:
				c.mask = code.next.mask
				c.rshiftNum = code.next.rshiftNum
				return opStructFieldHeadIntNPtr
			case opUint:
				c.mask = code.next.mask
				return opStructFieldHeadUintNPtr
			case opFloat32:
				return opStructFieldHeadFloat32NPtr
			case opFloat64:
				return opStructFieldHeadFloat64NPtr
			case opString:
				return opStructFieldHeadStringNPtr
			case opBool:
				return opStructFieldHeadBoolNPtr
			}
		} else {
			switch code.next.op {
			case opInt:
				c.mask = code.next.mask
				c.rshiftNum = code.next.rshiftNum
				return opStructFieldHeadIntPtr
			case opUint:
				c.mask = code.next.mask
				return opStructFieldHeadUintPtr
			case opFloat32:
				return opStructFieldHeadFloat32Ptr
			case opFloat64:
				return opStructFieldHeadFloat64Ptr
			case opString:
				return opStructFieldHeadStringPtr
			case opBool:
				return opStructFieldHeadBoolPtr
			}
		}
	case opInt:
		return opStructFieldHeadInt
	case opUint:
		return opStructFieldHeadUint
	case opFloat32:
		return opStructFieldHeadFloat32
	case opFloat64:
		return opStructFieldHeadFloat64
	case opString:
		return opStructFieldHeadString
	case opBool:
		return opStructFieldHeadBool
	case opMapHead:
		return opStructFieldHeadMap
	case opMapHeadLoad:
		return opStructFieldHeadMapLoad
	case opArrayHead:
		return opStructFieldHeadArray
	case opSliceHead:
		return opStructFieldHeadSlice
	case opStructFieldHead:
		return opStructFieldHeadStruct
	case opMarshalJSON:
		return opStructFieldHeadMarshalJSON
	case opMarshalText:
		return opStructFieldHeadMarshalText
	}
	return opStructFieldHead
}

func encodeTypeToFieldType(ctx *encodeCompileContext, code *opcode) opType {
	switch code.op {
	case opPtr:
		ptrNum := 1
		ctx.decIndex()
		c := code
		for {
			if code.next.op == opPtr {
				ptrNum++
				code = code.next
				ctx.decIndex()
				continue
			}
			break
		}
		c.ptrNum = ptrNum
		if ptrNum > 1 {
			switch code.next.op {
			case opInt:
				c.mask = code.next.mask
				c.rshiftNum = code.next.rshiftNum
				return opStructFieldIntNPtr
			case opUint:
				c.mask = code.next.mask
				return opStructFieldUintNPtr
			case opFloat32:
				return opStructFieldFloat32NPtr
			case opFloat64:
				return opStructFieldFloat64NPtr
			case opString:
				return opStructFieldStringNPtr
			case opBool:
				return opStructFieldBoolNPtr
			}
		} else {
			switch code.next.op {
			case opInt:
				c.mask = code.next.mask
				c.rshiftNum = code.next.rshiftNum
				return opStructFieldIntPtr
			case opUint:
				c.mask = code.next.mask
				return opStructFieldUintPtr
			case opFloat32:
				return opStructFieldFloat32Ptr
			case opFloat64:
				return opStructFieldFloat64Ptr
			case opString:
				return opStructFieldStringPtr
			case opBool:
				return opStructFieldBoolPtr
			}
		}
	case opInt:
		return opStructFieldInt
	case opUint:
		return opStructFieldUint
	case opFloat32:
		return opStructFieldFloat32
	case opFloat64:
		return opStructFieldFloat64
	case opString:
		return opStructFieldString
	case opBool:
		return opStructFieldBool
	case opMapHead:
		return opStructFieldMap
	case opMapHeadLoad:
		return opStructFieldMapLoad
	case opArrayHead:
		return opStructFieldArray
	case opSliceHead:
		return opStructFieldSlice
	case opStructFieldHead:
		return opStructFieldStruct
	case opMarshalJSON:
		return opStructFieldMarshalJSON
	case opMarshalText:
		return opStructFieldMarshalText
	}
	return opStructField
}

func encodeOptimizeStructHeader(ctx *encodeCompileContext, code *opcode, tag *structTag) opType {
	headType := encodeTypeToHeaderType(ctx, code)
	switch {
	case tag.isOmitEmpty:
		headType = headType.headToOmitEmptyHead()
	case tag.isString:
		headType = headType.headToStringTagHead()
	}
	return headType
}

func encodeOptimizeStructField(ctx *encodeCompileContext, code *opcode, tag *structTag) opType {
	fieldType := encodeTypeToFieldType(ctx, code)
	switch {
	case tag.isOmitEmpty:
		fieldType = fieldType.fieldToOmitEmptyField()
	case tag.isString:
		fieldType = fieldType.fieldToStringTagField()
	}
	return fieldType
}

func encodeRecursiveCode(ctx *encodeCompileContext, jmp *compiledCode) *opcode {
	code := newRecursiveCode(ctx, jmp)
	ctx.incIndex()
	return code
}

func encodeCompiledCode(ctx *encodeCompileContext) *opcode {
	typ := ctx.typ
	typeptr := uintptr(unsafe.Pointer(typ))
	if compiledCode, exists := ctx.structTypeToCompiledCode[typeptr]; exists {
		return encodeRecursiveCode(ctx, compiledCode)
	}
	return nil
}

func encodeStructHeader(ctx *encodeCompileContext, fieldCode *opcode, valueCode *opcode, tag *structTag) *opcode {
	fieldCode.indent--
	op := encodeOptimizeStructHeader(ctx, valueCode, tag)
	fieldCode.op = op
	fieldCode.mask = valueCode.mask
	fieldCode.rshiftNum = valueCode.rshiftNum
	fieldCode.ptrNum = valueCode.ptrNum
	switch op {
	case opStructFieldHead,
		opStructFieldHeadSlice,
		opStructFieldHeadArray,
		opStructFieldHeadMap,
		opStructFieldHeadMapLoad,
		opStructFieldHeadStruct,
		opStructFieldHeadOmitEmpty,
		opStructFieldHeadOmitEmptySlice,
		opStructFieldHeadOmitEmptyArray,
		opStructFieldHeadOmitEmptyMap,
		opStructFieldHeadOmitEmptyMapLoad,
		opStructFieldHeadOmitEmptyStruct,
		opStructFieldHeadStringTag:
		return valueCode.beforeLastCode()
	}
	ctx.decOpcodeIndex()
	return (*opcode)(unsafe.Pointer(fieldCode))
}

func encodeStructField(ctx *encodeCompileContext, fieldCode *opcode, valueCode *opcode, tag *structTag) *opcode {
	code := (*opcode)(unsafe.Pointer(fieldCode))
	op := encodeOptimizeStructField(ctx, valueCode, tag)
	fieldCode.op = op
	fieldCode.ptrNum = valueCode.ptrNum
	fieldCode.mask = valueCode.mask
	fieldCode.rshiftNum = valueCode.rshiftNum
	switch op {
	case opStructField,
		opStructFieldSlice,
		opStructFieldArray,
		opStructFieldMap,
		opStructFieldMapLoad,
		opStructFieldStruct,
		opStructFieldOmitEmpty,
		opStructFieldOmitEmptySlice,
		opStructFieldOmitEmptyArray,
		opStructFieldOmitEmptyMap,
		opStructFieldOmitEmptyMapLoad,
		opStructFieldOmitEmptyStruct,
		opStructFieldStringTag:
		return valueCode.beforeLastCode()
	}
	ctx.decIndex()
	return code
}

func encodeIsNotExistsField(head *opcode) bool {
	if head == nil {
		return false
	}
	if head.op != opStructFieldAnonymousHead {
		return false
	}
	if head.next == nil {
		return false
	}
	if head.nextField == nil {
		return false
	}
	if head.nextField.op != opStructAnonymousEnd {
		return false
	}
	if head.next.op == opStructAnonymousEnd {
		return true
	}
	if head.next.op.codeType() != codeStructField {
		return false
	}
	return encodeIsNotExistsField(head.next)
}

func encodeOptimizeAnonymousFields(head *opcode) {
	code := head
	var prev *opcode
	removedFields := map[*opcode]struct{}{}
	for {
		if code.op == opStructEnd {
			break
		}
		if code.op == opStructField {
			codeType := code.next.op.codeType()
			if codeType == codeStructField {
				if encodeIsNotExistsField(code.next) {
					code.next = code.nextField
					diff := code.next.displayIdx - code.displayIdx
					for i := 0; i < diff; i++ {
						code.next.decOpcodeIndex()
					}
					encodeLinkPrevToNextField(code, removedFields)
					code = prev
				}
			}
		}
		prev = code
		code = code.nextField
	}
}

type structFieldPair struct {
	prevField   *opcode
	curField    *opcode
	isTaggedKey bool
	linked      bool
}

func encodeAnonymousStructFieldPairMap(tags structTags, named string, valueCode *opcode) map[string][]structFieldPair {
	anonymousFields := map[string][]structFieldPair{}
	f := valueCode
	var prevAnonymousField *opcode
	removedFields := map[*opcode]struct{}{}
	for {
		existsKey := tags.existsKey(f.displayKey)
		op := f.op.headToAnonymousHead()
		if existsKey && (f.next.op == opStructFieldPtrAnonymousHeadRecursive || f.next.op == opStructFieldAnonymousHeadRecursive) {
			// through
		} else if op != f.op {
			if existsKey {
				f.op = opStructFieldAnonymousHead
			} else if named == "" {
				f.op = op
			}
		} else if named == "" && f.op == opStructEnd {
			f.op = opStructAnonymousEnd
		} else if existsKey {
			diff := f.nextField.displayIdx - f.displayIdx
			for i := 0; i < diff; i++ {
				f.nextField.decOpcodeIndex()
			}
			encodeLinkPrevToNextField(f, removedFields)
		}

		if f.displayKey == "" {
			if f.nextField == nil {
				break
			}
			prevAnonymousField = f
			f = f.nextField
			continue
		}

		key := fmt.Sprintf("%s.%s", named, f.displayKey)
		anonymousFields[key] = append(anonymousFields[key], structFieldPair{
			prevField:   prevAnonymousField,
			curField:    f,
			isTaggedKey: f.isTaggedKey,
		})
		if f.next != nil && f.nextField != f.next && f.next.op.codeType() == codeStructField {
			for k, v := range encodeAnonymousFieldPairRecursively(named, f.next) {
				anonymousFields[k] = append(anonymousFields[k], v...)
			}
		}
		if f.nextField == nil {
			break
		}
		prevAnonymousField = f
		f = f.nextField
	}
	return anonymousFields
}

func encodeAnonymousFieldPairRecursively(named string, valueCode *opcode) map[string][]structFieldPair {
	anonymousFields := map[string][]structFieldPair{}
	f := valueCode
	var prevAnonymousField *opcode
	for {
		if f.displayKey != "" && strings.Contains(f.op.String(), "Anonymous") {
			key := fmt.Sprintf("%s.%s", named, f.displayKey)
			anonymousFields[key] = append(anonymousFields[key], structFieldPair{
				prevField:   prevAnonymousField,
				curField:    f,
				isTaggedKey: f.isTaggedKey,
			})
			if f.next != nil && f.nextField != f.next && f.next.op.codeType() == codeStructField {
				for k, v := range encodeAnonymousFieldPairRecursively(named, f.next) {
					anonymousFields[k] = append(anonymousFields[k], v...)
				}
			}
		}
		if f.nextField == nil {
			break
		}
		prevAnonymousField = f
		f = f.nextField
	}
	return anonymousFields
}

func encodeOptimizeConflictAnonymousFields(anonymousFields map[string][]structFieldPair) {
	removedFields := map[*opcode]struct{}{}
	for _, fieldPairs := range anonymousFields {
		if len(fieldPairs) == 1 {
			continue
		}
		// conflict anonymous fields
		taggedPairs := []structFieldPair{}
		for _, fieldPair := range fieldPairs {
			if fieldPair.isTaggedKey {
				taggedPairs = append(taggedPairs, fieldPair)
			} else {
				if !fieldPair.linked {
					if fieldPair.prevField == nil {
						// head operation
						fieldPair.curField.op = opStructFieldAnonymousHead
					} else {
						diff := fieldPair.curField.nextField.displayIdx - fieldPair.curField.displayIdx
						for i := 0; i < diff; i++ {
							fieldPair.curField.nextField.decOpcodeIndex()
						}
						removedFields[fieldPair.curField] = struct{}{}
						encodeLinkPrevToNextField(fieldPair.curField, removedFields)
					}
					fieldPair.linked = true
				}
			}
		}
		if len(taggedPairs) > 1 {
			for _, fieldPair := range taggedPairs {
				if !fieldPair.linked {
					if fieldPair.prevField == nil {
						// head operation
						fieldPair.curField.op = opStructFieldAnonymousHead
					} else {
						diff := fieldPair.curField.nextField.displayIdx - fieldPair.curField.displayIdx
						removedFields[fieldPair.curField] = struct{}{}
						for i := 0; i < diff; i++ {
							fieldPair.curField.nextField.decOpcodeIndex()
						}
						encodeLinkPrevToNextField(fieldPair.curField, removedFields)
					}
					fieldPair.linked = true
				}
			}
		} else {
			for _, fieldPair := range taggedPairs {
				fieldPair.curField.isTaggedKey = false
			}
		}
	}
}

func encodeCompileStruct(ctx *encodeCompileContext, isPtr bool) (*opcode, error) {
	ctx.root = false
	if code := encodeCompiledCode(ctx); code != nil {
		return code, nil
	}
	typ := ctx.typ
	typeptr := uintptr(unsafe.Pointer(typ))
	compiled := &compiledCode{}
	ctx.structTypeToCompiledCode[typeptr] = compiled
	// header => code => structField => code => end
	//                        ^          |
	//                        |__________|
	fieldNum := typ.NumField()
	fieldIdx := 0
	var (
		head      *opcode
		code      *opcode
		prevField *opcode
	)
	ctx = ctx.incIndent()
	tags := structTags{}
	anonymousFields := map[string][]structFieldPair{}
	for i := 0; i < fieldNum; i++ {
		field := typ.Field(i)
		if isIgnoredStructField(field) {
			continue
		}
		tags = append(tags, structTagFromField(field))
	}
	for i, tag := range tags {
		field := tag.field
		fieldType := type2rtype(field.Type)
		if isPtr && i == 0 {
			// head field of pointer structure at top level
			// if field type is pointer and implements MarshalJSON or MarshalText,
			// it need to operation of dereference of pointer.
			if field.Type.Kind() == reflect.Ptr &&
				(field.Type.Implements(marshalJSONType) || field.Type.Implements(marshalTextType)) {
				fieldType = rtype_ptrTo(fieldType)
			}
		}
		fieldOpcodeIndex := ctx.opcodeIndex
		fieldPtrIndex := ctx.ptrIndex
		ctx.incIndex()
		valueCode, err := encodeCompile(ctx.withType(fieldType))
		if err != nil {
			return nil, err
		}

		if field.Anonymous {
			if valueCode.op == opPtr && valueCode.next.op == opStructFieldRecursive {
				valueCode = valueCode.next
				valueCode.decOpcodeIndex()
				ctx.decIndex()
				valueCode.op = opStructFieldPtrHeadRecursive
			}
			tagKey := ""
			if tag.isTaggedKey {
				tagKey = tag.key
			}
			for k, v := range encodeAnonymousStructFieldPairMap(tags, tagKey, valueCode) {
				anonymousFields[k] = append(anonymousFields[k], v...)
			}
		}
		key := fmt.Sprintf(`"%s":`, tag.key)
		escapedKey := fmt.Sprintf(`%s:`, string(encodeEscapedString([]byte{}, tag.key)))
		fieldCode := &opcode{
			typ:          valueCode.typ,
			displayIdx:   fieldOpcodeIndex,
			idx:          opcodeOffset(fieldPtrIndex),
			next:         valueCode,
			indent:       ctx.indent,
			anonymousKey: field.Anonymous,
			key:          []byte(key),
			escapedKey:   []byte(escapedKey),
			isTaggedKey:  tag.isTaggedKey,
			displayKey:   tag.key,
			offset:       field.Offset,
		}
		if fieldIdx == 0 {
			fieldCode.headIdx = fieldCode.idx
			code = encodeStructHeader(ctx, fieldCode, valueCode, tag)
			head = fieldCode
			prevField = fieldCode
		} else {
			fieldCode.headIdx = head.headIdx
			code.next = fieldCode
			code = encodeStructField(ctx, fieldCode, valueCode, tag)
			prevField.nextField = fieldCode
			fieldCode.prevField = prevField
			prevField = fieldCode
		}
		fieldIdx++
	}
	ctx = ctx.decIndent()

	structEndCode := &opcode{
		op:     opStructEnd,
		typ:    nil,
		indent: ctx.indent,
		next:   newEndOp(ctx),
	}

	// no struct field
	if head == nil {
		head = &opcode{
			op:         opStructFieldHead,
			typ:        typ,
			displayIdx: ctx.opcodeIndex,
			idx:        opcodeOffset(ctx.ptrIndex),
			headIdx:    opcodeOffset(ctx.ptrIndex),
			indent:     ctx.indent,
			nextField:  structEndCode,
		}
		structEndCode.prevField = head
		ctx.incIndex()
		code = head
	}

	structEndCode.displayIdx = ctx.opcodeIndex
	structEndCode.idx = opcodeOffset(ctx.ptrIndex)
	ctx.incIndex()

	if prevField != nil && prevField.nextField == nil {
		prevField.nextField = structEndCode
		structEndCode.prevField = prevField
	}

	head.end = structEndCode
	code.next = structEndCode
	encodeOptimizeConflictAnonymousFields(anonymousFields)
	encodeOptimizeAnonymousFields(head)
	ret := (*opcode)(unsafe.Pointer(head))
	compiled.code = ret

	delete(ctx.structTypeToCompiledCode, typeptr)

	return ret, nil
}
