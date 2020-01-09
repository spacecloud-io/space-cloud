package utils

import "sync"

// Array is an atomic array
type Array struct {
	lock  sync.RWMutex
	array []interface{}
}

// NewArray creates a new array
func NewArray(length int) *Array {
	return &Array{array: make([]interface{}, length)}
}

// Set sets an item in the array at a particular index
func (a *Array) Set(index int, value interface{}) {
	a.lock.Lock()
	a.array[index] = value
	a.lock.Unlock()
}

// GetAll gets the array
func (a *Array) GetAll() []interface{} {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.array
}

// Object is an atomic array
type Object struct {
	lock   sync.RWMutex
	object map[string]interface{}
}

// NewObject creates a new atomic object
func NewObject() *Object {
	return &Object{object: map[string]interface{}{}}
}

// Set sets a key in the object with a value
func (obj *Object) Set(key string, value interface{}) {
	obj.lock.Lock()
	obj.object[key] = value
	obj.lock.Unlock()
}

// Get gets the value of a key from the object
func (obj *Object) Get(key string) (value interface{}, present bool) {
	obj.lock.RLock()
	defer obj.lock.RUnlock()

	value, present = obj.object[key]
	return
}

// GetAll returns the entire object
func (obj *Object) GetAll() map[string]interface{} {
	obj.lock.RLock()
	defer obj.lock.RUnlock()

	return obj.object
}
