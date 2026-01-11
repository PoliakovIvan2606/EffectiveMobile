package subscription

import (
	"EffectiveMobile/pkg/handler"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func(router SubscriptionRouter) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

    IdStr := vars["id"]
	id, err := strconv.Atoi(IdStr)
	if err != nil {
		handler.ErrResponse(w, "Id должен быть положительным и содержать только цифры", err, http.StatusBadRequest)
		return
	}

	if err = router.UC.DeleteSubscription(r.Context(), id); err != nil {
		handler.ErrResponse(w, "Ошибка удаления подписки", err, http.StatusBadRequest)
		return
	}

	handler.OkResponse(w, "OK", http.StatusOK)
}