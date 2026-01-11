package subscription

import (
	models "EffectiveMobile/internal/models/subscription"
	repoSub "EffectiveMobile/internal/repository/subscription"
	"EffectiveMobile/pkg/handler"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// UpdateSubscription godoc
// @Summary      Обновить подписку
// @Description  Обновляет существующую подписку по ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id            path      int                      true  "ID подписки"
// @Param        subscription  body      models.SubscriptionRequest  true  "Данные для обновления"
// @Success      201           {string}  handler.ApiOkResponse                   "OK"
// @Failure      400           {object}  handler.ApiErrResponse  "Неверный JSON, ошибка валидации или ID"
// @Failure      404           {object}  handler.ApiErrResponse  "Подписка не найдена"
// @Failure      500           {object}  handler.ApiErrResponse  "Ошибка сервера"
// @Router       /subscription/{id} [put]
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

	err = router.UC.UpdateSubscription(r.Context(), &in, id)
	if err != nil {
		if errors.Is(err, repoSub.ErrNoRows) {
			handler.ErrResponse(w, "Не была найдена подписка", err, http.StatusNotFound)
			return
		}
		handler.ErrResponse(w, "Ошибка обновления", err, http.StatusBadRequest)
		return
	}



	handler.OkResponse(w, "OK", http.StatusCreated)
}

