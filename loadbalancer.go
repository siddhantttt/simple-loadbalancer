package main

import (
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type LoadBalancer struct {
	backends   []*Backend
	httpClient *http.Client
	Algorithm  func([]*Backend) *Backend
}

func NewLoadBalancer(urls ...string) *LoadBalancer {
	lb := &LoadBalancer{make([]*Backend, 0), &http.Client{}, SelectRandomBackend}
	for _, u := range urls {
		backend := NewBackend(u)
		lb.backends = append(lb.backends, backend)
	}
	return lb
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := lb.Algorithm(lb.backends)
	b.ReverseProxy.ServeHTTP(w, r)
}

func SelectRandomBackend(backends []*Backend) *Backend {
	return backends[rand.Intn(2)]
}

func (lb *LoadBalancer) StartHealthcheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, backend := range lb.backends {
				go func(b *Backend) {
					b.SetLivenessStatus(checkHealth(b))
				}(backend)
			}
		}
	}
}

func checkHealth(b *Backend) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", b.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	_ = conn.Close()
	return true
}
