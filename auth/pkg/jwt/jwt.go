package jwt

import (
	"errors"
	"strings"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

const UserClaimKey = "User-Claims"

func CreateToken(secret []byte, claims *UserClaims) (string, error) {
	claims.Role = "user"
	return jwt2.NewWithClaims(jwt2.SigningMethodHS512, claims).SignedString(secret)
}

func ParseToken(tokenString string, secret []byte) (*UserClaims, error) {
	// Check for valid JWT format
	if strings.Count(tokenString, ".") != 2 {
		return nil, errors.New("token contains an invalid number of segments")
	}
	token, err := jwt2.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt2.Token) (interface{}, error) {
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
		return claim, errors.New("token is not valid")
	}

	return claim, nil
}


package jwt

import (
    "time"

    "github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("your-secret-key") // Replace with a secure secret key

type Claims struct {
    UserID uint   `json:"user_id"`
    Roles  []string `json:"roles"`
    jwt.StandardClaims
}

func GenerateToken(userID uint, roles []string) (string, error) {
    claims := &Claims{
        UserID: userID,
        Roles:  roles,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrSignatureInvalid
}