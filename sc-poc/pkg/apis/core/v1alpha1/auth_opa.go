package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// OPAPolicy is the schema for the opapolicy api
type OPAPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OPAPolicySpec   `json:"spec,omitempty"`
	Status OPAPolicyStatus `json:"status,omitempty"`
}

// OPAPolicySpec describes the specification of the opapolicy api
type OPAPolicySpec struct {
	// Rego describes the opa policy in rego
	Rego string `json:"rego"`
}

// OPAPolicyStatus defines the observed state of the opapolicy api
type OPAPolicyStatus struct {
	// TODO: Add state to show if sync succeeded or if there was an error
}

// +kubebuilder:object:root=true

// OPAPolicyList contains a list of OPAPolicy
type OPAPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OPAPolicy `json:"items"`
}
