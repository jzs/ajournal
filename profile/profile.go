package profile

import (
	"context"
	"errors"
)

// ErrProfileNotExist custom error
var ErrProfileNotExist error

func init() {
	ErrProfileNotExist = errors.New("Profile does not exist")
}

// Repository describes the methods on the Profile repository
type Repository interface {
	Create(ctx context.Context, p *Profile) (*Profile, error)
	Update(ctx context.Context, p *Profile) (*Profile, error)
	FindByID(ctx context.Context, id int64) (*Profile, error)
	FindByShortName(ctx context.Context, sn string) (*Profile, error)
}

// SubscriptionRepository describes the methods on the subscription repository
type SubscriptionRepository interface {
	Create(ctx context.Context, s *Subscription) (*Subscription, error)
}

// Plan Subscription plan
type Plan int64

func (p Plan) String() string {
	switch p {
	case PlanFree:
		return "Free"
	case PlanPaid:
		return "Paid"
	}
	panic("Plan set to an invalid value!")
}

const (
	// PlanFree free plan constant
	PlanFree = iota + 1
	// PlanPaid paid plan constant
	PlanPaid
)

// Profile is a user profile
type Profile struct {
	ID          int64
	Name        string
	Email       string
	Plan        Plan
	ShortName   string // Short name of the users profile. Used for public links
	Description string
}

// Subscription is a user subscription
type Subscription struct {
	Token   string
	Profile *Profile
	Plan    Plan
}
