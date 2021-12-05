package context

import (
	"github.com/wyy-go/go-web-template/internal/common/jwt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/wyy-go/go-web-template/internal/common/constant"
	"github.com/wyy-go/go-web-template/internal/common/errors"
	"github.com/wyy-go/go-web-template/internal/common/typ"
)

func New(ctx *gin.Context) *Context {
	return &Context{gctx: ctx}
}

type Context struct {
	gctx *gin.Context
}

func (c *Context) Context() *gin.Context {
	return c.gctx
}

func (c *Context) Param(key string) string {
	return c.gctx.Param(key)
}

func (c *Context) Query(key string) string {
	return c.gctx.Query(key)
}

func (c *Context) ShouldBind(obj interface{}) error {
	if err := c.gctx.ShouldBind(obj); err != nil {
		return err
	}
	return nil
}

func (c *Context) Abort() {
	c.gctx.Abort()
}

func (c *Context) Next() {
	c.gctx.Next()
}

func (c *Context) Response(obj interface{}) {
	m := make(map[string]interface{})
	m["code"] = 0
	m["data"] = obj
	m["ts"] = time.Now().UnixNano()
	c.response(http.StatusOK, m)
}

func (c *Context) ResponseOK() {
	c.Response(nil)
}

func (c *Context) ResponseError(err error) {
	ce := errors.Parse(err.Error())

	m := make(map[string]interface{})
	m["code"] = ce.Code
	if ce.Message != "" {
		m["message"] = ce.Message
	}
	if ce.Detail != "" {
		m["detail"] = ce.Detail
	}
	m["ts"] = time.Now().UnixNano()
	if ce.Code == -1 {
		c.response(500, m)
		return
	}
	c.response(499, m)
}

func (c *Context) ResponseErrorEx(obj interface{}, err error) {
	ce := errors.Parse(err.Error())

	m := make(map[string]interface{})
	m["code"] = ce.Code
	if ce.Message != "" {
		m["message"] = ce.Message
	}
	if ce.Detail != "" {
		m["detail"] = ce.Detail
	}

	if obj != nil {
		m["data"] = obj
	}

	m["ts"] = time.Now().UnixNano()
	if ce.Code == -1 {
		c.response(500, m)
		return
	}
	c.response(499, m)
}

func (c *Context) response(status int, obj interface{}) {

	//deployEnv := config.GetConfig().Server.DeployEnv
	//if deployEnv == env.DeployEnvLocal || deployEnv == env.DeployEnvDev || deployEnv == env.DeployEnvTest {
	//	b, _ := json.Marshal(obj)
	//	log.Debug("应答包: " + string(b))
	//}

	c.gctx.JSON(status, obj)
	c.gctx.Abort()
}

var (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

func (c *Context) GetAppHeader() *typ.AppHeader {
	//log.Debug(c.gctx.Request.Header)
	h := &typ.AppHeader{}
	c.gctx.ShouldBindHeader(h)
	return h
}

func (c *Context) GetWsHeader() *typ.WsHeader {
	h := &typ.WsHeader{}
	c.gctx.ShouldBindHeader(h)
	return h
}

func (c *Context) GetForm() url.Values {
	c.gctx.Request.ParseForm()
	return c.gctx.Request.PostForm
}

func (c *Context) GetUidStr() string {
	return cast.ToString(c.GetUid())
}

func (c *Context) GetUid() int64 {
	if v, exists := c.gctx.Get(constant.AccountKey); !exists {
		return 0
	} else {
		if acc, ok := v.(jwt.Account); !ok {
			return 0
		} else {
			return acc.Uid
		}
	}
}

func (c *Context) GetPlatform() string {
	if v, exists := c.gctx.Get(constant.AccountKey); !exists {
		return ""
	} else {
		if acc, ok := v.(jwt.Account); !ok {
			return ""
		} else {
			return acc.Platform
		}
	}
}

func (c *Context) GetDeviceName() string {
	if v, exists := c.gctx.Get(constant.AccountKey); !exists {
		return ""
	} else {
		if acc, ok := v.(jwt.Account); !ok {
			return ""
		} else {
			return acc.DeviceName
		}
	}
}

func (c *Context) GetToken() string {
	return c.gctx.GetHeader("Token")
}

func (c *Context) GetPage() int {
	if v := c.Query("page"); len(v) > 0 {
		if n := cast.ToInt(v); n > 0 {
			return n
		}
	}

	return 1
}

func (c *Context) GetPageSize() int {
	if v := c.Query("per_page"); len(v) > 0 {
		if n := cast.ToInt(v); n > 0 {
			if n > MaxPageSize {
				return MaxPageSize
			}
			return n
		}
	}

	return DefaultPageSize
}

type list struct {
	List    interface{} `json:"list,omitempty"`
	Total   int         `json:"total,omitempty"`
	Page    int         `json:"page,omitempty"`
	PerPage int         `json:"per_page,omitempty"`
	//Pagination *pagination `json:"pagination,omitempty"`
}

type pagination struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func (c *Context) ResponseList(obj interface{}) {
	c.Response(list{List: obj})
}

func (c *Context) ResponsePage(total int, obj interface{}) {
	c.Response(list{
		List:    obj,
		Total:   total,
		Page:    c.GetPage(),
		PerPage: c.GetPageSize(),
		//Pagination: &pagination{
		//	Total:   total,
		//	Page:    c.GetPage(),
		//	PerPage: c.GetPageSize(),
		//},
	})
}
