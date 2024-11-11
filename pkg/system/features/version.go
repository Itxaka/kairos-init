package features

import (
	"github.com/joho/godotenv"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"os"
)

// Implement the version feature that generates a /etc/kairos-release file with the proper versions

// Version represents the Version feature.
type Version struct {
}

func (g Version) Name() string {
	return "Version"
}

// Install installs the Version feature.
func (g Version) Install(s System, l sdkTypes.KairosLogger) error {
	values := map[string]string{
		"KAIROS_FLAVOR": string(s.Distro),
	}
	err := godotenv.Write(values, "/etc/kairos-release")
	if err != nil {
		l.Logger.Error().Err(err).Msg("Error writing /etc/kairos-release")
		return err
	}
	return nil
}

// Remove removes the Version feature.
func (g Version) Remove(s System, l sdkTypes.KairosLogger) error {
	return nil
}

// Info logs information about the Version feature.
func (g Version) Info(s System, l sdkTypes.KairosLogger) {
	l.Info("Version feature.")
}

// HasServices returns true if the Version feature has services.
func (g Version) HasServices() bool {
	return false
}

// InstallsPackages returns true if the Version feature installs packages.
func (g Version) InstallsPackages() bool {
	return true
}

// Installed returns true if the Version feature is installed.
func (g Version) Installed(s System, l sdkTypes.KairosLogger) bool {
	// Check if the initrd file exists
	if _, err := os.Stat("/etc/kairos-release"); err == nil {
		l.Logger.Info().Msg("Version is already generated")
		return true
	}
	return false
}
