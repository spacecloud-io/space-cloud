package json

import (
	"strconv"
	"unsafe"
)

type floatDecoder struct {
	op         func(unsafe.Pointer, float64)
	structName string
	fieldName  string
}

func newFloatDecoder(structName, fieldName string, op func(unsafe.Pointer, float64)) *floatDecoder {
	return &floatDecoder{op: op, structName: structName, fieldName: fieldName}
}

var (
	floatTable = [256]bool{
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
		'.': true,
		'e': true,
		'E': true,
		'+': true,
		'-': true,
	}

	validEndNumberChar = [256]bool{
		nul:  true,
		' ':  true,
		'\t': true,
		'\r': true,
		'\n': true,
		',':  true,
		':':  true,
		'}':  true,
		']':  true,
	}
)

func floatBytes(s *stream) []byte {
	start := s.cursor
	for {
		s.cursor++
		if floatTable[s.char()] {
			continue
		} else if s.char() == nul {
			if s.read() {
				s.cursor-- // for retry current character
				continue
			}
		}
		break
	}
	return s.buf[start:s.cursor]
}

func (d *floatDecoder) decodeStreamByte(s *stream) ([]byte, error) {
	for {
		switch s.char() {
		case ' ', '\n', '\t', '\r':
			s.cursor++
			continue
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return floatBytes(s), nil
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
			goto ERROR
		}
	}
ERROR:
	return nil, errUnexpectedEndOfJSON("float", s.totalOffset())
}

func (d *floatDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	buflen := int64(len(buf))
	for ; cursor < buflen; cursor++ {
		switch buf[cursor] {
		case ' ', '\n', '\t', '\r':
			continue
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			start := cursor
			cursor++
			for ; cursor < buflen; cursor++ {
				if floatTable[buf[cursor]] {
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
			return nil, 0, errUnexpectedEndOfJSON("float", cursor)
		}
	}
	return nil, 0, errUnexpectedEndOfJSON("float", cursor)
}

func (d *floatDecoder) decodeStream(s *stream, depth int64, p unsafe.Pointer) error {
	bytes, err := d.decodeStreamByte(s)
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	str := *(*string)(unsafe.Pointer(&bytes))
	f64, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return &SyntaxError{msg: err.Error(), Offset: s.totalOffset()}
	}
	d.op(p, f64)
	return nil
}

func (d *floatDecoder) decode(buf []byte, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeByte(buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c
	if !validEndNumberChar[buf[cursor]] {
		return 0, errUnexpectedEndOfJSON("float", cursor)
	}
	s := *(*string)(unsafe.Pointer(&bytes))
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, &SyntaxError{msg: err.Error(), Offset: cursor}
	}
	d.op(p, f64)
	return cursor, nil
}
