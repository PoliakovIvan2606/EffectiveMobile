package usecase

import (
	models "EffectiveMobile/internal/models/subscription"
	"EffectiveMobile/internal/utils"
	"context"
	"time"
)

type RepositorySubscription interface {
	AddSubscription(ctx context.Context, s *models.Subscription) (int, error)
	GetSubscription(ctx context.Context, id int) (*models.GetSubscription, error)
	UpdateSubscription(ctx context.Context, s *models.Subscription, id int) error
	DeleteSubscription(ctx context.Context, id int) error
	GetTotalCost(ctx context.Context, userID string, serviceName string, startDate, endDate time.Time) (float64, error)
	GetListSubscription(ctx context.Context, UUID string) ([]models.GetSubscription, error)
}

type UseCaseSubscription struct {
	repo RepositorySubscription
}

func NewUseCaseSubscription(repo RepositorySubscription) *UseCaseSubscription {
	return &UseCaseSubscription{
		repo: repo,
	}
}

func(UC *UseCaseSubscription) AddSubscription(ctx context.Context, s *models.SubscriptionRequest) (int, error) {
	domain, err := s.ToDomain()
	if err != nil {
		return 0, err
	}
	id, err := UC.repo.AddSubscription(ctx, domain)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func(UC *UseCaseSubscription) GetSubscription(ctx context.Context, id int) (*models.SubscriptionResponse, error) {
	domain, err := UC.repo.GetSubscription(ctx, id)
	if err != nil {
		return nil, err
	}
	return domain.FromDomain(), nil
}

func(UC *UseCaseSubscription) UpdateSubscription(ctx context.Context, s *models.SubscriptionRequest, id int) error {
	domain, err := s.ToDomain()
	if err != nil {
		return err
	}
	if err := UC.repo.UpdateSubscription(ctx, domain, id); err != nil {
		return err
	}
	return nil
}

func(UC *UseCaseSubscription) DeleteSubscription(ctx context.Context, id int) error {
	err := UC.repo.DeleteSubscription(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func(UC *UseCaseSubscription) GetTotalCost(ctx context.Context, userID string, serviceName string, startDate, endDate time.Time) (float64, error) {
	totalCost, err := UC.repo.GetTotalCost(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		return 0, err
	}
	return totalCost, err
}

func(UC *UseCaseSubscription) GetListSubscription(ctx context.Context, UUID string) ([]models.SubscriptionResponse, error) {
	domains, err := UC.repo.GetListSubscription(ctx, UUID)
	if err != nil {
		return nil, err
	}

	var subscriptions []models.SubscriptionResponse
	for _, sub := range domains {
		subscriptions = append(subscriptions, models.SubscriptionResponse{
			Id: sub.Id,
			ServiceName: sub.ServiceName,
			Price: sub.Price,
			StartDate: utils.ParseDate(sub.StartDate),
			EndDate: utils.ParseDate(sub.EndDate),
		})
	}

	return subscriptions, nil
}