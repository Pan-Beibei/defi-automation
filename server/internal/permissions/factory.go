package permissions

import (
	"fmt"
	"server/internal/types"
)

// PermissionHandlerFactory 对应 TS: PermissionHandlerFactory
// Factory 本身不感知任何具体权限类型，完全通过注册表路由
type PermissionHandlerFactory struct {

}

func NewPermissionHandlerFactory() *PermissionHandlerFactory {
    return &PermissionHandlerFactory{}
}

// Create 根据权限类型从注册表查找并创建对应的 handler
// 对应 TS: PermissionHandlerFactory.createPermissionHandler()
func (f *PermissionHandlerFactory) Create(req *types.SubmitPermissionReq, deps HandlerDeps) (PermissionHandler, error) {
    ctor, ok := registry[req.Permission.Type]
    if !ok {
        return nil, fmt.Errorf("unsupported permission type: %s", req.Permission.Type)
    }
    return ctor(*req, deps), nil
}