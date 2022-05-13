package version

var (
	Litekube     = "v0.2.0"
	GitBranch    = "default-main" // branch of git
	GitVersion   = "v2.25.1"
	GitCommit    = "$HEAD"                 // sha1 from git, output of $(git rev-parse HEAD)
	GitTreeState = "clean"                 // state of git tree, either "clean" or "dirty"
	BuildDate    = "2022-04-10 T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	Kubernetes   = "v1.24.0"
	Kine         = "v0.9.0"
)
