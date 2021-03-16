package types

// Storage is a struct to store the Live Query Data (For Internal Use Only)
type Storage struct {
	Id        string
	Time      int64
	Payload   []byte
	IsDeleted bool
}
