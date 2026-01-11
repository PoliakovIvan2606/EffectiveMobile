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

type Subscription struct {
	ServiceName string
	Price int
	UserId string
	StartDate time.Time 
	EndDate time.Time
}

type SubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price int `json:"price"`
	UserId string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}

func(s *SubscriptionRequest) Validate() error {
	errList := []error{}
	if s.ServiceName == "" || s.UserId == "" || s.StartDate == "" || s.EndDate == "" {
		errList = append(errList, ErrEmptyField)
	}

	if len(s.ServiceName) > 100 {
		errList = append(errList, ErrServiceNameLong)
	}

	if s.Price < 0 {
		errList = append(errList, ErrNegativPrice)
	}

	_, err := uuid.Parse(s.UserId)
	if err != nil {
		errList = append(errList, ErrUserIdNotValid)
	}

	if len(errList) != 0 {
		return errors.Join(errList...)
	}

	return nil
}

func(s *SubscriptionRequest) ToDomain() (*Subscription, error) {
	errList := []error{}
	StartDateTime, err := utils.ValidDate(s.StartDate)
	if err != nil {
		errList = append(errList, fmt.Errorf("StartDate %w", ErrDataNotValid))
	}

	EndDateTime, err := utils.ValidDate(s.EndDate)
	if err != nil {
		errList = append(errList, fmt.Errorf("EndDate %w", ErrDataNotValid))
	}

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