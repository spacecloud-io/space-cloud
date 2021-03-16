package json

import (
	"encoding/base64"
	"unsafe"
)

type bytesDecoder struct {
	typ          *rtype
	sliceDecoder decoder
	structName   string
	fieldName    string
}

func byteUnmarshalerSliceDecoder(typ *rtype, structName string, fieldName string) decoder {
	var unmarshalDecoder decoder
	switch {
	case rtype_ptrTo(typ).Implements(unmarshalJSONType):
		unmarshalDecoder = newUnmarshalJSONDecoder(rtype_ptrTo(typ), structName, fieldName)
	case rtype_ptrTo(typ).Implements(unmarshalTextType):
		unmarshalDecoder = newUnmarshalTextDecoder(rtype_ptrTo(typ), structName, fieldName)
	}
	if unmarshalDecoder == nil {
		return nil
	}
	return newSliceDecoder(unmarshalDecoder, typ, 1, structName, fieldName)
}

func newBytesDecoder(typ *rtype, structName string, fieldName string) *bytesDecoder {
	return &bytesDecoder{
		typ:          typ,
		sliceDecoder: byteUnmarshalerSliceDecoder(typ, structName, fieldName),
		structName:   structName,
		fieldName:    fieldName,
	}
}

func (d *bytesDecoder) decodeStream(s *stream, depth int64, p unsafe.Pointer) error {
	bytes, err := d.decodeStreamBinary(s, depth, p)
	if err != nil {
		return err
	}
	if bytes == nil {
		s.reset()
		return nil
	}
	decodedLen := base64.StdEncoding.DecodedLen(len(bytes))
	buf := make([]byte, decodedLen)
	if _, err := base64.StdEncoding.Decode(buf, bytes); err != nil {
		return err
	}
	*(*[]byte)(p) = buf
	s.reset()
	return nil
}

func (d *bytesDecoder) decode(buf []byte, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeBinary(buf, cursor, depth, p)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c
	decodedLen := base64.StdEncoding.DecodedLen(len(bytes))
	b := make([]byte, decodedLen)
	n, err := base64.StdEncoding.Decode(b, bytes)
	if err != nil {
		return 0, err
	}
	*(*[]byte)(p) = b[:n]
	return cursor, nil
}

func binaryBytes(s *stream) ([]byte, error) {
	s.cursor++
	start := s.cursor
	for {
		switch s.char() {
		case '"':
			literal := s.buf[start:s.cursor]
			s.cursor++
			return literal, nil
		case nul:
			if s.read() {
				continue
			}
			goto ERROR
		}
		s.cursor++
	}
ERROR:
	return nil, errUnexpectedEndOfJSON("[]byte", s.totalOffset())
}

func (d *bytesDecoder) decodeStreamBinary(s *stream, depth int64, p unsafe.Pointer) ([]byte, error) {
	for {
		switch s.char() {
		case ' ', '\n', '\t', '\r':
			s.cursor++
			continue
		case '"':
			return binaryBytes(s)
		case 'n':
			if err := nullBytes(s); err != nil {
				return nil, err
			}
			return nil, nil
		case '[':
			if d.sliceDecoder == nil {
				return nil, &UnmarshalTypeError{
					Type:   rtype2type(d.typ),
					Offset: s.totalOffset(),
				}
			}
			if err := d.sliceDecoder.decodeStream(s, depth, p); err != nil {
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

func (d *bytesDecoder) decodeBinary(buf []byte, cursor, depth int64, p unsafe.Pointer) ([]byte, int64, error) {
	for {
		switch buf[cursor] {
		case ' ', '\n', '\t', '\r':
			cursor++
		case '"':
			cursor++
			start := cursor
			for {
				switch buf[cursor] {
				case '"':
					literal := buf[start:cursor]
					cursor++
					return literal, cursor, nil
				case nul:
					return nil, 0, errUnexpectedEndOfJSON("[]byte", cursor)
				}
				cursor++
			}
		case '[':
			if d.sliceDecoder == nil {
				return nil, 0, &UnmarshalTypeError{
					Type:   rtype2type(d.typ),
					Offset: cursor,
				}
			}
			c, err := d.sliceDecoder.decode(buf, cursor, depth, p)
			if err != nil {
				return nil, 0, err
			}
			return nil, c, nil
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
			return []byte{}, cursor, nil
		default:
			goto ERROR
		}
	}
ERROR:
	return nil, 0, errNotAtBeginningOfValue(cursor)
}
