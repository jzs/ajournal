package profile

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/sketchground/ajournal/user"
	"github.com/sketchground/ajournal/utils"
)

// Service describes the methods on a profile service
type Service interface {
	Create(ctx context.Context, p *Profile) (*Profile, error)
	Profile(ctx context.Context) (*Profile, error)
	UserProfile(ctx context.Context, userid int64) (*Profile, error)
	UpdateProfile(ctx context.Context, p *Profile) (*Profile, error)
	Subscribe(ctx context.Context, sub *Subscription) error
}

// NewService returns a new implementation of the Service interface
func NewService(pr Repository, sr SubscriptionRepository) Service {
	return &service{pr: pr, sr: sr}
}

type service struct {
	pr Repository
	sr SubscriptionRepository
}

func (s *service) Create(ctx context.Context, p *Profile) (*Profile, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, utils.NewAPIError(nil, http.StatusForbidden, "Cannot create profile without a user context")
	}
	if usr.ID != p.ID {
		return nil, utils.NewAPIError(nil, http.StatusBadRequest, "Cannot create profile for another user")
	}

	if p.Plan == 0 {
		p.Plan = PlanFree
	}

	prof, err := s.pr.Create(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "CreateProfile")
	}
	return prof, nil

}

func (s *service) UserProfile(ctx context.Context, userid int64) (*Profile, error) {
	pro, err := s.pr.FindByID(ctx, userid)
	if err == ErrProfileNotExist {
		return nil, errors.Wrap(err, "Profile doesn't exist")
	}
	pro.Email = ""
	return pro, nil
}

func (s *service) Profile(ctx context.Context) (*Profile, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, utils.NewAPIError(nil, http.StatusForbidden, "Cannot create a journal without a user context")
	}

	pro, err := s.pr.FindByID(ctx, usr.ID)
	if err == ErrProfileNotExist {
		// Create profile and return that.
		pro, err = s.pr.Create(ctx, &Profile{ID: usr.ID, Plan: PlanFree})
		if err != nil {
			return nil, errors.Wrap(err, "Could not create profile for user")
		}
		return pro, nil
	}
	return pro, nil
}

func (s *service) UpdateProfile(ctx context.Context, p *Profile) (*Profile, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, utils.NewAPIError(nil, http.StatusForbidden, "Cannot create a journal without a user context")
	}
	if usr.ID != p.ID {
		return nil, utils.NewAPIError(nil, http.StatusBadRequest, "Cannot update another users profile")
	}

	prof, err := s.pr.Update(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateProfile")
	}
	return prof, nil
}

// SubscriptionArgs args for signing up for a subscription
type SubscriptionArgs struct {
	CardName string
	Number   string
	Month    string
	Year     string
	CVC      string
}

func (s *service) Subscribe(ctx context.Context, sub *Subscription) error {
	_, err := s.sr.Create(ctx, sub)
	return err
}
