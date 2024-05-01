package main

import (
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Backend struct {
	Host         string
	Active       atomic.Bool
	ReverseProxy *httputil.ReverseProxy
}

func NewBackend(host string) *Backend {
	rp := httputil.NewSingleHostReverseProxy(&url.URL{Host: host, Scheme: "http"})
	return &Backend{host, atomic.Bool{}, rp}
}

func (b *Backend) IsAlive() bool {
	return b.Active.Load()
}

func (b *Backend) SetLivenessStatus(s bool) {
	b.Active.Store(s)
}
