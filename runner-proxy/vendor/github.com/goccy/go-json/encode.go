package json

import (
	"bytes"
	"encoding/base64"
	"io"
	"math"
	"strconv"
	"sync"
	"unsafe"
)

// An Encoder writes JSON values to an output stream.
type Encoder struct {
	w                 io.Writer
	enabledIndent     bool
	enabledHTMLEscape bool
	prefix            string
	indentStr         string
}

const (
	bufSize = 1024
)

type EncodeOption int

const (
	EncodeOptionHTMLEscape EncodeOption = 1 << iota
	EncodeOptionIndent
	EncodeOptionUnorderedMap
)

var (
	encRuntimeContextPool = sync.Pool{
		New: func() interface{} {
			return &encodeRuntimeContext{
				buf:      make([]byte, 0, bufSize),
				ptrs:     make([]uintptr, 128),
				keepRefs: make([]unsafe.Pointer, 0, 8),
			}
		},
	}
)

func takeEncodeRuntimeContext() *encodeRuntimeContext {
	return encRuntimeContextPool.Get().(*encodeRuntimeContext)
}

func releaseEncodeRuntimeContext(ctx *encodeRuntimeContext) {
	encRuntimeContextPool.Put(ctx)
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w, enabledHTMLEscape: true}
}

// Encode writes the JSON encoding of v to the stream, followed by a newline character.
//
// See the documentation for Marshal for details about the conversion of Go values to JSON.
func (e *Encoder) Encode(v interface{}) error {
	return e.EncodeWithOption(v)
}

// EncodeWithOption call Encode with EncodeOption.
func (e *Encoder) EncodeWithOption(v interface{}, optFuncs ...EncodeOptionFunc) error {
	ctx := takeEncodeRuntimeContext()

	err := e.encodeWithOption(ctx, v, optFuncs...)

	releaseEncodeRuntimeContext(ctx)
	return err
}

func (e *Encoder) encodeWithOption(ctx *encodeRuntimeContext, v interface{}, optFuncs ...EncodeOptionFunc) error {
	var opt EncodeOption
	if e.enabledHTMLEscape {
		opt |= EncodeOptionHTMLEscape
	}
	for _, optFunc := range optFuncs {
		opt = optFunc(opt)
	}
	var (
		buf []byte
		err error
	)
	if e.enabledIndent {
		buf, err = encodeIndent(ctx, v, e.prefix, e.indentStr, opt)
	} else {
		buf, err = encode(ctx, v, opt)
	}
	if err != nil {
		return err
	}
	if e.enabledIndent {
		buf = buf[:len(buf)-2]
	} else {
		buf = buf[:len(buf)-1]
	}
	buf = append(buf, '\n')
	if _, err := e.w.Write(buf); err != nil {
		return err
	}
	return nil
}

// SetEscapeHTML specifies whether problematic HTML characters should be escaped inside JSON quoted strings.
// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e to avoid certain safety problems that can arise when embedding JSON in HTML.
//
// In non-HTML settings where the escaping interferes with the readability of the output, SetEscapeHTML(false) disables this behavior.
func (e *Encoder) SetEscapeHTML(on bool) {
	e.enabledHTMLEscape = on
}

// SetIndent instructs the encoder to format each subsequent encoded value as if indented by the package-level function Indent(dst, src, prefix, indent).
// Calling SetIndent("", "") disables indentation.
func (e *Encoder) SetIndent(prefix, indent string) {
	if prefix == "" && indent == "" {
		e.enabledIndent = false
		return
	}
	e.prefix = prefix
	e.indentStr = indent
	e.enabledIndent = true
}

func marshal(v interface{}, opt EncodeOption) ([]byte, error) {
	ctx := takeEncodeRuntimeContext()

	buf, err := encode(ctx, v, opt|EncodeOptionHTMLEscape)
	if err != nil {
		releaseEncodeRuntimeContext(ctx)
		return nil, err
	}

	// this line exists to escape call of `runtime.makeslicecopy` .
	// if use `make([]byte, len(buf)-1)` and `copy(copied, buf)`,
	// dst buffer size and src buffer size are differrent.
	// in this case, compiler uses `runtime.makeslicecopy`, but it is slow.
	buf = buf[:len(buf)-1]
	copied := make([]byte, len(buf))
	copy(copied, buf)

	releaseEncodeRuntimeContext(ctx)
	return copied, nil
}

func marshalNoEscape(v interface{}, opt EncodeOption) ([]byte, error) {
	ctx := takeEncodeRuntimeContext()

	buf, err := encodeNoEscape(ctx, v, opt|EncodeOptionHTMLEscape)
	if err != nil {
		releaseEncodeRuntimeContext(ctx)
		return nil, err
	}

	// this line exists to escape call of `runtime.makeslicecopy` .
	// if use `make([]byte, len(buf)-1)` and `copy(copied, buf)`,
	// dst buffer size and src buffer size are differrent.
	// in this case, compiler uses `runtime.makeslicecopy`, but it is slow.
	buf = buf[:len(buf)-1]
	copied := make([]byte, len(buf))
	copy(copied, buf)

	releaseEncodeRuntimeContext(ctx)
	return copied, nil
}

func marshalIndent(v interface{}, prefix, indent string, opt EncodeOption) ([]byte, error) {
	ctx := takeEncodeRuntimeContext()

	buf, err := encodeIndent(ctx, v, prefix, indent, opt|EncodeOptionHTMLEscape)
	if err != nil {
		releaseEncodeRuntimeContext(ctx)
		return nil, err
	}

	buf = buf[:len(buf)-2]
	copied := make([]byte, len(buf))
	copy(copied, buf)

	releaseEncodeRuntimeContext(ctx)
	return copied, nil
}

