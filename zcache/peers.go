package zcache

// HTTPGetter ...
type HTTPGetter interface {
  Get(group, key string) (ByteView, error)
}

// PeerPicker ...
type PeerPicker interface {
  PickPeer(key string) (HTTPGetter, bool)
}
