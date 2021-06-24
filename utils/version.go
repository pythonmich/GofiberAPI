package utils

// GetVersion returns the build version of this binary
func GetVersion(config Config) string {
	var version = config.BuildVersion
	return version
}

