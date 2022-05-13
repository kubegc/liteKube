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
		LiteKube:     Litekube,
		GitVersion:   GitVersion,
		GitBranch:    GitBranch,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Kubernetes:   Kubernetes,
		Kine:         Kine,
	}
}

func GetSimple() string {
	return fmt.Sprintf("Version: LiteKube %s, kubernetes %s, kine %s", Litekube, Kubernetes, Kine)
}
