package pubsub_channel

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/modules/pubsub"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var pubsubchannelsResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "pubsubchannels"}

func init() {
	source.RegisterSource(PubsubChannelSource{}, pubsubchannelsResource)
}

// PubsubChannelSource describes a PubsubChannel source
type PubsubChannelSource struct {
	v1alpha1.PubsubChannel
}

// CaddyModule returns the Caddy module information.
func (PubsubChannelSource) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(pubsubchannelsResource)),
		New: func() caddy.Module { return new(PubsubChannelSource) },
	}
}

// GetPriority returns the priority of the source.
func (s *PubsubChannelSource) GetPriority() int {
	return 0
}

// GetProviders returns the providers this source is applicable for
func (s *PubsubChannelSource) GetProviders() []string {
	return []string{"pubsub"}
}

// GetChannel returns the channel of this source.
func (s *PubsubChannelSource) GetChannel() v1alpha1.PubsubChannelSpec {
	return s.Spec
}

// Interface guard
var (
	_ source.Source = (*PubsubChannelSource)(nil)
	_ pubsub.Source = (*PubsubChannelSource)(nil)
)
