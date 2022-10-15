package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HSASecret is the schema for the hsasecret api
type HSASecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HSASecretSpec   `json:"spec,omitempty"`
	Status HSASecretStatus `json:"status,omitempty"`
}

// HSASecretSpec defines the desired state of HSASecretSpec
type HSASecretSpec struct {
	AuthSecret `json:",inline"`

	// Value holds the content of the secret or jwk url
	// +kubebuilder:validation:Optional
	Value string `json:"value,omitempty"`

	// Source for the AuthSecret's value. Cannot be used if value is not empty.
	// +kubebuilder:validation:Optional
	ValueFrom *SecretSource `json:"valueFrom,omitempty"`
}

// HSASecretStatus defines the observed state of HSASecret
type HSASecretStatus struct {
	// TODO: Add state to show if sync succeeded or if there was an error
}

// +kubebuilder:object:root=true

// HSASecretList contains a list of HSASecret
type HSASecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HSASecret `json:"items"`
}
