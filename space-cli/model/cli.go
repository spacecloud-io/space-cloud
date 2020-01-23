package model

type LoginResponse struct {
	Token    string      `json:"token" yaml:"token"`
	Projects []*Projects `json:"projects" yaml:"projects"`
	Error    string      `json:"error"`
}

type Credential struct {
	Accounts        []Account `json:"accounts" yaml:"accounts"`
	SelectedAccount string    `json:"selectedAccount" yaml:"selectedAccount"`
}

type Account struct {
	ID        string `json:"id" yaml:"id"`
	UserName  string `json:"username" yaml:"username"`
	Key       string `json:"key" yaml:"key"`
	ServerUrl string `json:"serverurl" yaml:"serverurl"`
}

type Projects struct {
	Name         string         `json:"name" yaml:"name"`
	ID           string         `json:"id" yaml:"id"`
	Environments []*Environment `json:"environment,omitempty" yaml:"environment"`
}

type Environment struct {
	Name     string    `json:"name" yaml:"name"`
	ID       string    `json:"id" yaml:"id"`
	Clusters []Cluster `json:"clusters" yaml:"clusters"`
}

type Cluster struct {
	ID  string `json:"id" yaml:"id"`
	URL string `json:"url" yaml:"url"`
}
