package oauth

import "context"

// Provider represents a given provider
type Provider string

const (
	// ProviderGoogle google oauth provider
	ProviderGoogle = "google"
)

// UserProvider represents a relation between a given provider and a user.
type UserProvider struct {
	UserID           int64    // UserID of the user
	Provider         Provider // Identifier of the provider used
	ProviderUsername string   // The email used by the provider as authentication
}

// GoogleUserInfo represents user info for a user
type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

// Repository interface for Oauth support
type Repository interface {
	Find(ctx context.Context, userID int64) ([]UserProvider, error)
	Create(ctx context.Context, p UserProvider) error
}
