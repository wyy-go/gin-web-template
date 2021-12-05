package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wyy-go/go-web-template/internal/common/constant"
	"github.com/wyy-go/go-web-template/internal/common/context"
	"github.com/wyy-go/go-web-template/internal/common/errors"
	"github.com/wyy-go/go-web-template/internal/common/util"
	log "github.com/wyy-go/go-web-template/pkg/logger"
	"io/ioutil"
	"mime"
	"strings"
)

func CheckSign() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		title, keyString := context.GetRouterTitleAndKey(gctx.Request.Method, gctx.Request.URL.Path)
		urlData := strings.Split(keyString, "-")
		log.Debugf("title: %s method: %s url: %s", title, urlData[0], urlData[1])
		ctx := context.New(gctx)
		h := ctx.GetAppHeader()
		mapBody := make(map[string]interface{})

		ct := gctx.Request.Header.Get("Content-Type")
		ct, _, _ = mime.ParseMediaType(ct)
		if ct == "application/json" {
			body, _ := ioutil.ReadAll(gctx.Request.Body)
			gctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			log.Debug("请求包: " + string(body))
			dec := json.NewDecoder(bytes.NewReader(body))
			dec.UseNumber()
			dec.Decode(&mapBody)
			//log.Debug(mapBody)
		} else {
			for k, vs := range ctx.GetForm() {
				v := ""
				if len(vs) > 0 {
					v = vs[0]
				}
				if v != "" {
					mapBody[k] = v
				}
			}
		}

		sign := util.Sign(h, mapBody, constant.SignKey)
		log.Debugf("计算得到的sign:%s", sign)
		log.Debugf("客户端上传的sign:%s", h.Sign)
		if h.Sign != sign {
			ctx.ResponseError(errors.ErrSign)
			gctx.Abort()
		}

		gctx.Next()
	}
}
