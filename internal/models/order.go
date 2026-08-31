package models

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	BaseEntity

	UserID      uint        `json:"user_id" gorm:"not null"`
	TotalAmount float64     `json:"total_amount" gorm:"not null"`
	Status      OrderStatus `json:"status" gorm:"default:pending"`

	// Relationships
	User       User        `json:"user"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	BaseEntity

	OrderID   uint    `json:"order_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`

	// Relationships
	Order   Order   `json:"-"`
	Product Product `json:"product"`
}
