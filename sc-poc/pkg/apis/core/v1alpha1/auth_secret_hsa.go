package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JwtHSASecret is the schema for the hsasecret api
type JwtHSASecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JwtHSASecretSpec   `json:"spec,omitempty"`
	Status JwtHSASecretStatus `json:"status,omitempty"`
}

// JwtHSASecretSpec defines the desired state of JwtHSASecretSpec
type JwtHSASecretSpec struct {
	AuthSecret `json:",inline"`

	// Secret holds the content of the secret or jwk url
	Secret SecretSource `json:"secret"`
}

// JwtHSASecretStatus defines the observed state of HSASecret
type JwtHSASecretStatus struct {
	// TODO: Add state to show if sync succeeded or if there was an error
}

// +kubebuilder:object:root=true

// JwtHSASecretList contains a list of JwtHSASecret
type JwtHSASecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JwtHSASecret `json:"items"`
}
