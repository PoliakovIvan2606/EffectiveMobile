package subscription

import (
	repoSub "EffectiveMobile/internal/repository/subscription"
	"EffectiveMobile/pkg/handler"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func(router SubscriptionRouter) GetSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

    IdStr := vars["id"]
	id, err := strconv.Atoi(IdStr)
	if err != nil {
		handler.ErrResponse(w, "Id должен быть положительным и содержать только цифры", err, http.StatusBadRequest)
		return
	}

	out, err := router.UC.GetSubscription(r.Context(), id)
	if err != nil {
		if errors.Is(err, repoSub.ErrNoRows) {
			handler.ErrResponse(w, "Не была найдена подписка", err, http.StatusNotFound)
			return
		}
		handler.ErrResponse(w, "Ошибка получения подписки", err, http.StatusBadRequest)
		return
	}

	handler.OkResponse(w, out, http.StatusOK)
}