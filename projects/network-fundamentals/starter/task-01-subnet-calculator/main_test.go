package main

import (
	"net"
	"testing"
)

func TestSubnetInfo_Slash24(t *testing.T) {
	network, broadcast, first, last, hostCount, err := SubnetInfo("192.168.1.0/24")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !network.Equal(net.ParseIP("192.168.1.0")) {
		t.Errorf("network = %v, want 192.168.1.0", network)
	}
	if !broadcast.Equal(net.ParseIP("192.168.1.255")) {
		t.Errorf("broadcast = %v, want 192.168.1.255", broadcast)
	}
	if !first.Equal(net.ParseIP("192.168.1.1")) {
		t.Errorf("first = %v, want 192.168.1.1", first)
	}
	if !last.Equal(net.ParseIP("192.168.1.254")) {
		t.Errorf("last = %v, want 192.168.1.254", last)
	}
	if hostCount != 254 {
		t.Errorf("hostCount = %d, want 254", hostCount)
	}
}

func TestSubnetInfo_Slash30(t *testing.T) {
	network, broadcast, first, last, hostCount, err := SubnetInfo("10.0.0.0/30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !network.Equal(net.ParseIP("10.0.0.0")) {
		t.Errorf("network = %v, want 10.0.0.0", network)
	}
	if !broadcast.Equal(net.ParseIP("10.0.0.3")) {
		t.Errorf("broadcast = %v, want 10.0.0.3", broadcast)
	}
	if !first.Equal(net.ParseIP("10.0.0.1")) {
		t.Errorf("first = %v, want 10.0.0.1", first)
	}
	if !last.Equal(net.ParseIP("10.0.0.2")) {
		t.Errorf("last = %v, want 10.0.0.2", last)
	}
	if hostCount != 2 {
		t.Errorf("hostCount = %d, want 2", hostCount)
	}
}

func TestSubnetInfo_Slash28(t *testing.T) {
	network, broadcast, first, last, hostCount, err := SubnetInfo("172.16.5.130/28")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !network.Equal(net.ParseIP("172.16.5.128")) {
		t.Errorf("network = %v, want 172.16.5.128", network)
	}
	if !broadcast.Equal(net.ParseIP("172.16.5.143")) {
		t.Errorf("broadcast = %v, want 172.16.5.143", broadcast)
	}
	if !first.Equal(net.ParseIP("172.16.5.129")) {
		t.Errorf("first = %v, want 172.16.5.129", first)
	}
	if !last.Equal(net.ParseIP("172.16.5.142")) {
		t.Errorf("last = %v, want 172.16.5.142", last)
	}
	if hostCount != 14 {
		t.Errorf("hostCount = %d, want 14", hostCount)
	}
}

func TestSubnetInfo_HostIPNotAligned(t *testing.T) {
	network, broadcast, first, last, hostCount, err := SubnetInfo("192.168.1.100/24")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !network.Equal(net.ParseIP("192.168.1.0")) {
		t.Errorf("network = %v, want 192.168.1.0", network)
	}
	if !broadcast.Equal(net.ParseIP("192.168.1.255")) {
		t.Errorf("broadcast = %v, want 192.168.1.255", broadcast)
	}
	if !first.Equal(net.ParseIP("192.168.1.1")) {
		t.Errorf("first = %v, want 192.168.1.1", first)
	}
	if !last.Equal(net.ParseIP("192.168.1.254")) {
		t.Errorf("last = %v, want 192.168.1.254", last)
	}
	if hostCount != 254 {
		t.Errorf("hostCount = %d, want 254", hostCount)
	}
}

func TestSubnetInfo_Invalid(t *testing.T) {
	_, _, _, _, _, err := SubnetInfo("not-a-cidr")
	if err == nil {
		t.Error("expected error for invalid CIDR, got nil")
	}
}
