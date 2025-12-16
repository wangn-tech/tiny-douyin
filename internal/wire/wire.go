//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/wangn-tech/tiny-douyin/internal/api/handler"
	"github.com/wangn-tech/tiny-douyin/internal/dao"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/pkg/upload"
	"github.com/wangn-tech/tiny-douyin/internal/service"
	"gorm.io/gorm"
)

// ProvideDB 提供数据库连接
func ProvideDB() *gorm.DB {
	return global.DB
}

// UploadSet Upload 层 Provider Set
var UploadSet = wire.NewSet(
	upload.NewUploadService,
	upload.NewWorker,
)

// DAOSet DAO 层 Provider Set（只注入 DB）
var DAOSet = wire.NewSet(
	dao.NewUserDAO,
	dao.NewVideoDAO,
	dao.NewFavoriteDAO,
	dao.NewCommentDAO,
	dao.NewRelationDAO,
	dao.NewMessageDAO,
)

// ServiceSet Service 层 Provider Set（只注入 DAO）
var ServiceSet = wire.NewSet(
	service.NewUserService,
	service.NewVideoService,
	service.NewFavoriteService,
	service.NewCommentService,
	service.NewRelationService,
	service.NewMessageService,
	DAOSet,
)

// HandlerSet Handler 层 Provider Set（只注入 Service 和 Upload）
var HandlerSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewVideoHandler,
	handler.NewFavoriteHandler,
	handler.NewCommentHandler,
	handler.NewRelationHandler,
	handler.NewMessageHandler,
	ServiceSet,
	UploadSet,
)

// InitUserHandler 初始化 UserHandler（Wire 自动生成实现）
func InitUserHandler() *handler.UserHandler {
	wire.Build(
		ProvideDB,
		dao.NewUserDAO,
		dao.NewRelationDAO,
		service.NewUserService,
		handler.NewUserHandler,
	)
	return nil
}

// InitVideoHandler 初始化 VideoHandler（Wire 自动生成实现）
func InitVideoHandler() *handler.VideoHandler {
	wire.Build(
		ProvideDB,
		dao.NewUserDAO,
		dao.NewVideoDAO,
		dao.NewFavoriteDAO,
		dao.NewRelationDAO,
		service.NewVideoService,
		upload.NewUploadService,
		handler.NewVideoHandler,
	)
	return nil
}

// InitUploadWorker 初始化 UploadWorker（Wire 自动生成实现）
func InitUploadWorker() upload.IUploadWorker {
	wire.Build(
		ProvideDB,
		dao.NewVideoDAO,
		upload.NewUploadService,
		upload.NewWorker,
	)
	return nil
}

// InitFavoriteHandler 初始化 FavoriteHandler（Wire 自动生成实现）
func InitFavoriteHandler() *handler.FavoriteHandler {
	wire.Build(
		ProvideDB,
		dao.NewUserDAO,
		dao.NewVideoDAO,
		dao.NewFavoriteDAO,
		dao.NewRelationDAO,
		service.NewFavoriteService,
		handler.NewFavoriteHandler,
	)
	return nil
}

// InitCommentHandler 初始化 CommentHandler（Wire 自动生成实现）
func InitCommentHandler() *handler.CommentHandler {
	wire.Build(
		ProvideDB,
		dao.NewUserDAO,
		dao.NewVideoDAO,
		dao.NewCommentDAO,
		service.NewCommentService,
		handler.NewCommentHandler,
	)
	return nil
}

// InitRelationHandler 初始化 RelationHandler（Wire 自动生成实现）
func InitRelationHandler() *handler.RelationHandler {
	wire.Build(
		ProvideDB,
		dao.NewUserDAO,
		dao.NewRelationDAO,
		dao.NewMessageDAO,
		service.NewRelationService,
		handler.NewRelationHandler,
	)
	return nil
}

// InitMessageHandler 初始化 MessageHandler（Wire 自动生成实现）
func InitMessageHandler() *handler.MessageHandler {
	wire.Build(
		ProvideDB,
		dao.NewUserDAO,
		dao.NewRelationDAO,
		dao.NewMessageDAO,
		service.NewRelationService,
		service.NewMessageService,
		handler.NewMessageHandler,
	)
	return nil
}
