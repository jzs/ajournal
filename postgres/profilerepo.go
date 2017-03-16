package postgres

import (
	"context"

	"bitbucket.org/sketchground/ajournal/profile"
	"bitbucket.org/sketchground/ajournal/utils/logger"
	"github.com/jmoiron/sqlx"
)

type DBProfile struct {
	UserID int64
	Name   string
	Email  string
}

type profileRepo struct {
	db *sqlx.DB
}

func NewProfileRepo(db *sqlx.DB) profile.Repository {
	return &profileRepo{db: db}
}

func (pr *profileRepo) Create(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	var id int64
	err := pr.db.Get(&id, "INSERT INTO Profile(UserID, Name, Email) VALUES($1, $2, $3) RETURNING UserID", p.ID, p.Name, p.Email)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}
	p.ID = id
	return p, nil
}

func (pr *profileRepo) FindByID(ctx context.Context, id int64) (*profile.Profile, error) {
	prof := &DBProfile{}
	err := pr.db.Get(prof, "SELECT * FROM Profile WHERE UserID=$1", id)
	if err != nil {
		logger.Error(ctx, err)
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
