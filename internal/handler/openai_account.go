package handler

import (
	v1 "PandoraFuclaudePlusHelper/api/v1"
	"PandoraFuclaudePlusHelper/internal/model"
	"PandoraFuclaudePlusHelper/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type OpenaiAccountHandler struct {
	*Handler
	openaiAccountService service.OpenaiAccountService
}

func NewOpenaiAccountHandler(
	handler *Handler,
	openaiAccountService service.OpenaiAccountService,
) *OpenaiAccountHandler {
	return &OpenaiAccountHandler{
		Handler:              handler,
		openaiAccountService: openaiAccountService,
	}
}

func (h *OpenaiAccountHandler) GetAccount(ctx *gin.Context) {

}

func (h *OpenaiAccountHandler) CreateAccount(ctx *gin.Context) {
	req := new(v1.AddOpenaiAccountRequest)

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
		v1.HandleError(ctx, http.StatusBadRequest, fmt.Errorf("时间格式错误: %v", err), nil)
		return
	}

	// 构建新的 OpenaiAccount 对象
	account := &model.OpenaiAccount{
		ID:                req.ID,
		Account:           req.Account,
		ExpirationTime:    expTime,
		Status:            req.Status,
		Gpt35Limit:        req.Gpt35Limit,
		Gpt4Limit:         req.Gpt4Limit,
		Gpt4oLimit:        req.Gpt4oLimit,
		Gpt4oMiniLimit:    req.Gpt4oMiniLimit,
		O1Limit:           req.O1Limit,
		O1MiniLimit:       req.O1MiniLimit,
		ShowConversations: req.ShowConversations,
		TemporaryChat:     req.TemporaryChat,
		TokenID:           req.TokenID,
		ShareToken:        req.ShareToken,
	}

	if err := h.openaiAccountService.Create(ctx, account); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *OpenaiAccountHandler) UpdateAccount(ctx *gin.Context) {
	req := new(v1.AddOpenaiAccountRequest)

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
		v1.HandleError(ctx, http.StatusBadRequest, fmt.Errorf("时间格式错误: %v", err), nil)
		return
	}

	// 构建新的 Account 对象
	account := &model.OpenaiAccount{
		ID:                req.ID,
		Account:           req.Account,
		ExpirationTime:    expTime,
		Status:            req.Status,
		Gpt35Limit:        req.Gpt35Limit,
		Gpt4Limit:         req.Gpt4Limit,
		Gpt4oLimit:        req.Gpt4oLimit,
		Gpt4oMiniLimit:    req.Gpt4oMiniLimit,
		O1Limit:           req.O1Limit,
		O1MiniLimit:       req.O1MiniLimit,
		ShowConversations: req.ShowConversations,
		TemporaryChat:     req.TemporaryChat,
		TokenID:           req.TokenID,
		ShareToken:        req.ShareToken,
	}

	if err := h.openaiAccountService.Update(ctx, account); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *OpenaiAccountHandler) DeleteAccount(ctx *gin.Context) {
	req := new(v1.DeleteOpenaiAccountRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.openaiAccountService.DeleteAccount(ctx, req.Id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *OpenaiAccountHandler) SearchAccount(ctx *gin.Context) {

	req := new(v1.SearchOpenaiAccountRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		h.logger.Error("SearchAccount", zap.Any("err", err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	accountList, err := h.openaiAccountService.SearchAccount(ctx, req.TokenId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, accountList)
}

func (h *OpenaiAccountHandler) StatisticAccount(context *gin.Context) {
	req := new(v1.StatisticOpenaiAccountRequest)

	if err := context.ShouldBindJSON(req); err != nil {
		v1.HandleError(context, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	statistic, err := h.openaiAccountService.StatisticAccount(context, req.TokenId)
	if err != nil {
		v1.HandleError(context, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(context, statistic)
}

func (h *OpenaiAccountHandler) DisableAccount(ctx *gin.Context) {
	req := new(v1.DisableOpenaiAccountRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.openaiAccountService.DisableAccount(ctx, req.Id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *OpenaiAccountHandler) EnableAccount(ctx *gin.Context) {
	req := new(v1.EnableOpenaiAccountRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.openaiAccountService.EnableAccount(ctx, req.Id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
