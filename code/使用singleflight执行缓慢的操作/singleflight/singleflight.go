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
