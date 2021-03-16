package json

import (
	"bytes"
	"encoding"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"unsafe"
)

const startDetectingCyclesAfter = 1000

func load(base uintptr, idx uintptr) uintptr {
	addr := base + idx
	return **(**uintptr)(unsafe.Pointer(&addr))
}

func store(base uintptr, idx uintptr, p uintptr) {
	addr := base + idx
	**(**uintptr)(unsafe.Pointer(&addr)) = p
}

func ptrToUint64(p uintptr) uint64      { return **(**uint64)(unsafe.Pointer(&p)) }
func ptrToFloat32(p uintptr) float32    { return **(**float32)(unsafe.Pointer(&p)) }
func ptrToFloat64(p uintptr) float64    { return **(**float64)(unsafe.Pointer(&p)) }
func ptrToBool(p uintptr) bool          { return **(**bool)(unsafe.Pointer(&p)) }
func ptrToBytes(p uintptr) []byte       { return **(**[]byte)(unsafe.Pointer(&p)) }
func ptrToString(p uintptr) string      { return **(**string)(unsafe.Pointer(&p)) }
func ptrToSlice(p uintptr) *sliceHeader { return *(**sliceHeader)(unsafe.Pointer(&p)) }
func ptrToPtr(p uintptr) uintptr {
	return uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
}
func ptrToUnsafePtr(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}
func ptrToInterface(code *opcode, p uintptr) interface{} {
	return *(*interface{})(unsafe.Pointer(&interfaceHeader{
		typ: code.typ,
		ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
	}))
}

func errUnsupportedValue(code *opcode, ptr uintptr) *UnsupportedValueError {
	v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
		typ: code.typ,
		ptr: *(*unsafe.Pointer)(unsafe.Pointer(&ptr)),
	}))
	return &UnsupportedValueError{
		Value: reflect.ValueOf(v),
		Str:   fmt.Sprintf("encountered a cycle via %s", code.typ),
	}
}

func errUnsupportedFloat(v float64) *UnsupportedValueError {
	return &UnsupportedValueError{
		Value: reflect.ValueOf(v),
		Str:   strconv.FormatFloat(v, 'g', -1, 64),
	}
}

func errMarshaler(code *opcode, err error) *MarshalerError {
	return &MarshalerError{
		Type: rtype2type(code.typ),
		Err:  err,
	}
}

