// File: packages.go

package features

import (
	"github.com/Masterminds/semver/v3"
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
	// Get the packages to install for this system
	packages, err := getPackages(s, l)

	// Now parse the packages with the templating engine
	finalMergedPkgs, err := values.PackageListToTemplate(packages, s.GetTemplateParams(), l)
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

// getPackages returns the packages to install for the Immutability feature.
// It parses the package maps and returns the packages that match the system version with semver
func getPackages(s values.System, l sdkTypes.KairosLogger) ([]string, error) {
	mergedPkgs := values.CommonPackages
	version, err := semver.NewVersion(s.Version)
	if err != nil {
		l.Logger.Error().Err(err).Str("version", s.Version).Msg("Error parsing version.")
		return mergedPkgs, err
	}

	// Go over all packages maps
	for _, packages := range []values.VersionMap{
		values.BasePackages[s.Distro][s.Arch],
		values.ImmucorePackages[s.Distro][s.Arch],
		values.KernelPackages[s.Distro][s.Arch],
		values.GrubPackages[s.Distro][s.Arch],
		values.SystemdPackages[s.Distro][s.Arch],
	} {
		// for each package map, check if the version matches the constraint
		for k, v := range packages {
			// Add them if they are common
			l.Logger.Debug().Str("constraint", k).Str("version", k).Msg("Checking constraint")
			if k == values.Common {
				mergedPkgs = append(mergedPkgs, v...)
				continue
			}
			constraint, err := semver.NewConstraint(k)
			if err != nil {
				l.Logger.Error().Err(err).Str("constraint", k).Msg("Error parsing constraint.")
				continue
			}
			// Also add them if the constraint matches
			if constraint.Check(version) {
				mergedPkgs = append(mergedPkgs, v...)
			}
		}
	}

	return mergedPkgs, nil
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
