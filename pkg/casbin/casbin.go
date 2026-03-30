package casbin

import (
	"fmt"
	"log"
	"sync"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Casbin 权限管理组件
// 使用 RBAC 模型，策略通过 GORM Adapter 持久化到数据库

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

// rbacModel RBAC 权限模型定义
// sub = 用户/角色, obj = 资源标识, act = 操作
// g = _, _ 支持角色继承（用户→角色→角色...）
const rbacModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

// GetEnforcer 获取 Casbin Enforcer 实例
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}

// InitCasbin 初始化 Casbin（传入已有的 GORM DB 实例）
// GORM Adapter 会自动创建 casbin_rule 表
func InitCasbin(db *gorm.DB) error {
	var err error
	once.Do(func() {
		err = initCasbin(db)
	})
	return err
}

func initCasbin(db *gorm.DB) error {
	// 创建 GORM 适配器（自动建表 sys_casbin_rule）
	a, err := gormadapter.NewAdapterByDBUseTableName(db, "", "sys_casbin_rule")
	if err != nil {
		return fmt.Errorf("创建 Casbin 适配器失败: %v", err)
	}

	// 解析 RBAC 模型
	m, err := model.NewModelFromString(rbacModel)
	if err != nil {
		return fmt.Errorf("创建 Casbin 模型失败: %v", err)
	}

	// 创建 Enforcer
	enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		return fmt.Errorf("创建 Casbin Enforcer 失败: %v", err)
	}

	// 从数据库加载已有策略到内存
	if err := enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载 Casbin 策略失败: %v", err)
	}

	log.Println("Casbin 初始化成功")
	return nil
}
