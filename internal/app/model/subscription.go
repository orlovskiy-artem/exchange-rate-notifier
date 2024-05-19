package model

import (
	"context"

	"github.com/goware/emailx"
)

type Email string

func (e Email) String() string {
	return string(e)
}

func (e *Email) Validate() error {
	return emailx.Validate(e.String())
}

type Subscription struct {
	ID    int
	Email Email
}

func NewSubscription(e Email) Subscription {
	return Subscription{Email: e}
}

type SubscriptionHandle func(context.Context, Subscription) error
