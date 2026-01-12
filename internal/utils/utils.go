package utils

import "time"

// layout — формат даты: месяц-год (MM-YYYY)
var layout = "01-2006"

// ValidDate
// Принимает строку с датой в формате MM-YYYY
// Парсит её в time.Time и возвращает указатель
func ValidDate(date string) (*time.Time, error) {

	// Парсим строку в time.Time по заданному layout
	t, err := time.Parse(layout, date)
	if err != nil {
		// Ошибка возникает, если формат строки не совпадает с layout
		return nil, err
	}

	return &t, nil
}

// ParseDate
// Преобразует time.Time в строку формата MM-YYYY
func ParseDate(t time.Time) string {
	return t.Format(layout)
}