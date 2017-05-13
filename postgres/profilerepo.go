package postgres

import (
	"context"

	"bitbucket.org/sketchground/ajournal/profile"
	"bitbucket.org/sketchground/ajournal/utils/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type dbProfile struct {
	UserID      int64
	Name        string
	Email       string
	Description string
}

type profileRepo struct {
	db     *sqlx.DB
	logger logger.Logger
}

// NewProfileRepo returns a new profile repository
func NewProfileRepo(db *sqlx.DB, logger logger.Logger) profile.Repository {
	return &profileRepo{db: db, logger: logger}
}

func (pr *profileRepo) Create(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	var id int64
	err := pr.db.Get(&id, "INSERT INTO Profile(UserID, Name, Email, Description) VALUES($1, $2, $3, $4) RETURNING UserID", p.ID, p.Name, p.Email, p.Description)
	if err != nil {
		return nil, errors.Wrap(err, "ProfileRepo:Create failed")
	}
	p.ID = id
	return p, nil
}

func (pr *profileRepo) Update(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	res, err := pr.db.Exec("UPDATE Profile SET Name = $1, Email = $2, Description = $3 WHERE UserID=$4", p.Name, p.Email, p.Description, p.ID)
	if err != nil {
		return nil, errors.Wrap(err, "ProfileRepo:Update failed")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, errors.Wrap(err, "ProfileRepo:rows affected")
	}
	if rows != 1 {
		return nil, errors.New("ProfileRepo: rows affected is different from 1")
	}
	return p, nil
}

func (pr *profileRepo) FindByID(ctx context.Context, id int64) (*profile.Profile, error) {
	prof := &dbProfile{}
	err := pr.db.Get(prof, "SELECT * FROM Profile WHERE UserID=$1", id)
	if err != nil {
		pr.logger.Error(ctx, err)
		return nil, profile.ErrProfileNotExist
	}
	return toProfile(prof), nil
}

func toProfile(p *dbProfile) *profile.Profile {
	return &profile.Profile{
		ID:          p.UserID,
		Name:        p.Name,
		Email:       p.Email,
		Description: p.Description,
	}
}
