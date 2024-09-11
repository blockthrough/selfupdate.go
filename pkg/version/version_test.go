package version_test

import (
	"testing"

	"selfupdate.blockthrough.com/pkg/version"
)

func TestVersionCompare(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		want bool
	}{
		{"1.0.0", "1.0.0", false},
		{"1.0.0", "1.0.1", false},
		{"1.9.0", "1.10.0", false},
		{"1.10.0", "1.9.0", true},
	}

	for _, tt := range tests {
		got := version.Compare(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("%s > %s = %v; want %v", tt.a, tt.b, got, tt)
		}
	}
}
