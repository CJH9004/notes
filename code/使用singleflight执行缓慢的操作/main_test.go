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

func TestSlowPing(t *testing.T) {
	once.Do(mainSetup)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/slow", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestOptmizedPing(t *testing.T) {
	once.Do(mainSetup)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/optmized", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func testFunc(path string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	router.ServeHTTP(w, req)
}

func average(times []int64) float64 {
	var sum int64
	for _, v := range times {
		sum += v
	}
	return float64(sum) / float64(len(times))
}

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
