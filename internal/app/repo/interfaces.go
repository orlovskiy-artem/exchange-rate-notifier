package repo

import (
	"context"
	"genesis_test_task/internal/app/model"
)

type ISubscriptionRepo interface {
	SaveSubscription(ctx context.Context, s model.Subscription) error
	FindSubscription(ctx context.Context, e model.Email) (*model.Subscription, bool, error)
	ForEachSubscription(ctx context.Context, hf model.SubscriptionHandle) error
}

type INotificationRepo interface {
	SendEmailExchangeRate(ctx context.Context, addressee model.Email, exchangeRate model.ExchangeRate) error
}

type IExchangeRateRepo interface {
	GetExchangeRate(ctx context.Context) (model.ExchangeRate, error)
}
