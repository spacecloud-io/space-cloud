package handlers

import (
	"github.com/gorilla/mux"
	"github.com/miolini/datacounter"
	"github.com/spaceuptech/space-cloud/utils/metrics"
	"io"
	"net/http"
	"sync/atomic"
)

func HandleMetricMiddleWare(next http.Handler, metrics *metrics.Module) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		projectID, ok := vars["project"]
		if ok {
			readerCounter := NewReaderCounter(r.Body)
			writerCounter := datacounter.NewResponseWriterCounter(w)
			r.Body = readerCounter
			next.ServeHTTP(writerCounter, r)
			metrics.AddIngress(projectID, readerCounter.Count()+200)
			metrics.AddEgress(projectID, writerCounter.Count()+200)
			return
		}
		next.ServeHTTP(w, r)

	})
}

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
