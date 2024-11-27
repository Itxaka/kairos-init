package system

import (
	"github.com/kairos-io/kairos-init/pkg/values"
	"os"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kairos-io/kairos-init/pkg/system/features"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
)

func DetectSystem(l sdkTypes.KairosLogger) features.System {
	// Detects the system
	s := features.System{
		Distro: values.Unknown,
		Family: features.UnknownFamily,
	}

	file, err := os.Open("/etc/os-release")
	if err != nil {
		return s
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	val, err := godotenv.Parse(file)
	if err != nil {
		return s
	}
	l.Logger.Trace().Interface("values", val).Msg("Read values from os-release")
	// Match values to distros
	switch values.Distro(val["ID"]) {
	case values.Debian:
		s.Distro = values.Debian
		s.Family = features.DebianFamily
		s.Installer = features.APTInstaller
	case values.Ubuntu:
		s.Distro = values.Ubuntu
		s.Family = features.DebianFamily
		s.Installer = features.APTInstaller
	case values.Fedora:
		s.Distro = values.Fedora
		s.Family = features.RedHatFamily
		s.Installer = features.DNFInstaller
	case values.RockyLinux:
		s.Distro = values.RockyLinux
		s.Family = features.RedHatFamily
		s.Installer = features.DNFInstaller
	case values.AlmaLinux:
		s.Distro = values.AlmaLinux
		s.Family = features.RedHatFamily
		s.Installer = features.DNFInstaller
	case values.RedHat:
		s.Distro = values.RedHat
		s.Family = features.RedHatFamily
		s.Installer = features.DNFInstaller
	case values.Arch:
		s.Distro = values.Arch
		s.Family = features.ArchFamily
		s.Installer = features.PacmanInstaller
	case values.Alpine:
		s.Distro = values.Alpine
		s.Family = features.AlpineFamily
		s.Installer = features.AlpineInstaller
	case values.OpenSUSELeap:
		s.Distro = values.OpenSUSELeap
		s.Family = features.SUSEFamily
		s.Installer = features.SUSEInstaller
	case values.OpenSUSETumbleweed:
		s.Distro = values.OpenSUSETumbleweed
		s.Family = features.SUSEFamily
		s.Installer = features.SUSEInstaller
	}

	// Match architecture
	switch values.Architecture(runtime.GOARCH) {
	case values.ArchAMD64:
		s.Arch = values.ArchAMD64
	case values.ArchARM64:
		s.Arch = values.ArchARM64
	}

	// Check if we are still unknown
	if s.Distro == values.Unknown {
		// Check ID_LIKE value
		// For some derivatives they ID will be their own but the ID_LIKE will be the parent
		// So we may be able to detect the parent and use the same family and such
		switch features.Family(val["ID_LIKE"]) {
		case features.DebianFamily:
			s.Distro = values.Debian
			s.Family = features.DebianFamily
			s.Installer = features.APTInstaller
		case features.RedHatFamily, features.Family(values.Fedora):
			s.Distro = values.Fedora
			s.Family = features.RedHatFamily
			s.Installer = features.DNFInstaller
		case features.ArchFamily:
			s.Distro = values.Arch
			s.Family = features.ArchFamily
			s.Installer = features.PacmanInstaller
		case features.SUSEFamily:
			s.Distro = values.OpenSUSELeap
			s.Family = features.SUSEFamily
			s.Installer = features.SUSEInstaller
		}
	}

	// Store the version
	s.Version = val["VERSION_ID"]

	// Store the name
	s.Name = val["PRETTY_NAME"]
	// Fallback to normal name
	if s.Name == "" {
		s.Name = val["NAME"]
	}

	s.Features = []features.Feature{}
	l.Logger.Debug().Interface("system", s).Msg("Detected system")

	return s
}
