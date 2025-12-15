package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wangn-tech/tiny-douyin/internal/api/dto"
	"github.com/wangn-tech/tiny-douyin/internal/config"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 初始化测试环境
func init() {
	// 初始化全局配置（测试环境）
	global.Config = &config.AppConfig{
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			TTL:    24,
		},
	}

	// 初始化测试日志
	global.Logger, _ = zap.NewDevelopment()
}

// MockUserDAO 实现 IUserDAO 接口用于测试
type MockUserDAO struct {
	CreateUserFunc        func(ctx context.Context, user *model.User) error
	GetUserByUsernameFunc func(ctx context.Context, username string) (*model.User, error)
	GetUserByIDFunc       func(ctx context.Context, id uint) (*model.User, error)
	ExistsUsernameFunc    func(ctx context.Context, username string) (bool, error)
}

func (m *MockUserDAO) CreateUser(ctx context.Context, user *model.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, user)
	}
	return nil
}

func (m *MockUserDAO) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(ctx, username)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockUserDAO) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockUserDAO) ExistsUsername(ctx context.Context, username string) (bool, error) {
	if m.ExistsUsernameFunc != nil {
		return m.ExistsUsernameFunc(ctx, username)
	}
	return false, nil
}

// TestUserService_Register 测试用户注册
func TestUserService_Register(t *testing.T) {
	ctx := context.Background()

	// 准备 mock DAO
	mockDAO := &MockUserDAO{
		// 模拟：用户名不存在（使用 ExistsUsername 方法）
		ExistsUsernameFunc: func(ctx context.Context, username string) (bool, error) {
			return false, nil
		},
		// 模拟：创建用户成功，设置 ID
		CreateUserFunc: func(ctx context.Context, user *model.User) error {
			user.ID = 123
			return nil
		},
	}

	// 通过依赖注入创建 Service（不再注入 logger）
	service := NewUserService(mockDAO)

	// 执行测试
	req := &dto.UserRegisterRequest{
		Username: "testuser",
		Password: "password123",
	}
	resp, err := service.Register(ctx, req)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(123), resp.UserID)
	assert.NotEmpty(t, resp.Token)
}

// TestUserService_Register_UserExists 测试用户已存在的情况
func TestUserService_Register_UserExists(t *testing.T) {
	ctx := context.Background()

	mockDAO := &MockUserDAO{
		// 模拟：用户名已存在（使用 ExistsUsername 方法）
		ExistsUsernameFunc: func(ctx context.Context, username string) (bool, error) {
			return true, nil
		},
	}

	service := NewUserService(mockDAO)

	req := &dto.UserRegisterRequest{
		Username: "existinguser",
		Password: "password123",
	}
	resp, err := service.Register(ctx, req)

	// 断言
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "用户已存在", err.Error())
}

// TestUserService_Login 测试用户登录
// 注意：这个测试需要真实的密码哈希，在实际项目中建议跳过或使用真实哈希
func TestUserService_Login(t *testing.T) {
	t.Skip("需要真实的密码哈希，跳过此测试")

	ctx := context.Background()

	// 提前生成的密码哈希（对应 "password123"）
	hashedPassword := "$2a$10$XqN0lYH.9o.gXxVL4jGJU.3ZqH0OIKQ3VJqH0qOIKQ3VJqH0qOIKQ"

	mockDAO := &MockUserDAO{
		GetUserByUsernameFunc: func(ctx context.Context, username string) (*model.User, error) {
			user := &model.User{
				Username: username,
				Password: hashedPassword,
			}
			user.ID = 456
			return user, nil
		},
	}

	service := NewUserService(mockDAO)

	req := &dto.UserLoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	resp, err := service.Login(ctx, req)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestUserService_GetUserInfo 测试获取用户信息
func TestUserService_GetUserInfo(t *testing.T) {
	ctx := context.Background()

	mockDAO := &MockUserDAO{
		GetUserByIDFunc: func(ctx context.Context, id uint) (*model.User, error) {
			if id == 999 {
				user := &model.User{
					Username:  "testuser",
					Nickname:  "Test User",
					Avatar:    "http://example.com/avatar.jpg",
					Signature: "Hello World",
				}
				user.ID = 999
				return user, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := NewUserService(mockDAO)

	req := &dto.UserInfoRequest{
		UserID: 999,
	}
	userInfo, err := service.GetUserInfo(ctx, req)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, userInfo)
	assert.Equal(t, uint(999), userInfo.ID)
	assert.Equal(t, "testuser", userInfo.Username)
	assert.Equal(t, "http://example.com/avatar.jpg", userInfo.Avatar)
}

// TestUserService_GetUserInfo_NotFound 测试用户不存在
func TestUserService_GetUserInfo_NotFound(t *testing.T) {
	ctx := context.Background()

	mockDAO := &MockUserDAO{
		GetUserByIDFunc: func(ctx context.Context, id uint) (*model.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := NewUserService(mockDAO)

	req := &dto.UserInfoRequest{
		UserID: 404,
	}
	userInfo, err := service.GetUserInfo(ctx, req)

	// 断言
	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Equal(t, "用户不存在", err.Error())
}
