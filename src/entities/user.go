package entities

import (
	"fmt"
	"time"
)

type User struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Password  string `json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) VerifyFields() error {

	if u.Name == "" {
		return fmt.Errorf("field name is empty")
	}

	if u.Email == "" {
		return fmt.Errorf("field email is empty")
	}

	if u.Password == "" {
		return fmt.Errorf("field password is empty")
	}
	return nil
}
