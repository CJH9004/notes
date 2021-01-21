# 使用singleflight执行缓慢的操作

singleflight是groupcache中用于防止缓存更新时出现的瞬时负载升高问题，同样也适用于优化普通程序

## 实现

[代码连接](../../code/使用singleflight执行缓慢的操作/singleflight)

```go
package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	ret interface{}
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (ret interface{})) (ret interface{}) {
	g.mu.Lock()
	// 初始化g.m
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 如果已经存在调用，那么等待调用完成后返回结果
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.ret
	}

	// 否则开始调用
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.ret = fn()
	c.wg.Done()

	// 调用完成，清除调用后返回
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.ret
}
```

## 应用

### 1.合并缓慢函数的执行

[代码连接](../../code/使用singleflight执行缓慢的操作)

```go
func longOp() interface{} {
	resp, err := http.Get("https://godoc.org/github.com/golang/groupcache")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	return "pong"
}
```

假设这个函数是一个http handler函数的一部分，那么该接口并发时执行多次函数longOp，如果使用singleflight模式来优化这个函数，那么在第一次该函数调用完成之前，都不会再次触发该函数调用。这里使用gin做httpserver。

```go
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
```

测试代码

```go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

var (
	once   sync.Once
	router *gin.Engine
)

func mainSetup() {
	router = setupRouter()
}

// 发出一次get请求
func testFunc(path string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	router.ServeHTTP(w, req)
}

// 算平均值
func average(times []int64) float64 {
	var sum int64
	for _, v := range times {
		sum += v
	}
	return float64(sum) / float64(len(times))
}

// 简单实现压力测试，输出总时间，平均时间
func TestTimesOfSlow(t *testing.T) {
	once.Do(mainSetup)

	times := make([]int64, 0, 100)
	start := time.Now().Unix()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now().UnixNano()
			testFunc("/slow")
			times = append(times, time.Now().UnixNano()-start)
		}()
	}
	wg.Wait()
	fmt.Println("slow:", time.Now().Unix()-start, "s", "; average:", average(times)/1000000, "ms")
}

func TestTimesOfOptmized(t *testing.T) {
	once.Do(mainSetup)

	times := make([]int64, 0, 100)
	start := time.Now().Unix()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now().UnixNano()
			testFunc("/optmized")
			times = append(times, time.Now().UnixNano()-start)
		}()
	}
	wg.Wait()
	fmt.Println("optmized:", time.Now().Unix()-start, "s", "; average:", average(times)/1000000, "ms")
}
```

结果如下，优化后的，相当于一次调用longOp的时间

```go
slow: 8 s ; average: 3792.3356421 ms
optmized: 0 s ; average: 236.6876463738509 ms
```