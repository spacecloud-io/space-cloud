package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// OpenAPISource is the schema for the openapisource API
type OpenAPISource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenAPISourceSpec   `json:"spec,omitempty"`
	Status OpenAPISourceStatus `json:"status,omitempty"`
}

// OpenAPISourceSpec defines the desired state of OpenAPISourceSpec
type OpenAPISourceSpec struct {
	// Sources descrives the remote OpenAPI endpoints that we want to federate
	Source *HTTPSource `json:"source"`
}

// OpenAPISourceStatus defines the observed state of OpenAPISource
type OpenAPISourceStatus struct {
	// TODO: Add state to show if sync succeeded or if there was an error
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// OpenAPISourceList contains a list of OpenAPISource
type OpenAPISourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenAPISource `json:"items"`
}
