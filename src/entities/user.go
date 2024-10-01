package entities

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primarykey"`
	Name      string    `json:"name" gorm:"not null"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Bio       string    `json:"bio"`
	Password  string    `json:"password" gorm:"not null"`
	Salt      string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) VerifyFields() error {
	fields := map[string]string{
		"name":     u.Name,
		"username": u.Username,
		"email":    u.Email,
		"password": u.Password,
	}

	for fieldName, fieldValue := range fields {
		if fieldValue == "" {
			return fmt.Errorf("field %s is empty", fieldName)
		}
	}

	return nil
}
