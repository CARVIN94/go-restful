# go-resrful

够浪 koa 风格的路由控制

## usage

[go-example]("https://github.com/CARVIN94/go-example")

```go
package main

import (
	"os"

	"github.com/CARVIN94/hello/controller"

	"github.com/CARVIN94/go-restful"
)

 func init() {
	restful.Start(&restful.Config{
		Route: setRoute(),
		Port:  3000,
		// ReadTimeout:  time.Second * 10,
		// WriteTimeout: time.Second * 4,
	})
}

func setRoute() (route restful.Route) {
	route.Get("/111", controller.Pipe, controller.Midd)
	route.Post("/111", controller.IndexHandler)
	return route
}

func main() {
	quit := make(chan os.Signal)
	<-quit
}
```
