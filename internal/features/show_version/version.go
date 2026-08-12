package show_version

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// ldflags から注入されるビルド情報変数
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// Info はバージョン情報を保持する構造体です。
type Info struct {
	Version string
	Commit  string
	Date    string
}

// GetInfo はバージョン情報を取得します。
// ldflags による注入値がない場合は runtime/debug.ReadBuildInfo から補完を試みます。
func GetInfo() Info {
	buildInfo, _ := debug.ReadBuildInfo()
	v, c, d := parseBuildInfo(buildInfo, Version, Commit, Date)
	return Info{
		Version: v,
		Commit:  c,
		Date:    d,
	}
}

func parseBuildInfo(buildInfo *debug.BuildInfo, v, c, d string) (string, string, string) {
	if buildInfo != nil {
		if v == "dev" && buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
			v = buildInfo.Main.Version
		}

		var revision string
		var modTime string
		var modified bool

		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.time":
				modTime = setting.Value
			case "vcs.modified":
				if setting.Value == "true" {
					modified = true
				}
			}
		}

		if c == "none" && revision != "" {
			if len(revision) > 7 {
				c = revision[:7]
			} else {
				c = revision
			}
			if modified {
				c += "-dirty"
			}
		}

		if d == "unknown" && modTime != "" {
			d = modTime
		}
	}

	return v, c, d
}

// String はバージョン情報を表示用の文字列に整形します。
func (i Info) String() string {
	var details []string
	if i.Commit != "" && i.Commit != "none" {
		details = append(details, fmt.Sprintf("commit %s", i.Commit))
	}
	if i.Date != "" && i.Date != "unknown" {
		details = append(details, fmt.Sprintf("built at %s", i.Date))
	}

	if len(details) > 0 {
		return fmt.Sprintf("pjdoc version %s (%s)", i.Version, strings.Join(details, ", "))
	}
	return fmt.Sprintf("pjdoc version %s", i.Version)
}
