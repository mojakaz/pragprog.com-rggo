package scan

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single TCP port
type PortState struct {
	Port int
	Open state
}

type state bool

// String converts the boolean value of state to a human-readable string
func (s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

// scanPort performs a port scan on a single TCP port
func scanPort(host string, port int, useUDP bool, timeout time.Duration) PortState {
	p := PortState{
		Port: port,
	}
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	var network string
	if !useUDP {
		network = "tcp"
	} else {
		network = "udp"
	}
	scanConn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return p
	}
	scanConn.Close()
	p.Open = true
	return p
}

// Results represents the scan results for a single host
type Results struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

// Run performs a port scan on the hosts list
func Run(hl *HostsList, ports []int, useUDP bool, timeout time.Duration) []Results {
	res := make([]Results, 0, len(hl.Hosts))
	for _, h := range hl.Hosts {
		r := Results{
			Host: h,
		}
		if _, err := net.LookupHost(h); err != nil {
			r.NotFound = true
			res = append(res, r)
			continue
		}
		for _, p := range ports {
			if p < 1 || p > 65535 {
				continue
			}
			r.PortStates = append(r.PortStates, scanPort(h, p, useUDP, timeout))
		}
		res = append(res, r)
	}
	return res
}
