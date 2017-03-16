package profile_test

import (
	"context"
	"testing"

	"bitbucket.org/sketchground/ajournal/profile"
	"bitbucket.org/sketchground/ajournal/user"
)

func TestService(t *testing.T) {
	pr := NewInmemRepo()
	ps := profile.NewService(pr)

	u := &user.User{
		ID:       201,
		Username: "jzs",
	}
	ctx := user.TestContextWithUser(u)

	p, err := ps.Profile(ctx)
	if err != nil {
		t.Fatalf("Expected fetching profile, got: %v", err.Error())
	}
	if p.Name != "" {
		t.Fatalf("Expected profile name to be empty, got: %v", p.Name)
	}
}

type profileRepo struct {
	profiles []*profile.Profile
	id       int64
}

func NewInmemRepo() profile.Repository {
	repo := &profileRepo{
		id:       1,
		profiles: []*profile.Profile{},
	}
	return repo
}

func (pr *profileRepo) Create(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	p.ID = pr.id
	pr.profiles = append(pr.profiles, p)
	pr.id = pr.id + 1
	return p, nil
}

func (pr *profileRepo) FindByID(ctx context.Context, id int64) (*profile.Profile, error) {
	for _, p := range pr.profiles {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, profile.ErrProfileNotExist
}
