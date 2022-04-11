package version

var (
	litekube     string = "v0.2.0"
	gitBranch    string = "default-main" // branch of git
	gitVersion          = "v2.25.1"
	gitCommit           = "$HEAD"                 // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState        = "clean"                 // state of git tree, either "clean" or "dirty"
	buildDate           = "2022-04-10 T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	kubernetes          = "v1.23.5"
	kine                = "v0.9.0"
)
