package model

// LoginResponse is the object for storing payload for login response
type LoginResponse struct {
	Token    string      `json:"token" yaml:"token"`
	Projects []*Projects `json:"projects" yaml:"projects"`
	Error    string      `json:"error"`
}

// Credential is the object for representing all the account information in accounts.yaml file
type Credential struct {
	Accounts        []*Account `json:"accounts" yaml:"accounts"`
	SelectedAccount string     `json:"selectedAccount" yaml:"selectedAccount"`
}

// Account is the object for representing individual account information
type Account struct {
	ID        string `json:"id" yaml:"id"`
	UserName  string `json:"username" yaml:"username"`
	Key       string `json:"key" yaml:"key"`
	ServerURL string `json:"serverurl" yaml:"serverurl"`
}

// Projects describes the configuration of a single project
type Projects struct {
	Name string `json:"name" yaml:"name"`
	ID   string `json:"id" yaml:"id"`
}

// Environment describes the configuration of a single environment
type Environment struct {
	Name     string    `json:"name" yaml:"name"`
	ID       string    `json:"id" yaml:"id"`
	Clusters []Cluster `json:"clusters" yaml:"clusters"`
}

// Cluster describes the configuration of a single cluster
type Cluster struct {
	ID  string `json:"id" yaml:"id"`
	URL string `json:"url" yaml:"url"`
}
