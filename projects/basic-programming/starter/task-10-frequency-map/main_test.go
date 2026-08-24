package main

import (
	"reflect"
	"testing"
)

func TestFrequencyMap(t *testing.T) {
	words := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
	got := FrequencyMap(words)
	want := map[string]int{"apple": 3, "banana": 2, "cherry": 1}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("FrequencyMap[%q] = %d, want %d", k, got[k], v)
		}
	}
}

func TestFrequencyMapEmpty(t *testing.T) {
	got := FrequencyMap([]string{})
	if len(got) != 0 {
		t.Errorf("FrequencyMap([]) should be empty, got %v", got)
	}
}

func TestTopN(t *testing.T) {
	freq := map[string]int{"apple": 3, "banana": 2, "cherry": 1}
	if got := TopN(freq, 2); !reflect.DeepEqual(got, []string{"apple", "banana"}) {
		t.Errorf("TopN(freq, 2) = %v, want [apple banana]", got)
	}
	if got := TopN(freq, 1); !reflect.DeepEqual(got, []string{"apple"}) {
		t.Errorf("TopN(freq, 1) = %v, want [apple]", got)
	}
	if got := TopN(freq, 5); len(got) != 3 {
		t.Errorf("TopN(freq, 5) = %v, want 3 items", got)
	}
}

func TestTopNTieBreak(t *testing.T) {
	freq := map[string]int{"cat": 2, "ant": 2, "bat": 2}
	got := TopN(freq, 2)
	want := []string{"ant", "bat"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TopN tie-break: got %v, want %v", got, want)
	}
}
