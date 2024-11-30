// File: packages.go

package features

import (
	"github.com/kairos-io/kairos-init/pkg/values"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	sdkUtils "github.com/kairos-io/kairos-sdk/utils"
	"os"
)

// Immutability represents the Immutability feature.
// This install immucore and its required packages to run.
type Immutability struct {
	Order int
}

func (g Immutability) GetOrder() int {
	return g.Order
}

func (g Immutability) Name() string {
	return "Immutability"
}

// Install installs the Immutability feature.
func (g Immutability) Install(s values.System, l sdkTypes.KairosLogger) error {
	// First base packages so certs are updated
	pkg := values.BasePackages[s.Distro]
	// Add packages in which immucore depends
	pkg = append(pkg, values.ImmucorePackages[s.Distro]...)
	// Add kernel packages
	pkg = append(pkg, values.KernelPackages[s.Distro]...)
	// TODO: Somehow we need to know here if we are installing grub or systemd-boot
	// Add grub packages
	pkg = append(pkg, values.GrubPackages[s.Distro][s.Arch]...)
	// Add systemd packages
	pkg = append(pkg, values.SystemdPackages[s.Distro][s.Arch]...)

	err := s.Installer.Install(pkg, l)
	if err != nil {
		return err
	}

	l.Logger.Debug().Msg("Installing framework")
	frameworkImage, err := sdkUtils.GetImage("quay.io/kairos/framework:v2.14.4", "", nil, nil)
	err = sdkUtils.ExtractOCIImage(frameworkImage, "/")
	l.Logger.Debug().Msg("Installed framework")

	// Install config files that affect initramfs and rootfs only, which are the ones that affect immucore?
	err = os.MkdirAll("/etc/kairos", os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}
	_, err = os.Create(values.ImmutabilitySentinel)
	if err != nil {
		return err
	}
	return nil
}

// Remove removes the Immutability feature.
func (g Immutability) Remove(s values.System, l sdkTypes.KairosLogger) error {
	return nil
}

// Info logs information about the Immutability feature.
func (g Immutability) Info(s values.System, l sdkTypes.KairosLogger) {
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
func (g Immutability) Installed(s values.System, l sdkTypes.KairosLogger) bool {
	if s.Force {
		return false
	}
	// Check sentinel
	_, err := os.Stat(values.ImmutabilitySentinel)
	if err != nil {
		return false
	}
	return true
}
