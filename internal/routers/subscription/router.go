package subscription

import (
	models "EffectiveMobile/internal/models/subscription"
	"context"
	"time"

	"github.com/gorilla/mux"
)

type UseCaseSubscription interface {
	AddSubscription(ctx context.Context, s *models.SubscriptionRequest) (int, error)
	GetSubscription(ctx context.Context, id int) (*models.SubscriptionResponse, error)
	UpdateSubscription(ctx context.Context, s *models.SubscriptionRequest, id int) error
	DeleteSubscription(ctx context.Context, id int) error
	GetTotalCost(ctx context.Context, userID string, serviceName string, startDate, endDate time.Time) (float64, error)
}


type SubscriptionRouter struct {
	UC UseCaseSubscription
}

func InitRouter(r *mux.Router, UC UseCaseSubscription) {
	serviceRouter := SubscriptionRouter{UC: UC}
	chat := r.PathPrefix("/subscription").Subrouter()
	chat.HandleFunc("/stats", serviceRouter.GetTotalCost).Methods("GET")
	chat.HandleFunc("/{id}", serviceRouter.GetSubscription).Methods("GET")
	chat.HandleFunc("", serviceRouter.AddSubscription).Methods("POST")
	chat.HandleFunc("/{id}", serviceRouter.UpdateSubscription).Methods("PUT")
	chat.HandleFunc("/{id}", serviceRouter.DeleteSubscription).Methods("DELETE")
}