package version

import (
	"fmt"
	"runtime"
)

type Info struct {
	LiteKube     string `json:"litekube"`
	GitVersion   string `json:"gitVersion"`
	GitBranch    string `json:"gitBranch"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
	Kubernetes   string `json:"kubernetes"`
	Kine         string `json:"kine"`
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in pkg/version/base.go
	return Info{
		LiteKube:     litekube,
		GitVersion:   gitVersion,
		GitBranch:    gitBranch,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Kubernetes:   kubernetes,
		Kine:         kine,
	}
}

func GetSimple() string {
	return fmt.Sprintf("Version: LiteKube %s, kubernetes %s, kine %s", litekube, kubernetes, kine)
}
