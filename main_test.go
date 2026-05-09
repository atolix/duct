package main

import (
	"testing"

	"github.com/atolix/duct/internal/scan"
)

func TestTakeTopN(t *testing.T) {
	entries := []scan.Entry{
		{Size: 10},
		{Size: 20},
		{Size: 30},
	}

	tests := []struct {
		name string
		n    int
		want int
	}{
		{"top 2", 2, 2},
		{"zero returns all", 0, 3},
		{"over returns all", 10, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := takeTopN(entries, tt.n)
			if len(got) != tt.want {
				t.Errorf("got %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestFilterByMinSize(t *testing.T) {
	entries := []scan.Entry{
		{Size: 100},
		{Size: 200},
		{Size: 300},
	}

	tests := []struct {
		name    string
		minSize int64
		wantLen int
	}{
		{"min 200", 200, 2},
		{"min 0 returns all", 0, 3},
		{"min high", 1000, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterByMinSize(entries, tt.minSize)
			if len(got) != tt.wantLen {
				t.Errorf("got %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestSortEntries(t *testing.T) {
	entries := []scan.Entry{
		{Size: 10},
		{Size: 30},
		{Size: 20},
	}

	sortEntries(entries)

	if entries[0].Size != 30 {
		t.Errorf("expected first element to be 30, got %d", entries[0].Size)
	}
}
