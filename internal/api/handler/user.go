package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/tiny-douyin/internal/api/dto"
	"github.com/wangn-tech/tiny-douyin/internal/common/errc"
	"github.com/wangn-tech/tiny-douyin/internal/common/response"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/service"
	"go.uber.org/zap"
)

// UserHandler 用户处理器
// 职责：处理 HTTP 请求和响应，参数验证，调用 Service 层
// 不应包含：业务逻辑、数据库操作、数据结构定义
type UserHandler struct {
	userService service.IUserService
}

// NewUserHandler 创建 UserHandler 实例（依赖注入）
func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// NewUserHandlerDefault 创建默认配置的 UserHandler（便捷方法）
func NewUserHandlerDefault() *UserHandler {
	return NewUserHandler(service.NewUserServiceDefault())
}

// Register 用户注册接口
// POST /douyin/user/register/
func (h *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	// 参数绑定和验证
	var req dto.UserRegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.Warn("handler.Register.bind_error",
			zap.String("error", err.Error()),
			zap.String("client_ip", c.ClientIP()),
		)
		response.Error(c, errc.ErrInvalidParams, "参数错误: "+err.Error())
		return
	}

	global.Logger.Info("handler.Register.request",
		zap.String("username", req.Username),
		zap.String("client_ip", c.ClientIP()),
	)

	// 调用 Service 层处理业务逻辑
	serviceResp, err := h.userService.Register(ctx, &req)
	if err != nil {
		global.Logger.Error("handler.Register.service_error",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		response.Error(c, errc.ErrUserAlreadyExists, err.Error())
		return
	}

	// 构造响应并返回
	response.SuccessWithData(c, dto.UserRegisterResponse{
		Response: response.Response{
			StatusCode: errc.Success,
			StatusMsg:  errc.GetMsg(errc.Success),
		},
		UserID: serviceResp.UserID,
		Token:  serviceResp.Token,
	})
}

// Login 用户登录接口
// POST /douyin/user/login/
func (h *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. 参数绑定和验证
	var req dto.UserLoginRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.Warn("handler.Login.bind_error",
			zap.String("error", err.Error()),
			zap.String("client_ip", c.ClientIP()),
		)
		response.Error(c, errc.ErrInvalidParams, "参数错误: "+err.Error())
		return
	}

	global.Logger.Info("handler.Login.request",
		zap.String("username", req.Username),
		zap.String("client_ip", c.ClientIP()),
	)

	// 2. 调用 Service 层处理业务逻辑
	serviceResp, err := h.userService.Login(ctx, &req)
	if err != nil {
		global.Logger.Error("handler.Login.service_error",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		if err.Error() == "用户不存在" {
			response.Error(c, errc.ErrUserNotFound, err.Error())
		} else {
			response.Error(c, errc.ErrInvalidPassword, err.Error())
		}
		return
	}

	// 3. 构造响应并返回
	response.SuccessWithData(c, dto.UserLoginResponse{
		Response: response.Response{
			StatusCode: errc.Success,
		},
		UserID: serviceResp.UserID,
		Token:  serviceResp.Token,
	})
}

// GetUserInfo 获取用户信息接口
// GET /douyin/user/?user_id=xxx
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. 参数解析和验证
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		global.Logger.Warn("handler.GetUserInfo.missing_param",
			zap.String("client_ip", c.ClientIP()),
		)
		response.Error(c, errc.ErrInvalidParams, "缺少user_id参数")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		global.Logger.Warn("handler.GetUserInfo.invalid_param",
			zap.String("user_id", userIDStr),
			zap.String("client_ip", c.ClientIP()),
		)
		response.Error(c, errc.ErrInvalidParams, "user_id格式错误")
		return
	}

	global.Logger.Info("handler.GetUserInfo.request",
		zap.Uint64("user_id", userID),
		zap.String("client_ip", c.ClientIP()),
	)

	// 2. 调用 Service 层处理业务逻辑
	req := &dto.UserInfoRequest{UserID: uint(userID)}
	user, err := h.userService.GetUserInfo(ctx, req)
	if err != nil {
		global.Logger.Error("handler.GetUserInfo.service_error",
			zap.Uint64("user_id", userID),
			zap.Error(err),
		)
		response.Error(c, errc.ErrUserNotFound, err.Error())
		return
	}

	// 3. 构造响应并返回
	response.SuccessWithData(c, dto.UserInfoResponse{
		Response: response.Response{
			StatusCode: errc.Success,
		},
		User: *user,
	})
}
