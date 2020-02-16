package handlers

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

// ReaderCounter is counter for io.Reader
type ReaderCounter struct {
	r io.ReadCloser
	count uint64
}

// NewReaderCounter function for create new ReaderCounter
func NewReaderCounter(r io.ReadCloser) *ReaderCounter {
	return &ReaderCounter{
		r: r,
	}
}

func (counter *ReaderCounter) Read(buf []byte) (int, error) {
	n, err := counter.r.Read(buf)
	atomic.AddUint64(&counter.count, uint64(n))
	return n, err
}

// Count function return counted bytes
func (counter *ReaderCounter) Count() uint64 {
	return atomic.LoadUint64(&counter.count)
}

func (counter *ReaderCounter) Close() error {
	return counter.r.Close()
}

// ResponseWriterCounter is counter for http.ResponseWriter
type ResponseWriterCounter struct {
	http.ResponseWriter
	count   uint64
	started time.Time
}

// NewResponseWriterCounter function create new ResponseWriterCounter
func NewResponseWriterCounter(rw http.ResponseWriter) *ResponseWriterCounter {
	return &ResponseWriterCounter{
		ResponseWriter: rw,
		started:        time.Now(),
	}
}

func (counter *ResponseWriterCounter) Write(buf []byte) (int, error) {
	n, err := counter.ResponseWriter.Write(buf)
	atomic.AddUint64(&counter.count, uint64(n))
	return n, err
}

func (counter *ResponseWriterCounter) Header() http.Header {
	return counter.ResponseWriter.Header()
}

func (counter *ResponseWriterCounter) WriteHeader(statusCode int) {
	counter.Header().Set("X-Runtime", fmt.Sprintf("%.6f", time.Since(counter.started).Seconds()))
	counter.ResponseWriter.WriteHeader(statusCode)
}

func (counter *ResponseWriterCounter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return counter.ResponseWriter.(http.Hijacker).Hijack()
}

// Count function return counted bytes
func (counter *ResponseWriterCounter) Count() uint64 {
	return atomic.LoadUint64(&counter.count)
}

func (counter *ResponseWriterCounter) Started() time.Time {
	return counter.started
}
