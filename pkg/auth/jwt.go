/*
JWT 认证组件包
===========================================
提供 JWT Token 生成、解析、刷新功能

主要功能：
- Token 生成（登录时）
- Token 解析（验证身份）
- Token 刷新（续期）
- 全局配置管理

配置说明（config.yaml）：

	jwt:
	  secret: your-secret-key
	  expire_time: 24
	  issuer: go-mvc
*/
package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

/*
JWT 组件
===========================================
配置结构体定义在这里，自己解析配置
*/

// Config JWT配置
type Config struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"`
	Issuer     string `mapstructure:"issuer"`
}

var (
	jwtSecret []byte
	jwtConfig Config
	inited    bool
)

// Claims 自定义 claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// getDefaultConfig 获取默认配置
func getDefaultConfig() Config {
	return Config{
		Secret:     "default-secret-key-please-change-in-production",
		ExpireTime: 24,
		Issuer:     "go-mvc",
	}
}

// InitJWT 初始化 JWT
func InitJWT(v *viper.Viper) error {
	if inited {
		return nil
	}

	if err := v.UnmarshalKey("jwt", &jwtConfig); err != nil {
		log.Printf("解析 JWT 配置失败，使用默认配置: %v", err)
		jwtConfig = getDefaultConfig()
	}

	defaultCfg := getDefaultConfig()
	if jwtConfig.Secret == "" {
		jwtConfig.Secret = defaultCfg.Secret
		log.Println("警告: JWT secret 未配置，使用默认值（生产环境请修改）")
	}
	if jwtConfig.ExpireTime <= 0 {
		jwtConfig.ExpireTime = defaultCfg.ExpireTime
	}
	if jwtConfig.Issuer == "" {
		jwtConfig.Issuer = defaultCfg.Issuer
	}

	jwtSecret = []byte(jwtConfig.Secret)
	inited = true
	log.Println("JWT 初始化成功")
	return nil
}

func ensureInitialized() error {
	if !inited || len(jwtSecret) == 0 {
		return errors.New("JWT 未初始化")
	}
	return nil
}

// GenerateToken 生成 Token
func GenerateToken(userID int64, username string) (string, error) {
	if err := ensureInitialized(); err != nil {
		return "", err
	}

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
	if err := ensureInitialized(); err != nil {
		return nil, err
	}

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

	return GenerateToken(claims.UserID, claims.Username)
}

// GetExpireTime 获取过期时间（小时）
func GetExpireTime() int {
	return jwtConfig.ExpireTime
}

// MustBeReady 检查 JWT 是否已完成初始化
func MustBeReady() error {
	if err := ensureInitialized(); err != nil {
		return fmt.Errorf("JWT 组件不可用: %w", err)
	}
	return nil
}
