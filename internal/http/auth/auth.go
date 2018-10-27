package auth

// Auth defines an interface for authenticating API calls
type Auth interface {
	Authenticate(authHeader string) (User, error)
}
