package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/jzs/ajournal/blob"
	"github.com/jzs/ajournal/profile"
	"github.com/jzs/ajournal/utils/logger"
)

type dbProfile struct {
	UserID      int64
	Name        string
	Email       string
	Description string
	ShortName   sql.NullString // UNIQUE
	Picture     sql.NullString
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
	err := pr.db.Get(&id, "INSERT INTO Profile(UserID, Name, ShortName, Email, Description, Picture) VALUES($1, $2, $3, $4, $5, $6) RETURNING UserID", p.ID, p.Name, p.ShortName, p.Email, p.Description, p.Picture.Key)
	if err != nil {
		return nil, errors.Wrap(err, "ProfileRepo:Create failed")
	}
	p.ID = id
	return p, nil
}

func (pr *profileRepo) Update(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	res, err := pr.db.Exec("UPDATE Profile SET Name = $1, Email = $2, Description = $3, ShortName = $5, Picture = $6 WHERE UserID=$4", p.Name, p.Email, p.Description, p.ID, p.ShortName, p.Picture.Key)
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
		if err == sql.ErrNoRows {
			return nil, profile.ErrProfileNotExist
		}
		pr.logger.Error(ctx, err)
		return nil, errors.Wrap(err, "ProfileRepo:FindByID")
	}
	return toProfile(prof), nil
}

func (pr *profileRepo) FindByShortName(ctx context.Context, sn string) (*profile.Profile, error) {
	prof := &dbProfile{}
	err := pr.db.Get(prof, "SELECT * FROM Profile WHERE ShortName=$1", sn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, profile.ErrProfileNotExist
		}
		pr.logger.Error(ctx, err)
		return nil, errors.Wrap(err, "ProfileRepo:FindByShortName")
	}
	return toProfile(prof), nil
}

func toProfile(p *dbProfile) *profile.Profile {
	prof := &profile.Profile{
		ID:          p.UserID,
		Name:        p.Name,
		Email:       p.Email,
		Description: p.Description,
	}
	if p.ShortName.Valid {
		prof.ShortName = p.ShortName.String
	}
	if p.Picture.Valid {
		prof.Picture = blob.FileFromKey(p.Picture.String)
	}

	return prof
}
