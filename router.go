package restful

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CARVIN94/go-reply"
	"github.com/CARVIN94/go-util/log"
)

// Config Http服务器基础配置
type Config struct {
	Route
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Route 定义路由
type Route struct {
	Addr       routeMap
	BeforeFunc []middleware
}

// handler Http处理模型
type handler struct {
	Route
}

// middleware 中间件
type middleware func(*Context)

// routeMap 地址别名
type routeMap map[string][]middleware

var (
	server *http.Server
)

// Start 启动服务
func Start(config *Config) {
	port := ":" + strconv.Itoa(config.Port)
	defer log.Success("服务器启动成功" + " " + "http://localhost" + port)

	if config.Port == 0 {
		log.Fatal("请检查服务器端口配置")
	}
	if len(config.Route.Addr) == 0 {
		log.Fatal("请检查服务器路由配置")
	}
	server = &http.Server{
		Addr:         port,
		Handler:      &handler{config.Route},
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	go func() {
		err := server.ListenAndServe()
		log.FailOnError(err, "服务器启动失败")
	}()
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 初始化参数
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	ctx := &Context{Res: w, Req: r, Pipe: Pipe{}, Finish: false}
	ware := []middleware{}
	params := []string{}
	reg := ""
	ok := false
	// 解析地址
	url := r.URL.String()
	method := r.Method
	urlSplit := strings.Split(url, "?")
	if len(urlSplit) == 1 {
		urlSplit = append(urlSplit, "")
	}
	// 解析路由表
	for addr, middlewares := range h.Route.Addr {
		addrSplit := strings.Split(addr, " ")
		route := addrSplit[0]
		mode := addrSplit[1]
		switch mode {
		case "Gateway":
			ctx.Pipe["micro"] = route
			ctx.Pipe["proxy"] = url[len(route)+1:]
			reg = "^/" + route + "/*"
		case method:
			reg, params = analysisAddr(route)
		}
		if reg != "" {
			match, err := regexp.MatchString(reg, urlSplit[0])
			if match && err == nil {
				ware = middlewares
				ok = true
				goto Next
			}
		}
	}
Next:
	if ok {
		analysisURLData(urlSplit[0], urlSplit[1], reg, params, ctx)
		for _, v := range ware {
			if !ctx.Finish {
				v(ctx)
			}
		}
	} else {
		defer ctx.ErrorHandler()
		log.Connect("HTTP", "404", method+" "+url)
		ctx.ReplyJSON(reply.RouteNotExist())
	}
}

// Before 把事件放进 Route
func (r *Route) Before(args ...middleware) {
	r.BeforeFunc = args
}

// Get 把事件放进 Route
func (r *Route) Get(address string, args ...middleware) {
	handleRouteMethod("GET", r, address, args)
}

// Post 把事件放进 Route
func (r *Route) Post(address string, args ...middleware) {
	handleRouteMethod("POST", r, address, args)
}

// Put 把事件放进 Route
func (r *Route) Put(address string, args ...middleware) {
	handleRouteMethod("PUT", r, address, args)
}

// Delete 把事件放进 Route
func (r *Route) Delete(address string, args ...middleware) {
	handleRouteMethod("DELETE", r, address, args)
}

// Gateway 把事件放进 Route
func (r *Route) Gateway(address string, args ...middleware) {
	handleRouteMethod("Gateway", r, address, args)
}

// format 格式化 Route.Addr 修复 map nil 问题
func (r *Route) format() {
	if r.Addr == nil {
		r.Addr = make(routeMap)
	}
}

func handleRouteMethod(method string, r *Route, address string, args []middleware) {
	r.format()
	r.Addr[address+" "+method] = append(r.BeforeFunc, args...)
}
func analysisAddr(addr string) (reg string, params []string) {
	flysnowRegexp, _ := regexp.Compile(`/(:\S*?)(/|$)`)
	stringArr := flysnowRegexp.FindAllString(addr, -1)
	reg = "^" + addr
	for _, str := range stringArr {
		sub := flysnowRegexp.FindStringSubmatch(str)
		params = append(params, sub[1])
		reg = strings.Replace(reg, sub[1], `(\S*?)`, 1)
	}
	reg = reg + "$"
	return reg, params
}
func analysisURLData(url string, queryString string, reg string, paramKeys []string, ctx *Context) {
	// param
	flysnowRegexp, _ := regexp.Compile(reg)
	sub := flysnowRegexp.FindStringSubmatch(url)
	for index, value := range sub {
		if index != 0 {
			key := paramKeys[index-1][1:]
			ctx.Req.Form[key] = []string{value}
		}
	}
	// query
	if strings.HasSuffix(queryString, "&") {
		queryString = queryString[:len(queryString)-1]
	}
	queryArray := strings.Split(queryString, "&")
	if len(queryArray) != 0 && queryArray[0] != "" {
		for _, item := range queryArray {
			obj := strings.Split(item, "=")
			if len(obj) == 1 {
				ctx.Req.Form[obj[0]] = []string{obj[1]}
			}
		}
	}
}
