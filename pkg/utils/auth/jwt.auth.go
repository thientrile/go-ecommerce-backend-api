package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go-ecommerce-backend-api.com/global"
)

type PayloadClaims struct {
	jwt.RegisteredClaims
}
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateJWT(payload jwt.Claims) (string, error) {
	// Create a new token object, specifying signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(global.Config.JWT.ApiSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateToken(uuid string) (out *Token, err error) {
	timeEx := global.Config.JWT.JwtExpiration
	if timeEx == "" {
		timeEx = "1h"
	}
	expiration, err := time.ParseDuration(timeEx)
	if err != nil {
		return out, err
	}
	now := time.Now()
	expiresAt := now.Add(expiration)
	accessToken, err := GenerateJWT(PayloadClaims{RegisteredClaims: jwt.RegisteredClaims{
		Subject:   uuid,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    global.Config.JWT.Issuer,
		ID:        uuid,
	}})
	if err != nil {
		return nil, err
	}
	// Set refresh token expiration to 7 days
	refreshExpiration := 7 * 24 * time.Hour
	refreshExpiresAt := now.Add(refreshExpiration)
	refreshToken, err := GenerateJWT(PayloadClaims{RegisteredClaims: jwt.RegisteredClaims{
		Subject:   uuid,
		ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    global.Config.JWT.Issuer,
		ID:        uuid,
	}})
	if err != nil {
		return nil, err
	}
	out = &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return out, nil
}

func ParseJwtTokenSubject(token string) (string, error) {
	// Parse the token
	parsedToken, err := jwt.ParseWithClaims(token, &PayloadClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.ApiSecret), nil
	})
	if err != nil {
		return "", err
	}

	// Check if the token is valid
	if claims, ok := parsedToken.Claims.(*PayloadClaims); ok && parsedToken.Valid {
		return claims.Subject, nil
	}
	return "", jwt.ErrSignatureInvalid
}

// validateToken checks if the token is valid
func ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
	// Parse the token
	parsedToken, err := jwt.ParseWithClaims(tokenString, &PayloadClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.ApiSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsedToken.Claims.(*PayloadClaims); ok && parsedToken.Valid {
		return &claims.RegisteredClaims, nil
	}
	return nil, err
}

func VerifyTokenSubject(token string) (*jwt.RegisteredClaims, error) {
	// Parse the token
	parsedToken, err := jwt.ParseWithClaims(token, &PayloadClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.ApiSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if claims, ok := parsedToken.Claims.(*PayloadClaims); ok && parsedToken.Valid {
		return &claims.RegisteredClaims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
