package json

import (
	"reflect"
	"unsafe"
)

// rtype representing reflect.rtype for noescape trick
type rtype struct{}

//go:linkname rtype_Align reflect.(*rtype).Align
//go:noescape
func rtype_Align(*rtype) int

func (t *rtype) Align() int {
	return rtype_Align(t)
}

//go:linkname rtype_FieldAlign reflect.(*rtype).FieldAlign
//go:noescape
func rtype_FieldAlign(*rtype) int

func (t *rtype) FieldAlign() int {
	return rtype_FieldAlign(t)
}

//go:linkname rtype_Method reflect.(*rtype).Method
//go:noescape
func rtype_Method(*rtype, int) reflect.Method

func (t *rtype) Method(a0 int) reflect.Method {
	return rtype_Method(t, a0)
}

//go:linkname rtype_MethodByName reflect.(*rtype).MethodByName
//go:noescape
func rtype_MethodByName(*rtype, string) (reflect.Method, bool)

func (t *rtype) MethodByName(a0 string) (reflect.Method, bool) {
	return rtype_MethodByName(t, a0)
}

//go:linkname rtype_NumMethod reflect.(*rtype).NumMethod
//go:noescape
func rtype_NumMethod(*rtype) int

func (t *rtype) NumMethod() int {
	return rtype_NumMethod(t)
}

//go:linkname rtype_Name reflect.(*rtype).Name
//go:noescape
func rtype_Name(*rtype) string

func (t *rtype) Name() string {
	return rtype_Name(t)
}

//go:linkname rtype_PkgPath reflect.(*rtype).PkgPath
//go:noescape
func rtype_PkgPath(*rtype) string

func (t *rtype) PkgPath() string {
	return rtype_PkgPath(t)
}

//go:linkname rtype_Size reflect.(*rtype).Size
//go:noescape
func rtype_Size(*rtype) uintptr

func (t *rtype) Size() uintptr {
	return rtype_Size(t)
}

//go:linkname rtype_String reflect.(*rtype).String
//go:noescape
func rtype_String(*rtype) string

func (t *rtype) String() string {
	return rtype_String(t)
}

//go:linkname rtype_Kind reflect.(*rtype).Kind
//go:noescape
func rtype_Kind(*rtype) reflect.Kind

func (t *rtype) Kind() reflect.Kind {
	return rtype_Kind(t)
}

//go:linkname rtype_Implements reflect.(*rtype).Implements
//go:noescape
func rtype_Implements(*rtype, reflect.Type) bool

func (t *rtype) Implements(u reflect.Type) bool {
	return rtype_Implements(t, u)
}

//go:linkname rtype_AssignableTo reflect.(*rtype).AssignableTo
//go:noescape
func rtype_AssignableTo(*rtype, reflect.Type) bool

func (t *rtype) AssignableTo(u reflect.Type) bool {
	return rtype_AssignableTo(t, u)
}

//go:linkname rtype_ConvertibleTo reflect.(*rtype).ConvertibleTo
//go:noescape
func rtype_ConvertibleTo(*rtype, reflect.Type) bool

func (t *rtype) ConvertibleTo(u reflect.Type) bool {
	return rtype_ConvertibleTo(t, u)
}

//go:linkname rtype_Comparable reflect.(*rtype).Comparable
//go:noescape
func rtype_Comparable(*rtype) bool

func (t *rtype) Comparable() bool {
	return rtype_Comparable(t)
}

//go:linkname rtype_Bits reflect.(*rtype).Bits
//go:noescape
func rtype_Bits(*rtype) int

func (t *rtype) Bits() int {
	return rtype_Bits(t)
}

//go:linkname rtype_ChanDir reflect.(*rtype).ChanDir
//go:noescape
func rtype_ChanDir(*rtype) reflect.ChanDir

func (t *rtype) ChanDir() reflect.ChanDir {
	return rtype_ChanDir(t)
}

//go:linkname rtype_IsVariadic reflect.(*rtype).IsVariadic
//go:noescape
func rtype_IsVariadic(*rtype) bool

func (t *rtype) IsVariadic() bool {
	return rtype_IsVariadic(t)
}

//go:linkname rtype_Elem reflect.(*rtype).Elem
//go:noescape
func rtype_Elem(*rtype) reflect.Type

func (t *rtype) Elem() *rtype {
	return type2rtype(rtype_Elem(t))
}

//go:linkname rtype_Field reflect.(*rtype).Field
//go:noescape
func rtype_Field(*rtype, int) reflect.StructField

func (t *rtype) Field(i int) reflect.StructField {
	return rtype_Field(t, i)
}

//go:linkname rtype_FieldByIndex reflect.(*rtype).FieldByIndex
//go:noescape
func rtype_FieldByIndex(*rtype, []int) reflect.StructField

func (t *rtype) FieldByIndex(index []int) reflect.StructField {
	return rtype_FieldByIndex(t, index)
}

//go:linkname rtype_FieldByName reflect.(*rtype).FieldByName
//go:noescape
func rtype_FieldByName(*rtype, string) (reflect.StructField, bool)

func (t *rtype) FieldByName(name string) (reflect.StructField, bool) {
	return rtype_FieldByName(t, name)
}

//go:linkname rtype_FieldByNameFunc reflect.(*rtype).FieldByNameFunc
//go:noescape
func rtype_FieldByNameFunc(*rtype, func(string) bool) (reflect.StructField, bool)

func (t *rtype) FieldByNameFunc(match func(string) bool) (reflect.StructField, bool) {
	return rtype_FieldByNameFunc(t, match)
}

//go:linkname rtype_In reflect.(*rtype).In
//go:noescape
func rtype_In(*rtype, int) reflect.Type

func (t *rtype) In(i int) reflect.Type {
	return rtype_In(t, i)
}

//go:linkname rtype_Key reflect.(*rtype).Key
//go:noescape
func rtype_Key(*rtype) reflect.Type

func (t *rtype) Key() *rtype {
	return type2rtype(rtype_Key(t))
}

//go:linkname rtype_Len reflect.(*rtype).Len
//go:noescape
func rtype_Len(*rtype) int

func (t *rtype) Len() int {
	return rtype_Len(t)
}

//go:linkname rtype_NumField reflect.(*rtype).NumField
//go:noescape
func rtype_NumField(*rtype) int

func (t *rtype) NumField() int {
	return rtype_NumField(t)
}

//go:linkname rtype_NumIn reflect.(*rtype).NumIn
//go:noescape
func rtype_NumIn(*rtype) int

func (t *rtype) NumIn() int {
	return rtype_NumIn(t)
}

//go:linkname rtype_NumOut reflect.(*rtype).NumOut
//go:noescape
func rtype_NumOut(*rtype) int

func (t *rtype) NumOut() int {
	return rtype_NumOut(t)
}

//go:linkname rtype_Out reflect.(*rtype).Out
//go:noescape
func rtype_Out(*rtype, int) reflect.Type

//go:linkname rtype_ptrTo reflect.(*rtype).ptrTo
//go:noescape
func rtype_ptrTo(*rtype) *rtype

func (t *rtype) Out(i int) reflect.Type {
	return rtype_Out(t, i)
}

//go:linkname rtype2type reflect.toType
//go:noescape
func rtype2type(t *rtype) reflect.Type

type interfaceHeader struct {
	typ *rtype
	ptr unsafe.Pointer
}

func type2rtype(t reflect.Type) *rtype {
	return (*rtype)(((*interfaceHeader)(unsafe.Pointer(&t))).ptr)
}
