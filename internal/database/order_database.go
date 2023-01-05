package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/medenzel/orders-rest-api/internal/order"
)

func TimeNowToString() string {
	t := time.Now()
	return fmt.Sprintf("%d/%d/%d %02d:%02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
}

// GetOrder - retrieves order from the database by ID
func (db *Database) GetOrder(ctx context.Context, ID int) (order.Order, error) {
	var ord order.Order
	row := db.DB.QueryRowContext(
		ctx,
		`SELECT id, description, state, create_at
		FROM orders 
		WHERE id = $1;`, ID)
	err := row.Scan(&ord.ID, &ord.Description, &ord.State, &ord.CreateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return order.Order{}, order.ErrNoOrderFound
		}
		return order.Order{}, fmt.Errorf("get order from database: %w", err)
	}

	return ord, nil
}

func (db *Database) GetAllOrders(ctx context.Context) ([]order.Order, error) {
	orders := make([]order.Order, 0)
	rows, err := db.DB.QueryContext(ctx,
		`SELECT id, description, state, create_at FROM orders;`)
	if err != nil {
		return nil, fmt.Errorf("get all orders from database: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var ord order.Order
		if err := rows.Scan(&ord.ID, &ord.Description, &ord.State, &ord.CreateAt); err != nil {
			return nil, fmt.Errorf("get all orders from database: %w", err)
		}
		orders = append(orders, ord)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get all orders from database: %w", err)
	}
	return orders, nil
}

// PostOrder - adds a new order to the database
func (db *Database) PostOrder(ctx context.Context, ord order.Order) (order.Order, error) {
	if ord.CreateAt == "" {
		ord.CreateAt = TimeNowToString()
	}
	row := db.DB.QueryRowContext(
		ctx,
		`INSERT INTO orders 
		(description, state, create_at) VALUES
		($1, $2, $3) 
		RETURNING id;`, ord.Description, ord.State, ord.CreateAt,
	)
	err := row.Scan(&ord.ID)
	if err != nil {
		return order.Order{}, fmt.Errorf("post order to database: %w", err)
	}
	return ord, nil
}

// UpdateOrder - update an order in the database
func (db *Database) UpdateOrder(ctx context.Context, ID int, newOrd order.Order) (order.Order, error) {
	row := db.DB.QueryRowContext(
		ctx,
		`UPDATE orders SET
		description = $2,
		state = $3,
		create_at = $4
		WHERE id = $1
		RETURNING id;`, ID, newOrd.Description, newOrd.State, newOrd.CreateAt,
	)
	var updatedID string
	err := row.Scan(&updatedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return order.Order{}, order.ErrNoOrderFound
		}
		return order.Order{}, fmt.Errorf("update order in database: %w", err)
	}
	return newOrd, nil
}

// DeleteOrder - delete an order from the database
func (db *Database) DeleteOrder(ctx context.Context, id int) error {
	res, err := db.DB.ExecContext(ctx,
		`DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete order from database: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete order from database: %w", err)
	}
	if rows == 0 {
		return order.ErrNoOrderFound
	}
	return nil
}
