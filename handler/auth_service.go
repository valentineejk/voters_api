package handler

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func jwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

// make bearertoken
// generate access
// takehome - chcek for token type
func GenerateAccess(UserID, Role string) (string, error) {

	Claims := Claims{
		UserID: UserID,
		Role:   Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   UserID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	return token.SignedString(jwtSecret())

}

// generate refresh
// takehome - chcek for token type
// change refreh to nanoid(12)
func GenerateRefresh(UserID, Role string) (string, error) {

	Claims := Claims{
		UserID: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   UserID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	return token.SignedString(jwtSecret())

}

// validate
func Vaidate(tokenStr string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(t *jwt.Token) (interface{}, error) {
			//verify Algorithm
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("wrong signing method")
			}
			return jwtSecret(), nil

		})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil

}
