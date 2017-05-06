package services

import (
	"context"

	"bitbucket.org/sketchground/ajournal/profile"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

type subRepo struct {
	db     *sqlx.DB
	client *client.API
}

// NewStripeSubscriptionRepo returns a new subscription repo that integrates with stripe
func NewStripeSubscriptionRepo(skkey string, db *sqlx.DB) profile.SubscriptionRepository {
	repo := &subRepo{db: db}
	repo.client = &client.API{}
	repo.client.Init(skkey, nil)
	return repo
}

func (sr *subRepo) Create(ctx context.Context, s *profile.Subscription) (*profile.Subscription, error) {

	params := &stripe.CustomerParams{
		Desc:  s.Profile.Name,
		Email: s.Profile.Email,
	}
	err := params.SetSource(&stripe.CardParams{
		Token: s.Token,
	})
	if err != nil {
		return nil, errors.Wrap(err, "SubscriptionRepo: Failed setting stripe params")
	}

	// Create customer!
	cust, err := sr.client.Customers.New(params)
	if err != nil {
		return nil, errors.Wrap(err, "SubscriptionRepo:Create failed")
	}

	// TODO Consider if plan and trial period should be specified elsewhere or as args to program
	subparams := &stripe.SubParams{
		Customer:    cust.ID,
		Plan:        "ajournal_basic",
		TrialPeriod: 14,
	}
	subscription, err := sr.client.Subs.New(subparams)
	if err != nil {
		return nil, errors.Wrap(err, "SubscriptionRepo:Create failed")
	}

	// Store subscription info in db.
	_, err = sr.db.Exec("INSERT INTO Subscription(userid, stripecustomerid, stripesubscriptionid) VALUES($1, $2, $3)", s.Profile.ID, cust.ID, subscription.ID)
	if err != nil {
		return nil, errors.Wrap(err, "SubscriptionRepo:Create failed")
	}

	return s, nil
}
