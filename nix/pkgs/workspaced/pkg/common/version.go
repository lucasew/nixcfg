package common

// BuildID is set at compile time via -ldflags
var BuildID = "dev"

// GetBuildID returns the build ID
func GetBuildID() string {
	return BuildID
}
