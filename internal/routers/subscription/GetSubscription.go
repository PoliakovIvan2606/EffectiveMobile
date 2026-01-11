package subscription

import (
	repoSub "EffectiveMobile/internal/repository/subscription"
	"EffectiveMobile/pkg/handler"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetSubscription godoc
// @Summary      Получить подписку по ID
// @Description  Возвращает информацию о подписке по её идентификатору
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      int    true   "ID подписки"
// @Success      200   {object}  models.SubscriptionResponse   "Данные подписки"
// @Failure      400   {object}  handler.ApiErrResponse  "Неверный формат ID"
// @Failure      404   {object}  handler.ApiErrResponse  "Подписка не найдена"
// @Router       /subscription/{id} [get]
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