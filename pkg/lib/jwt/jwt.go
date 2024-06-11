package jwt

import (
	"fmt"
	"time"
	"vault/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(vaultID int, secret string, expires time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["vault_id"] = vaultID
	claims["exp"] = time.Now().Add(expires * time.Second).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeToken(token string, secret string) (int, error) {
	var model models.TokenModel

	jwtToken, err := jwt.ParseWithClaims(token, &model, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !jwtToken.Valid {
		return 0, fmt.Errorf("failed to decode token: %w", err)
	}

	return model.ID, nil
}
