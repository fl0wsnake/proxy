package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/elazarl/goproxy"
)

type entry struct {
	ip   string
	time int64
}

type byTime []entry

func (s byTime) Len() int {
	return len(s)
}
func (s byTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byTime) Less(i, j int) bool {
	return s[i].time < s[j].time
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	ips := make(map[string]int64)

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err == nil {
			ips[ip] = time.Now().Unix()
		}
		return req, nil
	})

	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		l := make([]entry, 0, len(ips))
		for k, v := range ips {
			l = append(l, entry{k, v})
		}
		sort.Sort(byTime(l))
		for _, e := range l {
			fmt.Fprintf(w, "%s: %s\n", e.ip, time.Unix(e.time, 0).Format(time.RFC850))
		}
	})

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), proxy))
}
