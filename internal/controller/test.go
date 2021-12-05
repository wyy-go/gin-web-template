package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wyy-go/go-web-template/internal/common/context"
	"github.com/wyy-go/wview/plugin/ginview"
	"net/http"
)

type IndexController struct {
	
}

func (t *IndexController)Index(ctx *context.Context) {
	ginview.HTML(ctx.Context(), http.StatusOK, "index", gin.H{
		"title": "Index title!",
		"add": func(a int, b int) int {
			return a + b
		},
	})
}

func (t *IndexController)Page(ctx *context.Context) {
	ginview.HTML(ctx.Context(), http.StatusOK, "page.html", gin.H{
		"title": "Page file title!!",
	})
}