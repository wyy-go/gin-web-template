package middleware

import (
	"github.com/wyy-go/go-web-template/internal/common/constant"
	"github.com/wyy-go/go-web-template/internal/common/jwt"
	"github.com/wyy-go/go-web-template/internal/service/passport"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wyy-go/go-web-template/internal/common/context"
	log "github.com/wyy-go/go-web-template/pkg/logger"
)

func authToken(token string) (acc *jwt.Account, err error) {
	acc, err = passport.GetService().AuthToken(token)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func CheckLogin(excludePrefixes ...string) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		//title, keyString := context.GetRouterTitleAndKey(gctx.Request.Method, gctx.Request.URL.Path)
		//urlData := strings.Split(keyString, "-")
		//log.Debugf("title: %s method: %s url: %s", title, urlData[0], urlData[1])

		//if env.DeployEnv != "prod" {
		//	// 不需要验证登录状态
		//	gctx.Next()
		//	return
		//}
		ctx := context.New(gctx)
		if checkPrefix(gctx.Request.URL.Path, excludePrefixes...) {
			// 不需要验证登录状态
			gctx.Next()
			return
		}
		info, err := authToken(ctx.GetToken())
		if err != nil {
			log.Error(err)
			ctx.ResponseError(err)
			return
		}
		gctx.Set(constant.AccountKey, *info)
		gctx.Next()
	}
}

func checkPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
