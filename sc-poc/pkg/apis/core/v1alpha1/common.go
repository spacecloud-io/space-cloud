package v1alpha1

// SecretSource represents a source for a confidential value.
// Taken from apiv1.EnvVarSource.
type SecretSource struct {
	// SecretKeyRef selects a key of a secret in the pod's namespace
	// +kubebuilder:validation:Optional
	SecretKeyRef *SecretKeyRef `json:"secretKeyRef,omitempty"`

	// EnvRef selects value from an environment variable
	// +kubebuilder:validation:Optional
	EnvRef *EnvRef `json:"envRef,omitempty"`
}

type SecretKeyRef struct {
	// Name of the secret
	Name string `json:"name"`

	// Key of the secret
	Key string `json:"key"`

	// Source describes which secret manager to use. By default its k8s
	// +kubebuilder:validation:Optional
	Source string `json:"source,omitempty"`
}

// EnvRef describes a reference to an environment variable
type EnvRef struct {
	// Name of the environment variable
	Name string `json:"name"`
}

// AuthSecret describes the state of common properties required in every auth secret
type AuthSecret struct {
	// IsPrimary denotes if this secret is to be used as the default secret
	IsPrimary bool `json:"isPrimary"`

	// The kid value of this secret
	KID string `json:"kid"`

	// AllowedAudiences is describes the allowed values in the "aud" field of the jwt token.
	// +kubebuilder:validation:Optional
	AllowedAudiences []string `json:"allowedAudiences,omitempty"`

	// AllowedIssuers describes the allowed values in the "iss" field of the jwt token.
	// +kubebuilder:validation:Optional
	AllowedIssuers []string `json:"allowedIssuers,omitempty"`
}
