// File: packages.go

package features

import (
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	sdkUtils "github.com/kairos-io/kairos-sdk/utils"
	"os"
)

// Immutability represents the Immutability feature.
// This install immucore and its required packages to run.
type Immutability struct {
}

func (g Immutability) Name() string {
	return "Immutability"
}

// Install installs the Immutability feature.
func (g Immutability) Install(s System, l sdkTypes.KairosLogger) error {
	// First packages so certs are updated

	err := s.Installer.Install(append(kernelPackages[string(s.Distro)], immucorePackages[string(s.Distro)]...), l)
	if err != nil {
		return err
	}

	l.Logger.Debug().Msg("Installing framework")
	frameworkImage, err := sdkUtils.GetImage("quay.io/kairos/framework:v2.14.1", "", nil, nil)
	err = sdkUtils.ExtractOCIImage(frameworkImage, "/")
	l.Logger.Debug().Msg("Installed framework")

	// Install config files that affect initramfs and rootfs only, which are the ones that affect immucore?
	return nil
}

// Remove removes the Immutability feature.
func (g Immutability) Remove(s System, l sdkTypes.KairosLogger) error {
	return nil
}

// Info logs information about the Immutability feature.
func (g Immutability) Info(s System, l sdkTypes.KairosLogger) {
	l.Info("Immutability feature. This installs immucore and the cloud configs files to support immutability")
}

// HasServices returns true if the Immutability feature has services.
func (g Immutability) HasServices() bool {
	return false
}

// InstallsPackages returns true if the Immutability feature installs packages.
func (g Immutability) InstallsPackages() bool {
	return true
}

// Installed returns true if the Immutability feature is installed.
func (g Immutability) Installed(s System, l sdkTypes.KairosLogger) bool {
	// Check if immucore binary is on the system
	_, err := os.Stat("/usr/bin/immucore")
	if err != nil {
		return false
	}
	return true
}
