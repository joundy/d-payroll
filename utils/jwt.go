package utils

import (
	"d-payroll/entity"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(secret string, payload *entity.AuthTokenPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   payload.ID,
		"role": payload.Role,
	})

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func VerifyToken(secret string, token string) (*entity.AuthTokenPayload, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if parsed.Valid {
		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			return nil, err
		}

		idFloat, ok := claims["id"].(float64)
		if !ok {
			return nil, err
		}
		roleStr, ok := claims["role"].(string)
		if !ok {
			return nil, err
		}

		payload := &entity.AuthTokenPayload{
			ID:   uint(idFloat),
			Role: entity.UserRole(roleStr),
		}

		return payload, nil
	}

	return nil, nil
}
