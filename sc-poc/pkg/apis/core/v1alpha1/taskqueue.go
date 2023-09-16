package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TaskQueue is the schema for the TaskQueue source
type TaskQueue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskQueueSpec   `json:"spec,omitempty"`
	Status TaskQueueStatus `json:"status,omitempty"`
}

// TaskQueueSpec describes the specification of a task queue
type TaskQueueSpec struct {
	// Source describes which source to select as a source for the task queue
	Source string `json:"source"`

	// Operations describes the options for all the operations
	// +kubebuilder:validation:Optional
	Operations *TaskQueueOperations `json:"operations,omitempty"`
}

type TaskQueueOperations struct {
	// AddTask describes the endpoint options for the "AddTask" operation
	// +kubebuilder:validation:Optional
	AddTask *TaskQueueEndpointOptions `json:"addTask,omitempty"`

	// ReceiveTasks describes the endpoint options for the "ReceiveTasks" operation
	// +kubebuilder:validation:Optional
	ReceiveTasks *TaskQueueEndpointOptions `json:"getTasks,omitempty"`

	// AckTask describes the endpoint options for the "AckTask" operation
	// +kubebuilder:validation:Optional
	AckTask *TaskQueueEndpointOptions `json:"ackTask,omitempty"`

	// DeleteTask describes the endpoint options for the "DeleteTask" operation
	// +kubebuilder:validation:Optional
	DeleteTask *TaskQueueEndpointOptions `json:"deleteTask,omitempty"`

	// GetPendingTasks describes the endpoint options for the "GetPendingTasks" operation
	// +kubebuilder:validation:Optional
	GetPendingTasks *TaskQueueEndpointOptions `json:"getPendingTasks,omitempty"`
}

// TaskQueueEndpointOptions describes the options present for each API
type TaskQueueEndpointOptions struct {
	// Plugins describes the plugins to be used for this endpoint
	// +kubebuilder:validation:Optional
	Plugins []HTTPPlugin `json:"plugins,omitempty"`
}

// TaskQueueStatus defines the observed state of the TaskQueue source
type TaskQueueStatus struct {
	// TODO: Add state to show if sync succeeded or if there was an error
}

// +kubebuilder:object:root=true

// TaskQueueList contains a list of TaskQueues
type TaskQueueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TaskQueue `json:"items"`
}
