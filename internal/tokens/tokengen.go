package tokens

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	secretKey             = []byte(os.Getenv("SECRET_KEY"))
	accessTokenDuration   = time.Minute * 15
	refreshTokenDuration  = time.Hour * 24 * 7
)

type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Type   string    `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// アクセス・リフレッシュトークン生成
func GenerateTokens(userID uuid.UUID) (accessToken string, refreshToken string, err error) {
	now := time.Now()

	// アクセストークン
	accessClaims := CustomClaims{
		UserID: userID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	at, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	// リフレッシュトークン
	refreshClaims := CustomClaims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	rt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

// トークン検証
func ValidateToken(tokenString string, expectedType string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Type != expectedType {
		return nil, errors.New("token type mismatch")
	}

	return claims, nil
}

// リフレッシュトークン → 新しいアクセストークン
func RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := ValidateToken(refreshToken, "refresh")
	if err != nil {
		return "", err
	}
	accessToken, _, err := GenerateTokens(claims.UserID)
	return accessToken, err
}
