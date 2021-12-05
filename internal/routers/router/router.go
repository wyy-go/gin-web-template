package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wyy-go/go-web-template/internal/controller"
)

func Register(r *gin.Engine) {
	g := r.Group("/",
		//middleware.Log(),
		//middleware.CheckSign(),
		// 登录状态校验即身份校验
		//middleware.CheckLogin(),
	)

	RegisterIndexRouter(g, &controller.IndexController{})
}
