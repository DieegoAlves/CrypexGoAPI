package repositories

import (
	"errors"
	"github.com/DieegoAlves/CrypexGoAPI/src/entities"
	"gorm.io/gorm"
)

type UserRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		database: db,
	}
}

func (db *UserRepository) AddNewUser(user *entities.User) error {
	err := db.database.Create(user)
	if err != nil {
		return err.Error
	}

	return nil
}

func (db *UserRepository) FindByUsername(username string) (entities.User, error) {
	var user entities.User
	if err := db.database.Where("username = ?", username).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (db *UserRepository) UpdateUsername(user entities.User, newUsername string) error {
	var actualUser entities.User
	if err := db.database.Where("username = ?", newUsername).First(&actualUser).Error; err == nil {
		return errors.New("username already exists")
	}

	user.Username = newUsername

	return db.database.Save(user).Error
}

func (db *UserRepository) UpdateBio(user entities.User, newBio string) error {
	user.Bio = newBio

	return db.database.Save(user).Error
}

func (db *UserRepository) DeleteUser(user entities.User) error {

	return db.database.Delete(user).Error
}
