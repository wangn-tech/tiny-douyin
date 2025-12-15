//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/wangn-tech/tiny-douyin/internal/api/handler"
	"github.com/wangn-tech/tiny-douyin/internal/dao"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/service"
	"gorm.io/gorm"
)

// ProvideDB 提供数据库连接
func ProvideDB() *gorm.DB {
	return global.DB
}

// DAOSet DAO 层 Provider Set（只注入 DB）
var DAOSet = wire.NewSet(
	dao.NewUserDAO,
)

// ServiceSet Service 层 Provider Set（只注入 DAO）
var ServiceSet = wire.NewSet(
	service.NewUserService,
	DAOSet,
)

// HandlerSet Handler 层 Provider Set（只注入 Service）
var HandlerSet = wire.NewSet(
	handler.NewUserHandler,
	ServiceSet,
)

// InitUserHandler 初始化 UserHandler（Wire 自动生成实现）
func InitUserHandler() *handler.UserHandler {
	wire.Build(
		ProvideDB,
		HandlerSet,
	)
	return nil
}
