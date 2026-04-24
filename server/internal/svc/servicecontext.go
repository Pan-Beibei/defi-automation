// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"server/internal/config"
	"server/internal/permissions/erc20tokenperiodic"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {

	// 显式注册所有权限类型处理器
  erc20tokenperiodic.Register()

	return &ServiceContext{
		Config: c,
	}
}
