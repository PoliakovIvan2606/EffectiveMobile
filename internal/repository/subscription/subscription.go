package repository

import (
	models "EffectiveMobile/internal/models/subscription"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)


var (
	ErrAddSubscription = errors.New("ошибка добавления подписки")
	ErrNoRows = errors.New("подписка не была найдена")
)

type RepositorySubscription struct {
	Db *sql.DB
}

func NewRepositorySubscription(Db *sql.DB) *RepositorySubscription {
	return &RepositorySubscription{
		Db: Db,
	}
}

// Метод добавления подписки получает контекст и модель подписки, если что-то пошло не так отправляется ошибка или
// id модели если всё прошло хорошо
func (repo *RepositorySubscription) AddSubscription(ctx context.Context, s *models.Subscription) (int, error) {
    // переменная для логов
	const op = "repository.subscription.AddSubscription"

    // запрос на добавление подписки
	query := `INSERT INTO subscription (service_name, price, user_id, start_date, end_date) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

    // переменная в которую будем складывать пришедшие данные
	var id int
    // отправляем запрос в БД
	err := repo.Db.QueryRowContext(ctx, query, 
		s.ServiceName, 
		s.Price, 
		s.UserId, 
		s.StartDate, 
		s.EndDate,
	).Scan(&id)

    // если произошлоа ошибка отправляем на слой выше самау ошибку и op (переменная где произошла ошибка)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	
	return id, nil
}

// Метод для получения подписки по id, получает контекст и id подписки, если что-то пошло не так отправляется ошибка или
// модель подписки если всё прошло хорошо
func (repo *RepositorySubscription) GetSubscription(ctx context.Context, id int) (*models.GetSubscription, error) {
    // переменная для логов
	const op = "repository.subscription.GetSubscription"
    // запрос на получение подписки
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscription WHERE id = $1`

    // переменная в которую будем складывать пришедшие данные
	s := models.GetSubscription{}
    // отправляем запрос в БД
	err := repo.Db.QueryRowContext(ctx, query, id).Scan(
		&s.Id,
		&s.ServiceName,
		&s.Price,
		&s.UserId,
		&s.StartDate,
		&s.EndDate,
	)

    // если произошлоа ошибка отправляем на слой выше самау ошибку и op (переменная где произошла ошибка)
	if err != nil {
        // если пришла ошибка с отсутсвием строк из БД отправляем свою ошибку которую отлавливаем
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows 
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &s, nil
}

// Метод обновления подписки, получает контекст, модель и id подписки, если что-то пошло не так отправляется ошибка
func (repo *RepositorySubscription) UpdateSubscription(ctx context.Context, s *models.Subscription, id int) error {
    // переменная для логов
    const op = "repository.subscription.UpdateSubscription"
    // запрос на обновление подписки
    query := `
        UPDATE subscription 
        SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
        WHERE id = $6`

    // отправляем запрос в БД
    res, err := repo.Db.ExecContext(ctx, query, s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate, id)
    // обрабатываем ошибку если что-то произошло не так
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }

    // получаем количество строк затронутых обновлением если даннные число равно 0 
    // то ничего не обновлялось и отправляем ошибку ErrNoRows
    rows, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }
    // если 
    if rows == 0 {
        return ErrNoRows
    }

    return nil
}

// Метод удаления подписки, получает контекст и id подписки, если что-то пошло не так отправляется ошибка
func (repo *RepositorySubscription) DeleteSubscription(ctx context.Context, id int) error {
    // переменная для логов
    const op = "repository.subscription.DeleteSubscription"
    // запрос на удалении подписки
    query := `DELETE FROM subscription WHERE id = $1`
    // отправляем запрос в БД
    res, err := repo.Db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }
    
    // получаем количество строк затронутых обновлением если даннные число равно 0 
    // то ничего не обновлялось и отправляем ошибку ErrNoRows
    rows, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }
    if rows == 0 {
        return ErrNoRows
    }

    return nil
}

// Метод который получает сумму всех подписок определённого пользователя с фильтрами по названию, дате начала и конца подписки
// получает контекст, uuid пользователя, название подписки и даты начала и конца подписки, отправляет ошибку если она произошла и
// сумму если всё прошло усупешно
func (repo *RepositorySubscription) GetTotalCost(ctx context.Context, userID string, serviceName string, startDate, endDate time.Time) (float64, error) {
    // переменная для логов
    const op = "repository.subscription.GetTotalCost"
    // sql запрос
    query := `
        SELECT COALESCE(SUM(price), 0) 
        FROM subscription 
        WHERE user_id = $1 
          AND service_name ILIKE $2 
          AND start_date >= $3 
          AND end_date <= $4`

    // переменная в которую будем складывать сумму
    var total float64
    nameFilter := "%" + serviceName + "%"
    // отправляем запрос в БД
    err := repo.Db.QueryRowContext(ctx, query, userID, nameFilter, startDate, endDate).Scan(&total)
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    return total, nil
}

// Метод получения списка подписок пользоваетля, метод получает контекст и uuid польщователя, если ошибка отправляет её в
// противоположном случае отправляет список моделей подписок
func(repo *RepositorySubscription) GetListSubscription(ctx context.Context, UUID string) ([]models.GetSubscription, error) {
    // переменная для логов
    const op = "repository.subscription.GetListSubscription"
    // sql запрос для получения списка
    query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscription WHERE user_id = $1`
    // отправляем запрос в БД
    rows, err := repo.Db.QueryContext(ctx, query, UUID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNoRows
        }
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer rows.Close()

    // перебираем ответ и кладём в переменную subscriptions
    var subscriptions []models.GetSubscription
    for rows.Next() {
        var s models.GetSubscription
        err := rows.Scan(&s.Id, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate)
        if err != nil {
            return nil, fmt.Errorf("%s: %w", op, err)
        }
        subscriptions = append(subscriptions, s)
    }

    // Проверка на ошибки после итерации
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    // если ответа нет отправляем ErrNoRows
    if len(subscriptions) == 0 {
        return nil, ErrNoRows
    }

    return subscriptions, nil
}