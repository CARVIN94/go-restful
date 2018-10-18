package restful

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CARVIN94/go-util/reply"
)

// Context 上下文
type Context struct {
	Res  http.ResponseWriter
	Req  *http.Request
	Pipe interface{}
}

// ReplyJSON 以 json 方式返回数据
func (ctx *Context) ReplyJSON(something interface{}) {
	json.NewEncoder(ctx.Res).Encode(something)
}

// IsUrlencoded 解析
func (ctx *Context) IsUrlencoded() (cb bool) {
	ct := ctx.Req.Header["Content-Type"]
	for _, v := range ct {
		ok := strings.Contains(v, "application/x-www-form-urlencoded")
		if ok {
			cb = true
		}
	}
	if !cb {
		ctx.ReplyJSON(reply.Urlencoded())
	}
	return cb
}

// IsMultipart 解析
func (ctx *Context) IsMultipart() (cb bool) {
	ct := ctx.Req.Header["Content-Type"]
	for _, v := range ct {
		ok := strings.Contains(v, "multipart/form-data")
		if ok {
			cb = true
		}
	}
	if !cb {
		ctx.ReplyJSON(reply.Multipart())
	}
	return cb
}
