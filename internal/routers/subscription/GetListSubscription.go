package subscription

import (
	repoSub "EffectiveMobile/internal/repository/subscription"
	"EffectiveMobile/pkg/handler"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetListSubscription godoc
// @Summary      Получить список подписок пользователя
// @Description  Возвращает список подписок по UUID пользователя
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        uuid  path     string  true  "UUID пользователя"
// @Success      200   {array}  models.SubscriptionResponse  "Список подписок"
// @Failure      400   {object} handler.ApiErrResponse       "Неверный формат UUID"
// @Failure      404   {object} handler.ApiErrResponse       "Подписки не найдены"
// @Router       /subscription/list/{uuid} [get]
func(router SubscriptionRouter) GetListSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

    UUIDstr := vars["uuid"]
	_, err := uuid.Parse(UUIDstr)
	if err != nil {
		handler.ErrResponse(w, "Не валидный uuid", err, http.StatusBadRequest)
		return
	}

	out, err := router.UC.GetListSubscription(r.Context(), UUIDstr)
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