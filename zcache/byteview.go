package zcache

// ByteView ...
type ByteView struct {
  b []byte
}

// Len ...
func (v ByteView) Len() int {
  return len(v.b)
}

// ByteSlice ...
func (v ByteView) ByteSlice() []byte {
  return cloneBytes(v.b)
}

func (v ByteView) String() string {
  return string(v.b)
}

func cloneBytes(b []byte) []byte {
  c := make([]byte, len(b))
  copy(c, b)
  return c
}
