package models

import (
	"EffectiveMobile/internal/utils"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyField = errors.New("поле не должно быть нулевым")
	ErrServiceNameLong = errors.New("поле ServiceName не должно длинной более 100 символов")
	ErrNegativPrice = errors.New("поле price не должно быть отрицательным")
	ErrUserIdNotValid = errors.New("поле user_id не валидно")
	ErrDataNotValid = errors.New("поле data не валидно")
)

// Модель для добавления и обновления подписки в БД
type Subscription struct {
	ServiceName string
	Price int
	UserId string
	StartDate time.Time 
	EndDate time.Time
}

// Модель для получения от клиента данных для добавлении подписки
type SubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price int `json:"price"`
	UserId string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}

// Метод валидации пришедших данных
func(s *SubscriptionRequest) Validate() error {
	// Переменная для аккумулирования ошибок пришедших данных
	errList := []error{}
	// Проверка пустых полей
	if s.ServiceName == "" || s.UserId == "" || s.StartDate == "" || s.EndDate == "" {
		errList = append(errList, ErrEmptyField)
	}

	// Проверка длинны названи подписки
	if len(s.ServiceName) > 100 {
		errList = append(errList, ErrServiceNameLong)
	}

	// Проверка правильности цены подписки
	if s.Price < 0 {
		errList = append(errList, ErrNegativPrice)
	}

	// Проверка UUID пользователя
	_, err := uuid.Parse(s.UserId)
	if err != nil {
		errList = append(errList, ErrUserIdNotValid)
	}

	// Если ошибок есть джойним их и отправлем на слой выше
	if len(errList) != 0 {
		return errors.Join(errList...)
	}

	return nil
}

// Метод получения молдели добавлени и обновления в БД из модели которая получает данные от пользователя
func(s *SubscriptionRequest) ToDomain() (*Subscription, error) {
	// переменная аккумулирующая ошибки 
	errList := []error{}

	// Получаем валидную переменную time.Time для добавления в БД
	StartDateTime, err := utils.ValidDate(s.StartDate)
	if err != nil {
		errList = append(errList, fmt.Errorf("StartDate %w", ErrDataNotValid))
	}

	EndDateTime, err := utils.ValidDate(s.EndDate)
	if err != nil {
		errList = append(errList, fmt.Errorf("EndDate %w", ErrDataNotValid))
	}

	// Если ошибок есть джойним их и отправлем на слой выше
	if len(errList) != 0 {
		return nil, errors.Join(errList...)
	}

    return &Subscription{
        ServiceName: s.ServiceName,
        Price:       s.Price,
        UserId:      s.UserId,
        StartDate:   *StartDateTime,
        EndDate:     *EndDateTime,
    }, nil
}