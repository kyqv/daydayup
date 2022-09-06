package singleflight

import "sync"

//call 代表正在进行中，或已经结束的请求。使用 sync.WaitGroup 锁避免重入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group 是 singleflight 的主数据结构，管理不同 key 的请求(call)
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
}

//接收 2 个参数，第一个参数是 key，第二个参数是一个函数 fn。
//针对相同的 key，无论 Do 被调用多少次，函数 fn 都只会被调用一次，等待 fn 调用结束了，返回返回值或错误。
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	//发起请求前加锁
	c.wg.Add(1)
	//添加到 g.m，表明 key 已经有对应的请求在处理
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	// 更新 g.m
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
