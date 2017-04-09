package postgres

import (
	"context"

	"bitbucket.org/sketchground/ajournal/profile"
	"bitbucket.org/sketchground/ajournal/utils/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DBProfile struct {
	UserID int64
	Name   string
	Email  string
}

type profileRepo struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewProfileRepo(db *sqlx.DB, logger logger.Logger) profile.Repository {
	return &profileRepo{db: db, logger: logger}
}

func (pr *profileRepo) Create(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	var id int64
	err := pr.db.Get(&id, "INSERT INTO Profile(UserID, Name, Email) VALUES($1, $2, $3) RETURNING UserID", p.ID, p.Name, p.Email)
	if err != nil {
		return nil, errors.Wrap(err, "ProfileRepo:Create failed")
	}
	p.ID = id
	return p, nil
}

func (pr *profileRepo) Update(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	_, err := pr.db.Exec("UPDATE Profile SET Name = $1, Email = $2", p.Name, p.Email)
	if err != nil {
		return nil, errors.Wrap(err, "ProfileRepo:Update failed")
	}
	return p, nil
}

func (pr *profileRepo) FindByID(ctx context.Context, id int64) (*profile.Profile, error) {
	prof := &DBProfile{}
	err := pr.db.Get(prof, "SELECT * FROM Profile WHERE UserID=$1", id)
	if err != nil {
		pr.logger.Error(ctx, err)
		return nil, profile.ErrProfileNotExist
	}
	return toProfile(prof), nil
}

func toProfile(p *DBProfile) *profile.Profile {
	return &profile.Profile{
		ID:    p.UserID,
		Name:  p.Name,
		Email: p.Email,
	}
}
