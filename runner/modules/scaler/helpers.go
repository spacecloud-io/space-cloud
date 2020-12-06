package scaler

import (
	"context"
	"fmt"
	"time"
)

type isActiveStream struct {
	ch chan bool
}

func (s *isActiveStream) Notify() {
	s.ch <- true
}

func (s *Scaler) isStreamActive(project, service, version, minReplicas string) bool {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Return true if min replicas is not equal to zero.
	if minReplicas != "0" {
		return true
	}

	// Check if stream exists
	count := 0
	for {
		// Return if count is 3
		if count == 3 {
			return false
		}

		// Check if key exists
		exists, err := s.pubsubClient.CheckIfKeyExists(ctx, generateKey(project, service, version))
		if err != nil {
			count++
			time.Sleep(2 * time.Second)
			continue
		}

		return exists
	}
}

func (s *Scaler) addIsActiveStream(project, service, version string, ch chan bool) {
	// First remove the previous stream with the same name
	s.removeIsActiveStream(project, service, version)

	// Acquire the lock after the removing step since it too acquires a lock internally
	s.lock.Lock()
	defer s.lock.Unlock()

	s.isActiveStreams[generateKey(project, service, version)] = &isActiveStream{ch: ch}
}

func (s *Scaler) removeIsActiveStream(project, service, version string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Generate a key for the stream
	key := generateKey(project, service, version)

	// Close the stream if it already existed
	if stream, p := s.isActiveStreams[key]; p {
		close(stream.ch)
	}

	// Remove the stream from the map
	delete(s.isActiveStreams, key)
}

func generateKey(project, service, version string) string {
	return fmt.Sprintf("%s---%s---%s", project, service, version)
}
