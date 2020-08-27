package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/jzs/ajournal/oauth"
	"github.com/jzs/ajournal/user"
)

type oauthRepo struct {
	db *sqlx.DB
}

// NewOauthRepo returns a postgres implementation of the oauth repository
func NewOauthRepo(db *sqlx.DB) oauth.Repository {
	return &oauthRepo{db: db}
}

func (or *oauthRepo) Find(ctx context.Context, userID int64) ([]oauth.UserProvider, error) {
	u := &oauth.UserProvider{}
	err := or.db.Get(u, "SELECT * FROM oauthuser WHERE userid=$1", userID)
	if err != nil {
		return nil, user.ErrUserNotExist
	}
	return []oauth.UserProvider{*u}, nil
}

func (or *oauthRepo) Create(ctx context.Context, p oauth.UserProvider) error {
	_, err := or.db.Exec("INSERT INTO oauthuser(userid, provider, providerusername) VALUES($1, $2, $3)", p.UserID, p.Provider, p.ProviderUsername)
	if err != nil {
		return errors.Wrap(err, "OauthRepo.Create")
	}
	return nil
}
