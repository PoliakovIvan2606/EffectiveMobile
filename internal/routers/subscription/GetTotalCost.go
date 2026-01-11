package subscription

import (
	"EffectiveMobile/internal/utils"
	"EffectiveMobile/pkg/handler"
	"net/http"
)

func(router SubscriptionRouter) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	userID := query.Get("user_id")
	serviceName := query.Get("service_name") // Может быть пустым
	fromStr := query.Get("from")
	toStr := query.Get("to")

	// 2. Валидация обязательных полей
	if userID == "" || fromStr == "" || toStr == "" {
		handler.ErrResponse(w, "Отсутствуют обязательные параметры (user_id, from, to)", nil, http.StatusBadRequest)
		return
	}


	startDate, err := utils.ValidDate(fromStr)
	if err != nil {
		handler.ErrResponse(w, "Неверный формат даты 'from'. Используйте ММ-ГГГГ", err, http.StatusBadRequest)
		return
	}

	endDate, err := utils.ValidDate(toStr)
	if err != nil {
		handler.ErrResponse(w, "Неверный формат даты 'to'. Используйте ММ-ГГГГ", err, http.StatusBadRequest)
		return
	}

	// 4. Вызов бизнес-логики (Use Case / Service)
	// Передаем контекст для возможности отмены запроса
	total, err := router.UC.GetTotalCost(r.Context(), userID, serviceName, *startDate, *endDate)
	if err != nil {
		handler.ErrResponse(w, "Ошибка при расчете статистики", err, http.StatusInternalServerError)
		return
	}

	handler.OkResponse(w, map[string]float64{"total":total}, http.StatusOK)
}