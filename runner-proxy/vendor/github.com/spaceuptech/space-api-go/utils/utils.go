package utils

import "io"

// CloseTheCloser closes the closer :P
func CloseTheCloser(c io.Closer) {
	_ = c.Close()
}
