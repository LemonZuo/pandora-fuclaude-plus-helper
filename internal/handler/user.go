package handler

import (
	v1 "PandoraFuclaudePlusHelper/api/v1"
	"PandoraFuclaudePlusHelper/internal/model"
	"PandoraFuclaudePlusHelper/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(
	handler *Handler,
	userService service.UserService,
) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

func (h *UserHandler) RefreshUser(ctx *gin.Context) {
	req := new(v1.RefreshOpenaiTokenRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// err := h.openaiTokenService.RE(ctx, req.Id)
	// if err != nil {
	// 	v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
	// 	return
	// }
	v1.HandleSuccess(ctx, nil)
}

func (h *UserHandler) SearchUser(ctx *gin.Context) {
	req := new(v1.SearchUserRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	tokens, err := h.userService.SearchUser(ctx, req.UniqueName)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, tokens)
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	req := new(v1.AddUserRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 定义时间格式和时区
	const layout = "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Shanghai") // 加载东八区时区
	// 在指定时区环境下解析时间字符串
	expTime, err := time.ParseInLocation(layout, req.ExpirationTime, loc)
	if err != nil {
		// 如果时间解析错误，返回400错误
		v1.HandleError(ctx, http.StatusBadRequest, fmt.Errorf("invalid expiration time format: %v", err), nil)
		return
	}

	// 构建新的 Account 对象
	user := &model.User{
		UniqueName:     req.UniqueName,
		Password:       req.Password,
		ExpirationTime: expTime,
		Enable:         req.Enable,
		Openai:         req.Openai,
		OpenaiToken:    req.OpenaiToken,
		Claude:         req.Claude,
		ClaudeToken:    req.ClaudeToken,
	}

	if err := h.userService.Create(ctx, user); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	req := new(v1.UpdateUserRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 定义时间格式和时区
	const layout = "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Shanghai") // 加载东八区时区
	// 在指定时区环境下解析时间字符串
	expTime, err := time.ParseInLocation(layout, req.ExpirationTime, loc)
	if err != nil {
		// 如果时间解析错误，返回400错误
		v1.HandleError(ctx, http.StatusBadRequest, fmt.Errorf("invalid expiration time format: %v", err), nil)
		return
	}

	// 构建新的 Account 对象
	user := &model.User{
		ID:             req.ID,
		UniqueName:     req.UniqueName,
		Password:       req.Password,
		ExpirationTime: expTime,
		Enable:         req.Enable,
		Openai:         req.Openai,
		OpenaiToken:    req.OpenaiToken,
		Claude:         req.Claude,
		ClaudeToken:    req.ClaudeToken,
	}

	if err := h.userService.Update(ctx, user); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	req := new(v1.DeleteOpenaiTokenRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.DeleteUser(ctx, req.Id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
