package profile_test

import (
	"context"
	"testing"

	"bitbucket.org/sketchground/ajournal/profile"
	"bitbucket.org/sketchground/ajournal/user"
)

func TestService(t *testing.T) {
	pr := NewInmemRepo()
	sr := NewInmemSubRepo()
	ps := profile.NewService(pr, sr)

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
	if p.Plan != profile.PlanFree {
		t.Fatalf("Expected to have the plan set to free, got: %v", p.Plan)
	}
	if p.ID != u.ID {
		t.Fatalf("Profile ID %v, got %v", u.ID, p.ID)
	}

	// Test update profile
	_, err = ps.UpdateProfile(ctx, p)
	if err != nil {
		t.Fatalf("Expected update success, got: %v", err.Error())
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
		t.Fatal("Expected an error updating other persons profile, got %v", np)
	}

	// Test Create subscription!
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
