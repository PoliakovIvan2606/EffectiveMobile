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


// AddSubscription
// Создаёт новую подписку
func (UC *UseCaseSubscription) AddSubscription(ctx context.Context, s *models.SubscriptionRequest) (int, error) {
	// Преобразуем входной DTO в доменную модель
	domain, err := s.ToDomain()
	if err != nil {
		return 0, err
	}

	// Сохраняем подписку через репозиторий
	id, err := UC.repo.AddSubscription(ctx, domain)
	if err != nil {
		return 0, err
	}

	// Возвращаем ID созданной записи
	return id, nil
}

// GetSubscription
// Получает одну подписку по ID
func (UC *UseCaseSubscription) GetSubscription(ctx context.Context, id int) (*models.SubscriptionResponse, error) {

	// Получаем доменную модель из репозитория
	domain, err := UC.repo.GetSubscription(ctx, id)
	if err != nil {
		return nil, err
	}

	// Преобразуем доменную модель в DTO ответа
	return domain.FromDomain(), nil
}

// UpdateSubscription
// Обновляет существующую подписку по ID
func (UC *UseCaseSubscription) UpdateSubscription(ctx context.Context, s *models.SubscriptionRequest, id int) error {

	// DTO → Domain
	domain, err := s.ToDomain()
	if err != nil {
		return err
	}

	// Обновляем запись в хранилище
	if err := UC.repo.UpdateSubscription(ctx, domain, id); err != nil {
		return err
	}

	return nil
}

// DeleteSubscription
// Удаляет подписку по ID
func (UC *UseCaseSubscription) DeleteSubscription(ctx context.Context, id int) error {

	// Удаляем подписку из репозитория
	if err := UC.repo.DeleteSubscription(ctx, id); err != nil {
		return err
	}

	return nil
}

// GetTotalCost
// Считает общую стоимость подписок пользователя
// за период по конкретному сервису
func (UC *UseCaseSubscription) GetTotalCost(ctx context.Context, userID string, serviceName string, startDate, endDate time.Time) (float64, error) {

	// Делегируем расчёт репозиторию (SQL aggregation)
	totalCost, err := UC.repo.GetTotalCost(
		ctx,
		userID,
		serviceName,
		startDate,
		endDate,
	)
	if err != nil {
		return 0, err
	}

	return totalCost, nil
}

// GetListSubscription
// Возвращает список подписок пользователя по UUID
func (UC *UseCaseSubscription) GetListSubscription(ctx context.Context, UUID string) ([]models.SubscriptionResponse, error) {

	// Получаем список доменных моделей
	domains, err := UC.repo.GetListSubscription(ctx, UUID)
	if err != nil {
		return nil, err
	}

	// Подготавливаем slice под ответ
	subscriptions := make([]models.SubscriptionResponse, 0, len(domains))

	// Domain → Response DTO
	for _, sub := range domains {
		subscriptions = append(subscriptions, models.SubscriptionResponse{
			Id:          sub.Id,
			ServiceName: sub.ServiceName,
			Price:       sub.Price,
			StartDate:   utils.ParseDate(sub.StartDate),
			EndDate:     utils.ParseDate(sub.EndDate),
		})
	}

	return subscriptions, nil
}