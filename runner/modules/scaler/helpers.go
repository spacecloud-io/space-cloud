package scaler

import (
	"fmt"
)

func (s *Scaler) isStreamActive(project, service, version, minReplicas string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Return true if min replicas is not equal to zero.
	if minReplicas != "0" {
		return true
	}

	// Check if stream exists
	stream, p := s.isActiveStreams[generateKey(project, service, version)]
	if !p {
		return false
	}

	return stream.IsActive()
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
	return fmt.Sprintf("%s--%s--%s", project, service, version)
}
