package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wyy-go/go-web-template/internal/controller"
	"github.com/wyy-go/go-web-template/internal/routers/helper"
)

func RegisterIndexRouter(g *gin.RouterGroup, c *controller.IndexController) {
	helper.GET(g, "/", c.Index, "首页")
	helper.GET(g, "/page", c.Page, "页面")
}
