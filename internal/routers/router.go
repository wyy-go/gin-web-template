package routers

import (
	"embed"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/wyy-go/go-web-template/internal/config"
	"github.com/wyy-go/go-web-template/internal/middleware"
	"github.com/wyy-go/go-web-template/internal/routers/router"
	log "github.com/wyy-go/go-web-template/pkg/logger"
	"html/template"
	"net/http"
	"time"

	"github.com/wyy-go/wview"
	"github.com/wyy-go/wview/plugin/ginview"
)

var (
	Views embed.FS
	Static embed.FS
	Favicon []byte
)

func initSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = system.LoadRules([]*system.Rule{
		{
			MetricType:   system.Load,
			TriggerCount: 8.0,
			Strategy:     system.BBR,
		},
	})

	if err != nil {
		log.Fatal(err)
		return
	}
}

func Setup(engine *gin.Engine) {
	initSentinel()
	engine.NoMethod(func(ctx *gin.Context) {
		ctx.String(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	})

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})

	engine.Use(gin.Recovery())
	engine.Use(requestid.New())
	engine.Use(middleware.Logger(config.GetConfig().LoggerConfig.Level, nil))

	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST, OPTIONS, GET, PUT, PATCH, DELETE"},
		AllowHeaders: []string{"*"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	}))

	engine.Use(
		middleware.Sentinel(
			middleware.WithResourceExtractor(func(ctx *gin.Context) string {
				return ctx.GetHeader("X-Real-IP")
			}),
			middleware.WithBlockFallback(func(ctx *gin.Context) {
				ctx.AbortWithStatusJSON(400, map[string]interface{}{
					"code":    9999,
					"message": "服务器忙，请稍后重试",
				})
			}),
		),
	)
	//engine.Use(middleware.CORSMiddleware())

	//apiPrefixes := []string{"/mobile/", "/swagger/"}
	//engine.Use(middleware.StaticFile("www", apiPrefixes...))

	//api.RegisterV1(engine)
	web(engine)
	router.Register(engine)
}

func web(r *gin.Engine)  {
	fm := make(template.FuncMap)
	fm["copy"] = func() string {
		return time.Now().Format("2006")
	}

	// new template engine
	r.HTMLRender = ginview.New(wview.Config{
		Root:         "views",
		Extension:    ".html",
		Master:       "layouts/master",
		Partials:     []string{},
		Funcs:        fm,
		DisableCache: true,
		EnableEmbed:  true,
		Views:        Views,
	})

	// static file
	r.Any("/static/*filepath", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(Static))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})

	// favicon
	r.GET("/favicon.ico", func(context *gin.Context) {
		context.Writer.Write(Favicon)
	})
}