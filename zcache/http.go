package zcache

import (
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const defaultBasePath = "/_zcache"

// HTTPPool ...
type HTTPPool struct {
	mu          sync.Mutex
	self        string
	httpGetters map[string]*httpGetter
	basePath    string
	peers       []string
}

// NewHTTPPool ...
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log ...
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// Set ...
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = peers
	p.httpGetters = make(map[string]*httpGetter)
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{
			baseURL: fmt.Sprintf("%s%s", peer, p.basePath),
		}
	}
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.Log("%s %s", req.URL.Path, req.Method)
	parts := strings.SplitN(req.URL.Path[len(p.basePath)+1:], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName, err := url.PathUnescape(parts[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	key, err := url.PathUnescape(parts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

// PickPeer ...
func (p *HTTPPool) PickPeer(key string) (HTTPGetter, bool) {
	h := crc32.Checksum([]byte(key), crc32.IEEETable)
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.peers) == 0 {
		return nil, false
	}
	if peer := p.peers[int(h)%len(p.peers)]; peer != p.self {
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ http.Handler = (*HTTPPool)(nil)
var _ PeerPicker = (*HTTPPool)(nil)

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group, key string) (ByteView, error) {
	u := fmt.Sprintf("%s/%s/%s", h.baseURL, group, key)
	res, err := http.Get(u)
	if err != nil {
		return ByteView{}, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{data}, nil
}

var _ HTTPGetter = (*httpGetter)(nil)
