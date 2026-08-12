package show_version

import (
	"runtime/debug"
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
		{
			name: "Version and Date only",
			info: Info{
				Version: "v1.0.0",
				Commit:  "none",
				Date:    "2026-08-12T22:00:00Z",
			},
			expected: "pjdoc version v1.0.0 (built at 2026-08-12T22:00:00Z)",
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

func TestParseBuildInfo(t *testing.T) {
	t.Run("Nil BuildInfo preserves original values", func(t *testing.T) {
		v, c, d := parseBuildInfo(nil, "v1.0.0", "abc1234", "2026-08-12")
		if v != "v1.0.0" || c != "abc1234" || d != "2026-08-12" {
			t.Errorf("parseBuildInfo(nil) returned unexpected values: %s, %s, %s", v, c, d)
		}
	})

	t.Run("Populates version, commit, and date from BuildInfo settings", func(t *testing.T) {
		bi := &debug.BuildInfo{
			Main: debug.Module{
				Version: "v2.0.0",
			},
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "1234567890abcdef"},
				{Key: "vcs.time", Value: "2026-08-12T12:34:56Z"},
				{Key: "vcs.modified", Value: "true"},
			},
		}

		v, c, d := parseBuildInfo(bi, "dev", "none", "unknown")
		if v != "v2.0.0" {
			t.Errorf("expected version v2.0.0, got %s", v)
		}
		if c != "1234567-dirty" {
			t.Errorf("expected commit 1234567-dirty, got %s", c)
		}
		if d != "2026-08-12T12:34:56Z" {
			t.Errorf("expected date 2026-08-12T12:34:56Z, got %s", d)
		}
	})

	t.Run("Handles short revision without modified flag", func(t *testing.T) {
		bi := &debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "12345"},
				{Key: "vcs.modified", Value: "false"},
			},
		}

		_, c, _ := parseBuildInfo(bi, "dev", "none", "unknown")
		if c != "12345" {
			t.Errorf("expected commit 12345, got %s", c)
		}
	})
}

func TestGetInfo_Default(t *testing.T) {
	info := GetInfo()
	if info.Version == "" {
		t.Errorf("GetInfo().Version should not be empty")
	}
}
