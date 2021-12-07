package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"zcache"
)

var db = map[string]string{
	"zhangsan": "11",
	"lisi":     "12",
	"wangwu":   "13",
}

func createGroup(name string) *zcache.Group {
	return zcache.NewGroup(name, 1<<10, zcache.GetterFunc(func(key string) ([]byte, error) {
		log.Printf("[Database] search key: %q", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s is not exist", key)
	}))
}

func startCacheServer(addr string, peers []string, group *zcache.Group) {
	p := zcache.NewHTTPPool(addr)
	p.Set(peers...)
	group.RegisterPeers(p)
	log.Println("zcache server is runing at", addr)
	addr = strings.TrimPrefix(strings.TrimPrefix(addr, "http://"), "https://")
	log.Fatal(http.ListenAndServe(addr, p))
}

func startAPIServer(group *zcache.Group) {
	addr := "localhost:9999"
	http.Handle("/api", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		view, err := group.Get(key)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write(view.ByteSlice())
	}))
	log.Println("api server is runing at", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	var addr string
	var api bool
	var peers string
	flag.StringVar(&addr, "addr", "http://localhost:8001", "zcache server addr")
	flag.StringVar(&peers, "peers", "http://localhost:8001", "all peers")
	flag.BoolVar(&api, "api", false, "start a api server?")
	flag.Parse()
	addrs := strings.Split(peers, ",")
	group := createGroup("age")
	if api {
		go startAPIServer(group)
	}
	startCacheServer(addr, addrs, group)
}
