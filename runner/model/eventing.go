package model

// CloudEventPayload is the the JSON event spec by Cloud Events Specification
type CloudEventPayload struct {
	SpecVersion string `json:"specversion"`
	Type        string `json:"type"`
	Source      string `json:"source"`
	ID          string `json:"id"`
	Time        string `json:"time"`
	Data        struct {
		Path string         `json:"path"`
		Meta ServiceRequest `json:"meta"`
	} `json:"data"`
}

// ServiceRequest is the meta format of the meta data received from artifact store
type ServiceRequest struct {
	IsDeploy bool     `json:"isDeploy"`
	Service  *Service `json:"service"`
}
