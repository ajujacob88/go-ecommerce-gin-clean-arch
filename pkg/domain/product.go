package domain

type ProductCategory struct {
	ID           uint   `gorm:"primaryKey, uniqueIndex" json:"id"`
	CategoryName string `gorm:"not null, index, unique" json:"category_name"`
}

type Product struct {
	ID                uint            `gorm:"primaryKey, uniqueIndex"" json:"id"`
	ProductCategoryID uint            `gorm:"not null" json:"product_category_id" validate:"required"`
	ProductCategory   ProductCategory `gorm:"foreignKey:ProductCategoryID" json:"-"`
	Name              string          `gorm:"not null,uniqueIndex" json:"name" validate:"required"`
	//BrandID           uint            `gorm:"not null" json:"brand_id" validate:"required"`
	Description string `json:"description"`
	//ProductImage      string          `json:"product_image"`
}
