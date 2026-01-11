package subscription

import (
	models "EffectiveMobile/internal/models/subscription"
	"EffectiveMobile/pkg/handler"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func(router SubscriptionRouter) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	in := models.SubscriptionRequest{}

	// Парсим JSON из body в структуру
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		handler.ErrResponse(w, "Неверный JSON", err, http.StatusBadRequest)
        return
    }

	vars := mux.Vars(r)

    IdStr := vars["id"]
	id, err := strconv.Atoi(IdStr)
	if err != nil {
		handler.ErrResponse(w, "Id должен быть положительным и содержать только цифры", err, http.StatusBadRequest)
		return
	}
	
	if err := in.Validate(); err != nil {
		handler.ErrResponse(w, "Ошибка валидации: "+err.Error(), err, http.StatusBadRequest)
        return
	}

	if err := router.UC.UpdateSubscription(r.Context(), &in, id); err != nil {
		handler.ErrResponse(w, "Ошибка обновления", err, http.StatusBadRequest)
        return
	}

	handler.OkResponse(w, "OK", http.StatusCreated)
}

