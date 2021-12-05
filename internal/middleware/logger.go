package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/wyy-go/go-web-template/internal/common/context"
	"github.com/wyy-go/go-web-template/pkg/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type dupBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w dupBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w dupBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger(levelStr string, fields map[string]interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		writer := &dupBodyWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}
		ctx.Writer = writer
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			ctx.Abort()
			return
		}
		title, _ := context.GetRouterTitleAndKey(ctx.Request.Method, ctx.Request.URL.Path)

		ctx.Request.Body.Close()
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		ctx.Next()

		m := make(map[string]interface{})
		m["title"] = title
		m["status"] = ctx.Writer.Status()
		m["method"] = ctx.Request.Method
		m["path"] = ctx.Request.URL.Path
		m["client_ip"] = ctx.ClientIP()
		m["ua"] = ctx.Request.UserAgent()
		m["request_body"] = string(body)
		query, _ := url.QueryUnescape(ctx.Request.URL.RawQuery)
		m["request_query"] = query
		//m["response_body"] = writer.body.String()
		duration := time.Since(start)
		m["duration"] = duration

		for k, v := range fields {
			m[k] = v
		}

		level, _ := logger.GetLevel(levelStr)
		logger.Fields(m).Log(level, "logger")
	}
}
