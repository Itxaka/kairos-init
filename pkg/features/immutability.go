// File: packages.go

package features

import (
	"github.com/blang/semver/v4"
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
	// First base packages so certs are updated + immucore
	version, err := semver.ParseTolerant(s.Version)
	if err != nil {
		l.Logger.Error().Err(err).Str("version", s.Version).Msg("Error parsing version.")
		return err
	}
	mergedPkgs := values.CommonPackages
	if _, ok := values.BasePackages[s.Distro][s.Arch][values.Common]; ok {
		mergedPkgs = append(mergedPkgs, values.BasePackages[s.Distro][s.Arch][values.Common]...)
	}

	if _, ok := values.BasePackages[s.Distro][s.Arch][s.Version]; ok {
		mergedPkgs = append(mergedPkgs, values.BasePackages[s.Distro][s.Arch][s.Version]...)
	}

	// Add immucore required packages
	if _, ok := values.ImmucorePackages[s.Distro][s.Arch][values.Common]; ok {
		mergedPkgs = append(mergedPkgs, values.ImmucorePackages[s.Distro][s.Arch][values.Common]...)
	}
	// Add immucore required packages for the distro version with versioning
	for k, v := range values.ImmucorePackages[s.Distro][s.Arch] {
		if k == values.Common {
			continue
		}
		constraint, err := semver.ParseRange(k)
		if err != nil {
			l.Logger.Error().Err(err).Str("constraint", k).Msg("Error parsing constraint.")
			continue
		}
		if constraint(version) {
			mergedPkgs = append(mergedPkgs, v...)
		}
	}
	// Add kernel packages
	mergedPkgs = append(mergedPkgs, values.KernelPackages[s.Distro]...)
	// TODO: Somehow we need to know here if we are installing grub or systemd-boot
	mergedPkgs = append(mergedPkgs, values.GrubPackages[s.Distro][s.Arch]...)
	if _, ok := values.SystemdPackages[s.Distro][s.Arch][values.Common]; ok {
		// Add common systemd packages
		mergedPkgs = append(mergedPkgs, values.SystemdPackages[s.Distro][s.Arch][values.Common]...)
	}
	// Add specific systemd packages for the distro version
	if _, ok := values.SystemdPackages[s.Distro][s.Arch][s.Version]; ok {
		mergedPkgs = append(mergedPkgs, values.SystemdPackages[s.Distro][s.Arch][s.Version]...)
	}

	// Now parse the packages with the templating engine
	finalMergedPkgs, err := values.PackageListToTemplate(mergedPkgs, s.GetTemplateParams(), l)
	if err != nil {
		l.Logger.Error().Err(err).Msg("Error parsing base packages.")
		return err
	}
	err = s.Installer.Install(finalMergedPkgs, l)
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
