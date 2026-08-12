package show_version

import (
	"testing"
)

func TestInfo_String(t *testing.T) {
	tests := []struct {
		name     string
		info     Info
		expected string
	}{
		{
			name: "Full info",
			info: Info{
				Version: "v1.0.0",
				Commit:  "a1b2c3d",
				Date:    "2026-08-12T22:00:00Z",
			},
			expected: "pjdoc version v1.0.0 (commit a1b2c3d, built at 2026-08-12T22:00:00Z)",
		},
		{
			name: "Version only",
			info: Info{
				Version: "v1.0.0",
				Commit:  "none",
				Date:    "unknown",
			},
			expected: "pjdoc version v1.0.0",
		},
		{
			name: "Version and Commit only",
			info: Info{
				Version: "v1.0.0",
				Commit:  "a1b2c3d",
				Date:    "unknown",
			},
			expected: "pjdoc version v1.0.0 (commit a1b2c3d)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.info.String()
			if got != tt.expected {
				t.Errorf("Info.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGetInfo_Default(t *testing.T) {
	info := GetInfo()
	if info.Version == "" {
		t.Errorf("GetInfo().Version should not be empty")
	}
}
