package handler

import (
	"io"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/wangn-tech/tiny-douyin/internal/api/dto"
	"github.com/wangn-tech/tiny-douyin/internal/common/constant"
	"github.com/wangn-tech/tiny-douyin/internal/common/errc"
	"github.com/wangn-tech/tiny-douyin/internal/common/response"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/pkg/upload"
	"github.com/wangn-tech/tiny-douyin/internal/service"
)

// VideoHandler 视频处理器
type VideoHandler struct {
	videoService  service.IVideoService
	uploadService upload.IUploadService
}

// NewVideoHandler 创建 VideoHandler 实例（通过依赖注入）
func NewVideoHandler(videoService service.IVideoService, uploadService upload.IUploadService) *VideoHandler {
	return &VideoHandler{
		videoService:  videoService,
		uploadService: uploadService,
	}
}

// PublishVideo 发布视频
// POST /douyin/publish/action/
// 参数：token（必填），data（必填，视频文件），title（必填，视频标题）
func (h *VideoHandler) PublishVideo(c *gin.Context) {
	ctx := c.Request.Context()

	// 从 JWT 中间件获取用户ID（中间件已处理 token 参数）
	userID, exists := c.Get("user_id")
	if !exists {
		global.Logger.Warn("handler.PublishVideo.missing_user_id")
		response.ErrorWithCode(c, errc.ErrUnauthorized)
		return
	}

	// 参数绑定（验证 token 和 title）
	var req dto.VideoPublishRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.Warn("handler.PublishVideo.bind_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInvalidParams, err.Error())
		return
	}

	// 获取上传的视频文件
	file, err := c.FormFile("data")
	if err != nil {
		global.Logger.Warn("handler.PublishVideo.no_file",
			zap.Error(err),
		)
		response.Error(c, errc.ErrVideoFileInvalid, "请上传视频文件")
		return
	}

	// 如果没有提供标题，使用文件名
	if req.Title == "" {
		req.Title = file.Filename
	}

	global.Logger.Info("handler.PublishVideo.request",
		zap.Uint("user_id", userID.(uint)),
		zap.String("title", req.Title),
		zap.String("filename", file.Filename),
		zap.Int64("size", file.Size),
	)

	// 读取文件内容
	fileReader, err := file.Open()
	if err != nil {
		global.Logger.Error("handler.PublishVideo.open_file_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrVideoFileInvalid, "打开文件失败")
		return
	}
	defer fileReader.Close()

	fileData, err := io.ReadAll(fileReader)
	if err != nil {
		global.Logger.Error("handler.PublishVideo.read_file_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrVideoFileInvalid, "读取文件失败")
		return
	}

	// 保存到临时目录
	ext := filepath.Ext(file.Filename)
	tempFilePath, err := h.uploadService.SaveTempFile(fileData, ext)
	if err != nil {
		global.Logger.Error("handler.PublishVideo.save_temp_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrVideoUploadFailed, "保存临时文件失败")
		return
	}

	// 生成 MinIO 对象名称
	videoName := h.uploadService.GenerateObjectName(userID.(uint), ext)
	coverName := h.uploadService.GenerateCoverObjectName(userID.(uint))

	// 先创建视频记录（使用临时 URL）
	tempPlayURL := constant.VideoStatusUploading
	tempCoverURL := constant.VideoStatusUploading

	videoID, err := h.videoService.PublishVideo(ctx, &req, userID.(uint), tempPlayURL, tempCoverURL)
	if err != nil {
		global.Logger.Error("handler.PublishVideo.service_error",
			zap.Error(err),
		)
		h.uploadService.CleanupTempFile(tempFilePath)
		response.Error(c, errc.ErrVideoUploadFailed, err.Error())
		return
	}

	// 发布上传任务到消息队列
	task := &upload.VideoUploadTask{
		VideoID:     videoID,
		VideoPath:   tempFilePath,
		VideoName:   videoName,
		CoverName:   coverName,
		ContentType: file.Header.Get("Content-Type"),
		UserID:      userID.(uint),
		Title:       req.Title,
		Description: "", // VideoPublishRequest 没有 Description 字段
	}

	err = h.uploadService.PublishUploadTask(ctx, task)
	if err != nil {
		global.Logger.Error("handler.PublishVideo.publish_task_error",
			zap.Error(err),
		)
		// 任务发布失败，但视频记录已创建，需要回滚或标记失败
		response.Error(c, errc.ErrVideoProcessFailed, "发布上传任务失败")
		return
	}

	response.SuccessWithMsg(c, "视频发布成功，正在处理中")
}

// GetVideoFeed 获取视频流
// GET /douyin/feed/
// 可选参数：latest_time
func (h *VideoHandler) GetVideoFeed(c *gin.Context) {
	ctx := c.Request.Context()

	// 参数绑定
	var req dto.VideoFeedRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Warn("handler.GetVideoFeed.bind_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInvalidParams, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID（可选，用于判断是否点赞）
	var currentUserID uint = 0
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(uint)
	}

	global.Logger.Info("handler.GetVideoFeed.request",
		zap.Int64("latest_time", req.LatestTime),
		zap.Uint("current_user_id", currentUserID),
	)

	// 调用 Service 层
	data, err := h.videoService.GetVideoFeed(ctx, &req, currentUserID)
	if err != nil {
		global.Logger.Error("handler.GetVideoFeed.service_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInternalServer, "获取视频流失败")
		return
	}

	// 在 Handler 层封装响应
	resp := dto.VideoFeedResponse{
		Response: response.Response{
			StatusCode: errc.Success,
			StatusMsg:  errc.GetMsg(errc.Success),
		},
		NextTime: data.NextTime,
		Videos:   data.Videos,
	}

	response.SuccessWithData(c, resp)
}

// GetVideoList 获取用户发布的视频列表
// GET /douyin/publish/list/
// 必需参数：user_id
func (h *VideoHandler) GetVideoList(c *gin.Context) {
	ctx := c.Request.Context()

	// 参数解析
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		global.Logger.Warn("handler.GetVideoList.missing_param")
		response.Error(c, errc.ErrInvalidParams, "缺少user_id参数")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		global.Logger.Warn("handler.GetVideoList.invalid_param",
			zap.String("user_id", userIDStr),
		)
		response.Error(c, errc.ErrInvalidParams, "user_id格式错误")
		return
	}

	// 获取当前用户ID（可选，用于判断是否点赞）
	var currentUserID uint = 0
	if uid, exists := c.Get("user_id"); exists {
		currentUserID = uid.(uint)
	}

	global.Logger.Info("handler.GetVideoList.request",
		zap.Uint64("user_id", userID),
		zap.Uint("current_user_id", currentUserID),
	)

	// 调用 Service 层
	req := &dto.VideoListRequest{UserID: uint(userID)}
	data, err := h.videoService.GetVideoList(ctx, req, currentUserID)
	if err != nil {
		global.Logger.Error("handler.GetVideoList.service_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInternalServer, "获取视频列表失败")
		return
	}

	// 在 Handler 层封装响应
	resp := dto.VideoListResponse{
		Response: response.Response{
			StatusCode: errc.Success,
			StatusMsg:  errc.GetMsg(errc.Success),
		},
		Videos: data.Videos,
	}

	response.SuccessWithData(c, resp)
}
