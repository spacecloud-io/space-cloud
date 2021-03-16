// +build go1.13

package json

import "unsafe"

//go:linkname mapitervalue reflect.mapiterelem
func mapitervalue(it unsafe.Pointer) unsafe.Pointer
