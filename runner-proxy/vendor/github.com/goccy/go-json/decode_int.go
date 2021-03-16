package json

import (
	"fmt"
	"reflect"
	"unsafe"
)

type intDecoder struct {
	typ        *rtype
	kind       reflect.Kind
	op         func(unsafe.Pointer, int64)
	structName string
	fieldName  string
}

func newIntDecoder(typ *rtype, structName, fieldName string, op func(unsafe.Pointer, int64)) *intDecoder {
	return &intDecoder{
		typ:        typ,
		kind:       typ.Kind(),
		op:         op,
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *intDecoder) typeError(buf []byte, offset int64) *UnmarshalTypeError {
	return &UnmarshalTypeError{
		Value:  fmt.Sprintf("number %s", string(buf)),
		Type:   rtype2type(d.typ),
		Struct: d.structName,
		Field:  d.fieldName,
		Offset: offset,
	}
}

var (
	pow10i64 = [...]int64{
		1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
	}
)

func (d *intDecoder) parseInt(b []byte) int64 {
	isNegative := false
	if b[0] == '-' {
		b = b[1:]
		isNegative = true
	}
	maxDigit := len(b)
	sum := int64(0)
	for i := 0; i < maxDigit; i++ {
		c := int64(b[i]) - 48
		digitValue := pow10i64[maxDigit-i-1]
		sum += c * digitValue
	}
	if isNegative {
		return -1 * sum
	}
	return sum
}

var (
	numTable = [256]bool{
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
	}
)

func (d *intDecoder) decodeStreamByte(s *stream) ([]byte, error) {
	for {
		switch s.char() {
		case ' ', '\n', '\t', '\r':
			s.cursor++
			continue
		case '-':
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
			if len(num) < 2 {
				goto ERROR
			}
			return num, nil
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
			goto ERROR
		default:
			return nil, d.typeError([]byte{s.char()}, s.totalOffset())
		}
	}
ERROR:
	return nil, errUnexpectedEndOfJSON("number(integer)", s.totalOffset())
}

func (d *intDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	b := (*sliceHeader)(unsafe.Pointer(&buf)).data
	for {
		switch char(b, cursor) {
		case ' ', '\n', '\t', '\r':
			cursor++
			continue
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			start := cursor
			cursor++
		LOOP:
			if numTable[char(b, cursor)] {
				cursor++
				goto LOOP
			}
			num := buf[start:cursor]
			return num, cursor, nil
		case 'n':
			buflen := int64(len(buf))
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
			return nil, 0, d.typeError([]byte{char(b, cursor)}, cursor)
		}
	}
}

func (d *intDecoder) decodeStream(s *stream, depth int64, p unsafe.Pointer) error {
	bytes, err := d.decodeStreamByte(s)
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	i64 := d.parseInt(bytes)
	switch d.kind {
	case reflect.Int8:
		if i64 <= -1*(1<<7) || (1<<7) <= i64 {
			return d.typeError(bytes, s.totalOffset())
		}
	case reflect.Int16:
		if i64 <= -1*(1<<15) || (1<<15) <= i64 {
			return d.typeError(bytes, s.totalOffset())
		}
	case reflect.Int32:
		if i64 <= -1*(1<<31) || (1<<31) <= i64 {
			return d.typeError(bytes, s.totalOffset())
		}
	}
	d.op(p, i64)
	s.reset()
	return nil
}

func (d *intDecoder) decode(buf []byte, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeByte(buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c
	i64 := d.parseInt(bytes)
	switch d.kind {
	case reflect.Int8:
		if i64 <= -1*(1<<7) || (1<<7) <= i64 {
			return 0, d.typeError(bytes, cursor)
		}
	case reflect.Int16:
		if i64 <= -1*(1<<15) || (1<<15) <= i64 {
			return 0, d.typeError(bytes, cursor)
		}
	case reflect.Int32:
		if i64 <= -1*(1<<31) || (1<<31) <= i64 {
			return 0, d.typeError(bytes, cursor)
		}
	}
	d.op(p, i64)
	return cursor, nil
}
