package domain

type ProductCategory struct {
	ID           uint   `gorm:"primaryKey, uniqueIndex" json:"id"`
	CategoryName string `gorm:"not null, index, unique" json:"category_name"`
}
