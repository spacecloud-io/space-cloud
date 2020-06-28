package model

// PublicKeyPayload holds the pem encoded public key
type PublicKeyPayload struct {
	PemData string `json:"pem"`
}
