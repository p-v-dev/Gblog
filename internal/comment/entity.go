package comment

import (
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        string         `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Content   string         `gorm:"type:text;not null"`
	PostID    string         `gorm:"type:uuid;not null;index"`
	UserID    string         `gorm:"type:uuid;not null;index"`
	IsActive  bool           `gorm:"default:true;not null"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return err
	}
	c.ID = fmt.Sprintf("%x-%x-%x-%x-%x", id[0:4], id[4:6], id[6:8], id[8:10], id[10:])
	return nil
}
