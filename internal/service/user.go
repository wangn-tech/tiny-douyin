package service

import (
	"context"
	"errors"
	"time"

	"github.com/wangn-tech/tiny-douyin/internal/api/dto"
	"github.com/wangn-tech/tiny-douyin/internal/dao"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"github.com/wangn-tech/tiny-douyin/internal/pkg/hash"
	"github.com/wangn-tech/tiny-douyin/internal/pkg/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// IUserService 用户服务接口
type IUserService interface {
	Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
	GetUserInfo(ctx context.Context, req *dto.UserInfoRequest) (*dto.UserInfo, error)
}

// UserService 用户服务实现
type UserService struct {
	userDAO dao.IUserDAO
}

// NewUserService 创建 UserService 实例（依赖注入）
func NewUserService(userDAO dao.IUserDAO) IUserService {
	return &UserService{
		userDAO: userDAO,
	}
}

// NewUserServiceDefault 创建默认配置的 UserService（便捷方法）
func NewUserServiceDefault() IUserService {
	return NewUserService(dao.NewUserDAO(global.DB))
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	start := time.Now()

	global.Logger.Info("service.Register.start",
		zap.String("username", req.Username),
	)

	// 只查询用户名是否存在
	exists, err := s.userDAO.ExistsUsername(ctx, req.Username)
	if err != nil {
		global.Logger.Error("service.Register.check_exists_error",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		return nil, err
	}
	if exists {
		global.Logger.Warn("service.Register.user_already_exists",
			zap.String("username", req.Username),
		)
		return nil, errors.New("用户已存在")
	}

	// 密码加密
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		global.Logger.Error("service.Register.hash_password_error",
			zap.Error(err),
		)
		return nil, err
	}

	// 创建用户（GORM 自动生成自增 ID）
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Username,
	}

	if err := s.userDAO.CreateUser(ctx, user); err != nil {
		global.Logger.Error("service.Register.create_user_error",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		return nil, err
	}

	// 生成 Token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		global.Logger.Error("service.Register.generate_token_error",
			zap.Uint("user_id", user.ID),
			zap.Error(err),
		)
		return nil, err
	}

	global.Logger.Info("service.Register.success",
		zap.String("username", req.Username),
		zap.Uint("user_id", user.ID),
		zap.Duration("duration", time.Since(start)),
	)

	return &dto.UserRegisterResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	start := time.Now()

	global.Logger.Info("service.Login.start",
		zap.String("username", req.Username),
	)

	// 查询用户
	user, err := s.userDAO.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Warn("service.Login.user_not_found",
				zap.String("username", req.Username),
			)
			return nil, errors.New("用户不存在")
		}
		global.Logger.Error("service.Login.get_user_error",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		return nil, err
	}

	// 验证密码
	if !hash.CheckPassword(user.Password, req.Password) {
		global.Logger.Warn("service.Login.invalid_password",
			zap.String("username", req.Username),
		)
		return nil, errors.New("密码错误")
	}

	// 生成 Token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		global.Logger.Error("service.Login.generate_token_error",
			zap.Uint("user_id", user.ID),
			zap.Error(err),
		)
		return nil, err
	}

	global.Logger.Info("service.Login.success",
		zap.String("username", req.Username),
		zap.Uint("user_id", user.ID),
		zap.Duration("duration", time.Since(start)),
	)

	return &dto.UserLoginResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, req *dto.UserInfoRequest) (*dto.UserInfo, error) {
	start := time.Now()

	global.Logger.Info("service.GetUserInfo.start",
		zap.Uint("user_id", req.UserID),
	)

	user, err := s.userDAO.GetUserByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Warn("service.GetUserInfo.user_not_found",
				zap.Uint("user_id", req.UserID),
			)
			return nil, errors.New("用户不存在")
		}
		global.Logger.Error("service.GetUserInfo.get_user_error",
			zap.Uint("user_id", req.UserID),
			zap.Error(err),
		)
		return nil, err
	}

	global.Logger.Info("service.GetUserInfo.success",
		zap.Uint("user_id", req.UserID),
		zap.Duration("duration", time.Since(start)),
	)

	return &dto.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Signature: user.Signature,
	}, nil
}
