package main

type LoadBalancingAlgorithm interface {
	SelectBackend([]*Backend) *Backend
}

type RoundRobinAlgorithm struct {
}

func (r *RoundRobinAlgorithm) SelectBackend(backends []*Backend) *Backend {
	return backends[0]
}
