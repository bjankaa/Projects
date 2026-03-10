package utility

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const helper = "key"

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 8).Unix(),
	})
	return token.SignedString([]byte(helper))
}

func VerifyToken(token string) (int64, error) {
	parsedtoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(helper), nil
	})

	if err != nil {
		return 0, errors.New("could not parse token")
	}

	tokenisValid := parsedtoken.Valid

	if !tokenisValid {
		return 0, errors.New("invalid token")
	}

	claims, ok := parsedtoken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userId := int64(claims["userId"].(float64))

	return userId, nil
}
