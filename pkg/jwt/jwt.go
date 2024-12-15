package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pdkonovalov/user-registration-service/pkg/config"
)

type JwtGenerator struct {
	accessTokenTtl  time.Duration
	refreshTokenTtl time.Duration
	secret          []byte
}

func Init(config *config.Config) (*JwtGenerator, error) {
	if config.JwtSecret == "" {
		return nil, fmt.Errorf("jwt secret is empty")
	}
	accessTokenTtl := config.AccessTokenTtl
	refreshTokenTtl := config.RefreshTokenTtl
	secret := []byte(config.JwtSecret)
	return &JwtGenerator{accessTokenTtl, refreshTokenTtl, secret}, nil
}

type accessTokenClaims struct {
	Ip string `json:"ip"`
	jwt.RegisteredClaims
}

func (gen *JwtGenerator) GenerateAccessToken(email string, ip string) (string, error) {
	claims := accessTokenClaims{
		ip,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(gen.accessTokenTtl)),
			Subject:   email,
		}}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(gen.secret)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (gen *JwtGenerator) GenerateRefreshToken(email string) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(gen.refreshTokenTtl)),
		Subject:   email,
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(gen.secret)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (gen *JwtGenerator) ValidateAccessToken(tokenStr string) (string, string, bool) {
	accessToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return gen.secret, nil
	})
	if err != nil {
		return "", "", false
	}
	email, err := accessToken.Claims.GetSubject()
	if err != nil {
		return "", "", false
	}
	mapClaims := accessToken.Claims.(jwt.MapClaims)
	_, ok := mapClaims["ip"]
	if !ok {
		return "", "", false
	}
	ip := mapClaims["ip"].(string)
	return email, ip, true
}

func (gen *JwtGenerator) ValidateRefreshToken(tokenStr string) (string, bool) {
	refreshToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return gen.secret, nil
	})
	if err != nil {
		return "", false
	}
	email, err := refreshToken.Claims.GetSubject()
	if err != nil {
		return "", false
	}
	return email, true
}
