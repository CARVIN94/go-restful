package restful

import (
	"strconv"
	"strings"

	"github.com/CARVIN94/go-reply"
)

// Param 参数
type Param struct {
	Ctx   *Context
	Key   string
	Value string
}

// Param 检查参数
func (ctx *Context) Param(key string) *Param {
	value := ctx.Req.FormValue(key)
	return &Param{ctx, key, value}
}

// Exist 参数必须存在
func (p *Param) Exist() *Param {
	if p.Value == "" {
		p.Ctx.ReplyJSON(reply.NotExist(p.Key))
	}
	return p
}

// String 参数输出字符串类型
func (p *Param) String() string {
	return p.Value
}

// Int 参数输出数字类型
func (p *Param) Int() (number int) {
	if p.Value == "" {
		return number
	}
	number, err := strconv.Atoi(p.Value)
	if err != nil {
		p.Ctx.ReplyJSON(reply.NotInt(p.Key))
	}
	return number
}

// Array 参数输出数组类型
func (p *Param) Array() (arr []string) {
	if p.Value == "" {
		return arr
	}
	return strings.Split(p.Value, ",")
}
