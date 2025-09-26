package pkg

import (
	"errors"
	"lytemp/config"
	"lytemp/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	jwt.RegisteredClaims
	User models.User `json:"user"`
}

func GenerateToken(User models.User) (string, error) {
	claim := &JWTClaim{
		User: User,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(config.Get().App.AuthExpireHours))),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := t.SignedString([]byte(config.Get().App.JwtSecret))

	if err != nil {
		return token, err
	}

	return token, nil
}

func ValidateToken(signedToken string) (claim *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Get().App.JwtSecret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(*JWTClaim)

	if !ok {
		return nil, errors.New("geçersiz oturum")
	}

	if claim.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, errors.New("oturumunuz sonlandı lütfen tekrar giriş yapın")
	}

	return claim, nil
}
