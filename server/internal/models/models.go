package models

import (
	"time"
)

type Category string

const (
	Digital     Category = "digital"    // 0
	Wearable    Category = "wearable"   // 2
	Electronics Category = "electronic" // 3
)

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
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"price"` // unit in Euro
	Total     int     `json:"total"`
}

type Product struct {
	ID            int64    `json:"id"`
	Name          string   `json:"name"`
	Category      Category `json:"category"`
	Price         float64  `json:"price"`
	Description   string   `json:"description"`
	Weight        string   `json:"weight"`
	StockQuantity int      `json:"instock"`
}

type Customer struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Password  string    `json:"-"`
}
type Cart struct {
	ID       string     `json:"id"`
	ClientID string     `json:"client_id"`
	Item     []CartItem `json:"items"`
}

type CartItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
