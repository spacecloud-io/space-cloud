package v1alpha1

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PubsubChannel is the schema for the pubsub channel
type PubsubChannel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PubsubChannelSpec   `json:"spec,omitempty"`
	Status PubsubChannelStatus `json:"status,omitempty"`
}

// PubsubChannel describes the specification of the pubsub channel
type PubsubChannelSpec struct {
	// Channel describes the name of the pubsub channel
	Channel string `json:"channel"`

	// Payload describes the payload schema of the channel
	Payload *ChannelSchema `json:"payload"`

	// ProducerOptions describes the options for a producer
	// +kubebuilder:validation:Optional
	ProducerOptions *ChannelOptions `json:"producerOptions,omitempty"`

	// ConsumerOptions describes the options for a consumer
	// +kubebuilder:validation:Optional
	ConsumerOptions *ChannelOptions `json:"consumerOptions,omitempty"`
}

// ChannelSchema defines the schema of the payload that the channel accepts
type ChannelSchema struct {
	// Type defines the type of the data
	// +kubebuilder:validation:Optional
	Type string `json:"type,omitempty"`

	// Items list the items of the array
	// +kubebuilder:validation:Optional
	Items *ChannelSchema `json:"items,omitempty"`

	// Properties list additional properties of the object
	// +kubebuilder:validation:Optional
	Properties map[string]*ChannelSchema `json:"properties,omitempty"`

	// Required specifies the required properties of the object
	// +kubebuilder:validation:Optional
	Required []string `json:"required,omitempty"`

	// AdditionalProperties defines if the schema accepts properties other than the ones mentioned
	// +kubebuilder:validation:Optional
	AdditionalProperties json.RawMessage `json:"additionalProperties,omitempty"`
}

// ChannelOptions describes the options available for a producer/ consumer
type ChannelOptions struct {
	// Plugins describes the plugins to be used for this endpoint
	// +kubebuilder:validation:Optional
	Plugins []HTTPPlugin `json:"plugins,omitempty"`
}

// PubsubChannel defines the observed state of the pubsub channel
type PubsubChannelStatus struct {
	// TODO: Add state to show if sync succeeded or if there was an error
}

// +kubebuilder:object:root=true

// OPAPolicyList contains a list of PubsubChannel
type PubsubChannelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PubsubChannel `json:"items"`
}