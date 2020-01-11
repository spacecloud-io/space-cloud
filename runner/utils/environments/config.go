package environments

// Config describes the configuration required by the manager
type Config struct {
	Store TypeStore
}

// TypeStore describes the type of the store to use
type TypeStore string

const (
	// SC is used when a sc gateway is used as the store
	SC TypeStore = "sc"
)
