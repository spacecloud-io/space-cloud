package json

import (
	"fmt"
	"reflect"
	"unsafe"
)

type uintDecoder struct {
	typ        *rtype
	kind       reflect.Kind
	op         func(unsafe.Pointer, uint64)
	structName string
	fieldName  string
}

func newUintDecoder(typ *rtype, structName, fieldName string, op func(unsafe.Pointer, uint64)) *uintDecoder {
	return &uintDecoder{
		typ:        typ,
		kind:       typ.Kind(),
		op:         op,
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *uintDecoder) typeError(buf []byte, offset int64) *UnmarshalTypeError {
	return &UnmarshalTypeError{
		Value:  fmt.Sprintf("number %s", string(buf)),
		Type:   rtype2type(d.typ),
		Offset: offset,
	}
}

var pow10u64 = [...]uint64{
	1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
}

func (d *uintDecoder) parseUint(b []byte) uint64 {
	maxDigit := len(b)
	sum := uint64(0)
	for i := 0; i < maxDigit; i++ {
		c := uint64(b[i]) - 48
		digitValue := pow10u64[maxDigit-i-1]
		sum += c * digitValue
	}
	return sum
}

func (d *uintDecoder) decodeStreamByte(s *stream) ([]byte, error) {
	for {
		switch s.char() {
		case ' ', '\n', '\t', '\r':
			s.cursor++
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			start := s.cursor
			for {
				s.cursor++
				if numTable[s.char()] {
					continue
				} else if s.char() == nul {
					if s.read() {
						s.cursor-- // for retry current character
						continue
					}
				}
				break
			}
			num := s.buf[start:s.cursor]
			return num, nil
		case 'n':
			if err := nullBytes(s); err != nil {
				return nil, err
			}
			return nil, nil
		case nul:
			if s.read() {
				continue
			}
		default:
			return nil, d.typeError([]byte{s.char()}, s.totalOffset())
		}
		break
	}
	return nil, errUnexpectedEndOfJSON("number(unsigned integer)", s.totalOffset())
}

func (d *uintDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	buflen := int64(len(buf))
	for ; cursor < buflen; cursor++ {
		switch buf[cursor] {
		case ' ', '\n', '\t', '\r':
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			start := cursor
			cursor++
			for ; cursor < buflen; cursor++ {
				tk := int(buf[cursor])
				if int('0') <= tk && tk <= int('9') {
					continue
				}
				break
			}
			num := buf[start:cursor]
			return num, cursor, nil
		case 'n':
			if cursor+3 >= buflen {
				return nil, 0, errUnexpectedEndOfJSON("null", cursor)
			}
			if buf[cursor+1] != 'u' {
				return nil, 0, errInvalidCharacter(buf[cursor+1], "null", cursor)
			}
			if buf[cursor+2] != 'l' {
				return nil, 0, errInvalidCharacter(buf[cursor+2], "null", cursor)
			}
			if buf[cursor+3] != 'l' {
				return nil, 0, errInvalidCharacter(buf[cursor+3], "null", cursor)
			}
			cursor += 4
			return nil, cursor, nil
		default:
			return nil, 0, d.typeError([]byte{buf[cursor]}, cursor)
		}
	}
	return nil, 0, errUnexpectedEndOfJSON("number(unsigned integer)", cursor)
}

func (d *uintDecoder) decodeStream(s *stream, depth int64, p unsafe.Pointer) error {
	bytes, err := d.decodeStreamByte(s)
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	u64 := d.parseUint(bytes)
	switch d.kind {
	case reflect.Uint8:
		if (1 << 8) <= u64 {
			return d.typeError(bytes, s.totalOffset())
		}
	case reflect.Uint16:
		if (1 << 16) <= u64 {
			return d.typeError(bytes, s.totalOffset())
		}
	case reflect.Uint32:
		if (1 << 32) <= u64 {
			return d.typeError(bytes, s.totalOffset())
		}
	}
	d.op(p, u64)
	return nil
}

func (d *uintDecoder) decode(buf []byte, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeByte(buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c
	u64 := d.parseUint(bytes)
	switch d.kind {
	case reflect.Uint8:
		if (1 << 8) <= u64 {
			return 0, d.typeError(bytes, cursor)
		}
	case reflect.Uint16:
		if (1 << 16) <= u64 {
			return 0, d.typeError(bytes, cursor)
		}
	case reflect.Uint32:
		if (1 << 32) <= u64 {
			return 0, d.typeError(bytes, cursor)
		}
	}
	d.op(p, u64)
	return cursor, nil
}
