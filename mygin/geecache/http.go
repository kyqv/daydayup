package geecache

import (
	"daydayup/mygin/geecache/consistenthash"
	pb "daydayup/mygin/geecache/geecachepb"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
)

const (
	_defaultBasePath = "/_geecache/"
	_defaultReplicas = 50
)

//既具备了提供 HTTP 服务的能力，也具备了根据具体的 key，创建 HTTP 客户端从远程节点获取缓存值的能力。
type HTTPPool struct {
	//记录自己的地址，包括主机名/IP 和端口
	self string
	//作为节点间通讯地址的前缀，默认是 /_geecache/
	//那么 http://example.com/_geecache/ 开头的请求，就用于节点间的访问。
	//因为一个主机上还可能承载其他的服务，加一段 Path 是一个好习惯。比如，大部分网站的 API 接口，一般以 /api 作为前缀。
	basePath string
	mu       sync.Mutex
	//一致性哈希算法的 Map，用来根据具体的 key 选择节点
	peers *consistenthash.Map
	//映射远程节点与对应的 httpGetter，每一个远程节点对应一个 httpGetter，因为 httpGetter 与远程节点的地址 baseURL 有关
	httpGetters map[string]*httpGetter
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: _defaultBasePath,
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//先判断访问路径的前缀是否是 basePath，不是返回错误。
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	//我们约定访问路径格式为 /<basepath>/<groupname>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	//通过 groupname 得到 group 实例，
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	//再使用 group.Get(key) 获取缓存数据。
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 将缓存值作为 httpResponse 的 body 返回
	// protobuf 广泛地应用于远程过程调用(RPC) 的二进制传输，使用 protobuf 的目的非常简单，为了获得更高的性能。
	// 传输前使用 protobuf 编码，接收方再进行解码，可以显著地降低二进制传输的大小。
	// 另外一方面，protobuf 可非常适合传输结构化数据，便于通信字段的扩展
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

// 实例化了一致性哈希算法,每个远端创建httpGetter对象
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//实例化了一致性哈希算法，并且添加了传入的节点
	p.peers = consistenthash.New(_defaultReplicas, nil)
	p.peers.Add(peers...)

	//并为每一个节点创建了一个 HTTP 客户端 httpGetter
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

//key-->gather
var _ PeerPicker = (*HTTPPool)(nil)

// 包装了一致性哈希算法的 Get() 方法，根据具体的 key，选择节点，返回节点对应的 HTTP 客户端。
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

type httpGetter struct {
	//表示将要访问的远程节点的地址
	baseURL string
}

var _ PeerGetter = (*httpGetter)(nil)

func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}
