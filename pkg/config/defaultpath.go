package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// I personally hate that XDG for MacOS has become $HOME/Library/...
// and this is my party so I'll cry if I want to

func DefaultFilePath() string {
	if runtime.GOOS == "darwin" {
		return defaultFilePathDarwin()
	}

	return defaultFilePathNonDarwin()
}

func defaultFilePathDarwin() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".config", "sickle", "config.dhall")
}

func defaultFilePathNonDarwin() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	return filepath.Join(configDir, "sickle", "config.dhall")
}
