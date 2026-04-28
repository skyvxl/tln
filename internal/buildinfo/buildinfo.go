// Package buildinfo provides build-time information such as version and commit hash.
package buildinfo

// Build-time variables injected by the linker.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// String returns a formatted build info string containing the version, commit, and build date.
func String() string {
	return Version + " (" + Commit + ", built " + BuildDate + ")"
}
