package client

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	BaseUrl  string `json:"baseUrl"`
}
