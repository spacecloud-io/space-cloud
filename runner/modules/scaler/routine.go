package scaler

import (
	"context"
)

func (s *Scaler) routineScaleUp() {
	messages, err := s.pubsubClient.Subscribe(context.Background(), "scale-up")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		// Notify if the stream is present internally
		s.lock.RLock()
		if isActive, p := s.isActiveStreams[msg.Payload]; p {
			isActive.Notify()
		}
		s.lock.RUnlock()
	}
}
