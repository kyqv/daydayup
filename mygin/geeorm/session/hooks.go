package session

import (
	"daydayup/mygin/geeorm/log"
	"reflect"
)

// Hooks constants
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// Hook，翻译为钩子，其主要思想是提前在可能增加功能的地方埋好(预设)一个钩子，
// 当我们需要重新修改或者增加这个地方的逻辑的时候，把扩展的类或者方法挂载到这个点即可。钩子的应用非常广泛，
// 例如 Github 支持的 travis 持续集成服务，当有 git push 事件发生时，会触发 travis 拉取新的代码进行构建。
// IDE 中钩子也非常常见，比如，当按下 Ctrl + s 后，自动格式化代码。再比如前端常用的 hot reload 机制，
// 前端代码发生变更时，自动编译打包，通知浏览器自动刷新页面，实现所写即所得。

// 钩子机制设计的好坏，取决于扩展点选择的是否合适。例如对于持续集成来说，代码如果不发生变更，
// 反复构建是没有意义的，因此钩子应设计在代码可能发生变更的地方，比如 MR、PR 合并前后。

// 那对于 ORM 框架来说，合适的扩展点在哪里呢？很显然，记录的增删查改前后都是非常合适的。

// CallMethod calls the registered hooks
func (s *Session) CallMethod(method string, value interface{}) {
	// s.RefTable().Model 或 value 即当前会话正在操作的对象，使用 MethodByName 方法反射得到该对象的方法。
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}
