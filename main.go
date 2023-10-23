package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

type Service struct {
	dialer    proxy.Dialer
	targetURL *url.URL
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(s.targetURL)
	proxy.Transport = &http.Transport{Dial: s.dialer.Dial}
	originalDirector := proxy.Director

	startTime := time.Now()
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = s.targetURL.Host
	}
	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("response: [%d] %s %s %v", resp.StatusCode, resp.Request.Method, resp.Request.URL.String(), time.Since(startTime))
		return nil
	}

	proxy.ServeHTTP(w, r)
}

var socks5proxy, upstream, listen string

func main() {
	flag.StringVar(&socks5proxy, "socks5", "127.0.0.1:1080", "Socks5 代理地址")
	flag.StringVar(&upstream, "upstream", "https://api.openai.com", "目标地址")
	flag.StringVar(&listen, "listen", ":8081", "监听地址")
	flag.Parse()

	dialer, err := proxy.SOCKS5("tcp", socks5proxy, nil, proxy.Direct)
	if err != nil {
		panic(err)
	}

	target, err := url.Parse(upstream)
	if err != nil {
		panic(err)
	}

	service := &Service{
		dialer:    dialer,
		targetURL: target,
	}

	http.ListenAndServe(listen, service)
}
