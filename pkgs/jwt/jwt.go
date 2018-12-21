package jwt

import (
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

//CustomClaims jwt claim struct
type CustomClaims struct {
	ID      string
	Name    string
	Account string
	Email   string
	jwt.StandardClaims
}

//NewCustomClaims create claims
func NewCustomClaims(id, name, account, email string) CustomClaims {
	return CustomClaims{
		id,
		name,
		account,
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(60*60*24*365) * time.Second).Unix(),
			Id:        id,
		},
	}
}

//NewToken generate token string
func NewToken(c CustomClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	return token.SignedString([]byte(secret))
}

//ValidateToken validate token and return claims
func ValidateToken(tokenstring, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &CustomClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
