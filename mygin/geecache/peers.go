package geecache

import pb "daydayup/mygin/geecache/geecachepb"

//用于根据传入的 key 选择相应节点 PeerGetter。
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

//用于从对应 group 查找缓存值
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
