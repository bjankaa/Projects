package utility

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func getEnvOrFile(key string) string {
	fileKey := key + "_FILE"
	if path := os.Getenv(fileKey); path != "" {
		b, err := os.ReadFile(path)
		if err == nil {
			return strings.TrimSpace(string(b))
		}
		log.Printf("Failed reading secrt: %v", err)
	}
	return os.Getenv(key)
}

// automatic call
func init() {

	sec := getEnvOrFile("JWT_SECRET")
	sec = strings.TrimSpace(sec)
	if sec == "" {
		sec = "key"
	}
	jwtSecret = []byte(sec)
}

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 8).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func VerifyToken(token string) (int64, error) {
	parsedtoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
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
