package casbin

import (
	"fmt"
	"log"
	"strings"
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

	// urlCodeMap 内存映射：URL+Method → code
	// 用于例外查询时快速找到请求 URL 对应的 code，避免每次 deny 都查数据库
	urlCodeMap sync.Map
)

// urlCodeKey 生成 URL+Method 的映射 key
func urlCodeKey(url, method string) string {
	return url + "||" + strings.ToUpper(method)
}

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

// GetCodeByURL 根据 URL+Method 查询对应的权限 code。
// 从内存映射中查找，不涉及数据库查询。
// 返回 code 和是否存在。用于例外查询场景：Casbin deny 后查此映射拿到 code，再查例外表。
func GetCodeByURL(url, method string) (string, bool) {
	val, ok := urlCodeMap.Load(urlCodeKey(url, method))
	if !ok {
		return "", false
	}
	return val.(string), true
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

	rebuildURLCodeMap(instance)
	return instance, nil
}

// rebuildURLCodeMap 从 Casbin 策略中重建 URL+Method → code 内存映射。
// 遍历所有 p 规则，提取 v1（URL）、v2（Method）、v3（code）建立映射。
func rebuildURLCodeMap(e *casbinlib.Enforcer) {
	urlCodeMap = sync.Map{}

	policies, err := e.GetPolicy()
	if err != nil {
		log.Printf("Casbin 获取策略失败: %v", err)
		return
	}

	for _, rule := range policies {
		// p 规则格式：[sub, obj, act, code]
		if len(rule) < 4 {
			continue
		}
		url := rule[1]   // obj = URL
		method := rule[2] // act = Method
		code := rule[3]   // code
		if code == "" {
			continue
		}
		urlCodeMap.Store(urlCodeKey(url, method), code)
	}
	log.Printf("Casbin URL→code 映射已加载: %d 条", func() int {
		count := 0
		urlCodeMap.Range(func(_, _ interface{}) bool { count++; return true })
		return count
	}())
}
