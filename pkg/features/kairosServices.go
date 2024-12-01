package features

import (
	"context"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/kairos-io/kairos-init/pkg/values"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
)

// Implement the Kairos services feature that installs required services in the system
// It implements the Feature interface.

// KairosServices represents the KairosServices feature.
type KairosServices struct {
	Order int
}

func (g KairosServices) GetOrder() int {
	return g.Order
}

func (g KairosServices) Name() string {
	return "KairosServices"
}

// Install installs the KairosServices feature.
func (g KairosServices) Install(s values.System, l sdkTypes.KairosLogger) error {
	conn, err := dbus.NewWithContext(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	servicesNames := []string{
		"systemd-pcrlock-make-policy",
	}

	// TODO: return the error?
	// Disable the service
	_, err = conn.DisableUnitFilesContext(context.Background(), servicesNames, false)
	if err != nil {
		l.Logger.Error().Err(err).Msg("Disabling services")
	}
	// Mask the service
	_, err = conn.MaskUnitFilesContext(context.Background(), servicesNames, false, true)
	if err != nil {
		l.Logger.Error().Err(err).Msg("Masking services")
	}
	return nil
}

// Remove removes the KairosServices feature.
func (g KairosServices) Remove(s values.System, l sdkTypes.KairosLogger) error {
	return nil
}

// Info logs information about the KairosServices feature.
func (g KairosServices) Info(s values.System, l sdkTypes.KairosLogger) {
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
func (g KairosServices) Installed(s values.System, l sdkTypes.KairosLogger) bool {
	return true
}

type EnableUnitFileChange struct {
	Type        string // Type of the change (one of symlink or unlink)
	Filename    string // File name of the symlink
	Destination string // Destination of the symlink
}
