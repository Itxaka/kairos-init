package values

import (
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"github.com/rs/zerolog"
	"github.com/sanity-io/litter"
	"strings"
)

// BinariesCheck returns the list of expected binaries to be in a kairos system
func BinariesCheck() []string {
	return []string{
		"immucore",
		"kairos-agent",
		"grub-install|grub2-install", // same binary, different names across OSes
	}
}

func FilesToRemove() []string {
	return []string{
		"/var/lib/dbus/machine-id",
		"/etc/hostname",
	}
}

type Architecture string

func (a Architecture) String() string {
	return string(a)
}

const (
	ArchAMD64 Architecture = "amd64"
	ArchARM64 Architecture = "arm64"
)

type Distro string

func (d Distro) String() string {
	return string(d)
}

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

const (
	ImmutabilitySentinel = "/etc/kairos/.inmmutability_installed"
	KernelSentinel       = "/etc/kairos/.kernel_installed"
	InitrdSentinel       = "/etc/kairos/.initrd_installed"
)

type Family string

func (f Family) String() string {
	return string(f)
}

// generic families that have things in common and we can apply to all of them
const (
	UnknownFamily Family = "unknown"
	DebianFamily  Family = "debian"
	RedHatFamily  Family = "redhat"
	ArchFamily    Family = "arch"
	AlpineFamily  Family = "alpine"
	SUSEFamily    Family = "suse"
)

type Feature interface {
	Install(System, sdkTypes.KairosLogger) error
	Remove(System, sdkTypes.KairosLogger) error
	Info(System, sdkTypes.KairosLogger)
	Installed(System, sdkTypes.KairosLogger) bool

	/*  Simple functions that return "fixed" values	thus they require no inputs */
	HasServices() bool
	InstallsPackages() bool
	Name() string
	GetOrder() int
}

type Features []Feature

// MarshalZerologObject For zerolog to be able to log the features in a nicer way
func (f Features) MarshalZerologObject(e *zerolog.Event) {
	for _, feature := range f {
		e.Str("name", feature.Name())
		e.Int("order", feature.GetOrder())
	}
}

type Workarounds []Workaround

// MarshalZerologObject For zerolog to be able to log the features in a nicer way
func (w Workarounds) MarshalZerologObject(e *zerolog.Event) {
	for _, workaround := range w {
		e.Str("name", litter.Sdump(workaround))
	}
}

type Workaround func(s *System, l sdkTypes.KairosLogger) error

// Installer is an interface that defines the methods to install and remove packages
type Installer interface {
	Install(packages []string, l sdkTypes.KairosLogger) error
	Remove(packages []string, l sdkTypes.KairosLogger) error
}

// System Represents a kairos-to-be system
type System struct {
	Name        string
	Distro      Distro
	Family      Family
	Version     string
	Arch        Architecture
	Features    Features
	Workarounds Workarounds `json:"-,omitempty" yaml:"-,omitempty"`
	Installer   Installer
	Force       bool // Force will force the installation of the features without checking the Installed() method
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
func (s *System) ApplyWorkarounds(l sdkTypes.KairosLogger) error {
	for _, w := range s.Workarounds {
		err := w(s, l)
		if err != nil {
			return err
		}
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

// GetTemplateParams returns a map of parameters that can be used in a template
func (s *System) GetTemplateParams() map[string]string {
	return map[string]string{
		"distro":  s.Distro.String(),
		"version": s.Version,
		"arch":    s.Arch.String(),
		"family":  s.Family.String(),
	}
}

func (s System) MarshalZerologObject(e *zerolog.Event) {
	e.Str("name", s.Name).
		Str("distro", s.Distro.String()).
		Str("family", s.Family.String()).
		Str("version", s.Version).
		Str("arch", s.Arch.String())

	e.Object("features", s.Features)
	e.Object("workarounds", s.Workarounds)
}
