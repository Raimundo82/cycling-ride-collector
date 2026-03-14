package auth_interfaces

type TokenProvider interface {
	GetValidToken() (string, error)
}
