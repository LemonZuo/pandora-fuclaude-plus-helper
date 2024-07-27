package handler

import (
	v1 "PandoraFuclaudePlusHelper/api/v1"
	"PandoraFuclaudePlusHelper/internal/model"
	"PandoraFuclaudePlusHelper/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OpenaiTokenHandler struct {
	*Handler
	openaiTokenService service.OpenaiTokenService
}

func NewOpenaiTokenHandler(
	handler *Handler,
	openaiTokenService service.OpenaiTokenService,
) *OpenaiTokenHandler {
	return &OpenaiTokenHandler{
		Handler:            handler,
		openaiTokenService: openaiTokenService,
	}
}

func (h *OpenaiTokenHandler) RefreshToken(ctx *gin.Context) {
	req := new(v1.RefreshOpenaiTokenRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	err := h.openaiTokenService.RefreshToken(ctx, req.Id)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *OpenaiTokenHandler) SearchToken(ctx *gin.Context) {
	req := new(v1.SearchOpenaiTokenRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	tokens, err := h.openaiTokenService.SearchToken(ctx, req.TokenName)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, tokens)
}

func (h *OpenaiTokenHandler) CreateToken(ctx *gin.Context) {
	req := new(model.OpenaiToken)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.openaiTokenService.Create(ctx, req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *OpenaiTokenHandler) UpdateToken(ctx *gin.Context) {
	req := new(model.OpenaiToken)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.openaiTokenService.Update(ctx, req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *OpenaiTokenHandler) DeleteToken(ctx *gin.Context) {
	req := new(v1.DeleteOpenaiTokenRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.openaiTokenService.DeleteToken(ctx, req.Id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
