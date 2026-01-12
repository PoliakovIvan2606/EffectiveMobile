package models

import "time"

// Модель для получения подписок из БД
type GetSubscription struct {
	Id int
	ServiceName string
	Price int
	UserId string
	StartDate time.Time 
	EndDate time.Time
}

// Модель для отправки клиенту моедели
type SubscriptionResponse struct {
	Id int `json:"id"`
	ServiceName string `json:"service_name"`
	Price int `json:"price"`
	UserId string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}

// Метод для получения модели отправки клиенту из модели получения из БД
func(s *GetSubscription) FromDomain() *SubscriptionResponse {
    return &SubscriptionResponse{
		Id: 		 s.Id,
        ServiceName: s.ServiceName,
        Price:       s.Price,
        UserId:      s.UserId,
        StartDate: s.StartDate.Format("01-2006"),
		EndDate:   s.EndDate.Format("01-2006"),
    }
}