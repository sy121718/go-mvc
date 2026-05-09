package casbin

import (
	"fmt"
	"log"
	"sync"

	"go-mvc/pkg/database"

	casbinlib "github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Casbin 权限管理组件
// 使用 RBAC 模型，策略通过 GORM Adapter 持久化到数据库

var (
	enforcer *casbinlib.Enforcer
	mu       sync.RWMutex
)

// rbacModel RBAC 权限模型定义
const rbacModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, code

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

// GetEnforcer 获取 Casbin Enforcer 实例
func GetEnforcer() *casbinlib.Enforcer {
	mu.RLock()
	defer mu.RUnlock()
	return enforcer
}

// Init 初始化 Casbin 组件。
func Init(_ *viper.Viper) error {
	db, err := database.GetDB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}
	return InitCasbin(db)
}

// InitCasbin 初始化 Casbin（传入已有的 GORM DB 实例）
func InitCasbin(db *gorm.DB) error {
	mu.Lock()
	defer mu.Unlock()

	if enforcer != nil {
		return nil
	}

	instance, err := initCasbin(db)
	if err != nil {
		return err
	}

	enforcer = instance
	log.Println("Casbin 初始化成功")
	return nil
}

// Close 关闭 Casbin 组件并清理运行时状态。
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	enforcer = nil
	return nil
}

// Ready 检查 Casbin 组件是否已完成初始化。
func Ready() error {
	if GetEnforcer() == nil {
		return fmt.Errorf("casbin 未初始化")
	}
	return nil
}

func initCasbin(db *gorm.DB) (*casbinlib.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDBUseTableName(db, "", "sys_casbin_rule")
	if err != nil {
		return nil, fmt.Errorf("创建 Casbin 适配器失败: %w", err)
	}

	m, err := model.NewModelFromString(rbacModel)
	if err != nil {
		return nil, fmt.Errorf("创建 Casbin 模型失败: %w", err)
	}

	instance, err := casbinlib.NewEnforcer(m, a)
	if err != nil {
		return nil, fmt.Errorf("创建 Casbin Enforcer 失败: %w", err)
	}

	if err := instance.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("加载 Casbin 策略失败: %w", err)
	}

	return instance, nil
}
