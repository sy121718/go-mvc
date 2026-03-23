package auth

import (
	"errors"
	"go-mvc/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret []byte
	jwtConfig config.JWTConfig
)

// Claims 自定义 claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// InitJWT 初始化 JWT
func InitJWT() error {
	jwtConfig = config.GetJWT()

	// 检查是否懒加载
	if jwtConfig.LazyInit {
		return nil
	}

	if jwtConfig.Secret == "" {
		return errors.New("JWT secret 不能为空")
	}

	jwtSecret = []byte(jwtConfig.Secret)
	return nil
}

// GenerateToken 生成 Token
func GenerateToken(userID int64, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtConfig.ExpireTime) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConfig.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 Token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 生成新 Token
	return GenerateToken(claims.UserID, claims.Username)
}

// GetExpireTime 获取过期时间（小时）
func GetExpireTime() int {
	return jwtConfig.ExpireTime
}
