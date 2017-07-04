package user

import (
	"context"
	"net/http"
	"time"

	"github.com/sketchground/ajournal/utils"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"golang.org/x/crypto/bcrypt"
)

// Service interface for users
type Service interface {
	// Register creates a new user
	Register(ctx context.Context, u *User) error
	// Activate activates an account by processing activate token
	//Activate(ctx context.Context, activateToken string) error
	// Login creates a new token for a given user
	Login(ctx context.Context, username string, password string) (*Token, error)
	// Logout invalidates a user token
	Logout(ctx context.Context, token string) error
	// User fetches a user with the given username
	User(ctx context.Context, username string) (*User, error)
	// UserWithToken fetches a user from a valid token
	UserWithToken(ctx context.Context, token string) (*User, error)
}

// NewService returns a service implementation of the user service
func NewService(t *utils.Translator, repo Repository) Service {
	return &service{repo: repo, t: t}
}

type service struct {
	repo Repository
	t    *utils.Translator
}

func (s *service) Register(ctx context.Context, u *User) error {
	tr := s.t.T(ctx)
	// TODO: Check for password complexity
	if u.Username == "" {
		return utils.NewAPIError(nil, http.StatusBadRequest, tr("user.missing-username").String())
	}
	if u.Password == "" {
		return utils.NewAPIError(nil, http.StatusBadRequest, tr("user.missing-password").String())
	}
	existing, err := s.repo.FindByUsername(ctx, u.Username)
	if err != nil && err != ErrUserNotExist {
		return err
	}
	if existing != nil {
		return errors.New("Username already exists")
	}

	u.Active = true
	u.Created = time.Now()
	enc, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return errors.New("Failed encrypting password")
	}
	u.Password = string(enc)

	_, err = s.repo.Create(ctx, u)
	return err
}

func (s *service) Login(ctx context.Context, username string, password string) (*Token, error) {
	tr := s.t.T(ctx)
	u, err := s.repo.FindByUsername(ctx, username)
	switch {
	case err == ErrUserNotExist:
		return nil, utils.NewAPIError(err, http.StatusNotFound, tr("user.userpasswrong").String())
	case err != nil:
		return nil, errors.Wrap(err, "Login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	switch {
	case err == bcrypt.ErrMismatchedHashAndPassword:
		return nil, utils.NewAPIError(err, http.StatusNotFound, tr("user.userpasswrong").String())
	case err != nil:
		return nil, errors.Wrap(err, "Login")
	}

	u1 := uuid.NewV4()
	token := &Token{
		Token:   u1.String(),
		Expires: time.Now().Add(24 * 7 * time.Hour),
		UserID:  u.ID,
	}

	err = s.repo.CreateToken(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "Login")
	}
	return token, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	err := s.repo.DeleteToken(ctx, token)
	return err
}

func (s *service) User(ctx context.Context, username string) (*User, error) {
	// TODO: Make auth check to see how much of the user info we should return!
	u, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	u.Password = "" // Setting password to blank such that we are not leaking the hash by mistake.
	return u, nil
}

func (s *service) UserWithToken(ctx context.Context, token string) (*User, error) {
	u, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	u.Password = "" // Setting password to blank such that we are not leaking the hash by mistake.
	return u, nil
}
