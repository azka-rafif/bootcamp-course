package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type JWT struct {
	secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{secret: secret}
}

func (j *JWT) GenerateJwt(userId, userName, role string) (string, error) {
	claims := Claims{
		UserId:   userId,
		UserName: userName,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour + 1).Unix(),
			Issuer:    "Bootcamp-auth",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func (j *JWT) ValidateJwt(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("JWT not valid")
}
