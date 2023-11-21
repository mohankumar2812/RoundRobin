package main

import (
	loadbalancer "loadbalancer/LoadBalancer"
	// "loadbalancer/Servers"
)

func main() {
	loadbalancer.MakeLoadBalancer(6)
	// servers.RunServers(6)
}
