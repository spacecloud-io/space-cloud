package utils

import (
	"fmt"
	"hash/fnv"
)

// Hash returns a hash value for provided string
func Hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%v", h.Sum32())
}
