package json

import (
	"unsafe"
)

type numberDecoder struct {
	*floatDecoder
	op         func(unsafe.Pointer, Number)
	structName string
	fieldName  string
}

func newNumberDecoder(structName, fieldName string, op func(unsafe.Pointer, Number)) *numberDecoder {
	return &numberDecoder{
		floatDecoder: newFloatDecoder(structName, fieldName, nil),
		op:           op,
		structName:   structName,
		fieldName:    fieldName,
	}
}

func (d *numberDecoder) decodeStream(s *stream, depth int64, p unsafe.Pointer) error {
	bytes, err := d.floatDecoder.decodeStreamByte(s)
	if err != nil {
		return err
	}
	d.op(p, Number(string(bytes)))
	s.reset()
	return nil
}

func (d *numberDecoder) decode(buf []byte, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.floatDecoder.decodeByte(buf, cursor)
	if err != nil {
		return 0, err
	}
	cursor = c
	s := *(*string)(unsafe.Pointer(&bytes))
	d.op(p, Number(s))
	return cursor, nil
}
