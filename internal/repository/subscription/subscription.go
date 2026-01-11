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

func (repo *RepositorySubscription) AddSubscription(ctx context.Context, s *models.Subscription) (int, error) {
	const op = "repository.subscription.AddSubscription"

	query := `INSERT INTO subscription (service_name, price, user_id, start_date, end_date) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int
	err := repo.Db.QueryRowContext(ctx, query, 
		s.ServiceName, 
		s.Price, 
		s.UserId, 
		s.StartDate, 
		s.EndDate,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	
	return id, nil
}

func (repo *RepositorySubscription) GetSubscription(ctx context.Context, id int) (*models.GetSubscription, error) {
	const op = "repository.subscription.GetSubscription"
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscription WHERE id = $1`

	s := models.GetSubscription{}
	err := repo.Db.QueryRowContext(ctx, query, id).Scan(
		&s.Id,
		&s.ServiceName,
		&s.Price,
		&s.UserId,
		&s.StartDate,
		&s.EndDate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows 
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &s, nil
}


func (repo *RepositorySubscription) UpdateSubscription(ctx context.Context, s *models.Subscription, id int) error {
    const op = "repository.subscription.UpdateSubscription"

    query := `
        UPDATE subscription 
        SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
        WHERE id = $6`

    res, err := repo.Db.ExecContext(ctx, query, s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate, id)
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }

    rows, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }
    if rows == 0 {
        return ErrNoRows
    }

    return nil
}

func (repo *RepositorySubscription) DeleteSubscription(ctx context.Context, id int) error {
    const op = "repository.subscription.DeleteSubscription"

    query := `DELETE FROM subscription WHERE id = $1`

    res, err := repo.Db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }

    rows, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }
    
    if rows == 0 {
        return ErrNoRows
    }

    return nil
}

func (repo *RepositorySubscription) GetTotalCost(ctx context.Context, userID string, serviceName string, startDate, endDate time.Time) (float64, error) {
    const op = "repository.subscription.GetTotalCost"

    query := `
        SELECT COALESCE(SUM(price), 0) 
        FROM subscription 
        WHERE user_id = $1 
          AND service_name ILIKE $2 
          AND start_date >= $3 
          AND end_date <= $4`

    var total float64
    nameFilter := "%" + serviceName + "%"

    err := repo.Db.QueryRowContext(ctx, query, userID, nameFilter, startDate, endDate).Scan(&total)
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    return total, nil
}

func(repo *RepositorySubscription) GetListSubscription(ctx context.Context, UUID string) ([]models.GetSubscription, error) {
    const op = "repository.subscription.GetListSubscription"

    query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscription WHERE user_id = $1`

    rows, err := repo.Db.QueryContext(ctx, query, UUID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNoRows
        }
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer rows.Close()

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

    if len(subscriptions) == 0 {
        return nil, ErrNoRows
    }

    return subscriptions, nil
}