package models

import "time"

type Order struct {
	ID              int64       `json:"id"`
	ClientID        string      `json:"client_id"`
	Status          string      `json:"status"` // enum typ
	PaymentIntentID string      `json:"payment_intent_id"`
	Total           int         `json:"total"`
	ShippingAddress string      `json:"shipping_address"`
	ProductsPlaced  []LineOrder `json:"products_placed"`
	PlacedAt        time.Time   `json:"placed_at"`
	PaidAt          *time.Time  `json:"paid_at,omitempty"` // ✅ pointer for nullable
}

type LineOrder struct {
	ID        int64
	Quantity  int
	UnitPrice float64 // unit in Euro
	Total     int
}

type Product struct {
	ID            int64
	Name          string
	Category      string // specialized Enum -> customized depending on the ecom business
	Price         float64
	StockQuantity int
}

type Customer struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
	Password  string
}
