package interfaces

type TokenProvider interface {
	GetValidToken() (string, error)
}
