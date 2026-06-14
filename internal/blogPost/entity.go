package blogPost

import (
	"Gblog/pkg"

	"gorm.io/gorm"
)

type BlogPost struct {
	gorm.Model

	Title    string         `gorm:"type:varchar(255);not null"`
	Content  string         `gorm:"type:text;not null"`
	Status   pkg.BlogStatus `gorm:"type:varchar(20);not null"`
	Slug     string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	IsActive bool           `gorm:"default:true;not null"` // Campo para soft delete
}