func encodeRun(ctx *encodeRuntimeContext, b []byte, codeSet *opcodeSet, opt EncodeOption) ([]byte, error) {
	recursiveLevel := 0
	ptrOffset := uintptr(0)
	ctxptr := ctx.ptr()
	code := codeSet.code

	for {
		switch code.op {
		default:
			return nil, fmt.Errorf("encoder: opcode %s has not been implemented", code.op)
		case opPtr:
			ptr := load(ctxptr, code.idx)
			code = code.next
			store(ctxptr, code.idx, ptrToPtr(ptr))
		case opInt:
			b = appendInt(b, ptrToUint64(load(ctxptr, code.idx)), code)
			b = encodeComma(b)
			code = code.next
		case opUint:
			b = appendUint(b, ptrToUint64(load(ctxptr, code.idx)), code)
			b = encodeComma(b)
			code = code.next
		case opIntString:
			b = append(b, '"')
			b = appendInt(b, ptrToUint64(load(ctxptr, code.idx)), code)
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opUintString:
			b = append(b, '"')
			b = appendUint(b, ptrToUint64(load(ctxptr, code.idx)), code)
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opFloat32:
			b = encodeFloat32(b, ptrToFloat32(load(ctxptr, code.idx)))
			b = encodeComma(b)
			code = code.next
		case opFloat64:
			v := ptrToFloat64(load(ctxptr, code.idx))
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return nil, errUnsupportedFloat(v)
			}
			b = encodeFloat64(b, v)
			b = encodeComma(b)
			code = code.next
		case opString:
			b = encodeNoEscapedString(b, ptrToString(load(ctxptr, code.idx)))
			b = encodeComma(b)
			code = code.next
		case opBool:
			b = encodeBool(b, ptrToBool(load(ctxptr, code.idx)))
			b = encodeComma(b)
			code = code.next
		case opBytes:
			ptr := load(ctxptr, code.idx)
			slice := ptrToSlice(ptr)
			if ptr == 0 || uintptr(slice.data) == 0 {
				b = encodeNull(b)
			} else {
				b = encodeByteSlice(b, ptrToBytes(ptr))
			}
			b = encodeComma(b)
			code = code.next
		case opInterface:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.next
				break
			}
			for _, seen := range ctx.seenPtr {
				if ptr == seen {
					return nil, errUnsupportedValue(code, ptr)
				}
			}
			ctx.seenPtr = append(ctx.seenPtr, ptr)
			iface := (*interfaceHeader)(ptrToUnsafePtr(ptr))
			if iface == nil || iface.ptr == nil {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.next
				break
			}
			ctx.keepRefs = append(ctx.keepRefs, unsafe.Pointer(iface))
			ifaceCodeSet, err := encodeCompileToGetCodeSet(uintptr(unsafe.Pointer(iface.typ)))
			if err != nil {
				return nil, err
			}

			totalLength := uintptr(codeSet.codeLength)
			nextTotalLength := uintptr(ifaceCodeSet.codeLength)

			curlen := uintptr(len(ctx.ptrs))
			offsetNum := ptrOffset / uintptrSize

			newLen := offsetNum + totalLength + nextTotalLength
			if curlen < newLen {
				ctx.ptrs = append(ctx.ptrs, make([]uintptr, newLen-curlen)...)
			}
			oldPtrs := ctx.ptrs

			newPtrs := ctx.ptrs[(ptrOffset+totalLength*uintptrSize)/uintptrSize:]
			newPtrs[0] = uintptr(iface.ptr)

			ctx.ptrs = newPtrs

			bb, err := encodeRun(ctx, b, ifaceCodeSet, opt)
			if err != nil {
				return nil, err
			}

			ctx.ptrs = oldPtrs
			ctxptr = ctx.ptr()
			ctx.seenPtr = ctx.seenPtr[:len(ctx.seenPtr)-1]

			b = bb
			code = code.next
		case opMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.next
				break
			}
			v := ptrToInterface(code, ptr)
			bb, err := v.(Marshaler).MarshalJSON()
			if err != nil {
				return nil, errMarshaler(code, err)
			}
			runtime.KeepAlive(v)
			if len(bb) == 0 {
				return nil, errUnexpectedEndOfJSON(
					fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
					0,
				)
			}
			buf := bytes.NewBuffer(b)
			//TODO: we should validate buffer with `compact`
			if err := compact(buf, bb, false); err != nil {
				return nil, err
			}
			b = buf.Bytes()
			b = encodeComma(b)
			code = code.next
		case opMarshalText:
			ptr := load(ctxptr, code.idx)
			isPtr := code.typ.Kind() == reflect.Ptr
			p := ptrToUnsafePtr(ptr)
			if p == nil || isPtr && *(*unsafe.Pointer)(p) == nil {
				b = append(b, '"', '"', ',')
			} else {
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
					typ: code.typ,
					ptr: p,
				}))
				bytes, err := v.(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return nil, errMarshaler(code, err)
				}
				b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
				b = encodeComma(b)
			}
			code = code.next
		case opSliceHead:
			p := load(ctxptr, code.idx)
			slice := ptrToSlice(p)
			if p == 0 || uintptr(slice.data) == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				store(ctxptr, code.elemIdx, 0)
				store(ctxptr, code.length, uintptr(slice.len))
				store(ctxptr, code.idx, uintptr(slice.data))
				if slice.len > 0 {
					b = append(b, '[')
					code = code.next
					store(ctxptr, code.idx, uintptr(slice.data))
				} else {
					b = append(b, '[', ']', ',')
					code = code.end.next
				}
			}
		case opSliceElem:
			idx := load(ctxptr, code.elemIdx)
			length := load(ctxptr, code.length)
			idx++
			if idx < length {
				store(ctxptr, code.elemIdx, idx)
				data := load(ctxptr, code.headIdx)
				size := code.size
				code = code.next
				store(ctxptr, code.idx, data+idx*size)
			} else {
				last := len(b) - 1
				b[last] = ']'
				b = encodeComma(b)
				code = code.end.next
			}
		case opArrayHead:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				if code.length > 0 {
					b = append(b, '[')
					store(ctxptr, code.elemIdx, 0)
					code = code.next
					store(ctxptr, code.idx, p)
				} else {
					b = append(b, '[', ']', ',')
					code = code.end.next
				}
			}
		case opArrayElem:
			idx := load(ctxptr, code.elemIdx)
			idx++
			if idx < code.length {
				store(ctxptr, code.elemIdx, idx)
				p := load(ctxptr, code.headIdx)
				size := code.size
				code = code.next
				store(ctxptr, code.idx, p+idx*size)
			} else {
				last := len(b) - 1
				b[last] = ']'
				b = encodeComma(b)
				code = code.end.next
			}
		case opMapHead:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				uptr := ptrToUnsafePtr(ptr)
				mlen := maplen(uptr)
				if mlen > 0 {
					b = append(b, '{')
					iter := mapiterinit(code.typ, uptr)
					ctx.keepRefs = append(ctx.keepRefs, iter)
					store(ctxptr, code.elemIdx, 0)
					store(ctxptr, code.length, uintptr(mlen))
					store(ctxptr, code.mapIter, uintptr(iter))
					if (opt & EncodeOptionUnorderedMap) == 0 {
						mapCtx := newMapContext(mlen)
						mapCtx.pos = append(mapCtx.pos, len(b))
						ctx.keepRefs = append(ctx.keepRefs, unsafe.Pointer(mapCtx))
						store(ctxptr, code.end.mapPos, uintptr(unsafe.Pointer(mapCtx)))
					}
					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))
					code = code.next
				} else {
					b = append(b, '{', '}', ',')
					code = code.end.next
				}
			}
		case opMapHeadLoad:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				// load pointer
				ptr = ptrToPtr(ptr)
				uptr := ptrToUnsafePtr(ptr)
				if ptr == 0 {
					b = encodeNull(b)
					b = encodeComma(b)
					code = code.end.next
					break
				}
				mlen := maplen(uptr)
				if mlen > 0 {
					b = append(b, '{')
					iter := mapiterinit(code.typ, uptr)
					ctx.keepRefs = append(ctx.keepRefs, iter)
					store(ctxptr, code.elemIdx, 0)
					store(ctxptr, code.length, uintptr(mlen))
					store(ctxptr, code.mapIter, uintptr(iter))
					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))
					if (opt & EncodeOptionUnorderedMap) == 0 {
						mapCtx := newMapContext(mlen)
						mapCtx.pos = append(mapCtx.pos, len(b))
						ctx.keepRefs = append(ctx.keepRefs, unsafe.Pointer(mapCtx))
						store(ctxptr, code.end.mapPos, uintptr(unsafe.Pointer(mapCtx)))
					}
					code = code.next
				} else {
					b = append(b, '{', '}', ',')
					code = code.end.next
				}
			}
		case opMapKey:
			idx := load(ctxptr, code.elemIdx)
			length := load(ctxptr, code.length)
			idx++
			if (opt & EncodeOptionUnorderedMap) != 0 {
				if idx < length {
					ptr := load(ctxptr, code.mapIter)
					iter := ptrToUnsafePtr(ptr)
					store(ctxptr, code.elemIdx, idx)
					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))
					code = code.next
				} else {
					last := len(b) - 1
					b[last] = '}'
					b = encodeComma(b)
					code = code.end.next
				}
			} else {
				ptr := load(ctxptr, code.end.mapPos)
				mapCtx := (*encodeMapContext)(ptrToUnsafePtr(ptr))
				mapCtx.pos = append(mapCtx.pos, len(b))
				if idx < length {
					ptr := load(ctxptr, code.mapIter)
					iter := ptrToUnsafePtr(ptr)
					store(ctxptr, code.elemIdx, idx)
					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))
					code = code.next
				} else {
					code = code.end
				}
			}
		case opMapValue:
			if (opt & EncodeOptionUnorderedMap) != 0 {
				last := len(b) - 1
				b[last] = ':'
			} else {
				ptr := load(ctxptr, code.end.mapPos)
				mapCtx := (*encodeMapContext)(ptrToUnsafePtr(ptr))
				mapCtx.pos = append(mapCtx.pos, len(b))
			}
			ptr := load(ctxptr, code.mapIter)
			iter := ptrToUnsafePtr(ptr)
			value := mapitervalue(iter)
			store(ctxptr, code.next.idx, uintptr(value))
			mapiternext(iter)
			code = code.next
		case opMapEnd:
			// this operation only used by sorted map.
			length := int(load(ctxptr, code.length))
			ptr := load(ctxptr, code.mapPos)
			mapCtx := (*encodeMapContext)(ptrToUnsafePtr(ptr))
			pos := mapCtx.pos
			for i := 0; i < length; i++ {
				startKey := pos[i*2]
				startValue := pos[i*2+1]
				var endValue int
				if i+1 < length {
					endValue = pos[i*2+2]
				} else {
					endValue = len(b)
				}
				mapCtx.slice.items = append(mapCtx.slice.items, mapItem{
					key:   b[startKey:startValue],
					value: b[startValue:endValue],
				})
			}
			sort.Sort(mapCtx.slice)
			buf := mapCtx.buf
			for _, item := range mapCtx.slice.items {
				buf = append(buf, item.key...)
				buf[len(buf)-1] = ':'
				buf = append(buf, item.value...)
			}
			buf[len(buf)-1] = '}'
			buf = append(buf, ',')
			b = b[:pos[0]]
			b = append(b, buf...)
			mapCtx.buf = buf
			releaseMapContext(mapCtx)
			code = code.next
		case opStructFieldPtrAnonymousHeadRecursive:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadRecursive:
			fallthrough
		case opStructFieldRecursive:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				if recursiveLevel > startDetectingCyclesAfter {
					for _, seen := range ctx.seenPtr {
						if ptr == seen {
							return nil, errUnsupportedValue(code, ptr)
						}
					}
				}
			}
			ctx.seenPtr = append(ctx.seenPtr, ptr)
			c := code.jmp.code
			curlen := uintptr(len(ctx.ptrs))
			offsetNum := ptrOffset / uintptrSize
			oldOffset := ptrOffset
			ptrOffset += code.jmp.curLen * uintptrSize

			newLen := offsetNum + code.jmp.curLen + code.jmp.nextLen
			if curlen < newLen {
				ctx.ptrs = append(ctx.ptrs, make([]uintptr, newLen-curlen)...)
			}
			ctxptr = ctx.ptr() + ptrOffset // assign new ctxptr

			store(ctxptr, c.idx, ptr)
			store(ctxptr, c.end.next.idx, oldOffset)
			store(ctxptr, c.end.next.elemIdx, uintptr(unsafe.Pointer(code.next)))
			code = c
			recursiveLevel++
		case opStructFieldRecursiveEnd:
			recursiveLevel--

			// restore ctxptr
			offset := load(ctxptr, code.idx)
			ctx.seenPtr = ctx.seenPtr[:len(ctx.seenPtr)-1]

			codePtr := load(ctxptr, code.elemIdx)
			code = (*opcode)(ptrToUnsafePtr(codePtr))
			ctxptr = ctx.ptr() + offset
			ptrOffset = offset
		case opStructFieldPtrHead:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHead:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				if !code.anonymousKey {
					b = append(b, code.key...)
				}
				p := ptr + code.offset
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldPtrHeadOmitEmpty:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmpty:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				p := ptr + code.offset
				if p == 0 || *(*uintptr)(*(*unsafe.Pointer)(unsafe.Pointer(&p))) == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					code = code.next
					store(ctxptr, code.idx, p)
				}
			}
		case opStructFieldHeadOnly, opStructFieldHeadStringTagOnly:
			ptr := load(ctxptr, code.idx)
			b = append(b, '{')
			if !code.anonymousKey {
				b = append(b, code.key...)
			}
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldHeadOmitEmptyOnly:
			ptr := load(ctxptr, code.idx)
			b = append(b, '{')
			if !code.anonymousKey {
				if ptr != 0 {
					b = append(b, code.key...)
					p := ptr + code.offset
					code = code.next
					store(ctxptr, code.idx, p)
				} else {
					code = code.nextField
				}
			}
		case opStructFieldAnonymousHead:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				code = code.next
				store(ctxptr, code.idx, ptr)
			}
		case opStructFieldPtrAnonymousHeadOmitEmpty:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmpty:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				p := ptr + code.offset
				if p == 0 || *(*uintptr)(*(*unsafe.Pointer)(unsafe.Pointer(&p))) == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					code = code.next
					store(ctxptr, code.idx, p)
				}
			}
		case opStructFieldPtrHeadInt:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				b = appendInt(b, ptrToUint64(ptr+code.offset), code)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadOmitEmptyInt:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				u64 := ptrToUint64(ptr + code.offset)
				v := u64 & code.mask
				if v == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = appendInt(b, u64, code)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrHeadStringTagInt:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				b = append(b, '"')
				b = appendInt(b, ptrToUint64(ptr+code.offset), code)
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadIntOnly, opStructFieldHeadIntOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			b = appendInt(b, ptrToUint64(p), code)
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyIntOnly, opStructFieldHeadOmitEmptyIntOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			u64 := ptrToUint64(p)
			v := u64 & code.mask
			if v != 0 {
				b = append(b, code.key...)
				b = appendInt(b, u64, code)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagIntOnly, opStructFieldHeadStringTagIntOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			b = append(b, '"')
			b = appendInt(b, ptrToUint64(p), code)
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadIntPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadIntPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p + code.offset)
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = appendInt(b, ptrToUint64(p), code)
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyIntPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadOmitEmptyIntPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				p = ptrToPtr(p + code.offset)
				if p != 0 {
					b = append(b, code.key...)
					b = appendInt(b, ptrToUint64(p), code)
					b = encodeComma(b)
				}
				code = code.next
			}
		case opStructFieldPtrHeadStringTagIntPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadStringTagIntPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p + code.offset)
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = append(b, '"')
					b = appendInt(b, ptrToUint64(p), code)
					b = append(b, '"')
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadIntPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadIntPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendInt(b, ptrToUint64(p+code.offset), code)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyIntPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadOmitEmptyIntPtrOnly:
			b = append(b, '{')
			p := load(ctxptr, code.idx)
			if p != 0 {
				b = append(b, code.key...)
				b = appendInt(b, ptrToUint64(p), code)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagIntPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadStringTagIntPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = appendInt(b, ptrToUint64(p), code)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldHeadIntNPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				for i := 0; i < code.ptrNum; i++ {
					if p == 0 {
						break
					}
					p = ptrToPtr(p)
				}
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = appendInt(b, ptrToUint64(p), code)
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadInt:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = appendInt(b, ptrToUint64(ptr+code.offset), code)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyInt:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				u64 := ptrToUint64(ptr + code.offset)
				v := u64 & code.mask
				if v == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = appendInt(b, u64, code)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagInt:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = append(b, '"')
				b = appendInt(b, ptrToUint64(ptr+code.offset), code)
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadIntOnly, opStructFieldAnonymousHeadIntOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = appendInt(b, ptrToUint64(ptr+code.offset), code)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyIntOnly, opStructFieldAnonymousHeadOmitEmptyIntOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				u64 := ptrToUint64(ptr + code.offset)
				v := u64 & code.mask
				if v == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = appendInt(b, u64, code)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagIntOnly, opStructFieldAnonymousHeadStringTagIntOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = append(b, '"')
				b = appendInt(b, ptrToUint64(ptr+code.offset), code)
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadIntPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadIntPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendInt(b, ptrToUint64(p), code)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyIntPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyIntPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			p = ptrToPtr(p + code.offset)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				b = appendInt(b, ptrToUint64(p), code)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagIntPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadStringTagIntPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = appendInt(b, ptrToUint64(p), code)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadIntPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadIntPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendInt(b, ptrToUint64(p+code.offset), code)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyIntPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyIntPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				b = appendInt(b, ptrToUint64(p+code.offset), code)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagIntPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadStringTagIntPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = appendInt(b, ptrToUint64(p+code.offset), code)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadUint:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				b = appendUint(b, ptrToUint64(ptr+code.offset), code)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadOmitEmptyUint:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				u64 := ptrToUint64(ptr + code.offset)
				v := u64 & code.mask
				if v == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = appendUint(b, u64, code)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrHeadStringTagUint:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				b = append(b, '"')
				b = appendUint(b, ptrToUint64(ptr+code.offset), code)
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadUintOnly, opStructFieldHeadUintOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			b = appendUint(b, ptrToUint64(p), code)
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyUintOnly, opStructFieldHeadOmitEmptyUintOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			u64 := ptrToUint64(p)
			v := u64 & code.mask
			if v != 0 {
				b = append(b, code.key...)
				b = appendUint(b, u64, code)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagUintOnly, opStructFieldHeadStringTagUintOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			b = append(b, '"')
			b = appendUint(b, ptrToUint64(p), code)
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadUintPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadUintPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p + code.offset)
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = appendUint(b, ptrToUint64(p), code)
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyUintPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadOmitEmptyUintPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				p = ptrToPtr(p + code.offset)
				if p != 0 {
					b = append(b, code.key...)
					b = appendUint(b, ptrToUint64(p), code)
					b = encodeComma(b)
				}
				code = code.next
			}
		case opStructFieldPtrHeadStringTagUintPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadStringTagUintPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p + code.offset)
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = append(b, '"')
					b = appendUint(b, ptrToUint64(p), code)
					b = append(b, '"')
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadUintPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadUintPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendUint(b, ptrToUint64(p+code.offset), code)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyUintPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadOmitEmptyUintPtrOnly:
			b = append(b, '{')
			p := load(ctxptr, code.idx)
			if p != 0 {
				b = append(b, code.key...)
				b = appendUint(b, ptrToUint64(p+code.offset), code)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagUintPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadStringTagUintPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = appendUint(b, ptrToUint64(p+code.offset), code)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldHeadUintNPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				for i := 0; i < code.ptrNum; i++ {
					if p == 0 {
						break
					}
					p = ptrToPtr(p)
				}
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = appendUint(b, ptrToUint64(p), code)
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadUint:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = appendUint(b, ptrToUint64(ptr+code.offset), code)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyUint:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				u64 := ptrToUint64(ptr + code.offset)
				v := u64 & code.mask
				if v == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = appendUint(b, u64, code)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagUint:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = append(b, '"')
				b = appendUint(b, ptrToUint64(ptr+code.offset), code)
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadUintOnly, opStructFieldAnonymousHeadUintOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = appendUint(b, ptrToUint64(ptr+code.offset), code)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyUintOnly, opStructFieldAnonymousHeadOmitEmptyUintOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				u64 := ptrToUint64(ptr + code.offset)
				v := u64 & code.mask
				if v == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = appendUint(b, u64, code)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagUintOnly, opStructFieldAnonymousHeadStringTagUintOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = append(b, '"')
				b = appendUint(b, ptrToUint64(ptr+code.offset), code)
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadUintPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadUintPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendUint(b, ptrToUint64(p), code)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyUintPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyUintPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			p = ptrToPtr(p + code.offset)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				b = appendUint(b, ptrToUint64(p), code)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagUintPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadStringTagUintPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = appendUint(b, ptrToUint64(p), code)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadUintPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadUintPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendUint(b, ptrToUint64(p+code.offset), code)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyUintPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyUintPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				b = appendUint(b, ptrToUint64(p+code.offset), code)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagUintPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadStringTagUintPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = appendUint(b, ptrToUint64(p+code.offset), code)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadFloat32:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				b = encodeFloat32(b, ptrToFloat32(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadOmitEmptyFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				v := ptrToFloat32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = encodeFloat32(b, v)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrHeadStringTagFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				b = append(b, '"')
				b = encodeFloat32(b, ptrToFloat32(ptr+code.offset))
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadFloat32Only, opStructFieldHeadFloat32Only:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			b = encodeFloat32(b, ptrToFloat32(p))
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyFloat32Only, opStructFieldHeadOmitEmptyFloat32Only:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			v := ptrToFloat32(p)
			if v != 0 {
				b = append(b, code.key...)
				b = encodeFloat32(b, v)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagFloat32Only, opStructFieldHeadStringTagFloat32Only:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			b = append(b, '"')
			b = encodeFloat32(b, ptrToFloat32(p))
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadFloat32Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadFloat32Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p)
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = encodeFloat32(b, ptrToFloat32(p+code.offset))
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyFloat32Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadOmitEmptyFloat32Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				p = ptrToPtr(p)
				if p != 0 {
					b = append(b, code.key...)
					b = encodeFloat32(b, ptrToFloat32(p))
					b = encodeComma(b)
				}
				code = code.next
			}
		case opStructFieldPtrHeadStringTagFloat32Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadStringTagFloat32Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p)
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = append(b, '"')
					b = encodeFloat32(b, ptrToFloat32(p+code.offset))
					b = append(b, '"')
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadFloat32PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadFloat32PtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeFloat32(b, ptrToFloat32(p+code.offset))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyFloat32PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadOmitEmptyFloat32PtrOnly:
			b = append(b, '{')
			p := load(ctxptr, code.idx)
			if p != 0 {
				b = append(b, code.key...)
				b = encodeFloat32(b, ptrToFloat32(p+code.offset))
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagFloat32PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadStringTagFloat32PtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = encodeFloat32(b, ptrToFloat32(p+code.offset))
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldHeadFloat32NPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				for i := 0; i < code.ptrNum; i++ {
					if p == 0 {
						break
					}
					p = ptrToPtr(p)
				}
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = encodeFloat32(b, ptrToFloat32(p+code.offset))
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadFloat32:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = encodeFloat32(b, ptrToFloat32(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := ptrToFloat32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = encodeFloat32(b, v)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = append(b, '"')
				b = encodeFloat32(b, ptrToFloat32(ptr+code.offset))
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadFloat32Only, opStructFieldAnonymousHeadFloat32Only:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = encodeFloat32(b, ptrToFloat32(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyFloat32Only, opStructFieldAnonymousHeadOmitEmptyFloat32Only:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := ptrToFloat32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = encodeFloat32(b, v)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagFloat32Only, opStructFieldAnonymousHeadStringTagFloat32Only:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = append(b, '"')
				b = encodeFloat32(b, ptrToFloat32(ptr+code.offset))
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadFloat32Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadFloat32Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeFloat32(b, ptrToFloat32(p+code.offset))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyFloat32Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyFloat32Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			p = ptrToPtr(p)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				b = encodeFloat32(b, ptrToFloat32(p))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagFloat32Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadStringTagFloat32Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = encodeFloat32(b, ptrToFloat32(p+code.offset))
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadFloat32PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadFloat32PtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeFloat32(b, ptrToFloat32(p+code.offset))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyFloat32PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyFloat32PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				b = encodeFloat32(b, ptrToFloat32(p+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagFloat32PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadStringTagFloat32PtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = encodeFloat32(b, ptrToFloat32(p+code.offset))
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadFloat64:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				v := ptrToFloat64(ptr + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = append(b, '{')
				b = append(b, code.key...)
				b = encodeFloat64(b, v)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadOmitEmptyFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				v := ptrToFloat64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					if math.IsInf(v, 0) || math.IsNaN(v) {
						return nil, errUnsupportedFloat(v)
					}
					b = append(b, code.key...)
					b = encodeFloat64(b, v)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrHeadStringTagFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				v := ptrToFloat64(ptr + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = append(b, code.key...)
				b = append(b, '"')
				b = encodeFloat64(b, v)
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadFloat64Only, opStructFieldHeadFloat64Only:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			v := ptrToFloat64(p)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return nil, errUnsupportedFloat(v)
			}
			b = encodeFloat64(b, v)
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyFloat64Only, opStructFieldHeadOmitEmptyFloat64Only:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			v := ptrToFloat64(p)
			if v != 0 {
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = append(b, code.key...)
				b = encodeFloat64(b, v)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagFloat64Only, opStructFieldHeadStringTagFloat64Only:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			b = append(b, '"')
			v := ptrToFloat64(p)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return nil, errUnsupportedFloat(v)
			}
			b = encodeFloat64(b, v)
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadFloat64Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadFloat64Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p)
				if p == 0 {
					b = encodeNull(b)
				} else {
					v := ptrToFloat64(p + code.offset)
					if math.IsInf(v, 0) || math.IsNaN(v) {
						return nil, errUnsupportedFloat(v)
					}
					b = encodeFloat64(b, v)
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyFloat64Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadOmitEmptyFloat64Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				p = ptrToPtr(p)
				if p != 0 {
					b = append(b, code.key...)
					v := ptrToFloat64(p + code.offset)
					if math.IsInf(v, 0) || math.IsNaN(v) {
						return nil, errUnsupportedFloat(v)
					}
					b = encodeFloat64(b, v)
					b = encodeComma(b)
				}
				code = code.next
			}
		case opStructFieldPtrHeadStringTagFloat64Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadStringTagFloat64Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p)
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = append(b, '"')
					v := ptrToFloat64(p + code.offset)
					if math.IsInf(v, 0) || math.IsNaN(v) {
						return nil, errUnsupportedFloat(v)
					}
					b = encodeFloat64(b, v)
					b = append(b, '"')
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadFloat64PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadFloat64PtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				v := ptrToFloat64(p + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyFloat64PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadOmitEmptyFloat64PtrOnly:
			b = append(b, '{')
			p := load(ctxptr, code.idx)
			if p != 0 {
				b = append(b, code.key...)
				v := ptrToFloat64(p)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagFloat64PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadStringTagFloat64PtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				v := ptrToFloat64(p)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldHeadFloat64NPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				for i := 0; i < code.ptrNum; i++ {
					if p == 0 {
						break
					}
					p = ptrToPtr(p)
				}
				if p == 0 {
					b = encodeNull(b)
				} else {
					v := ptrToFloat64(p)
					if math.IsInf(v, 0) || math.IsNaN(v) {
						return nil, errUnsupportedFloat(v)
					}
					b = encodeFloat64(b, v)
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadFloat64:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := ptrToFloat64(ptr + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = append(b, code.key...)
				b = encodeFloat64(b, v)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := ptrToFloat64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					v := ptrToFloat64(ptr + code.offset)
					if math.IsInf(v, 0) || math.IsNaN(v) {
						return nil, errUnsupportedFloat(v)
					}
					b = encodeFloat64(b, v)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = append(b, '"')
				v := ptrToFloat64(ptr + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadFloat64Only, opStructFieldAnonymousHeadFloat64Only:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				v := ptrToFloat64(ptr + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyFloat64Only, opStructFieldAnonymousHeadOmitEmptyFloat64Only:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := ptrToFloat64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					v := ptrToFloat64(ptr + code.offset)
					if math.IsInf(v, 0) || math.IsNaN(v) {
						return nil, errUnsupportedFloat(v)
					}
					b = encodeFloat64(b, v)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagFloat64Only, opStructFieldAnonymousHeadStringTagFloat64Only:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = append(b, '"')
				v := ptrToFloat64(ptr + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadFloat64Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadFloat64Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p)
			if p == 0 {
				b = encodeNull(b)
			} else {
				v := ptrToFloat64(p + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyFloat64Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyFloat64Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			p = ptrToPtr(p)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				v := ptrToFloat64(p + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagFloat64Ptr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadStringTagFloat64Ptr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				v := ptrToFloat64(p + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadFloat64PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadFloat64PtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				v := ptrToFloat64(p + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyFloat64PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyFloat64PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				v := ptrToFloat64(p + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagFloat64PtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadStringTagFloat64PtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				v := ptrToFloat64(p + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadString:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, ptrToString(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadOmitEmptyString:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				v := ptrToString(ptr + code.offset)
				if v == "" {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = encodeNoEscapedString(b, v)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrHeadStringTagString:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				s := ptrToString(ptr + code.offset)
				b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, s)))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadStringOnly, opStructFieldHeadStringOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			b = encodeNoEscapedString(b, ptrToString(p))
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyStringOnly, opStructFieldHeadOmitEmptyStringOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			v := ptrToString(p)
			if v != "" {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, v)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagStringOnly, opStructFieldHeadStringTagStringOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, ptrToString(p))))
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadStringPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadStringPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p)
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = encodeNoEscapedString(b, ptrToString(p+code.offset))
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyStringPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadOmitEmptyStringPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				p = ptrToPtr(p)
				if p != 0 {
					b = append(b, code.key...)
					b = encodeNoEscapedString(b, ptrToString(p))
					b = encodeComma(b)
				}
				code = code.next
			}
		case opStructFieldPtrHeadStringTagStringPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadStringTagStringPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p)
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, ptrToString(p))))
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadStringPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadStringPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeNoEscapedString(b, ptrToString(p+code.offset))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyStringPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadOmitEmptyStringPtrOnly:
			b = append(b, '{')
			p := load(ctxptr, code.idx)
			if p != 0 {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, ptrToString(p+code.offset))
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagStringPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadStringTagStringPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, ptrToString(p))))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldHeadStringNPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				for i := 0; i < code.ptrNum; i++ {
					if p == 0 {
						break
					}
					p = ptrToPtr(p)
				}
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = encodeNoEscapedString(b, ptrToString(p+code.offset))
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadString:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, ptrToString(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyString:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := ptrToString(ptr + code.offset)
				if v == "" {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = encodeNoEscapedString(b, v)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagString:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, ptrToString(ptr))))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringOnly, opStructFieldAnonymousHeadStringOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, ptrToString(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyStringOnly, opStructFieldAnonymousHeadOmitEmptyStringOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := ptrToString(ptr + code.offset)
				if v == "" {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = encodeNoEscapedString(b, v)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagStringOnly, opStructFieldAnonymousHeadStringTagStringOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, ptrToString(ptr))))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadStringPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeNoEscapedString(b, ptrToString(p+code.offset))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyStringPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyStringPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			p = ptrToPtr(p)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, ptrToString(p))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagStringPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadStringTagStringPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, ptrToString(p))))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadStringPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadStringPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeNoEscapedString(b, ptrToString(p+code.offset))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyStringPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyStringPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, ptrToString(p+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagStringPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadStringTagStringPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, ptrToString(p))))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadBool:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				b = encodeBool(b, ptrToBool(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadOmitEmptyBool:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				v := ptrToBool(ptr + code.offset)
				if v {
					b = append(b, code.key...)
					b = encodeBool(b, v)
					b = encodeComma(b)
					code = code.next
				} else {
					code = code.nextField
				}
			}
		case opStructFieldPtrHeadStringTagBool:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				b = append(b, '"')
				b = encodeBool(b, ptrToBool(ptr+code.offset))
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadBoolOnly, opStructFieldHeadBoolOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			b = encodeBool(b, ptrToBool(p))
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyBoolOnly, opStructFieldHeadOmitEmptyBoolOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			v := ptrToBool(p)
			if v {
				b = append(b, code.key...)
				b = encodeBool(b, v)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagBoolOnly, opStructFieldHeadStringTagBoolOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			b = append(b, '"')
			b = encodeBool(b, ptrToBool(p))
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadBoolPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadBoolPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p)
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = encodeBool(b, ptrToBool(p+code.offset))
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyBoolPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadOmitEmptyBoolPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				p = ptrToPtr(p)
				if p != 0 {
					b = append(b, code.key...)
					b = encodeBool(b, ptrToBool(p))
					b = encodeComma(b)
				}
				code = code.next
			}
		case opStructFieldPtrHeadStringTagBoolPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadStringTagBoolPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				p = ptrToPtr(p)
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = append(b, '"')
					b = encodeBool(b, ptrToBool(p+code.offset))
					b = append(b, '"')
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadBoolPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadBoolPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeBool(b, ptrToBool(p+code.offset))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadOmitEmptyBoolPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadOmitEmptyBoolPtrOnly:
			b = append(b, '{')
			p := load(ctxptr, code.idx)
			if p != 0 {
				b = append(b, code.key...)
				b = encodeBool(b, ptrToBool(p+code.offset))
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldPtrHeadStringTagBoolPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadStringTagBoolPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, '{')
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = encodeBool(b, ptrToBool(p+code.offset))
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldHeadBoolNPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				for i := 0; i < code.ptrNum; i++ {
					if p == 0 {
						break
					}
					p = ptrToPtr(p)
				}
				if p == 0 {
					b = encodeNull(b)
				} else {
					b = encodeBool(b, ptrToBool(p+code.offset))
				}
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadBool:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = encodeBool(b, ptrToBool(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyBool:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := ptrToBool(ptr + code.offset)
				if v {
					b = append(b, code.key...)
					b = encodeBool(b, v)
					b = encodeComma(b)
					code = code.next
				} else {
					code = code.nextField
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagBool:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = append(b, '"')
				b = encodeBool(b, ptrToBool(ptr+code.offset))
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadBoolOnly, opStructFieldAnonymousHeadBoolOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = encodeBool(b, ptrToBool(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyBoolOnly, opStructFieldAnonymousHeadOmitEmptyBoolOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := ptrToBool(ptr + code.offset)
				if v {
					b = append(b, code.key...)
					b = encodeBool(b, v)
					b = encodeComma(b)
					code = code.next
				} else {
					code = code.nextField
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagBoolOnly, opStructFieldAnonymousHeadStringTagBoolOnly:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = append(b, '"')
				b = encodeBool(b, ptrToBool(ptr+code.offset))
				b = append(b, '"')
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadBoolPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadBoolPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeBool(b, ptrToBool(p+code.offset))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyBoolPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyBoolPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			p = ptrToPtr(p)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				b = encodeBool(b, ptrToBool(p))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagBoolPtr:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadStringTagBoolPtr:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			b = append(b, code.key...)
			p = ptrToPtr(p)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = encodeBool(b, ptrToBool(p+code.offset))
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadBoolPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadBoolPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeBool(b, ptrToBool(p+code.offset))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrAnonymousHeadOmitEmptyBoolPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyBoolPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				b = encodeBool(b, ptrToBool(p+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagBoolPtrOnly:
			p := load(ctxptr, code.idx)
			if p == 0 {
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldAnonymousHeadStringTagBoolPtrOnly:
			p := load(ctxptr, code.idx)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = encodeBool(b, ptrToBool(p+code.offset))
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldPtrHeadBytes:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadBytes {
					b = encodeNull(b)
					b = encodeComma(b)
				} else {
					b = append(b, '{', '}', ',')
				}
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				b = encodeByteSlice(b, ptrToBytes(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadBytes:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = encodeByteSlice(b, ptrToBytes(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadArray:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadArray:
			ptr := load(ctxptr, code.idx) + code.offset
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadArray {
					b = encodeNull(b)
					b = encodeComma(b)
				} else {
					b = append(b, '[', ']', ',')
				}
				code = code.end.next
			} else {
				b = append(b, '{')
				if !code.anonymousKey {
					b = append(b, code.key...)
				}
				code = code.next
				store(ctxptr, code.idx, ptr)
			}
		case opStructFieldPtrAnonymousHeadArray:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadArray:
			ptr := load(ctxptr, code.idx) + code.offset
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				store(ctxptr, code.idx, ptr)
				code = code.next
			}
		case opStructFieldPtrHeadSlice:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadSlice:
			ptr := load(ctxptr, code.idx)
			p := ptr + code.offset
			if p == 0 {
				if code.op == opStructFieldPtrHeadSlice {
					b = encodeNull(b)
					b = encodeComma(b)
				} else {
					b = append(b, '[', ']', ',')
				}
				code = code.end.next
			} else {
				b = append(b, '{')
				if !code.anonymousKey {
					b = append(b, code.key...)
				}
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldPtrAnonymousHeadSlice:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadSlice:
			ptr := load(ctxptr, code.idx)
			p := ptr + code.offset
			if p == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				store(ctxptr, code.idx, p)
				code = code.next
			}
		case opStructFieldPtrHeadMarshalJSON:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				ptr += code.offset
				v := ptrToInterface(code, ptr)
				rv := reflect.ValueOf(v)
				if rv.Type().Kind() == reflect.Interface && rv.IsNil() {
					b = encodeNull(b)
					code = code.end
					break
				}
				bb, err := rv.Interface().(Marshaler).MarshalJSON()
				if err != nil {
					return nil, errMarshaler(code, err)
				}
				if len(bb) == 0 {
					return nil, errUnexpectedEndOfJSON(
						fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
						0,
					)
				}
				buf := bytes.NewBuffer(b)
				//TODO: we should validate buffer with `compact`
				if err := compact(buf, bb, false); err != nil {
					return nil, err
				}
				b = buf.Bytes()
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadMarshalJSON:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				ptr += code.offset
				v := ptrToInterface(code, ptr)
				rv := reflect.ValueOf(v)
				if rv.Type().Kind() == reflect.Interface && rv.IsNil() {
					b = encodeNull(b)
					code = code.end.next
					break
				}
				bb, err := rv.Interface().(Marshaler).MarshalJSON()
				if err != nil {
					return nil, errMarshaler(code, err)
				}
				if len(bb) == 0 {
					return nil, errUnexpectedEndOfJSON(
						fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
						0,
					)
				}
				buf := bytes.NewBuffer(b)
				//TODO: we should validate buffer with `compact`
				if err := compact(buf, bb, false); err != nil {
					return nil, err
				}
				b = buf.Bytes()
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadMarshalText:
			p := load(ctxptr, code.idx)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, ptrToPtr(p))
			fallthrough
		case opStructFieldHeadMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				ptr += code.offset
				v := ptrToInterface(code, ptr)
				rv := reflect.ValueOf(v)
				if rv.Type().Kind() == reflect.Interface && rv.IsNil() {
					b = encodeNull(b)
					b = encodeComma(b)
					code = code.end
					break
				}
				bytes, err := rv.Interface().(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return nil, errMarshaler(code, err)
				}
				b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadMarshalText:
			store(ctxptr, code.idx, ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				ptr += code.offset
				v := ptrToInterface(code, ptr)
				rv := reflect.ValueOf(v)
				if rv.Type().Kind() == reflect.Interface && rv.IsNil() {
					b = encodeNull(b)
					b = encodeComma(b)
					code = code.end.next
					break
				}
				bytes, err := rv.Interface().(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return nil, errMarshaler(code, err)
				}
				b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadOmitEmptyBytes:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				v := ptrToBytes(ptr + code.offset)
				if len(v) == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = encodeByteSlice(b, v)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyBytes:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := ptrToBytes(ptr + code.offset)
				if len(v) == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					b = encodeByteSlice(b, v)
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				ptr += code.offset
				p := ptrToUnsafePtr(ptr)
				isPtr := code.typ.Kind() == reflect.Ptr
				if p == nil || (!isPtr && *(*unsafe.Pointer)(p) == nil) {
					code = code.nextField
				} else {
					v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
					bb, err := v.(Marshaler).MarshalJSON()
					if err != nil {
						return nil, &MarshalerError{
							Type: rtype2type(code.typ),
							Err:  err,
						}
					}
					if len(bb) == 0 {
						if isPtr {
							return nil, errUnexpectedEndOfJSON(
								fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
								0,
							)
						}
						code = code.nextField
					} else {
						b = append(b, code.key...)
						buf := bytes.NewBuffer(b)
						//TODO: we should validate buffer with `compact`
						if err := compact(buf, bb, false); err != nil {
							return nil, err
						}
						b = buf.Bytes()
						b = encodeComma(b)
						code = code.next
					}
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				ptr += code.offset
				p := ptrToUnsafePtr(ptr)
				isPtr := code.typ.Kind() == reflect.Ptr
				if p == nil || (!isPtr && *(*unsafe.Pointer)(p) == nil) {
					code = code.nextField
				} else {
					v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
					bb, err := v.(Marshaler).MarshalJSON()
					if err != nil {
						return nil, &MarshalerError{
							Type: rtype2type(code.typ),
							Err:  err,
						}
					}
					if len(bb) == 0 {
						if isPtr {
							return nil, errUnexpectedEndOfJSON(
								fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
								0,
							)
						}
						code = code.nextField
					} else {
						b = append(b, code.key...)
						buf := bytes.NewBuffer(b)
						//TODO: we should validate buffer with `compact`
						if err := compact(buf, bb, false); err != nil {
							return nil, err
						}
						b = buf.Bytes()
						b = encodeComma(b)
						code = code.next
					}
				}
			}
		case opStructFieldPtrHeadOmitEmptyMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				ptr += code.offset
				p := ptrToUnsafePtr(ptr)
				isPtr := code.typ.Kind() == reflect.Ptr
				if p == nil || (!isPtr && *(*unsafe.Pointer)(p) == nil) {
					code = code.nextField
				} else {
					v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
					bytes, err := v.(encoding.TextMarshaler).MarshalText()
					if err != nil {
						return nil, &MarshalerError{
							Type: rtype2type(code.typ),
							Err:  err,
						}
					}
					b = append(b, code.key...)
					b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				ptr += code.offset
				p := ptrToUnsafePtr(ptr)
				isPtr := code.typ.Kind() == reflect.Ptr
				if p == nil || (!isPtr && *(*unsafe.Pointer)(p) == nil) {
					code = code.nextField
				} else {
					v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
					bytes, err := v.(encoding.TextMarshaler).MarshalText()
					if err != nil {
						return nil, &MarshalerError{
							Type: rtype2type(code.typ),
							Err:  err,
						}
					}
					b = append(b, code.key...)
					b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrHeadStringTag:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTag:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				p := ptr + code.offset
				b = append(b, code.key...)
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldPtrAnonymousHeadStringTag:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTag:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				code = code.next
				store(ctxptr, code.idx, ptr+code.offset)
			}
		case opStructFieldPtrHeadStringTagBytes:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				b = append(b, code.key...)
				b = encodeByteSlice(b, ptrToBytes(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagBytes:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				b = append(b, code.key...)
				b = encodeByteSlice(b, ptrToBytes(ptr+code.offset))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrHeadStringTagMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				ptr += code.offset
				p := ptrToUnsafePtr(ptr)
				isPtr := code.typ.Kind() == reflect.Ptr
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
				bb, err := v.(Marshaler).MarshalJSON()
				if err != nil {
					return nil, &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				if len(bb) == 0 {
					if isPtr {
						return nil, errUnexpectedEndOfJSON(
							fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
							0,
						)
					}
					b = append(b, code.key...)
					b = append(b, '"', '"')
					b = encodeComma(b)
					code = code.nextField
				} else {
					var buf bytes.Buffer
					if err := compact(&buf, bb, false); err != nil {
						return nil, err
					}
					b = append(b, code.key...)
					b = encodeNoEscapedString(b, buf.String())
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				ptr += code.offset
				p := ptrToUnsafePtr(ptr)
				isPtr := code.typ.Kind() == reflect.Ptr
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
				bb, err := v.(Marshaler).MarshalJSON()
				if err != nil {
					return nil, &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				if len(bb) == 0 {
					if isPtr {
						return nil, errUnexpectedEndOfJSON(
							fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
							0,
						)
					}
					b = append(b, code.key...)
					b = append(b, '"', '"')
					b = encodeComma(b)
					code = code.nextField
				} else {
					var buf bytes.Buffer
					if err := compact(&buf, bb, false); err != nil {
						return nil, err
					}
					b = append(b, code.key...)
					b = encodeNoEscapedString(b, buf.String())
					b = encodeComma(b)
					code = code.next
				}
			}
		case opStructFieldPtrHeadStringTagMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.end.next
			} else {
				b = append(b, '{')
				ptr += code.offset
				p := ptrToUnsafePtr(ptr)
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
				bytes, err := v.(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return nil, &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
				b = encodeComma(b)
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				ptr += code.offset
				p := ptrToUnsafePtr(ptr)
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
				bytes, err := v.(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return nil, errMarshaler(code, err)
				}
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
				b = encodeComma(b)
				code = code.next
			}
		case opStructField:
			if !code.anonymousKey {
				b = append(b, code.key...)
			}
			ptr := load(ctxptr, code.headIdx) + code.offset
			code = code.next
			store(ctxptr, code.idx, ptr)
		case opStructFieldOmitEmpty:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 || **(**uintptr)(unsafe.Pointer(&p)) == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldStringTag:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			b = append(b, code.key...)
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldInt:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = appendInt(b, ptrToUint64(ptr+code.offset), code)
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyInt:
			ptr := load(ctxptr, code.headIdx)
			u64 := ptrToUint64(ptr + code.offset)
			v := u64 & code.mask
			if v != 0 {
				b = append(b, code.key...)
				b = appendInt(b, u64, code)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagInt:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = append(b, '"')
			b = appendInt(b, ptrToUint64(ptr+code.offset), code)
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opStructFieldIntPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendInt(b, ptrToUint64(p), code)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyIntPtr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				b = appendInt(b, ptrToUint64(p), code)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagIntPtr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = appendInt(b, ptrToUint64(p), code)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldIntNPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			for i := 0; i < code.ptrNum-1; i++ {
				if p == 0 {
					break
				}
				p = ptrToPtr(p)
			}
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendInt(b, ptrToUint64(p), code)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldUint:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = appendUint(b, ptrToUint64(ptr+code.offset), code)
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyUint:
			ptr := load(ctxptr, code.headIdx)
			u64 := ptrToUint64(ptr + code.offset)
			v := u64 & code.mask
			if v != 0 {
				b = append(b, code.key...)
				b = appendUint(b, u64, code)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagUint:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = append(b, '"')
			b = appendUint(b, ptrToUint64(ptr+code.offset), code)
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opStructFieldUintPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendUint(b, ptrToUint64(p), code)
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyUintPtr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				b = appendUint(b, ptrToUint64(p), code)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagUintPtr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = appendUint(b, ptrToUint64(p), code)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldFloat32:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = encodeFloat32(b, ptrToFloat32(ptr+code.offset))
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyFloat32:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToFloat32(ptr + code.offset)
			if v != 0 {
				b = append(b, code.key...)
				b = encodeFloat32(b, v)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagFloat32:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = append(b, '"')
			b = encodeFloat32(b, ptrToFloat32(ptr+code.offset))
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opStructFieldFloat32Ptr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeFloat32(b, ptrToFloat32(p))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyFloat32Ptr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				b = encodeFloat32(b, ptrToFloat32(p))
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagFloat32Ptr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = encodeFloat32(b, ptrToFloat32(p))
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldFloat64:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			v := ptrToFloat64(ptr + code.offset)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return nil, errUnsupportedFloat(v)
			}
			b = encodeFloat64(b, v)
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyFloat64:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToFloat64(ptr + code.offset)
			if v != 0 {
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = append(b, code.key...)
				b = encodeFloat64(b, v)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagFloat64:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToFloat64(ptr + code.offset)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return nil, errUnsupportedFloat(v)
			}
			b = append(b, code.key...)
			b = append(b, '"')
			b = encodeFloat64(b, v)
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opStructFieldFloat64Ptr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
				b = encodeComma(b)
				code = code.next
				break
			}
			v := ptrToFloat64(p)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return nil, errUnsupportedFloat(v)
			}
			b = encodeFloat64(b, v)
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyFloat64Ptr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				v := ptrToFloat64(p)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagFloat64Ptr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			b = append(b, code.key...)
			if p == 0 {
				b = encodeNull(b)
			} else {
				v := ptrToFloat64(p)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = append(b, '"')
				b = encodeFloat64(b, v)
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldString:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = encodeNoEscapedString(b, ptrToString(ptr+code.offset))
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyString:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToString(ptr + code.offset)
			if v != "" {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, v)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagString:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			s := ptrToString(ptr + code.offset)
			b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, s)))
			b = encodeComma(b)
			code = code.next
		case opStructFieldStringPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeNoEscapedString(b, ptrToString(p))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyStringPtr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, ptrToString(p))
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagStringPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, ptrToString(p))))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldBool:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = encodeBool(b, ptrToBool(ptr+code.offset))
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyBool:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToBool(ptr + code.offset)
			if v {
				b = append(b, code.key...)
				b = encodeBool(b, v)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagBool:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = append(b, '"')
			b = encodeBool(b, ptrToBool(ptr+code.offset))
			b = append(b, '"')
			b = encodeComma(b)
			code = code.next
		case opStructFieldBoolPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeBool(b, ptrToBool(p))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyBoolPtr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				b = encodeBool(b, ptrToBool(p))
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagBoolPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = encodeBool(b, ptrToBool(p))
				b = append(b, '"')
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldBytes:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = encodeByteSlice(b, ptrToBytes(ptr+code.offset))
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyBytes:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToBytes(ptr + code.offset)
			if len(v) > 0 {
				b = append(b, code.key...)
				b = encodeByteSlice(b, v)
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagBytes:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToBytes(ptr + code.offset)
			b = append(b, code.key...)
			b = encodeByteSlice(b, v)
			b = encodeComma(b)
			code = code.next
		case opStructFieldBytesPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeByteSlice(b, ptrToBytes(p))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyBytesPtr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				b = encodeByteSlice(b, ptrToBytes(p))
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldStringTagBytesPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeByteSlice(b, ptrToBytes(p))
			}
			b = encodeComma(b)
			code = code.next
		case opStructFieldMarshalJSON:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			p := ptr + code.offset
			v := ptrToInterface(code, p)
			bb, err := v.(Marshaler).MarshalJSON()
			if err != nil {
				return nil, errMarshaler(code, err)
			}
			buf := bytes.NewBuffer(b)
			//TODO: we should validate buffer with `compact`
			if err := compact(buf, bb, false); err != nil {
				return nil, err
			}
			b = buf.Bytes()
			b = encodeComma(b)
			code = code.next
		case opStructFieldStringTagMarshalJSON:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := ptrToInterface(code, p)
			bb, err := v.(Marshaler).MarshalJSON()
			if err != nil {
				return nil, errMarshaler(code, err)
			}
			var buf bytes.Buffer
			if err := compact(&buf, bb, false); err != nil {
				return nil, err
			}
			b = append(b, code.key...)
			b = encodeNoEscapedString(b, buf.String())
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyMarshalJSON:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if code.typ.Kind() == reflect.Ptr && code.typ.Elem().Implements(marshalJSONType) {
				p = ptrToPtr(p)
			}
			v := ptrToInterface(code, p)
			if v != nil && p != 0 {
				bb, err := v.(Marshaler).MarshalJSON()
				if err != nil {
					return nil, errMarshaler(code, err)
				}
				b = append(b, code.key...)
				buf := bytes.NewBuffer(b)
				//TODO: we should validate buffer with `compact`
				if err := compact(buf, bb, false); err != nil {
					return nil, err
				}
				b = buf.Bytes()
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldMarshalText:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			p := ptr + code.offset
			v := ptrToInterface(code, p)
			bytes, err := v.(encoding.TextMarshaler).MarshalText()
			if err != nil {
				return nil, errMarshaler(code, err)
			}
			b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
			b = encodeComma(b)
			code = code.next
		case opStructFieldStringTagMarshalText:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := ptrToInterface(code, p)
			bytes, err := v.(encoding.TextMarshaler).MarshalText()
			if err != nil {
				return nil, errMarshaler(code, err)
			}
			b = append(b, code.key...)
			b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
			b = encodeComma(b)
			code = code.next
		case opStructFieldOmitEmptyMarshalText:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := ptrToInterface(code, p)
			if v != nil {
				bytes, err := v.(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return nil, errMarshaler(code, err)
				}
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
				b = encodeComma(b)
			}
			code = code.next
		case opStructFieldArray:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldOmitEmptyArray:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			array := ptrToSlice(p)
			if p == 0 || uintptr(array.data) == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldSlice:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldOmitEmptySlice:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			slice := ptrToSlice(p)
			if p == 0 || uintptr(slice.data) == 0 {
				code = code.nextField
			} else {
				b = append(b, code.key...)
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldMap:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldOmitEmptyMap:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 {
				code = code.nextField
			} else {
				mlen := maplen(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
				if mlen == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					code = code.next
					store(ctxptr, code.idx, p)
				}
			}
		case opStructFieldMapLoad:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldOmitEmptyMapLoad:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 {
				code = code.nextField
			} else {
				mlen := maplen(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
				if mlen == 0 {
					code = code.nextField
				} else {
					b = append(b, code.key...)
					code = code.next
					store(ctxptr, code.idx, p)
				}
			}
		case opStructFieldStruct:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructEnd:
			last := len(b) - 1
			if b[last] == ',' {
				b[last] = '}'
			} else {
				b = append(b, '}')
			}
			b = encodeComma(b)
			code = code.next
		case opStructAnonymousEnd:
			code = code.next
		case opStructEndInt:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = appendInt(b, ptrToUint64(ptr+code.offset), code)
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyInt:
			ptr := load(ctxptr, code.headIdx)
			u64 := ptrToUint64(ptr + code.offset)
			v := u64 & code.mask
			if v != 0 {
				b = append(b, code.key...)
				b = appendInt(b, u64, code)
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagInt:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = append(b, '"')
			b = appendInt(b, ptrToUint64(ptr+code.offset), code)
			b = append(b, '"')
			b = appendStructEnd(b)
			code = code.next
		case opStructEndIntPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendInt(b, ptrToUint64(p), code)
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyIntPtr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				b = appendInt(b, ptrToUint64(p), code)
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagIntPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = appendInt(b, ptrToUint64(p), code)
				b = append(b, '"')
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndIntNPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			for i := 0; i < code.ptrNum-1; i++ {
				if p == 0 {
					break
				}
				p = ptrToPtr(p)
			}
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendInt(b, ptrToUint64(p), code)
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndUint:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = appendUint(b, ptrToUint64(ptr+code.offset), code)
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyUint:
			ptr := load(ctxptr, code.headIdx)
			u64 := ptrToUint64(ptr + code.offset)
			v := u64 & code.mask
			if v != 0 {
				b = append(b, code.key...)
				b = appendUint(b, u64, code)
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagUint:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = append(b, '"')
			b = appendUint(b, ptrToUint64(ptr+code.offset), code)
			b = append(b, '"')
			b = appendStructEnd(b)
			code = code.next
		case opStructEndUintPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendUint(b, ptrToUint64(p), code)
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyUintPtr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				b = appendUint(b, ptrToUint64(p), code)
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagUintPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = appendUint(b, ptrToUint64(p), code)
				b = append(b, '"')
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndUintNPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			for i := 0; i < code.ptrNum-1; i++ {
				if p == 0 {
					break
				}
				p = ptrToPtr(p)
			}
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = appendUint(b, ptrToUint64(p), code)
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndFloat32:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = encodeFloat32(b, ptrToFloat32(ptr+code.offset))
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyFloat32:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToFloat32(ptr + code.offset)
			if v != 0 {
				b = append(b, code.key...)
				b = encodeFloat32(b, v)
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagFloat32:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = append(b, '"')
			b = encodeFloat32(b, ptrToFloat32(ptr+code.offset))
			b = append(b, '"')
			b = appendStructEnd(b)
			code = code.next
		case opStructEndFloat32Ptr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeFloat32(b, ptrToFloat32(p))
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyFloat32Ptr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				b = encodeFloat32(b, ptrToFloat32(p))
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagFloat32Ptr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = encodeFloat32(b, ptrToFloat32(p))
				b = append(b, '"')
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndFloat32NPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			for i := 0; i < code.ptrNum-1; i++ {
				if p == 0 {
					break
				}
				p = ptrToPtr(p)
			}
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeFloat32(b, ptrToFloat32(p))
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndFloat64:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			v := ptrToFloat64(ptr + code.offset)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return nil, errUnsupportedFloat(v)
			}
			b = encodeFloat64(b, v)
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyFloat64:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToFloat64(ptr + code.offset)
			if v != 0 {
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = append(b, code.key...)
				b = encodeFloat64(b, v)
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagFloat64:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToFloat64(ptr + code.offset)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return nil, errUnsupportedFloat(v)
			}
			b = append(b, code.key...)
			b = append(b, '"')
			b = encodeFloat64(b, v)
			b = append(b, '"')
			b = appendStructEnd(b)
			code = code.next
		case opStructEndFloat64Ptr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
				b = appendStructEnd(b)
				code = code.next
				break
			}
			v := ptrToFloat64(p)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return nil, errUnsupportedFloat(v)
			}
			b = encodeFloat64(b, v)
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyFloat64Ptr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				v := ptrToFloat64(p)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagFloat64Ptr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				v := ptrToFloat64(p)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
				b = append(b, '"')
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndFloat64NPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			for i := 0; i < code.ptrNum-1; i++ {
				if p == 0 {
					break
				}
				p = ptrToPtr(p)
			}
			if p == 0 {
				b = encodeNull(b)
			} else {
				v := ptrToFloat64(p)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return nil, errUnsupportedFloat(v)
				}
				b = encodeFloat64(b, v)
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndString:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = encodeNoEscapedString(b, ptrToString(ptr+code.offset))
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyString:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToString(ptr + code.offset)
			if v != "" {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, v)
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagString:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			s := ptrToString(ptr + code.offset)
			b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, s)))
			b = appendStructEnd(b)
			code = code.next
		case opStructEndStringPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeNoEscapedString(b, ptrToString(p))
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyStringPtr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, ptrToString(p))
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagStringPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				v := ptrToString(p)
				b = encodeNoEscapedString(b, string(encodeNoEscapedString([]byte{}, v)))
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndStringNPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			for i := 0; i < code.ptrNum-1; i++ {
				if p == 0 {
					break
				}
				p = ptrToPtr(p)
			}
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeNoEscapedString(b, ptrToString(p))
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndBool:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = encodeBool(b, ptrToBool(ptr+code.offset))
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyBool:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToBool(ptr + code.offset)
			if v {
				b = append(b, code.key...)
				b = encodeBool(b, v)
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagBool:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = append(b, '"')
			b = encodeBool(b, ptrToBool(ptr+code.offset))
			b = append(b, '"')
			b = appendStructEnd(b)
			code = code.next
		case opStructEndBoolPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = encodeBool(b, ptrToBool(p))
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyBoolPtr:
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p != 0 {
				b = append(b, code.key...)
				b = encodeBool(b, ptrToBool(p))
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagBoolPtr:
			b = append(b, code.key...)
			ptr := load(ctxptr, code.headIdx)
			p := ptrToPtr(ptr + code.offset)
			if p == 0 {
				b = encodeNull(b)
			} else {
				b = append(b, '"')
				b = encodeBool(b, ptrToBool(p))
				b = append(b, '"')
			}
			b = appendStructEnd(b)
			code = code.next
		case opStructEndBytes:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			b = encodeByteSlice(b, ptrToBytes(ptr+code.offset))
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyBytes:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToBytes(ptr + code.offset)
			if len(v) > 0 {
				b = append(b, code.key...)
				b = encodeByteSlice(b, v)
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagBytes:
			ptr := load(ctxptr, code.headIdx)
			v := ptrToBytes(ptr + code.offset)
			b = append(b, code.key...)
			b = encodeByteSlice(b, v)
			b = appendStructEnd(b)
			code = code.next
		case opStructEndMarshalJSON:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			p := ptr + code.offset
			v := ptrToInterface(code, p)
			bb, err := v.(Marshaler).MarshalJSON()
			if err != nil {
				return nil, errMarshaler(code, err)
			}
			buf := bytes.NewBuffer(b)
			//TODO: we should validate buffer with `compact`
			if err := compact(buf, bb, false); err != nil {
				return nil, err
			}
			b = buf.Bytes()
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyMarshalJSON:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := ptrToInterface(code, p)
			if v != nil && (code.typ.Kind() != reflect.Ptr || ptrToPtr(p) != 0) {
				bb, err := v.(Marshaler).MarshalJSON()
				if err != nil {
					return nil, errMarshaler(code, err)
				}
				b = append(b, code.key...)
				buf := bytes.NewBuffer(b)
				//TODO: we should validate buffer with `compact`
				if err := compact(buf, bb, false); err != nil {
					return nil, err
				}
				b = buf.Bytes()
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagMarshalJSON:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := ptrToInterface(code, p)
			bb, err := v.(Marshaler).MarshalJSON()
			if err != nil {
				return nil, errMarshaler(code, err)
			}
			var buf bytes.Buffer
			if err := compact(&buf, bb, false); err != nil {
				return nil, err
			}
			b = append(b, code.key...)
			b = encodeNoEscapedString(b, buf.String())
			b = appendStructEnd(b)
			code = code.next
		case opStructEndMarshalText:
			ptr := load(ctxptr, code.headIdx)
			b = append(b, code.key...)
			p := ptr + code.offset
			v := ptrToInterface(code, p)
			bytes, err := v.(encoding.TextMarshaler).MarshalText()
			if err != nil {
				return nil, errMarshaler(code, err)
			}
			b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
			b = appendStructEnd(b)
			code = code.next
		case opStructEndOmitEmptyMarshalText:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := ptrToInterface(code, p)
			if v != nil {
				bytes, err := v.(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return nil, errMarshaler(code, err)
				}
				b = append(b, code.key...)
				b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
				b = appendStructEnd(b)
			} else {
				last := len(b) - 1
				if b[last] == ',' {
					b[last] = '}'
					b = encodeComma(b)
				} else {
					b = appendStructEnd(b)
				}
			}
			code = code.next
		case opStructEndStringTagMarshalText:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := ptrToInterface(code, p)
			bytes, err := v.(encoding.TextMarshaler).MarshalText()
			if err != nil {
				return nil, errMarshaler(code, err)
			}
			b = append(b, code.key...)
			b = encodeNoEscapedString(b, *(*string)(unsafe.Pointer(&bytes)))
			b = appendStructEnd(b)
			code = code.next
		case opEnd:
			goto END
		}
	}
END:
	return b, nil
}
