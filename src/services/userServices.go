package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/DieegoAlves/CrypexGoAPI/src/entities"
	"github.com/DieegoAlves/CrypexGoAPI/src/repositories"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type UserService struct {
	repository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return UserService{
		repository: userRepository,
	}
}

func (s *UserService) CreateUser(user *entities.User) error {
	Salt, err := generateSalt()
	if err != nil {
		return err
	}

	user.Salt = Salt
	user.Password = hashPassword(user.Password, user.Salt)

	err = s.repository.AddNewUser(user)
	if err != nil {
		return err
	}

	return nil
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)

	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(salt), nil
}

func hashPassword(password, salt string) string {
	sha := sha256.New()

	sha.Write([]byte(password + salt))

	return hex.EncodeToString(sha.Sum(nil))
}

func (s *UserService) VerifyCredentials(username, password string) (bool, error) {
	user, err := s.repository.FindByUsername(username)
	if err != nil {
		return false, err
	}

	hashedPassword := hashPassword(password, user.Salt)

	if user.Password != hashedPassword {
		return false, errors.New("password incorrect")
	}

	return true, nil
}

func (s *UserService) GenerateJWT(username string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) FindByUsername(username string) (entities.User, error) {
	var user entities.User

	user, err := s.repository.FindByUsername(username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) UpdateUsername(currentUsername, newUsername string) error {
	user, err := s.repository.FindByUsername(currentUsername)
	if err != nil {
		return errors.New("user not found")
	}

	if err := s.repository.UpdateUsername(user, newUsername); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateBio(username, newBio string) error {
	user, err := s.repository.FindByUsername(username)
	if err != nil {
		return errors.New("user not found")
	}

	err = s.repository.UpdateBio(user, newBio)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteUser(user entities.User) error {
	err := s.repository.DeleteUser(user)
	if err != nil {
		return err
	}

	return nil
}
