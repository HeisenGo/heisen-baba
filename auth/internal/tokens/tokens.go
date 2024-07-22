package tokens

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID       uint
	Email        string
	Role         string
	IsSuperAdmin bool
	Sections     []string
}

func CreateToken(secret []byte, claims *UserClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(secret)
}

func ParseToken(tokenString string, secret []byte) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	var claim *UserClaims
	if token.Claims != nil {
		cc, ok := token.Claims.(*UserClaims)
		if ok {
			claim = cc
		}
	}

	if err != nil {
		return claim, err
	}

	if !token.Valid {
		return claim, jwt.ErrSignatureInvalid
	}

	return claim, nil
}
