package main

type DBConfig struct {
	Host    string `json:"host"`
	Port    string `json:"port"`
	User    string `json:"user"`
	Pass    string `json:"password"`
	DB      string `json:"db"`
	SSLMode string `json:"sslMode"`
}
