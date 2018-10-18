package restful

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CARVIN94/go-util/logcolor"
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
	Addr routeMap
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
	defer log.Print(logcolor.Success("服务器启动成功" + " " + "http://localhost" + port))

	if config.Port == 0 {
		log.Fatal(logcolor.Error("请检查服务器端口配置"))
	}
	if len(config.Route.Addr) == 0 {
		log.Fatal(logcolor.Error("请检查服务器路由配置"))
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = time.Second * 10
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = time.Second * 4
	}
	server = &http.Server{
		Addr:         port,
		Handler:      &handler{config.Route},
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(logcolor.Error("服务器启动失败"), err.Error())
		}
	}()
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	urlSplit := strings.Split(url, "?")
	method := r.Method
	ware, ok := h.Route.Addr[urlSplit[0]+" "+method]
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	if ok {
		r.ParseForm()
		ctx := &Context{w, r, nil}
		for _, v := range ware {
			v(ctx)
		}
	} else {
		log.Print(logcolor.Magenta("[HTTP] ") + "404: " + url)
	}
}

// Get 把事件放进 Route
func (r *Route) Get(address string, args ...middleware) {
	r.format()
	r.Addr[address+" "+"GET"] = args
}

// Post 把事件放进 Route
func (r *Route) Post(address string, args ...middleware) {
	r.format()
	r.Addr[address+" "+"POST"] = args
}

// format 格式化 Route.Addr 修复 map nil 问题
func (r *Route) format() {
	if r.Addr == nil {
		r.Addr = make(routeMap)
	}
}
