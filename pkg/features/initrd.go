package features

import (
	"github.com/kairos-io/kairos-init/pkg/values"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"os"
)

// Implement the initrd feature that generates a initrd with the needed packages on it and configuration.
// It implements the Feature interface.

// Initrd represents the Initrd feature.
type Initrd struct {
	Order int
}

func (g Initrd) GetOrder() int {
	return g.Order
}

func (g Initrd) Name() string {
	return "Initrd"
}

// Install installs the Initrd feature.
func (g Initrd) Install(s values.System, l sdkTypes.KairosLogger) error {
	kernelVersion, err := GetLatestKernel(l)
	if err != nil {
		return err
	}

	cmd := "dracut"
	args := []string{"-v", "-f", "/boot/initrd", kernelVersion}
	l.Logger.Debug().Str("command", cmd).Strs("args", args).Msg("Running command")
	if err := CommandToLogger(cmd, args, l); err != nil {
		return err
	}
	return nil
}

// Remove removes the Initrd feature.
func (g Initrd) Remove(s values.System, l sdkTypes.KairosLogger) error {
	return nil
}

// Info logs information about the Initrd feature.
func (g Initrd) Info(s values.System, l sdkTypes.KairosLogger) {
	l.Info("Initrd feature.")
}

// HasServices returns true if the Initrd feature has services.
func (g Initrd) HasServices() bool {
	return false
}

// InstallsPackages returns true if the Initrd feature installs packages.
func (g Initrd) InstallsPackages() bool {
	return true
}

// Installed returns true if the Initrd feature is installed.
func (g Initrd) Installed(s values.System, l sdkTypes.KairosLogger) bool {
	// Check if the initrd file exists
	if _, err := os.Stat("/boot/initrd"); err == nil {
		l.Logger.Debug().Msg("Initrd is already generated")
		return true
	}
	return false
}
