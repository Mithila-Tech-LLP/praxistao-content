package main

import "testing"

func TestRouter_LongestPrefixMatch(t *testing.T) {
	r := NewRouter()
	r.AddRoute("0.0.0.0/0", "default-gw")
	r.AddRoute("10.0.0.0/8", "isp-a")
	r.AddRoute("10.1.0.0/16", "isp-b")
	r.AddRoute("10.1.2.0/24", "isp-c")

	cases := []struct {
		ip       string
		nextHop  string
	}{
		{"10.1.2.5", "isp-c"},
		{"10.1.5.5", "isp-b"},
		{"10.5.5.5", "isp-a"},
		{"192.168.1.1", "default-gw"},
	}

	for _, c := range cases {
		nextHop, found := r.Lookup(c.ip)
		if !found {
			t.Errorf("Lookup(%q): expected found=true", c.ip)
			continue
		}
		if nextHop != c.nextHop {
			t.Errorf("Lookup(%q) = %q, want %q", c.ip, nextHop, c.nextHop)
		}
	}
}

func TestRouter_NoDefaultRoute_NotFound(t *testing.T) {
	r := NewRouter()
	r.AddRoute("10.0.0.0/8", "isp-a")

	_, found := r.Lookup("192.168.1.1")
	if found {
		t.Error("expected found=false when no route matches and there is no default route")
	}
}
