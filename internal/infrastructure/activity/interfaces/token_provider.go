package activity_interfaces

type TokenProvider interface {
	GetValidToken() (string, error)
}
