package user

import (
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"type:varchar(255);not null"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string         `gorm:"type:varchar(255);not null"`
	IsActive  bool           `gorm:"default:true;not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return err
	}
	// ponytail: not RFC 4122, strict UUID format isn't needed here
	u.ID = fmt.Sprintf("%x-%x-%x-%x-%x", id[0:4], id[4:6], id[6:8], id[8:10], id[10:])
	return nil
}
