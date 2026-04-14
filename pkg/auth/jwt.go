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
	  secret: your-secret-key    # JWT 密钥（必须修改）
	  expire_time: 24            # 过期时间（小时）
	  issuer: go-mvc             # 签发者
	  lazy_init: false           # 是否懒加载

使用示例：

	// 在 main.go 中初始化
	auth.InitJWT(viper)

	// 登录时生成 Token
	token, err := auth.GenerateToken(userID, username)

	// 中间件中验证 Token
	claims, err := auth.ParseToken(tokenString)
	if err != nil {
	    // Token 无效
	}
	userID := claims.UserID

	// 刷新 Token
	newToken, err := auth.RefreshToken(oldToken)

PHP 对比：

	// Laravel JWT
	$token = JWTAuth::fromUser($user);
	$user = JWTAuth::toUser($token);

	// Go
	token, _ := auth.GenerateToken(user.ID, user.Name)
	claims, _ := auth.ParseToken(token)
*/
package auth

import (
	"errors"
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
	LazyInit   bool   `mapstructure:"lazy_init"`
}

var (
	jwtSecret []byte
	jwtConfig Config
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
		LazyInit:   false,
	}
}

// InitJWT 初始化 JWT
// 内部处理错误，致命错误会直接退出程序
func InitJWT(v *viper.Viper) {
	// 自己解析配置
	if err := v.UnmarshalKey("jwt", &jwtConfig); err != nil {
		log.Printf("解析 JWT 配置失败，使用默认配置: %v", err)
		jwtConfig = getDefaultConfig()
	}

	// 配置兜底：如果关键字段为空，使用默认值
	defaultCfg := getDefaultConfig()
	if jwtConfig.Secret == "" {
		jwtConfig.Secret = defaultCfg.Secret
		log.Println("警告: JWT secret 未配置，使用默认值（生产环境请修改）")
	}
	if jwtConfig.ExpireTime == 0 {
		jwtConfig.ExpireTime = defaultCfg.ExpireTime
	}
	if jwtConfig.Issuer == "" {
		jwtConfig.Issuer = defaultCfg.Issuer
	}

	// 检查是否懒加载
	if jwtConfig.LazyInit {
		return
	}

	jwtSecret = []byte(jwtConfig.Secret)
	log.Println("JWT 初始化成功")
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
