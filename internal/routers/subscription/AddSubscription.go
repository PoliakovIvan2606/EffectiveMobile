package subscription

import (
	models "EffectiveMobile/internal/models/subscription"
	"EffectiveMobile/pkg/handler"
	"encoding/json"
	"net/http"
)

// AddSubscription godoc
// @Summary      Создать подписку
// @Description  Добавляет новую подписку для пользователя
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      models.SubscriptionRequest  true  "Данные новой подписки"
// @Success      201           {object}  handler.ApiOkResponse              "ID созданной подписки"
// @Failure      400           {object}  handler.ApiErrResponse      "Неверный JSON или ошибка валидации"
// @Router       /subscription [post]
func(router SubscriptionRouter) AddSubscription(w http.ResponseWriter, r *http.Request) {
	// структура для входных данных
	in := models.SubscriptionRequest{}

	// Парсим JSON из body в структуру
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		handler.ErrResponse(w, "Неверный JSON", err, http.StatusBadRequest)
        return
    }
	
	// проверяем правильность пришедших данных
	if err := in.Validate(); err != nil {
		handler.ErrResponse(w, "Ошибка валидации: "+err.Error(), err, http.StatusBadRequest)
        return
	}

	// отправляем структуру на слой usecase для дальнейшей работы
	id, err := router.UC.AddSubscription(r.Context(), &in)
	if err != nil {
		handler.ErrResponse(w, "Ошибка добавления", err, http.StatusBadRequest)
        return
	}

	// при успешном добавлении подписки отправляем пользователю id
	handler.OkResponse(w, map[string]int{"id": id}, http.StatusCreated)
}

