package models

import (
	"time"
)

type User struct {
	BaseEntity

	Email     string   `json:"email" gorm:"uniqueIndex;not null"`
	Password  string   `json:"-" gorm:"not null"`
	FirstName string   `json:"first_name" gorm:"not null"`
	LastName  string   `json:"last_name" gorm:"not null"`
	Phone     string   `json:"phone"`
	IsActive  bool     `json:"is_active" gorm:"default:true"`
	Role      UserRole `json:"role" gorm:"default:customer"`

	// Relationships
	RefreshToken []RefreshToken `json:"-"`
	Orders       []Order        `json:"-"`
	Cart         Cart           `json:"-"`
}

type UserRole string

const (
	UserRoleCustomer UserRole = "customer"
	UserRoleAdmin    UserRole = "admin"
)

type RefreshToken struct {
	BaseEntity

	UserID    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`

	// Relationships
	User User `json:"-"`
}

type Cart struct {
	BaseEntity

	UserID uint `json:"user_id" gorm:"not null"`
}

type CartItem struct {
	BaseEntity

	CartID    uint `json:"cart_id" gorm:"not null"`
	ProductID uint `json:"product_id" gorm:"not null"`
	Quantity  int  `json:"quantity" gorm:"not null"`

	// Relationships
	Cart    Cart    `json:"-"`
	Product Product `json:"product"`
}
