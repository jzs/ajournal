package profile_test

import (
	"testing"

	"bitbucket.org/sketchground/journal/profile"
	"bitbucket.org/sketchground/journal/user"
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
	if p.Name != "Hanzi" {
		t.Fatalf("Expected profile name to be Hanzi, got: %v", p.Name)
	}
}

type profileRepo struct {
	id int64
}

func NewInmemRepo() profile.Repository {
	repo := &profileRepo{
		id: 1,
	}
	return repo
}
