package model

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sgoldenf/wb_l0/internal/interface/order"
)

var (
	ErrNoRecord = errors.New("no matching record found")
)

type OrderModel struct {
	Pool *pgxpool.Pool
}

func (m *OrderModel) AddOrder(o *order.Order) error {
	_, err := m.Pool.Exec(context.Background(),
		`INSERT INTO orders VALUES ($1, $2);`, o.OrderID, o.Data)
	if err != nil {
		return err
	}
	return nil
}

func (m *OrderModel) ReadOrder(orderID string) (*order.Order, error) {
	o := &order.Order{OrderID: orderID}
	err := m.Pool.QueryRow(context.Background(),
		`SELECT data FROM orders WHERE order_uid=$1;`, orderID).Scan(&o.Data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return o, nil
}

func (m *OrderModel) ReadAllOrders() ([]*order.Order, error) {
	rows, err := m.Pool.Query(context.Background(),
		`SELECT * FROM orders;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []*order.Order
	for rows.Next() {
		o := &order.Order{}
		if err := rows.Scan(&o.OrderID, &o.Data); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (m *OrderModel) Shutdown() {
	m.Pool.Close()
}
