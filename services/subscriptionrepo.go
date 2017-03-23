package services

import (
	"context"

	"bitbucket.org/sketchground/ajournal/profile"
	"bitbucket.org/sketchground/ajournal/utils/logger"
	"github.com/jmoiron/sqlx"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
)

type subRepo struct {
	db     *sqlx.DB
	client *client.API
}

func NewStripeSubscriptionRepo(skkey string, db *sqlx.DB) profile.SubscriptionRepository {
	repo := &subRepo{db: db}
	repo.client = &client.API{}
	repo.client.Init(skkey, nil)
	return &subRepo{}
}

func (sr *subRepo) Create(ctx context.Context, s *profile.Subscription) (*profile.Subscription, error) {

	params := &stripe.CustomerParams{
		Balance: -123,
		Desc:    "Stripe Developer",
		Email:   "gostripe@stripe.com",
	}
	params.SetSource(&stripe.CardParams{
		Name:   "Go Stripe",
		Number: "378282246310005",
		Month:  "06",
		Year:   "15",
	})

	// Create customer!
	cust, err := customer.New(params)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	subparams := &stripe.SubParams{
		Customer:    cust.ID,
		Plan:        "basic-monthly",
		TrialPeriod: 14,
	}
	subscription, err := sub.New(subparams)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	// Store subscription info in db.
	_, err = sr.db.Exec("INSERT INTO Subscription(userid, subscriptionid) VALUES($1, $2)", s.Profile.ID, subscription.ID)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	return s, nil
}
