package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// RedisSource is the schema for the RedisSource source
type RedisSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisSourceSpec   `json:"spec,omitempty"`
	Status RedisSourceStatus `json:"status,omitempty"`
}

// Redis describes the specification of the Redis source
type RedisSourceSpec struct {
	// Host describes the host of the Redis instance
	Host SecretSource `json:"host"`

	// Port describes the port of the Redis instance
	Port SecretSource `json:"port"`

	// Password describes the password to use to connect to the Redis instance
	// +kubebuilder:validation:Optional
	Password *SecretSource `json:"password,omitempty"`

	// TaskQueueConfig describes the configuration to use when using Redis as a task queue.
	// +kubebuilder:validation:Optional
	TaskQueueConfig *RedisTaskQueueConfig `json:"taskQueueConfig,omitempty"`
}

// RedisTaskQueueConfig define the task queue options exposed by the Redis source.
type RedisTaskQueueConfig struct {
	CommonTaskQueueConfig `json:"-"`

	// MaxQueueSize specifies the maximum length of the Redis queue.
	// +kubebuilder:validation:Optional
	MaxQueueSize int `json:"maxQueueSize,omitempty"`
}

// RedisSourceStatus defines the observed state of the RedisSource source
type RedisSourceStatus struct {
	// TODO: Add state to show if sync succeeded or if there was an error
}

// +kubebuilder:object:root=true

// RedisSourceList contains a list of Redis
type RedisSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RedisSource `json:"items"`
}
