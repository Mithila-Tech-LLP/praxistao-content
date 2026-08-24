package main

import "testing"

func TestARPCache_LRUEviction(t *testing.T) {
	c := NewARPCache(2)
	c.Set("10.0.0.1", "aa:aa")
	c.Set("10.0.0.2", "bb:bb")

	if mac, ok := c.Lookup("10.0.0.1"); !ok || mac != "aa:aa" {
		t.Fatalf("Lookup(10.0.0.1) = (%q, %v), want (aa:aa, true)", mac, ok)
	}

	c.Set("10.0.0.3", "cc:cc")

	if _, ok := c.Lookup("10.0.0.2"); ok {
		t.Error("expected 10.0.0.2 to be evicted (least recently used)")
	}
	if mac, ok := c.Lookup("10.0.0.1"); !ok || mac != "aa:aa" {
		t.Errorf("Lookup(10.0.0.1) = (%q, %v), want (aa:aa, true)", mac, ok)
	}
	if mac, ok := c.Lookup("10.0.0.3"); !ok || mac != "cc:cc" {
		t.Errorf("Lookup(10.0.0.3) = (%q, %v), want (cc:cc, true)", mac, ok)
	}
}

func TestARPCache_LookupMiss(t *testing.T) {
	c := NewARPCache(2)
	if _, ok := c.Lookup("192.168.1.1"); ok {
		t.Error("expected ok=false for a never-set IP")
	}
}
