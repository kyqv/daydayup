package geecache

import (
	pb "daydayup/mygin/geecache/geecachepb"
	"daydayup/mygin/geecache/singleflight"
	"fmt"
	"log"
	"sync"
)

var (
	_mu     sync.RWMutex
	_groups = make(map[string]*Group)
)

//Group 是 GeeCache 最核心的数据结构，负责与用户的交互，并且控制缓存值存储和获取的流程。
type Group struct {
	//一个 Group 可以认为是一个缓存的命名空间，每个 Group 拥有一个唯一的名称 name
	name string
	//缓存未命中时获取源数据的回调(callback)
	getter Getter
	//实现的并发缓存
	mainCache cache
	//分布式节点，缓存未命中根据key选择远程调用方式
	peers  PeerPicker
	loader *singleflight.Group
}

//我们思考一下，如果缓存不存在，应从数据源（文件，数据库等）获取数据并添加到缓存中。
//GeeCache 是否应该支持多种数据源的配置呢？不应该，一是数据源的种类太多，没办法一一实现；二是扩展性不好。
//如何从源头获取数据，应该是用户决定的事情，我们就把这件事交给用户好了。
//因此，我们设计了一个回调函数(callback)，在缓存不存在时，调用这个函数，得到源数据。
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

//函数类型实现某一个接口，称之为接口型函数，方便使用者在调用时既能够传入函数作为参数，也能够传入实现了该接口的结构体作为参数。
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

func GetGroup(name string) *Group {
	_mu.RLock()
	g := _groups[name]
	_mu.RUnlock()
	return g
}

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	_mu.Lock()
	defer _mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	_groups[name] = g
	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	//缓存获取
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	//远端获取
	return g.load(key)
}

//实现了 PeerPicker 接口的 HTTPPool 注入到 Group 中
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

//使用 PickPeer() 方法选择节点，若非本机节点，则调用 getFromPeer() 从远程获取。若是本机节点或失败，则回退到 getLocally()
func (g *Group) load(key string) (value ByteView, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}

		return g.getLocally(key)
	})

	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

//自己去获取，拿到后更新到缓存
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	//更新缓存
	g.populateCache(key, value)
	return value, nil
}

//使用实现了 PeerGetter 接口的 httpGetter 从访问远程节点，获取缓存值。
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}
