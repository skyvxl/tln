package buildinfo

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func String() string {
	return Version + " (" + Commit + ", built " + BuildDate + ")"
}
