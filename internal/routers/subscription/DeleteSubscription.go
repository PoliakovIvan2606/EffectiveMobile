package subscription

import (
	"EffectiveMobile/pkg/handler"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DeleteSubscription godoc
// @Summary      Удалить подписку
// @Description  Удаляет подписку по ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      int    true  "ID подписки"
// @Success      200  {string}  handler.ApiOkResponse  "OK"
// @Failure      400  {object}  handler.ApiErrResponse  "Неверный ID или ошибка удаления"
// @Failure      404  {object}  handler.ApiErrResponse  "Подписка не найдена"
// @Failure      500  {object}  handler.ApiErrResponse  "Ошибка сервера"
// @Router       /subscription/{id} [delete]
func(router SubscriptionRouter) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	// получаем переменные из URL
	vars := mux.Vars(r)

    IdStr := vars["id"]
	// Проверяем правильность пришедшего id
	id, err := strconv.Atoi(IdStr)
	if err != nil {
		handler.ErrResponse(w, "Id должен быть положительным и содержать только цифры", err, http.StatusBadRequest)
		return
	}

	// отправляем данные на слой usecase
	if err = router.UC.DeleteSubscription(r.Context(), id); err != nil {
		handler.ErrResponse(w, "Ошибка удаления подписки", err, http.StatusBadRequest)
		return
	}

	// при успешном выполнении отправляем положительный статус
	handler.OkResponse(w, "OK", http.StatusOK)
}