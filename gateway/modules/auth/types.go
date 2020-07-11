package auth

type adminMan interface {
	GetSecret() string
}
