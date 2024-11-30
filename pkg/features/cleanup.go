package features

import (
	"github.com/kairos-io/kairos-init/pkg/values"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"os"
)

// Cleanup represents the Cleanup feature.
// This feature is used to cleanup the system after the installation.
// Removes unnecessary files and directories. Cleans packages caches, etc...
type Cleanup struct {
	Order int
}

func (c Cleanup) Install(system values.System, logger sdkTypes.KairosLogger) error {
	// Empty machine-id
	f, err := os.OpenFile("/etc/machine-id", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_ = f.Close()
	// remove specific files
	for _, f := range values.FilesToRemove() {
		err = os.RemoveAll(f)
		if err != nil {
			logger.Logger.Error().Err(err).Str("file", f).Msg("Error removing file.")
			return err
		}
	}
	// Remove old initrds and kernels
	// We are only interested in keeping the one linked to /etc/initrd and /etc/vmlinuz
	// So we read the softlink at /boot/initrd and /boot/vmlinuz and remove the others
	return nil
}

func (c Cleanup) Remove(system values.System, logger sdkTypes.KairosLogger) error {
	return nil
}

func (c Cleanup) Info(system values.System, logger sdkTypes.KairosLogger) {
	logger.Logger.Info().Str("feature", c.Name()).Msg("Cleanup feature.")
}

func (c Cleanup) Installed(system values.System, logger sdkTypes.KairosLogger) bool {
	return false
}

func (c Cleanup) HasServices() bool {
	return false
}

func (c Cleanup) InstallsPackages() bool {
	return false
}

func (c Cleanup) Name() string {
	return "Cleanup"
}

func (c Cleanup) GetOrder() int {
	return c.Order
}
