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

type Repository interface {
	Create(ctx context.Context, p *Profile) (*Profile, error)
	Update(ctx context.Context, p *Profile) (*Profile, error)
	FindByID(ctx context.Context, id int64) (*Profile, error)
}

type SubscriptionRepository interface {
	Create(ctx context.Context, s *Subscription) (*Subscription, error)
}

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
	PlanFree = iota + 1
	PlanPaid
)

type Profile struct {
	ID    int64
	Name  string
	Email string
	Plan  Plan
}

type Subscription struct {
	Token   string
	Profile *Profile
	Plan    Plan
}
