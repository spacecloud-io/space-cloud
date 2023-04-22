package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JwtRSASecret is the schema for the rsasecret api
type JwtRSASecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JwtRSASecretSpec   `json:"spec,omitempty"`
	Status JwtRSASecretStatus `json:"status,omitempty"`
}

// JwtRSASecretSpec defines the desired state of JwtRSASecretSpec
type JwtRSASecretSpec struct {
	AuthSecret `json:",inline"`

	// Value holds the content of the secret or jwk url
	// +kubebuilder:validation:Optional
	PublicValue string `json:"public,omitempty"`

	PrivateValue string `json:"private,omitempty"`

	// Source for the AuthSecret's value. Cannot be used if value is not empty.
	// +kubebuilder:validation:Optional
	ValueFrom *SecretSource `json:"valueFrom,omitempty"`
}

// JwtRSASecretStatus defines the observed state of RSASecret
type JwtRSASecretStatus struct {
	// TODO: Add state to show if sync succeeded or if there was an error
}

// +kubebuilder:object:root=true

// JwtRSASecretList contains a list of JwtRSASecret
type JwtRSASecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JwtRSASecret `json:"items"`
}
