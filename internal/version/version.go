package version

import (
	"runtime"
	"strings"
)

var (
	// The version is of the format Major.Minor.Patch[-Prerelease][+BuildMetadata]
	version = "0.1.0"
	// metadata is extra build time data
	metadata = ""
	gitCommit = ""
)

type BuildInfo struct {
	// Version is the current semver.
	Version string `json:"version,omitempty"`
	// GitCommit is the git sha1.
	GitCommit string `json:"git_commit,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"go_version,omitempty"`
}

// GetVersion returns the semver string of the version
func GetVersion() string {
	if metadata == "" {
		 return version
	}
	return version + "+" + metadata
}

// GetUserAgent returns a user agent for user with an HTTP client
func GetUserAgent() string {
	return "osv/" + strings.TrimPrefix(GetVersion(), "v")
}


func Get() BuildInfo {
	v := BuildInfo{
		Version:   GetVersion(),
		GitCommit: gitCommit,
		GoVersion: runtime.Version(),
	}

	return v
}