package postgres

import (
	"context"

	"bitbucket.org/sketchground/ajournal/user"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

// NewUserRepo returns a postgres implementation of the user repository
func NewUserRepo(db *sqlx.DB) user.Repository {
	return &repo{db: db}
}

func (ur *repo) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	u := &user.User{}
	err := ur.db.Get(u, "SELECT * FROM _User WHERE Username=$1", username)
	if err != nil {
		return nil, user.ErrUserNotExist
	}
	return u, nil
}

func (ur *repo) FindByToken(ctx context.Context, token string) (*user.User, error) {
	u := &user.User{}
	err := ur.db.Get(u, "SELECT _User.* FROM _User JOIN UserToken on UserToken.UserID=_User.id WHERE UserToken.Token=$1", token)
	if err != nil {
		return nil, user.ErrUserNotExist
	}
	return u, nil
}

func (ur *repo) Create(ctx context.Context, u *user.User) (*user.User, error) {
	var id int64
	err := ur.db.Get(&id, "INSERT INTO _User(Username, Password, Active, Created) VALUES($1, $2, $3, $4) RETURNING id", u.Username, u.Password, u.Active, u.Created)
	if err != nil {
		return nil, err
	}
	u.ID = id
	return u, nil
}

func (ur *repo) CreateToken(ctx context.Context, t *user.Token) error {
	_, err := ur.db.Exec("INSERT INTO UserToken(Token, UserID, Expires) VALUES($1,$2,$3)", t.Token, t.UserID, t.Expires)
	if err != nil {
		return err
	}
	return nil
}

func (ur *repo) DeleteToken(ctx context.Context, token string) error {
	_, err := ur.db.Exec("DELETE FROM UserToken WHERE Token=$1", token)
	if err != nil {
		return err
	}
	return nil
}
