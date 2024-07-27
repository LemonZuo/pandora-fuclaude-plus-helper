package middleware

import (
	v1 "PandoraFuclaudePlusHelper/api/v1"
	commonConfig "PandoraFuclaudePlusHelper/config"
	"PandoraFuclaudePlusHelper/pkg/log"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"sort"
	"strings"
)

func SignMiddleware(logger *log.Logger, conf *viper.Viper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requiredHeaders := []string{"Timestamp", "Nonce", "Sign", "App-Version"}

		for _, header := range requiredHeaders {
			value, ok := ctx.Request.Header[header]
			if !ok || len(value) == 0 {
				v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
				ctx.Abort()
				return
			}
		}

		data := map[string]string{
			"AppKey":     commonConfig.GetConfig().AppKey,
			"Timestamp":  ctx.Request.Header.Get("Timestamp"),
			"Nonce":      ctx.Request.Header.Get("Nonce"),
			"AppVersion": ctx.Request.Header.Get("App-Version"),
		}

		var keys []string
		for k := range data {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return strings.ToLower(keys[i]) < strings.ToLower(keys[j]) })

		var str string
		for _, k := range keys {
			str += k + data[k]
		}
		str += commonConfig.GetConfig().AppSecurity

		if ctx.Request.Header.Get("Sign") != strings.ToUpper(cryptor.Md5String(str)) {
			v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
