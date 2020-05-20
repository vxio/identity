package authn

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey string

// loginSessionContextKey is the context key for the Login Session
const loginSessionContextKey contextKey = "login-session"
