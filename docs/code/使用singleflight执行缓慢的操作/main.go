package main

import (
	"fmt"
	"net/http"

	"test-gin/singleflight"

	"github.com/gin-gonic/gin"
)

var g singleflight.Group

func main() {
	// Creates a router without any middleware by default
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.New()
	r.GET("/slow", slowPong)
	r.GET("/optmized", optmizedPong)
	return r
}

func slowPong(c *gin.Context) {
	c.String(200, longOp().(string))
}

func optmizedPong(c *gin.Context) {
	c.String(200, g.Do("pong", longOp).(string))
}

func longOp() interface{} {
	resp, err := http.Get("https://godoc.org/github.com/golang/groupcache")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	return "pong"
}
