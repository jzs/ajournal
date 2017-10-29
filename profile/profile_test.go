package profile_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sketchground/ajournal/profile"
	"github.com/sketchground/ajournal/user"
)

func TestService(t *testing.T) {
	pr := NewInmemRepo()
	sr := NewInmemSubRepo()
	ps := profile.NewService(pr, sr, nil)

	u := &user.User{
		ID:       201,
		Username: "jzs",
	}
	ctx := user.TestContextWithUser(u)

	prof := &profile.Profile{ID: 201, Name: "Bobo", Email: "bobo@bobo.bobo"}
	p, err := ps.Create(ctx, prof)
	if err != nil {
		t.Fatalf("Expected creating profile, got: %v", err.Error())
	}
	if p.Name != prof.Name {
		t.Fatalf("Expected %v profile, got: %v", prof.Name, p.Name)
	}
	if p.ShortName != "bobo" {
		t.Fatalf("Expected %v shortname, got: %v", "bobo", p.ShortName)
	}

	p, err = ps.Profile(ctx)
	if err != nil {
		t.Fatalf("Expected fetching profile, got: %v", err.Error())
	}
	if p.Name != prof.Name {
		t.Fatalf("Expected profile name to be empty, got: %v", p.Name)
	}
	if p.Plan != profile.PlanFree {
		t.Fatalf("Expected to have the plan set to free, got: %v", p.Plan)
	}
	if p.ID != u.ID {
		t.Fatalf("Profile ID %v, got %v", u.ID, p.ID)
	}

	// Test update profile
	p.Description = "Hello World"
	pupdated, err := ps.UpdateProfile(ctx, p)
	if err != nil {
		t.Fatalf("Expected update success, got: %v", err.Error())
	}
	if pupdated.Description != p.Description {
		t.Fatalf("Expected update description: %v, got: %v", p.Description, pupdated.Description)
	}

	np, err := ps.Profile(ctx)
	if err != nil {
		t.Fatalf("Expected fetching profile, got: %v", err.Error())
	}
	if np.Name != p.Name {
		t.Fatalf("Expected %v, Got: %v", p.Name, np.Name)
	}

	np.ID = 1
	_, err = ps.UpdateProfile(ctx, np)
	if err == nil {
		t.Fatalf("Expected an error updating other persons profile, got %v", np)
	}

	// Test Profile GET when it doesn't already exist
	u = &user.User{
		ID:       202,
		Username: "ok",
	}
	ctx = user.TestContextWithUser(u)

	p, err = ps.Profile(ctx)
	if err != nil {
		t.Fatalf("Expected fetching profile, got: %v", err.Error())
	}
	if p.Name != "" {
		t.Fatalf("Expected profile name to be empty, got: %v", p.Name)
	}
	if p.Plan != profile.PlanFree {
		t.Fatalf("Expected to have the plan set to free, got: %v", p.Plan)
	}
	if p.ID != u.ID {
		t.Fatalf("Profile ID %v, got %v", u.ID, p.ID)
	}

	// Test Create subscription!
	sub := &profile.Subscription{}
	err = ps.Subscribe(ctx, sub)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err.Error())
	}
}

func TestGenerateShortName(t *testing.T) {
	pr := NewInmemRepo()
	sr := NewInmemSubRepo()
	ps := profile.NewService(pr, sr, nil)

	tests := map[string]string{
		"bob@cat.de":  "bob",
		"bob1@cat.de": "bob1",
		"bobx@cat.de": "bobx",
		"bob;@cat.de": "bob-",
		"b;ob@cat.de": "b-ob",
	}
	for i, o := range tests {
		r := ps.GenerateShortName(i)
		if r != o {
			t.Errorf("Generated short: %v does not match: %v", r, o)
		}
	}
}

func TestValidateShortName(t *testing.T) {
	pr := NewInmemRepo()
	sr := NewInmemSubRepo()
	ps := profile.NewService(pr, sr, nil)

	tests := map[string]bool{
		"bob@cat.de": false,
		"bob":        true,
		"bobcat-de":  true,
		"bob-cat.de": false,
	}
	for i, o := range tests {
		r := ps.ValidateShortName(context.Background(), 1, i)
		if r != o {
			t.Errorf("Validated short: %v is: %v expected: %v", i, r, o)
		}
	}
}

type profileRepo struct {
	profiles []*profile.Profile
}

func NewInmemRepo() profile.Repository {
	repo := &profileRepo{
		profiles: []*profile.Profile{},
	}
	return repo
}

func (pr *profileRepo) Create(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	for _, pp := range pr.profiles {
		// Uniqueness check
		if pp.ShortName == p.ShortName {
			return nil, errors.New("Profile already exists with short name: " + p.ShortName)
		}
	}
	pr.profiles = append(pr.profiles, p)
	return p, nil
}

func (pr *profileRepo) Update(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	for i, prof := range pr.profiles {
		if prof.ID == p.ID {
			pr.profiles[i] = p
			break
		}
	}
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

func (pr *profileRepo) FindByShortName(ctx context.Context, sn string) (*profile.Profile, error) {
	for _, p := range pr.profiles {
		if p.ShortName == sn {
			return p, nil
		}
	}
	return nil, profile.ErrProfileNotExist
}

type subRepo struct {
	subs map[string]*profile.Subscription
}

func NewInmemSubRepo() profile.SubscriptionRepository {
	return &subRepo{subs: map[string]*profile.Subscription{}}
}

func (sr *subRepo) Create(ctx context.Context, s *profile.Subscription) (*profile.Subscription, error) {
	sr.subs[s.Token] = s
	return s, nil
}
