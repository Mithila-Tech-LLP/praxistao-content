package main

import "testing"

func TestZAdd_ZRange_SortedOrder(t *testing.T) {
	s := NewStore()
	s.ZAdd("scores", 3.0, "c")
	s.ZAdd("scores", 1.0, "a")
	s.ZAdd("scores", 2.0, "b")

	members := s.ZRange("scores", 0, -1)
	if len(members) != 3 {
		t.Fatalf("expected 3 members, got %d", len(members))
	}
	if members[0] != "a" || members[1] != "b" || members[2] != "c" {
		t.Fatalf("expected sorted order [a b c], got %v", members)
	}
}

func TestZScore(t *testing.T) {
	s := NewStore()
	s.ZAdd("z", 5.5, "m")

	score, ok := s.ZScore("z", "m")
	if !ok || score != 5.5 {
		t.Fatalf("expected 5.5 ok=true, got %v ok=%v", score, ok)
	}
}

func TestZScore_Missing(t *testing.T) {
	s := NewStore()
	s.ZAdd("z", 1.0, "x")

	_, ok := s.ZScore("z", "missing")
	if ok {
		t.Fatal("expected ok=false for missing member")
	}
}

func TestZRank(t *testing.T) {
	s := NewStore()
	s.ZAdd("z", 10.0, "x")
	s.ZAdd("z", 20.0, "y")
	s.ZAdd("z", 30.0, "w")

	rank, ok := s.ZRank("z", "x")
	if !ok || rank != 0 {
		t.Fatalf("expected rank=0, got %d ok=%v", rank, ok)
	}
	rank, ok = s.ZRank("z", "y")
	if !ok || rank != 1 {
		t.Fatalf("expected rank=1, got %d ok=%v", rank, ok)
	}
	rank, ok = s.ZRank("z", "w")
	if !ok || rank != 2 {
		t.Fatalf("expected rank=2, got %d ok=%v", rank, ok)
	}
}

func TestZRank_Missing(t *testing.T) {
	s := NewStore()
	s.ZAdd("z", 1.0, "a")

	_, ok := s.ZRank("z", "missing")
	if ok {
		t.Fatal("expected ok=false for missing member")
	}
}

func TestZRangeWithScores(t *testing.T) {
	s := NewStore()
	s.ZAdd("z", 1.0, "a")
	s.ZAdd("z", 2.0, "b")

	members := s.ZRangeWithScores("z", 0, -1)
	if len(members) != 2 {
		t.Fatalf("expected 2, got %d", len(members))
	}
	if members[0].Member != "a" || members[0].Score != 1.0 {
		t.Fatalf("unexpected first member: %+v", members[0])
	}
	if members[1].Member != "b" || members[1].Score != 2.0 {
		t.Fatalf("unexpected second member: %+v", members[1])
	}
}

func TestZAdd_UpdateScore(t *testing.T) {
	s := NewStore()
	s.ZAdd("z", 1.0, "a")
	s.ZAdd("z", 2.0, "b")
	s.ZAdd("z", 0.5, "a") // update a's score to 0.5

	members := s.ZRange("z", 0, -1)
	if len(members) != 2 {
		t.Fatalf("expected 2 members after update, got %d", len(members))
	}
	if members[0] != "a" {
		t.Fatalf("after score update, a should be first (score=0.5), got %v", members)
	}
	score, _ := s.ZScore("z", "a")
	if score != 0.5 {
		t.Fatalf("expected updated score 0.5, got %v", score)
	}
}

func TestZRange_NegativeStop(t *testing.T) {
	s := NewStore()
	s.ZAdd("z", 1.0, "a")
	s.ZAdd("z", 2.0, "b")
	s.ZAdd("z", 3.0, "c")

	members := s.ZRange("z", 1, -1)
	if len(members) != 2 || members[0] != "b" || members[1] != "c" {
		t.Fatalf("expected [b c] with start=1 stop=-1, got %v", members)
	}
}
