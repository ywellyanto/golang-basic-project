package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// sebaiknya jangan hardcode, namun ditaruh di env / secret
var jwtKey = []byte("rahasia")

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     uint   `json:"role"`
	jwt.StandardClaims
}

// generate jwt token
func GenerateJWT(email, username string, role uint) (tokenString string, err error) {
	expTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaim{
		Username: username,
		Email:    email,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

// func validate jwt token
func ValidateToken(signedToken string) (email string, role uint, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)

	// jika claims gagal
	if !ok {
		err = errors.New("could not parse claims for token")
		return
	}

	// jika expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}

	role = claims.Role
	email = claims.Email
	return
}
