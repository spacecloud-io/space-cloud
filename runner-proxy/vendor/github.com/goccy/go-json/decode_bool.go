package json

import (
	"unsafe"
)

type boolDecoder struct {
	structName string
	fieldName  string
}

func newBoolDecoder(structName, fieldName string) *boolDecoder {
	return &boolDecoder{structName: structName, fieldName: fieldName}
}

func trueBytes(s *stream) error {
	if s.cursor+3 >= s.length {
		if !s.read() {
			return errInvalidCharacter(s.char(), "bool(true)", s.totalOffset())
		}
	}
	s.cursor++
	if s.char() != 'r' {
		return errInvalidCharacter(s.char(), "bool(true)", s.totalOffset())
	}
	s.cursor++
	if s.char() != 'u' {
		return errInvalidCharacter(s.char(), "bool(true)", s.totalOffset())
	}
	s.cursor++
	if s.char() != 'e' {
		return errInvalidCharacter(s.char(), "bool(true)", s.totalOffset())
	}
	s.cursor++
	return nil
}

func falseBytes(s *stream) error {
	if s.cursor+4 >= s.length {
		if !s.read() {
			return errInvalidCharacter(s.char(), "bool(false)", s.totalOffset())
		}
	}
	s.cursor++
	if s.char() != 'a' {
		return errInvalidCharacter(s.char(), "bool(false)", s.totalOffset())
	}
	s.cursor++
	if s.char() != 'l' {
		return errInvalidCharacter(s.char(), "bool(false)", s.totalOffset())
	}
	s.cursor++
	if s.char() != 's' {
		return errInvalidCharacter(s.char(), "bool(false)", s.totalOffset())
	}
	s.cursor++
	if s.char() != 'e' {
		return errInvalidCharacter(s.char(), "bool(false)", s.totalOffset())
	}
	s.cursor++
	return nil
}

func (d *boolDecoder) decodeStream(s *stream, depth int64, p unsafe.Pointer) error {
	s.skipWhiteSpace()
	for {
		switch s.char() {
		case 't':
			if err := trueBytes(s); err != nil {
				return err
			}
			**(**bool)(unsafe.Pointer(&p)) = true
			return nil
		case 'f':
			if err := falseBytes(s); err != nil {
				return err
			}
			**(**bool)(unsafe.Pointer(&p)) = false
			return nil
		case 'n':
			if err := nullBytes(s); err != nil {
				return err
			}
			return nil
		case nul:
			if s.read() {
				continue
			}
			goto ERROR
		}
		break
	}
ERROR:
	return errUnexpectedEndOfJSON("bool", s.totalOffset())
}

func (d *boolDecoder) decode(buf []byte, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buflen := int64(len(buf))
	cursor = skipWhiteSpace(buf, cursor)
	switch buf[cursor] {
	case 't':
		if cursor+3 >= buflen {
			return 0, errUnexpectedEndOfJSON("bool(true)", cursor)
		}
		if buf[cursor+1] != 'r' {
			return 0, errInvalidCharacter(buf[cursor+1], "bool(true)", cursor)
		}
		if buf[cursor+2] != 'u' {
			return 0, errInvalidCharacter(buf[cursor+2], "bool(true)", cursor)
		}
		if buf[cursor+3] != 'e' {
			return 0, errInvalidCharacter(buf[cursor+3], "bool(true)", cursor)
		}
		cursor += 4
		**(**bool)(unsafe.Pointer(&p)) = true
		return cursor, nil
	case 'f':
		if cursor+4 >= buflen {
			return 0, errUnexpectedEndOfJSON("bool(false)", cursor)
		}
		if buf[cursor+1] != 'a' {
			return 0, errInvalidCharacter(buf[cursor+1], "bool(false)", cursor)
		}
		if buf[cursor+2] != 'l' {
			return 0, errInvalidCharacter(buf[cursor+2], "bool(false)", cursor)
		}
		if buf[cursor+3] != 's' {
			return 0, errInvalidCharacter(buf[cursor+3], "bool(false)", cursor)
		}
		if buf[cursor+4] != 'e' {
			return 0, errInvalidCharacter(buf[cursor+4], "bool(false)", cursor)
		}
		cursor += 5
		**(**bool)(unsafe.Pointer(&p)) = false
		return cursor, nil
	case 'n':
		if cursor+3 >= buflen {
			return 0, errUnexpectedEndOfJSON("null", cursor)
		}
		if buf[cursor+1] != 'u' {
			return 0, errInvalidCharacter(buf[cursor+1], "null", cursor)
		}
		if buf[cursor+2] != 'l' {
			return 0, errInvalidCharacter(buf[cursor+2], "null", cursor)
		}
		if buf[cursor+3] != 'l' {
			return 0, errInvalidCharacter(buf[cursor+3], "null", cursor)
		}
		cursor += 4
		return cursor, nil
	}
	return 0, errUnexpectedEndOfJSON("bool", cursor)
}
