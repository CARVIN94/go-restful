package restful

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CARVIN94/go-reply"
)

// Context 上下文
type Context struct {
	Res    http.ResponseWriter
	Req    *http.Request
	Finish bool        // 控制中间件流程
	Pipe   interface{} // 传递数据
}

// ReplyJSON 以 json 方式返回数据
func (ctx *Context) ReplyJSON(something interface{}) {
	if ctx.Finish {
		return
	}
	defer ctx.Req.Body.Close()
	json.NewEncoder(ctx.Res).Encode(something)
	ctx.Finish = true
	panic("Close")
}

// Close 处理验证失败
func (ctx *Context) Close() {
	if err := recover(); err == "Close" {
		return
	}
}

// Check 检测运行时
func (ctx *Context) Check() {
	if ctx.Finish {
		panic("Close")
	}
}

// Urlencoded 解析
func Urlencoded(ctx *Context) {
	var ok = false
	ct := ctx.Req.Header["Content-Type"]
	for _, v := range ct {
		ok = strings.Contains(v, "application/x-www-form-urlencoded")
	}
	if !ok {
		ctx.ReplyJSON(reply.Urlencoded())
	}
}

// Multipart 解析
func Multipart(ctx *Context) {
	var ok = false
	ct := ctx.Req.Header["Content-Type"]
	for _, v := range ct {
		ok = strings.Contains(v, "multipart/form-data")
	}
	if !ok {
		ctx.ReplyJSON(reply.Multipart())
	}
}
