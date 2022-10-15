package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GraphqlSource is the schema for the graphqlsource API
type GraphqlSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GraphqlSourceSpec   `json:"spec,omitempty"`
	Status GraphqlSourceStatus `json:"status,omitempty"`
}

// GraphqlSourceSpec defines the desired state of GraphqlSourceSpec
type GraphqlSourceSpec struct {
	// Sources descrives the remote graphql endpoints that we want to federate
	Source *HTTPSource `json:"source"`
}

// GraphqlSourceStatus defines the observed state of GraphqlSource
type GraphqlSourceStatus struct {
	// TODO: Add state to show if sync succeeded or if there was an error
}

// HTTPSource defines the parameters to connect to the remote http source
type HTTPSource struct {
	// URL describes the full http url to communicate with remote
	URL string `json:"url"`

	// TODO: Add support for authentication with remote
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// GraphqlSourceList contains a list of GraphqlSource
type GraphqlSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GraphqlSource `json:"items"`
}