func encode(ctx *encodeRuntimeContext, v interface{}, opt EncodeOption) ([]byte, error) {
	b := ctx.buf[:0]
	if v == nil {
		b = encodeNull(b)
		b = encodeComma(b)
		return b, nil
	}
	header := (*interfaceHeader)(unsafe.Pointer(&v))
	typ := header.typ

	typeptr := uintptr(unsafe.Pointer(typ))
	codeSet, err := encodeCompileToGetCodeSet(typeptr)
	if err != nil {
		return nil, err
	}

	p := uintptr(header.ptr)
	ctx.init(p, codeSet.codeLength)
	buf, err := encodeRunCode(ctx, b, codeSet, opt)

	ctx.keepRefs = append(ctx.keepRefs, header.ptr)

	if err != nil {
		return nil, err
	}

	ctx.buf = buf
	return buf, nil
}

func encodeNoEscape(ctx *encodeRuntimeContext, v interface{}, opt EncodeOption) ([]byte, error) {
	b := ctx.buf[:0]
	if v == nil {
		b = encodeNull(b)
		b = encodeComma(b)
		return b, nil
	}
	header := (*interfaceHeader)(unsafe.Pointer(&v))
	typ := header.typ

	typeptr := uintptr(unsafe.Pointer(typ))
	codeSet, err := encodeCompileToGetCodeSet(typeptr)
	if err != nil {
		return nil, err
	}

	p := uintptr(header.ptr)
	ctx.init(p, codeSet.codeLength)
	buf, err := encodeRunCode(ctx, b, codeSet, opt)
	if err != nil {
		return nil, err
	}

	ctx.buf = buf
	return buf, nil
}

func encodeIndent(ctx *encodeRuntimeContext, v interface{}, prefix, indent string, opt EncodeOption) ([]byte, error) {
	b := ctx.buf[:0]
	if v == nil {
		b = encodeNull(b)
		b = encodeIndentComma(b)
		return b, nil
	}
	header := (*interfaceHeader)(unsafe.Pointer(&v))
	typ := header.typ

	typeptr := uintptr(unsafe.Pointer(typ))
	codeSet, err := encodeCompileToGetCodeSet(typeptr)
	if err != nil {
		return nil, err
	}

	p := uintptr(header.ptr)
	ctx.init(p, codeSet.codeLength)
	buf, err := encodeRunIndentCode(ctx, b, codeSet, prefix, indent, opt)

	ctx.keepRefs = append(ctx.keepRefs, header.ptr)

	if err != nil {
		return nil, err
	}

	ctx.buf = buf
	return buf, nil
}

func encodeRunCode(ctx *encodeRuntimeContext, b []byte, codeSet *opcodeSet, opt EncodeOption) ([]byte, error) {
	if (opt & EncodeOptionHTMLEscape) != 0 {
		return encodeRunEscaped(ctx, b, codeSet, opt)
	}
	return encodeRun(ctx, b, codeSet, opt)
}

func encodeRunIndentCode(ctx *encodeRuntimeContext, b []byte, codeSet *opcodeSet, prefix, indent string, opt EncodeOption) ([]byte, error) {
	ctx.prefix = []byte(prefix)
	ctx.indentStr = []byte(indent)
	if (opt & EncodeOptionHTMLEscape) != 0 {
		return encodeRunEscapedIndent(ctx, b, codeSet, opt)
	}
	return encodeRunIndent(ctx, b, codeSet, opt)
}

func encodeFloat32(b []byte, v float32) []byte {
	f64 := float64(v)
	abs := math.Abs(f64)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		f32 := float32(abs)
		if f32 < 1e-6 || f32 >= 1e21 {
			fmt = 'e'
		}
	}
	return strconv.AppendFloat(b, f64, fmt, -1, 32)
}

func encodeFloat64(b []byte, v float64) []byte {
	abs := math.Abs(v)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			fmt = 'e'
		}
	}
	return strconv.AppendFloat(b, v, fmt, -1, 64)
}

func encodeBool(b []byte, v bool) []byte {
	if v {
		return append(b, "true"...)
	}
	return append(b, "false"...)
}

func encodeNull(b []byte) []byte {
	return append(b, "null"...)
}

func encodeComma(b []byte) []byte {
	return append(b, ',')
}

func encodeIndentComma(b []byte) []byte {
	return append(b, ',', '\n')
}

func appendStructEnd(b []byte) []byte {
	return append(b, '}', ',')
}

func appendStructEndIndent(ctx *encodeRuntimeContext, b []byte, indent int) []byte {
	b = append(b, '\n')
	b = append(b, ctx.prefix...)
	b = append(b, bytes.Repeat(ctx.indentStr, ctx.baseIndent+indent)...)
	return append(b, '}', ',', '\n')
}

func encodeByteSlice(b []byte, src []byte) []byte {
	encodedLen := base64.StdEncoding.EncodedLen(len(src))
	b = append(b, '"')
	pos := len(b)
	remainLen := cap(b[pos:])
	var buf []byte
	if remainLen > encodedLen {
		buf = b[pos : pos+encodedLen]
	} else {
		buf = make([]byte, encodedLen)
	}
	base64.StdEncoding.Encode(buf, src)
	return append(append(b, buf...), '"')
}

func appendIndent(ctx *encodeRuntimeContext, b []byte, indent int) []byte {
	b = append(b, ctx.prefix...)
	return append(b, bytes.Repeat(ctx.indentStr, ctx.baseIndent+indent)...)
}
