// File: packages.go

package features

import (
	"fmt"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"github.com/sanity-io/litter"
	"os"
)

// Kernel represents the Kernel feature.
// This just links the latest kernel to /boot/vmlinuz
type Kernel struct {
}

func (g Kernel) Name() string {
	return "Kernel"
}

// Install installs the Immutability feature.
func (g Kernel) Install(s System, l sdkTypes.KairosLogger) error {
	kernelVersion, err := GetLatestKernel(l)
	err = os.Link("/boot/vmlinuz-"+kernelVersion, "/boot/vmlinuz")
	if err != nil {
		l.Logger.Error().Err(err).Msgf("Failed to link the kernel file: %s", err)
		return err
	}
	return nil
}

// Remove removes the Immutability feature.
func (g Kernel) Remove(s System, l sdkTypes.KairosLogger) error {
	return nil
}

// Info logs information about the Immutability feature.
func (g Kernel) Info(s System, l sdkTypes.KairosLogger) {
	l.Info("Kernel feature.")
}

// HasServices returns true if the Immutability feature has services.
func (g Kernel) HasServices() bool {
	return false
}

// InstallsPackages returns true if the Immutability feature installs packages.
func (g Kernel) InstallsPackages() bool {
	return true
}

// Installed returns true if the Immutability feature is installed.
func (g Kernel) Installed(s System, l sdkTypes.KairosLogger) bool {
	// Check if the kernel file exists
	if _, err := os.Stat("/boot/vmlinuz"); err == nil {
		l.Logger.Info().Msg("Kernel is already linked")
		stat, _ := os.Stat("/boot/vmlinuz")
		fmt.Println(litter.Sdump(stat))
		return true
	}
	return false
}
