// Package version exposes the build version of gitcat.
package version

import "runtime/debug"

// version is set at link time by the Makefile and by GoReleaser:
//
//	-X github.com/i-zaitsev/gitcat/internal/version.version=v0.1.0
var version string

// Get returns the version of the build. It is the version set at link time,
// which is a release tag like v0.1.0 or the output of git describe. Without
// it, as in a build by go install, it is the module version recorded in the
// build info. When neither is known, it returns "dev".
func Get() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}
