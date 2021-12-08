package zcache

import pb "zcache/zcachepb"

// HTTPGetter ...
type HTTPGetter interface {
  // Get(group, key string) (ByteView, error)
  Get(in *pb.Request, out *pb.Response) error
}

// PeerPicker ...
type PeerPicker interface {
  PickPeer(key string) (HTTPGetter, bool)
}
