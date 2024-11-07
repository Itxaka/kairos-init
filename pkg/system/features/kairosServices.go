package features

import (
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
)

// Implement the Kairos services feature that installs required services in the system
// It implements the Feature interface.

// KairosServices represents the KairosServices feature.
type KairosServices struct {
}

func (g KairosServices) Name() string {
	return "KairosServices"
}

// Install installs the KairosServices feature.
func (g KairosServices) Install(s System, l sdkTypes.KairosLogger) error {
	return nil
}

// Remove removes the KairosServices feature.
func (g KairosServices) Remove(s System, l sdkTypes.KairosLogger) error {
	return nil
}

// Info logs information about the KairosServices feature.
func (g KairosServices) Info(s System, l sdkTypes.KairosLogger) {
	l.Info("Kairos Services feature.")
}

// HasServices returns true if the KairosServices feature has services.
func (g KairosServices) HasServices() bool {
	return true
}

// InstallsPackages returns true if the KairosServices feature installs packages.
func (g KairosServices) InstallsPackages() bool {
	return false
}

// Installed returns true if the KairosServices feature is installed.
func (g KairosServices) Installed(s System, l sdkTypes.KairosLogger) bool {
	return true
}
