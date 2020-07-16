package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

var (
	SecretKey = []byte("ohmpatel")
)

func GenerateToken(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Println("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		log.Println("invalid token")
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if email, ok := claims["email"].(string); ok {
			return email, nil
		}
	}
	log.Println("invalid claims")
	return "", fmt.Errorf(fmt.Sprintf("Could not able to parse token. Token valid : %v", token.Valid))

}
