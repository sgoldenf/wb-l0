package order

type OrderModelInterface interface {
	AddOrder(*Order) error
	ReadOrder(orderID string) (*Order, error)
	ReadAllOrders() ([]*Order, error)
	Shutdown()
}

type Order struct {
	OrderID string `json:"order_uid"`
	Data    string `json:"-"`
}
