package features

import (
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"strings"
)

type Distro string

// Individual distros for when we need to be specific
const (
	Unknown            Distro = "unknown"
	Debian             Distro = "debian"
	Ubuntu             Distro = "ubuntu"
	RedHat             Distro = "redhat"
	RockyLinux         Distro = "rocky"
	AlmaLinux          Distro = "almalinux"
	Fedora             Distro = "fedora"
	Arch               Distro = "arch"
	Alpine             Distro = "alpine"
	OpenSUSELeap       Distro = "opensuse-leap"
	OpenSUSETumbleweed Distro = "opensuse-tumbleweed"
)

type Family string

// generic families that have things in common and we can apply to all of them
const (
	UnknownFamily Family = "unknown"
	DebianFamily  Family = "debian"
	RedHatFamily  Family = "redhat"
	ArchFamily    Family = "arch"
	AlpineFamily  Family = "alpine"
	SUSEFamily    Family = "suse"
)

type Architecture string

const (
	ArchAMD64 Architecture = "amd64"
	ArchARM64 Architecture = "arm64"
)

// System Represents a given system
type System struct {
	Name        string
	Distro      Distro
	Family      Family
	Version     string
	Arch        Architecture
	Features    []Feature
	Workarounds []func() error
	Installer   Installer
}

// ApplyFeatures will apply the features to the system
func (s *System) ApplyFeatures(l sdkTypes.KairosLogger) error {
	for _, f := range s.Features {
		if f.Installed(*s, l) {
			l.Logger.Info().Str("feature", f.Name()).Msg("Feature already installed.")
			continue
		} else {
			l.Logger.Info().Str("feature", f.Name()).Msg("Installing feature...")
			err := f.Install(*s, l)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RemoveFeatures will remove the features from the system
func (s *System) RemoveFeatures(l sdkTypes.KairosLogger) error {
	for _, f := range s.Features {
		err := f.Remove(*s, l)
		if err != nil {
			return err
		}
	}
	return nil
}

// ApplyWorkarounds will apply the workarounds to the system
func (s *System) ApplyWorkarounds() error {
	for _, w := range s.Workarounds {
		w()
	}
	return nil
}

func (s *System) GetFeature(name string, l sdkTypes.KairosLogger) Feature {
	for _, f := range s.Features {
		if strings.ToLower(f.Name()) == strings.ToLower(name) {
			return f
		}
	}
	return nil
}

func (s *System) AddFeature(f Feature) {
	s.Features = append(s.Features, f)
}

func (s *System) RemoveFeature(f Feature) {
	for i, feature := range s.Features {
		if feature == f {
			s.Features = append(s.Features[:i], s.Features[i+1:]...)
			return
		}
	}
}
