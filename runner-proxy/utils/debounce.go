package utils

import "sync"

// Debounce is a functionality which helps make sure a certain piece of logic doesn't get called multiple times.
type Debounce struct {
	serviceLocks sync.Map
}

// NewDebounce creates a new debounce object
func NewDebounce() *Debounce {
	return &Debounce{}
}

// Wait makes a caller wait on a certain key. The function provided will be called just once per unique key even though
// the caller may call it multiple times. It is important that the same callback is passes for each given key.
func (d *Debounce) Wait(key string, cb func() error) error {
	// Create an empty array and put it in the map atomically. This array will hold all the callers interested in calling the
	// debounced invocation.
	emptyArray, ch := NewDebounceArray()
	tempArray, _ := d.serviceLocks.LoadOrStore(key, emptyArray)
	array := tempArray.(*DebounceArray)

	// Add channel to array. If the current caller is the first one to use the key (i.e. size = 1), start the debounced logic.
	size := array.Add(ch)
	if size == 1 {
		// Notify all clients that the task is done so they may end their Wait. Also delete the key to free up space.
		// It is possible that the same key can have multiple callbacks running in parallel. It is unlikely.
		array.Notify(cb())
		d.serviceLocks.Delete(key)
	}

	// Wait for signal of completion
	return <-ch
}

// DebounceArray is the entity which holds the channels the callers are waiting on
type DebounceArray struct {
	lock  sync.Mutex
	array []chan error
}

// NewDebounceArray creates a new debounce array object
func NewDebounceArray() (*DebounceArray, chan error) {
	ch := make(chan error, 1)
	return &DebounceArray{array: []chan error{}}, ch
}

// Add adds a new caller to the wait list
func (a *DebounceArray) Add(ch chan error) int {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.array = append(a.array, ch)
	return len(a.array)
}

// Notify signals all interested callers that the event is completed
func (a *DebounceArray) Notify(err error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for _, ch := range a.array {
		ch <- err
	}

	a.array = []chan error{}
}
