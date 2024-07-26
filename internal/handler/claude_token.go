package handler

import (
	v1 "PandoraPlusHelper/api/v1"
	"PandoraPlusHelper/internal/model"
	"PandoraPlusHelper/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ClaudeTokenHandler struct {
	*Handler
	claudeTokenService service.ClaudeTokenService
}

func NewClaudeTokenHandler(
	handler *Handler,
	claudeTokenService service.ClaudeTokenService,
) *ClaudeTokenHandler {
	return &ClaudeTokenHandler{
		Handler:            handler,
		claudeTokenService: claudeTokenService,
	}
}

func (h *ClaudeTokenHandler) RefreshToken(ctx *gin.Context) {
	req := new(v1.RefreshOpenaiTokenRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	err := h.claudeTokenService.RefreshToken(ctx, req.Id)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *ClaudeTokenHandler) SearchToken(ctx *gin.Context) {
	req := new(v1.SearchOpenaiTokenRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	tokens, err := h.claudeTokenService.SearchToken(ctx, req.TokenName)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, tokens)
}

func (h *ClaudeTokenHandler) CreateToken(ctx *gin.Context) {
	req := new(model.ClaudeToken)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.claudeTokenService.Create(ctx, req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *ClaudeTokenHandler) UpdateToken(ctx *gin.Context) {
	req := new(model.ClaudeToken)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.claudeTokenService.Update(ctx, req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *ClaudeTokenHandler) DeleteToken(ctx *gin.Context) {
	req := new(v1.DeleteOpenaiTokenRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.claudeTokenService.DeleteToken(ctx, req.Id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
