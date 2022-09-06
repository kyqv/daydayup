package lru

import "container/list"

type Cache struct {
	//允许使用的最大内存
	maxBytes int64
	//当前已使用的内存
	nbytes int64
	ll     *list.List
	cache  map[string]*list.Element
	// 某条记录被移除时的回调函数，可以为 nil
	OnEvicted func(key string, value Value)
}

//双向链表节点的数据类型，
//在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射
type entry struct {
	key   string
	value Value
}

//为了通用性，我们允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小。
type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		//如果键存在，则更新对应节点的值，并将该节点移到队头。
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//不存在则是新增场景，首先队头添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系。
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	//更新 c.nbytes，如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点。
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

//查找功能
func (c *Cache) Get(key string) (value Value, ok bool) {
	//字典中找到对应的双向链表的节点
	if ele, ok := c.cache[key]; ok {
		//将该节点移动到队头
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//缓存淘汰。即移除最近最少访问的节点
func (c *Cache) RemoveOldest() {
	//取到队首节点
	ele := c.ll.Back()
	if ele != nil {
		//从链表中删除。
		c.ll.Remove(ele)
		//从字典中 c.cache 删除该节点的映射关系。
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		//更新当前所用的内存 c.nbytes。
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		//如果回调函数 OnEvicted 不为 nil，则调用回调函数。
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
