package domain

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var signature = []byte("myPrivateSignature")

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	No_hp     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRegister struct {
	Name     string
	Email    string
	Password string
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(name, email, password string) (*User, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	if len(password) < 7 {
		return nil, errors.New("password minimal 7 chararcters")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &User{
		Name:     name,
		Email:    email,
		Password: string(hash),
	}, nil
}

func (u User) GenerateJWT() (string, error) {
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iss":     "edspert",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString(signature)

	return stringToken, err
}

func (u User) DecryptJWT(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid Token")
		}
		return signature, nil
	})
	data := make(map[string]interface{})
	if err != nil {
		return data, err
	}
	if !parsedToken.Valid {
		return data, errors.New("invalid token")
	}

	return parsedToken.Claims.(jwt.MapClaims), nil
}

func (u User) Login(email, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
