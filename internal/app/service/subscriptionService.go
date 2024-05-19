package service

import (
	"context"
	"fmt"

	"genesis_test_task/internal/app/model"
	"genesis_test_task/internal/app/repo"
)

type SubscriptionService struct {
	subRepo repo.ISubscriptionRepo
	notRepo repo.INotificationRepo
	exrepo  repo.IExchangeRateRepo
}

func NewSubscriptionService(
	sr repo.ISubscriptionRepo,
	nr repo.INotificationRepo,
	er repo.IExchangeRateRepo) *SubscriptionService {

	return &SubscriptionService{subRepo: sr, notRepo: nr, exrepo: er}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, email model.Email) error {
	err := email.Validate()
	if err != nil {
		return fmt.Errorf("invalid email: %v", err)
	}
	sub := model.NewSubscription(email)
	return s.subRepo.SaveSubscription(ctx, sub)
}

func (s *SubscriptionService) CheckSubscriberExists(ctx context.Context, email model.Email) (found bool, err error) {
	err = email.Validate()
	if err != nil {
		return false, fmt.Errorf("invalid email: %v", err)
	}
	_, found, err = s.subRepo.FindSubscription(ctx, email)
	if err != nil {
		return false, err
	}
	return found, nil

}

// TODO
func (s *SubscriptionService) NotifyAll(ctx context.Context) error {
	exr, err := s.exrepo.GetExchangeRate(ctx)
	if err != nil {
		return err
	}
	return s.subRepo.ForEachSubscription(ctx,
		func(ctx context.Context, sub model.Subscription) error {
			return s.notRepo.SendEmailExchangeRate(ctx, sub.Email, exr)
		})
}
