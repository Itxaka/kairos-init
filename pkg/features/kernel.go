// File: packages.go

package features

import (
	"github.com/kairos-io/kairos-init/pkg/values"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"os"
)

// Kernel represents the Kernel feature.
// This just links the latest kernel to /boot/vmlinuz
type Kernel struct {
	Order int
}

func (g Kernel) GetOrder() int {
	return g.Order
}

func (g Kernel) Name() string {
	return "Kernel"
}

// Install installs the Immutability feature.
func (g Kernel) Install(s values.System, l sdkTypes.KairosLogger) error {
	kernelVersion, err := GetLatestKernel(l)
	err = os.Link("/boot/vmlinuz-"+kernelVersion, "/boot/vmlinuz")
	if err != nil {
		l.Logger.Error().Err(err).Msgf("Failed to link the kernel file: %s", err)
		return err
	}
	return nil
}

// Remove removes the Immutability feature.
func (g Kernel) Remove(s values.System, l sdkTypes.KairosLogger) error {
	return nil
}

// Info logs information about the Immutability feature.
func (g Kernel) Info(s values.System, l sdkTypes.KairosLogger) {
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
func (g Kernel) Installed(s values.System, l sdkTypes.KairosLogger) bool {
	// Check if the kernel file exists
	if _, err := os.Stat("/boot/vmlinuz"); err == nil {
		l.Logger.Debug().Msg("Kernel is already linked")
		return true
	}
	return false
}
