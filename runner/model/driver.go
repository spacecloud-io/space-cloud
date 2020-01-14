package model

// DriverType is used to describe which deployment target is to be used
type DriverType string

const (
	// TypeIstio is the driver type used to target istio on kubernetes
	TypeIstio DriverType = "istio"
)
