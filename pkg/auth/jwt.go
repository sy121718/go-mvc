package auth

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// Config JWT 配置。
type Config struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"`
	Issuer     string `mapstructure:"issuer"`
}

var (
	jwtSecret []byte
	jwtConfig Config
	inited    bool
	jwtMu     sync.RWMutex
)

// Claims 自定义 claims。
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func getDefaultConfig() Config {
	return Config{
		Secret:     "default-secret-key-please-change-in-production",
		ExpireTime: 24,
		Issuer:     "go-mvc",
	}
}

// ValidateConfig 校验 JWT 配置。
//
// 说明：
// - 用于框架启动前的 fail-fast 校验
// - strict=true 时，会拒绝默认 secret 和空 secret
func ValidateConfig(v *viper.Viper, strict bool) error {
	cfg := getDefaultConfig()
	if v != nil {
		if err := v.UnmarshalKey("jwt", &cfg); err != nil {
			return fmt.Errorf("解析 JWT 配置失败: %w", err)
		}
	}

	if !strict {
		return nil
	}

	if cfg.Secret == "" {
		return fmt.Errorf("jwt.secret 不能为空")
	}
	if cfg.Secret == getDefaultConfig().Secret {
		return fmt.Errorf("jwt.secret 不能使用默认值")
	}
	return nil
}

// InitJWT 初始化 JWT。
func InitJWT(v *viper.Viper) error {
	jwtMu.Lock()
	defer jwtMu.Unlock()

	if inited {
		return nil
	}

	cfg := getDefaultConfig()
	if v != nil {
		if err := v.UnmarshalKey("jwt", &cfg); err != nil {
			log.Printf("解析 JWT 配置失败，使用默认配置: %v", err)
			cfg = getDefaultConfig()
		}
	}

	defaultCfg := getDefaultConfig()
	if cfg.Secret == "" {
		cfg.Secret = defaultCfg.Secret
		log.Println("警告: JWT secret 未配置，使用默认值（生产环境请修改）")
	}
	if cfg.ExpireTime <= 0 {
		cfg.ExpireTime = defaultCfg.ExpireTime
	}
	if cfg.Issuer == "" {
		cfg.Issuer = defaultCfg.Issuer
	}

	jwtConfig = cfg
	jwtSecret = []byte(cfg.Secret)
	inited = true
	log.Println("JWT 初始化成功")
	return nil
}

// GenerateToken 生成 Token。
func GenerateToken(userID int64, username string) (string, error) {
	secret, cfg, err := snapshotState()
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.ExpireTime) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken 解析 Token。
func ParseToken(tokenString string) (*Claims, error) {
	secret, _, err := snapshotState()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("不支持的签名算法: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 Token。
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	return GenerateToken(claims.UserID, claims.Username)
}

// GetExpireTime 获取过期时间（小时）。
func GetExpireTime() int {
	jwtMu.RLock()
	defer jwtMu.RUnlock()

	if !inited {
		return getDefaultConfig().ExpireTime
	}
	return jwtConfig.ExpireTime
}

// MustBeReady 检查 JWT 是否可用。
func MustBeReady() error {
	if _, _, err := snapshotState(); err != nil {
		return fmt.Errorf("JWT 组件不可用: %w", err)
	}
	return nil
}

func snapshotState() ([]byte, Config, error) {
	jwtMu.RLock()
	defer jwtMu.RUnlock()

	if !inited || len(jwtSecret) == 0 {
		return nil, Config{}, errors.New("JWT 未初始化")
	}

	secret := make([]byte, len(jwtSecret))
	copy(secret, jwtSecret)
	return secret, jwtConfig, nil
}
