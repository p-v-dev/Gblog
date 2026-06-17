package tag

import (
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        string         `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"type:varchar(100);uniqueIndex;not null"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return err
	}
	t.ID = fmt.Sprintf("%x-%x-%x-%x-%x", id[0:4], id[4:6], id[6:8], id[8:10], id[10:])
	return nil
}
