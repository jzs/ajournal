package profile_test

import (
	"context"
	"net/http"
	"testing"

	"bitbucket.org/sketchground/ajournal/profile"
	"bitbucket.org/sketchground/ajournal/user"
)

type logger struct{}

func (l *logger) Error(ctx context.Context, err error)                                    {}
func (l *logger) Errorf(ctx context.Context, format string, args ...interface{})          {}
func (l *logger) Print(ctx context.Context, err error)                                    {}
func (l *logger) Printf(ctx context.Context, format string, args ...interface{})          {}
func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) { next(w, r) }

func TestService(t *testing.T) {
	pr := NewInmemRepo()
	sr := NewInmemSubRepo()
	ps := profile.NewService(pr, sr)

	u := &user.User{
		ID:       201,
		Username: "jzs",
	}
	ctx := user.TestContextWithUser(u)

	prof := &profile.Profile{ID: 201, Name: "Bobo"}
	p, err := ps.Create(ctx, prof)
	if err != nil {
		t.Fatalf("Expected creating profile, got: %v", err.Error())
	}
	if p.Name != prof.Name {
		t.Fatalf("Expected %v profile, got: %v", prof.Name, p.Name)
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
