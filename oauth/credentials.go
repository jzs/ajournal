package oauth

// Credentials describes credentials for a specific oauth provider
type Credentials struct {
	Provider     Provider
	ClientID     string
	ClientSecret string
	RedirectURL  string
}
