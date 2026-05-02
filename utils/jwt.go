package utils

import (
	"parking/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	RoleID   uint   `json:"role_id"`
	jwt.StandardClaims
}

func GenerateToken(userID uint, username string, roleID uint) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.AppConfig.JWT.Expire) * time.Second).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "parking",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
