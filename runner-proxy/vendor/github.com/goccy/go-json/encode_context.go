package json

import (
	"bytes"
	"sync"
	"unsafe"
)

type mapItem struct {
	key   []byte
	value []byte
}

type mapslice struct {
	items []mapItem
}

func (m *mapslice) Len() int {
	return len(m.items)
}

func (m *mapslice) Less(i, j int) bool {
	return bytes.Compare(m.items[i].key, m.items[j].key) < 0
}

func (m *mapslice) Swap(i, j int) {
	m.items[i], m.items[j] = m.items[j], m.items[i]
}

type encodeMapContext struct {
	pos   []int
	slice *mapslice
	buf   []byte
}

var mapContextPool = sync.Pool{
	New: func() interface{} {
		return &encodeMapContext{}
	},
}

func newMapContext(mapLen int) *encodeMapContext {
	ctx := mapContextPool.Get().(*encodeMapContext)
	if ctx.slice == nil {
		ctx.slice = &mapslice{
			items: make([]mapItem, 0, mapLen),
		}
	}
	if cap(ctx.pos) < (mapLen*2 + 1) {
		ctx.pos = make([]int, 0, mapLen*2+1)
		ctx.slice.items = make([]mapItem, 0, mapLen)
	} else {
		ctx.pos = ctx.pos[:0]
		ctx.slice.items = ctx.slice.items[:0]
	}
	ctx.buf = ctx.buf[:0]
	return ctx
}

func releaseMapContext(c *encodeMapContext) {
	mapContextPool.Put(c)
}

type encodeCompileContext struct {
	typ                      *rtype
	root                     bool
	opcodeIndex              int
	ptrIndex                 int
	indent                   int
	structTypeToCompiledCode map[uintptr]*compiledCode

	parent *encodeCompileContext
}

func (c *encodeCompileContext) context() *encodeCompileContext {
	return &encodeCompileContext{
		typ:                      c.typ,
		root:                     c.root,
		opcodeIndex:              c.opcodeIndex,
		ptrIndex:                 c.ptrIndex,
		indent:                   c.indent,
		structTypeToCompiledCode: c.structTypeToCompiledCode,
		parent:                   c,
	}
}

func (c *encodeCompileContext) withType(typ *rtype) *encodeCompileContext {
	ctx := c.context()
	ctx.typ = typ
	return ctx
}

func (c *encodeCompileContext) incIndent() *encodeCompileContext {
	ctx := c.context()
	ctx.indent++
	return ctx
}

func (c *encodeCompileContext) decIndent() *encodeCompileContext {
	ctx := c.context()
	ctx.indent--
	return ctx
}

func (c *encodeCompileContext) incIndex() {
	c.incOpcodeIndex()
	c.incPtrIndex()
}

func (c *encodeCompileContext) decIndex() {
	c.decOpcodeIndex()
	c.decPtrIndex()
}

func (c *encodeCompileContext) incOpcodeIndex() {
	c.opcodeIndex++
	if c.parent != nil {
		c.parent.incOpcodeIndex()
	}
}

func (c *encodeCompileContext) decOpcodeIndex() {
	c.opcodeIndex--
	if c.parent != nil {
		c.parent.decOpcodeIndex()
	}
}

func (c *encodeCompileContext) incPtrIndex() {
	c.ptrIndex++
	if c.parent != nil {
		c.parent.incPtrIndex()
	}
}

func (c *encodeCompileContext) decPtrIndex() {
	c.ptrIndex--
	if c.parent != nil {
		c.parent.decPtrIndex()
	}
}

type encodeRuntimeContext struct {
	buf        []byte
	ptrs       []uintptr
	keepRefs   []unsafe.Pointer
	seenPtr    []uintptr
	baseIndent int
	prefix     []byte
	indentStr  []byte
}

func (c *encodeRuntimeContext) init(p uintptr, codelen int) {
	if len(c.ptrs) < codelen {
		c.ptrs = make([]uintptr, codelen)
	}
	c.ptrs[0] = p
	c.keepRefs = c.keepRefs[:0]
	c.seenPtr = c.seenPtr[:0]
	c.baseIndent = 0
}

func (c *encodeRuntimeContext) ptr() uintptr {
	header := (*sliceHeader)(unsafe.Pointer(&c.ptrs))
	return uintptr(header.data)
}
