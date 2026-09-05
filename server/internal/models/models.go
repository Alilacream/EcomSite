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
	id         int64
	quantity   int
	unit_price float64
	total      int
}
type Product struct {
	id                int64
	name              string
	category          string // specialized Enum -> customized depending on the ecom business
	price             float64
	stock_quantity    int
	reserved_quantity int // for items in cart
}
