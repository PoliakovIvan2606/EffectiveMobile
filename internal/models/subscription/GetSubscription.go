package models

import "time"

type GetSubscription struct {
	Id int
	ServiceName string
	Price int
	UserId string
	StartDate time.Time 
	EndDate time.Time
}

type SubscriptionResponse struct {
	Id int `json:"id"`
	ServiceName string `json:"service_name"`
	Price int `json:"price"`
	UserId string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}

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