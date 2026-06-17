package blogPost

import (
	"Gblog/internal/tag"
	"Gblog/pkg/blogstatus"
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BlogPost struct {
	ID        string              `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt      `gorm:"index"`
	Title     string              `gorm:"type:varchar(255);not null"`
	Content   string              `gorm:"type:text;not null"`
	Status    blogstatus.BlogStatus `gorm:"type:varchar(20);not null"`
	Slug      string              `gorm:"type:varchar(255);uniqueIndex;not null"`
	IsActive  bool                `gorm:"default:true;not null"`
	UserID    string              `gorm:"type:uuid;not null;index"`
	Tags      []tag.Tag           `gorm:"many2many:post_tags;"`
}

func (p *BlogPost) BeforeCreate(tx *gorm.DB) error {
	if !p.Status.IsValid() {
		return fmt.Errorf("status inválido: %s", p.Status)
	}
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return err
	}
	// ponytail: not RFC 4122, strict UUID format isn't needed here
	p.ID = fmt.Sprintf("%x-%x-%x-%x-%x", id[0:4], id[4:6], id[6:8], id[8:10], id[10:])
	return nil
}
