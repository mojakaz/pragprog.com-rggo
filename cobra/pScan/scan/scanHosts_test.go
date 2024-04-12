package scan_test

import (
	"net"
	"pragprog.com/rggo/cobra/pScan/scan"
	"strconv"
	"testing"
	"time"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}
	if ps.Open.String() != "closed" {
		t.Errorf("Expected %q, got %q instead\n", "closed", ps.Open.String())
	}
	ps.Open = true
	if ps.Open.String() != "open" {
		t.Errorf("Expected %q, got %q instead\n", "open", ps.Open.String())
	}
}
func TestRunHostFoundTCP(t *testing.T) {
	testCases := []struct {
		name        string
		expectState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}
	host := "localhost"
	hl := &scan.HostsList{}
	hl.Add(host)
	ports := []int{}
	// Init ports, 1 open, 1 closed
	for _, tc := range testCases {
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()
		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}
		ports = append(ports, port)
		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}
	res := scan.Run(hl, ports, false, 1*time.Second)
	// Verify results for HostFound test
	if len(res) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}
	if res[0].NotFound {
		t.Errorf("Expected host %q to be found\n", host)
	}
	if len(res[0].PortStates) != 2 {
		t.Fatalf("Expected 2 port states, got %d instead\n", len(res[0].PortStates))
	}
	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Expected port %d, got %d instead\n", ports[i], res[0].PortStates[i].Port)
		}
		if res[0].PortStates[i].Open.String() != tc.expectState {
			t.Errorf("Expected port %d to be %s\n", ports[i], tc.expectState)
		}
	}
}

func TestRunHostFoundUDP(t *testing.T) {
	testCases := []struct {
		name        string
		expectState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "open"},
	}
	host := "127.0.0.1"
	hl := &scan.HostsList{}
	hl.Add(host)
	ports := []int{}
	// Init ports, 1 open, 1 closed
	for _, tc := range testCases {
		ln, err := net.ListenPacket("udp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()
		_, portStr, err := net.SplitHostPort(ln.LocalAddr().String())
		if err != nil {
			t.Fatal(err)
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}
		ports = append(ports, port)
		if tc.name == "ClosedPort" {
			if err := ln.Close(); err != nil {
				t.Fatal(err)
			}
		}
	}
	res := scan.Run(hl, ports, true, 1*time.Second)
	// Verify results for HostFound test
	if len(res) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}
	if res[0].NotFound {
		t.Errorf("Expected host %q to be found\n", host)
	}
	if len(res[0].PortStates) != 2 {
		t.Fatalf("Expected 2 port states, got %d instead\n", len(res[0].PortStates))
	}
	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Expected port %d, got %d instead\n", ports[i], res[0].PortStates[i].Port)
		}
		if res[0].PortStates[i].Open.String() != tc.expectState {
			t.Errorf("Expected port %d to be %s\n", ports[i], tc.expectState)
		}
	}
}

func TestRunHostNotFound(t *testing.T) {
	host := "389.389.389.389"
	hl := &scan.HostsList{}
	hl.Add(host)
	res := scan.Run(hl, []int{}, false, 1*time.Second)
	if len(res) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}
	if !res[0].NotFound {
		t.Errorf("Expected host %q NOT to be found\n", host)
	}
	if len(res[0].PortStates) != 0 {
		t.Fatalf("Expected 0 port states, got %d instead\n", len(res[0].PortStates))
	}
}
