package restful

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/CARVIN94/go-reply"
)

// Pipe 传递数据
type Pipe map[string]interface{}

// Context 上下文
type Context struct {
	Res    http.ResponseWriter
	Req    *http.Request
	Pipe   Pipe
	Finish bool // 控制中间件流程
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

// ErrorHandler 错误处理
func (ctx *Context) ErrorHandler() {
	err := recover()
	if err == "Close" {
		return
	} else if err != nil {
		panic(err)
	}
}

// Check 检测运行时
func (ctx *Context) Check() {
	if ctx.Finish {
		panic("Close")
	}
}

// SetFormData 改变 form 内参数的值
func (ctx *Context) SetFormData(key string, value string) {
	ctx.Req.Form[key] = []string{value}
}

// Urlencoded 解析
func Urlencoded(ctx *Context) {
	defer ctx.ErrorHandler()
	ctx.Req.ParseForm()
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
	defer ctx.ErrorHandler()
	ctx.Req.ParseMultipartForm(32 << 20)
	var ok = false
	ct := ctx.Req.Header["Content-Type"]
	for _, v := range ct {
		ok = strings.Contains(v, "multipart/form-data")
	}
	if !ok {
		ctx.ReplyJSON(reply.Multipart())
	}
}

// JSON 解析
func JSON(ctx *Context) {
	defer ctx.ErrorHandler()
	ctx.Req.ParseForm()
	var ok = false
	ct := ctx.Req.Header["Content-Type"]
	for _, v := range ct {
		ok = strings.Contains(v, "application/json")
	}
	bodyStr, err := ioutil.ReadAll(ctx.Req.Body)
	if !ok || err != nil {
		ctx.ReplyJSON(reply.JSON())
	}
	jsonErr := json.Unmarshal(bodyStr, &ctx.Pipe)
	if jsonErr != nil {
		ctx.ReplyJSON(reply.JSON())
	}
}
