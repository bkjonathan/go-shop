package models

type Category struct {
	BaseEntity

	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`

	// Relationships
	Products []Product `json:"-"`
}

type ProductImage struct {
	BaseEntity

	ProductID uint   `json:"product_id" gorm:"not null"`
	URL       string `json:"url" gorm:"not null"`
	AltText   string `json:"alt_text"`
	IsPrimary bool   `json:"is_primary" gorm:"default:false"`

	// Relationships
	Product Product `json:"-"`
}

type Product struct {
	BaseEntity

	CategoryID  uint    `json:"category_id" gorm:"not null"`
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"not null"`
	Stock       int     `json:"stock" gorm:"default:0"`
	IsActive    bool    `json:"is_active" gorm:"default:true"`

	// Relationships
	Category   Category       `json:"category"`
	Images     []ProductImage `json:"image"`
	OrderItems []OrderItem    `json:"-"`
	CartItems  []CartItem     `json:"-"`
}
