package permissions

import (
	"fmt"
	"server/internal/types"
)

// HandlerConstructor 是每种权限类型的构造函数签名
type HandlerConstructor func(req types.SubmitPermissionReq, deps HandlerDeps) PermissionHandler

// registry 是全局注册表，权限类型名 → 构造函数
// 各权限类型在自己包的 init() 中调用 Register() 注册自己
var registry = map[string]HandlerConstructor{}

// Register 注册一个权限类型的处理器构造函数
// 在各权限子包的 init() 中调用
func Register(permType string, ctor HandlerConstructor) {
    if _, exists := registry[permType]; exists {
        panic(fmt.Sprintf("permission handler already registered for type: %s", permType))
    }
    registry[permType] = ctor
}