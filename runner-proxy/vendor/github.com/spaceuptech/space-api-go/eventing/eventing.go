package eventing

import (
	"context"
	"time"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/types"
)

// Eventing stores eventing config
type Eventing struct {
	config *config.Config
	event  *event
}

type event struct {
	Type          string            `json:"type"`      // The type of the event
	Delay         int64             `json:"delay"`     // Time in seconds
	Timestamp     string            `json:"timestamp"` // Milliseconds from unix epoch (UTC)
	Payload       interface{}       `json:"payload"`   // payload contains necessary event dat
	Options       map[string]string `json:"options"`
	IsSynchronous bool              `json:"isSynchronous"` // if true then client will wait for response of event
}

// New returns a eventing object
func New(eventType string, payload map[string]interface{}, config *config.Config) *Eventing {
	return &Eventing{config: config, event: &event{
		Type:    eventType,
		Payload: payload,
	}}
}

// Apply triggers a custom event
func (d *Eventing) Apply(ctx context.Context) (*types.Response, error) {
	return d.config.Transport.TriggerEvent(ctx, &types.Meta{Project: d.config.Project, Token: d.config.Token}, d.event)
}

// Delay specified in seconds delays the custom event trigger from timestamp specified
func (d *Eventing) Delay(delay int64) *Eventing {
	d.event.Delay = delay
	return d
}

// TimeStamp schedule an event trigger at the given timestamp (in milliseconds)
func (d *Eventing) TimeStamp(timestamp time.Time) *Eventing {
	d.event.Timestamp = timestamp.Format(time.RFC3339)
	return d
}

// Options sets the eventing options
func (d *Eventing) Options(options map[string]string) *Eventing {
	d.event.Options = options
	return d
}

// Synchronous makes the event synchronous by waiting for the event response
func (d *Eventing) Synchronous() *Eventing {
	d.event.IsSynchronous = true
	return d
}
