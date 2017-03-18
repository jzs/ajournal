package services

import (
	"context"
	"errors"

	"bitbucket.org/sketchground/ajournal/profile"
)

type subRepo struct {
}

func NewStripeSubscriptionRepo() profile.SubscriptionRepository {
	return &subRepo{}
}

func (sr *subRepo) Create(ctx context.Context, s *profile.Subscription) (*profile.Subscription, error) {
	return nil, errors.New("Not implemented")
}
