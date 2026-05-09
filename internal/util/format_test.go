package util

import "testing"

func TestHumanize(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{0, "0B"},
		{512, "512B"},
		{1024, "1.0KB"},
		{1536, "1.5KB"},
		{1024 * 1024, "1.0MB"},
		{1024 * 1024 * 1024, "1.0GB"},
	}

	for _, tt := range tests {
		got := Humanize(tt.input)
		if got != tt.want {
			t.Errorf("HUmanize(%d) = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestShortenWithHome(t *testing.T) {
	tests := []struct {
		path string
		home string
		want string
	}{
		{"/Users/foo/file.txt", "/Users/foo", "~/file.txt"},
		{"/Users/foo", "/Users/foo", "~"},
		{"/other/path", "/Users/foo", "/other/path"},
	}

	for _, tt := range tests {
		got := shortenWithHome(tt.path, tt.home)
		if got != tt.want {
			t.Errorf("got %s, want %s", got, tt.want)
		}
	}
}
