package common

import (
	"os"
	"path"
	"runtime"
)

// UserHomeDir returns the home directory for the user the process is
// running under.
func UserHomeDir() string {
	if runtime.GOOS == "windows" { // Windows
		return os.Getenv("USERPROFILE")
	}

	// *nix
	return os.Getenv("HOME")
}

// DefaultAWSProfilePath return the default AWS profile path.
func DefaultAWSProfilePath() string {
	return path.Join(UserHomeDir(), ".aws")
}

// DefaultOLAWSConfigPath return the default OL-AWS configuration file.
func DefaultOLAWSConfigPath() string {
	return path.Join(UserHomeDir(), ".ol-aws.yml")
}
