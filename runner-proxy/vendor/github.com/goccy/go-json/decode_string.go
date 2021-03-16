package json

import (
	"reflect"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

type stringDecoder struct {
	structName string
	fieldName  string
}

func newStringDecoder(structName, fieldName string) *stringDecoder {
	return &stringDecoder{
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *stringDecoder) errUnmarshalType(typeName string, offset int64) *UnmarshalTypeError {
	return &UnmarshalTypeError{
		Value:  typeName,
		Type:   reflect.TypeOf(""),
		Offset: offset,
		Struct: d.structName,
		Field:  d.fieldName,
	}
}

func (d *stringDecoder) decodeStream(s *stream, depth int64, p unsafe.Pointer) error {
	bytes, err := d.decodeStreamByte(s)
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	**(**string)(unsafe.Pointer(&p)) = *(*string)(unsafe.Pointer(&bytes))
	s.reset()
	return nil
}

func (d *stringDecoder) decode(buf []byte, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeByte(buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c
	**(**string)(unsafe.Pointer(&p)) = *(*string)(unsafe.Pointer(&bytes))
	return cursor, nil
}

var (
	hexToInt = [256]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'A': 10,
		'B': 11,
		'C': 12,
		'D': 13,
		'E': 14,
		'F': 15,
		'a': 10,
		'b': 11,
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
	}
)

func unicodeToRune(code []byte) rune {
	var r rune
	for i := 0; i < len(code); i++ {
		r = r*16 + rune(hexToInt[code[i]])
	}
	return r
}

func decodeEscapeString(s *stream) error {
	s.cursor++
RETRY:
	switch s.buf[s.cursor] {
	case '"':
		s.buf[s.cursor] = '"'
	case '\\':
		s.buf[s.cursor] = '\\'
	case '/':
		s.buf[s.cursor] = '/'
	case 'b':
		s.buf[s.cursor] = '\b'
	case 'f':
		s.buf[s.cursor] = '\f'
	case 'n':
		s.buf[s.cursor] = '\n'
	case 'r':
		s.buf[s.cursor] = '\r'
	case 't':
		s.buf[s.cursor] = '\t'
	case 'u':
		if s.cursor+5 >= s.length {
			if !s.read() {
				return errInvalidCharacter(s.char(), "escaped string", s.totalOffset())
			}
		}
		r := unicodeToRune(s.buf[s.cursor+1 : s.cursor+5])
		if utf16.IsSurrogate(r) {
			if s.cursor+11 >= s.length || s.buf[s.cursor+5] != '\\' || s.buf[s.cursor+6] != 'u' {
				r = unicode.ReplacementChar
				unicode := []byte(string(r))
				s.buf = append(append(s.buf[:s.cursor-1], unicode...), s.buf[s.cursor+5:]...)
				s.cursor = s.cursor - 2 + int64(len(unicode))
				return nil
			}
			r2 := unicodeToRune(s.buf[s.cursor+7 : s.cursor+11])
			if r := utf16.DecodeRune(r, r2); r != unicode.ReplacementChar {
				// valid surrogate pair
				unicode := []byte(string(r))
				s.buf = append(append(s.buf[:s.cursor-1], unicode...), s.buf[s.cursor+11:]...)
				s.cursor = s.cursor - 2 + int64(len(unicode))
			} else {
				unicode := []byte(string(r))
				s.buf = append(append(s.buf[:s.cursor-1], unicode...), s.buf[s.cursor+5:]...)
				s.cursor = s.cursor - 2 + int64(len(unicode))
			}
		} else {
			unicode := []byte(string(r))
			s.buf = append(append(s.buf[:s.cursor-1], unicode...), s.buf[s.cursor+5:]...)
			s.cursor = s.cursor - 2 + int64(len(unicode))
		}
		return nil
	case nul:
		if !s.read() {
			return errInvalidCharacter(s.char(), "escaped string", s.totalOffset())
		}
		goto RETRY
	default:
		return errUnexpectedEndOfJSON("string", s.totalOffset())
	}
	s.buf = append(s.buf[:s.cursor-1], s.buf[s.cursor:]...)
	s.cursor--
	return nil
}

//nolint:deadcode,unused
func appendCoerceInvalidUTF8(b []byte, s []byte) []byte {
	c := [4]byte{}

	for _, r := range string(s) {
		b = append(b, c[:utf8.EncodeRune(c[:], r)]...)
	}

	return b
}

func stringBytes(s *stream) ([]byte, error) {
	buf, cursor, p := s.stat()

	cursor++ // skip double quote char
	start := cursor
	for {
		switch char(p, cursor) {
		case '\\':
			s.cursor = cursor
			if err := decodeEscapeString(s); err != nil {
				return nil, err
			}
			buf, cursor, p = s.stat()
		case '"':
			literal := buf[start:cursor]
			// TODO: this flow is so slow sequence.
			// literal = appendCoerceInvalidUTF8(make([]byte, 0, len(literal)), literal)
			cursor++
			s.cursor = cursor
			return literal, nil
		case nul:
			s.cursor = cursor
			if s.read() {
				buf, cursor, p = s.stat()
				continue
			}
			goto ERROR
		}
		cursor++
	}
ERROR:
	return nil, errUnexpectedEndOfJSON("string", s.totalOffset())
}

func nullBytes(s *stream) error {
	if s.cursor+3 >= s.length {
		if !s.read() {
			return errInvalidCharacter(s.char(), "null", s.totalOffset())
		}
	}
	s.cursor++
	if s.char() != 'u' {
		return errInvalidCharacter(s.char(), "null", s.totalOffset())
	}
	s.cursor++
	if s.char() != 'l' {
		return errInvalidCharacter(s.char(), "null", s.totalOffset())
	}
	s.cursor++
	if s.char() != 'l' {
		return errInvalidCharacter(s.char(), "null", s.totalOffset())
	}
	s.cursor++
	return nil
}

func (d *stringDecoder) decodeStreamByte(s *stream) ([]byte, error) {
	for {
		switch s.char() {
		case ' ', '\n', '\t', '\r':
			s.cursor++
			continue
		case '[':
			return nil, d.errUnmarshalType("array", s.totalOffset())
		case '{':
			return nil, d.errUnmarshalType("object", s.totalOffset())
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return nil, d.errUnmarshalType("number", s.totalOffset())
		case '"':
			return stringBytes(s)
		case 'n':
			if err := nullBytes(s); err != nil {
				return nil, err
			}
			return nil, nil
		case nul:
			if s.read() {
				continue
			}
		}
		break
	}
	return nil, errNotAtBeginningOfValue(s.totalOffset())
}

func (d *stringDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	for {
		switch buf[cursor] {
		case ' ', '\n', '\t', '\r':
			cursor++
		case '[':
			return nil, 0, d.errUnmarshalType("array", cursor)
		case '{':
			return nil, 0, d.errUnmarshalType("object", cursor)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return nil, 0, d.errUnmarshalType("number", cursor)
		case '"':
			cursor++
			start := cursor
			b := (*sliceHeader)(unsafe.Pointer(&buf)).data
			for {
				switch char(b, cursor) {
				case '\\':
					cursor++
					switch char(b, cursor) {
					case '"':
						buf[cursor] = '"'
						buf = append(buf[:cursor-1], buf[cursor:]...)
					case '\\':
						buf[cursor] = '\\'
						buf = append(buf[:cursor-1], buf[cursor:]...)
					case '/':
						buf[cursor] = '/'
						buf = append(buf[:cursor-1], buf[cursor:]...)
					case 'b':
						buf[cursor] = '\b'
						buf = append(buf[:cursor-1], buf[cursor:]...)
					case 'f':
						buf[cursor] = '\f'
						buf = append(buf[:cursor-1], buf[cursor:]...)
					case 'n':
						buf[cursor] = '\n'
						buf = append(buf[:cursor-1], buf[cursor:]...)
					case 'r':
						buf[cursor] = '\r'
						buf = append(buf[:cursor-1], buf[cursor:]...)
					case 't':
						buf[cursor] = '\t'
						buf = append(buf[:cursor-1], buf[cursor:]...)
					case 'u':
						buflen := int64(len(buf))
						if cursor+5 >= buflen {
							return nil, 0, errUnexpectedEndOfJSON("escaped string", cursor)
						}
						code := unicodeToRune(buf[cursor+1 : cursor+5])
						unicode := []byte(string(code))
						buf = append(append(buf[:cursor-1], unicode...), buf[cursor+5:]...)
					default:
						return nil, 0, errUnexpectedEndOfJSON("escaped string", cursor)
					}
					continue
				case '"':
					literal := buf[start:cursor]
					cursor++
					return literal, cursor, nil
				case nul:
					return nil, 0, errUnexpectedEndOfJSON("string", cursor)
				}
				cursor++
			}
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
			goto ERROR
		}
	}
ERROR:
	return nil, 0, errNotAtBeginningOfValue(cursor)
}
