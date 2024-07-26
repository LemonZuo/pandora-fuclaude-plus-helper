package handler

import (
	"PandoraPlusHelper/api/v1"
	"PandoraPlusHelper/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginHandler struct {
	*Handler
	loginService service.LoginService
}

func NewLoginHandler(handler *Handler, loginService service.LoginService) *LoginHandler {
	return &LoginHandler{
		Handler:      handler,
		loginService: loginService,
	}
}

// Login godoc
// @Summary 账号登录
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "params"
// @Success 200 {object} v1.LoginResponse
// @Router /login [post]
func (h *LoginHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	loginType, token, rules, loginUrl, err := h.loginService.Login(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrBadRequest, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.LoginResponseData{
		LoginType:   loginType,
		AccessToken: token,
		User:        rules,
		LoginUrl:    loginUrl,
	})
}
