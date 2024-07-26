package handler

import (
	v1 "PandoraPlusHelper/api/v1"
	"PandoraPlusHelper/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ClaudeAccountHandler struct {
	*Handler
	claudeAccountService service.ClaudeAccountService
}

func NewClaudeAccountHandler(
	handler *Handler,
	claudeAccountService service.ClaudeAccountService,
) *ClaudeAccountHandler {
	return &ClaudeAccountHandler{
		Handler:              handler,
		claudeAccountService: claudeAccountService,
	}
}

func (h *ClaudeAccountHandler) GetAccount(ctx *gin.Context) {

}

func (h *ClaudeAccountHandler) CreateAccount(ctx *gin.Context) {
	// req := new(v1.AddOpenaiAccountRequest)
	//
	// if err := ctx.ShouldBindJSON(req); err != nil {
	// 	v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
	// 	return
	// }
	//
	// // 定义时间格式和时区
	// const layout = "2006-01-02 15:04:05"
	// loc, _ := time.LoadLocation("Asia/Shanghai") // 加载东八区时区
	// // 在指定时区环境下解析时间字符串
	// expTime, err := time.ParseInLocation(layout, req.ExpirationTime, loc)
	// if err != nil {
	// 	// 如果时间解析错误，返回400错误
	// 	v1.HandleError(ctx, http.StatusBadRequest, fmt.Errorf("invalid expiration time format: %v", err), nil)
	// 	return
	// }
	//
	// // 构建新的 ClaudeAccount 对象
	// account := &model.ClaudeAccount{
	// 	ID:                req.ID,
	// 	Account:           req.Account,
	// 	ExpirationTime:    expTime,
	// 	Status:            req.Status,
	// 	Gpt35Limit:        req.Gpt35Limit,
	// 	Gpt4Limit:         req.Gpt4Limit,
	// 	ShowConversations: req.ShowConversations,
	// 	TemporaryChat:     req.TemporaryChat,
	// 	TokenID:           req.TokenID,
	// 	ShareToken:        req.ShareToken,
	// }
	//
	// if err := h.claudeAccountService.Create(ctx, account); err != nil {
	// 	v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
	// 	return
	// }
	// v1.HandleSuccess(ctx, nil)
}

func (h *ClaudeAccountHandler) UpdateAccount(ctx *gin.Context) {
	// req := new(v1.AddOpenaiAccountRequest)
	//
	// if err := ctx.ShouldBindJSON(req); err != nil {
	// 	v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
	// 	return
	// }
	//
	// // 定义时间格式和时区
	// const layout = "2006-01-02 15:04:05"
	// loc, _ := time.LoadLocation("Asia/Shanghai") // 加载东八区时区
	// // 在指定时区环境下解析时间字符串
	// expTime, err := time.ParseInLocation(layout, req.ExpirationTime, loc)
	// if err != nil {
	// 	// 如果时间解析错误，返回400错误
	// 	v1.HandleError(ctx, http.StatusBadRequest, fmt.Errorf("invalid expiration time format: %v", err), nil)
	// 	return
	// }
	//
	// // 构建新的 Account 对象
	// account := &model.ClaudeAccount{
	// 	ID:                req.ID,
	// 	Account:           req.Account,
	// 	ExpirationTime:    expTime,
	// 	Status:            req.Status,
	// 	Gpt35Limit:        req.Gpt35Limit,
	// 	Gpt4Limit:         req.Gpt4Limit,
	// 	ShowConversations: req.ShowConversations,
	// 	TemporaryChat:     req.TemporaryChat,
	// 	TokenID:           req.TokenID,
	// 	ShareToken:        req.ShareToken,
	// }
	//
	// if err := h.claudeAccountService.Update(ctx, account); err != nil {
	// 	v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
	// 	return
	// }
	// v1.HandleSuccess(ctx, nil)
}

func (h *ClaudeAccountHandler) DeleteAccount(ctx *gin.Context) {
	req := new(v1.DeleteOpenaiAccountRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.claudeAccountService.DeleteAccount(ctx, req.Id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *ClaudeAccountHandler) SearchAccount(ctx *gin.Context) {

	req := new(v1.SearchOpenaiAccountRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		h.logger.Error("SearchAccount", zap.Any("err", err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	accountList, err := h.claudeAccountService.SearchAccount(ctx, req.TokenId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, accountList)
}

func (h *ClaudeAccountHandler) StatisticAccount(context *gin.Context) {
	req := new(v1.StatisticOpenaiAccountRequest)

	if err := context.ShouldBindJSON(req); err != nil {
		v1.HandleError(context, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	statistic, err := h.claudeAccountService.StatisticAccount(context, req.TokenId)
	if err != nil {
		v1.HandleError(context, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(context, statistic)
}

func (h *ClaudeAccountHandler) DisableAccount(ctx *gin.Context) {
	req := new(v1.DisableOpenaiAccountRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.claudeAccountService.DisableAccount(ctx, req.Id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *ClaudeAccountHandler) EnableAccount(ctx *gin.Context) {
	req := new(v1.EnableOpenaiAccountRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.claudeAccountService.EnableAccount(ctx, req.Id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
