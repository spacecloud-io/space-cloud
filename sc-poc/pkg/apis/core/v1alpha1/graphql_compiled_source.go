package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CompiledGraphqlSource is the schema for the compiledgraphqlsource API
type CompiledGraphqlSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CompiledGraphqlSourceSpec   `json:"spec,omitempty"`
	Status CompiledGraphqlSourceStatus `json:"status,omitempty"`
}

// CompiledGraphqlSourceSpec defines the desired state of CompiledGraphqlSourceSpec
type CompiledGraphqlSourceSpec struct {
	// Graphql describes the graphql query to be compiled
	Graphql GraphqlQuery `json:"graphql"`

	// HTTP describes the http properties we want to override
	// +kubebuilder:validation:Optional
	HTTP *HTTPProperties `json:"http,omitempty"`

	// Plugins describes the plugins to be used for this endpoint
	// +kubebuilder:validation:Optional
	Plugins []*HTTPPlugin `json:"plugins,omitempty"`
}

// GraphqlQuery describes the graphql query to be compiled
type GraphqlQuery struct {
	// Query is the string query to be compiled
	Query string `json:"query"`

	// OperationName is the name of the operation to execute
	// +kubebuilder:validation:Optional
	OperationName string `json:"operationName"`

	// DefaultVariables defines the default variables we want to apply to the
	// graphql query
	// +kubebuilder:validation:Optional
	DefaultVariables runtime.RawExtension `json:"defaultVariables,omitempty"`
}

// HTTPProperties describes the rest api properties we want to attach to the source
type HTTPProperties struct {
	// Method describes the http method we want to use for the rest api
	// +kubebuilder:validation:Optional
	Method string `json:"method,omitempty"`

	// Urk describes the http url we want to use for the rest api
	// +kubebuilder:validation:Optional
	URL string `json:"url,omitempty"`
}

// CompiledGraphqlSourceStatus defines the observed state of CompiledGraphqlSource
type CompiledGraphqlSourceStatus struct {
	// TODO: Add state to show if sync succeeded or if there was an error
}

// +kubebuilder:object:root=true

// CompiledGraphqlSourceList contains a list of CompiledGraphqlSource
type CompiledGraphqlSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CompiledGraphqlSource `json:"items"`
}
