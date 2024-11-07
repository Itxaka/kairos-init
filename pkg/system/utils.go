package system

import (
	"os"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kairos-io/kairos-init/pkg/system/features"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
)

func DetectSystem(l sdkTypes.KairosLogger) features.System {
	// Detects the system
	s := features.System{
		Distro: features.Unknown,
		Family: features.UnknownFamily,
	}

	file, err := os.Open("/etc/os-release")
	if err != nil {
		return s
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	values, err := godotenv.Parse(file)
	if err != nil {
		return s
	}
	l.Logger.Trace().Interface("values", values).Msg("Read values from os-release")
	// Match values to distros
	switch features.Distro(values["ID"]) {
	case features.Debian:
		s.Distro = features.Debian
		s.Family = features.DebianFamily
		s.Installer = features.APTInstaller
	case features.Ubuntu:
		s.Distro = features.Ubuntu
		s.Family = features.DebianFamily
		s.Installer = features.APTInstaller
	case features.Fedora:
		s.Distro = features.Fedora
		s.Family = features.RedHatFamily
		s.Installer = features.DNFInstaller
	case features.RockyLinux:
		s.Distro = features.RockyLinux
		s.Family = features.RedHatFamily
		s.Installer = features.DNFInstaller
	case features.AlmaLinux:
		s.Distro = features.AlmaLinux
		s.Family = features.RedHatFamily
		s.Installer = features.DNFInstaller
	case features.RedHat:
		s.Distro = features.RedHat
		s.Family = features.RedHatFamily
		s.Installer = features.DNFInstaller
	case features.Arch:
		s.Distro = features.Arch
		s.Family = features.ArchFamily
		s.Installer = features.PacmanInstaller
	case features.Alpine:
		s.Distro = features.Alpine
		s.Family = features.AlpineFamily
		s.Installer = features.AlpineInstaller
	case features.OpenSUSELeap:
		s.Distro = features.OpenSUSELeap
		s.Family = features.SUSEFamily
		s.Installer = features.SUSEInstaller
	case features.OpenSUSETumbleweed:
		s.Distro = features.OpenSUSETumbleweed
		s.Family = features.SUSEFamily
		s.Installer = features.SUSEInstaller
	}

	// Match architecture
	switch features.Architecture(runtime.GOARCH) {
	case features.ArchAMD64:
		s.Arch = features.ArchAMD64
	case features.ArchARM64:
		s.Arch = features.ArchARM64
	}

	// Check if we are still unknown
	if s.Distro == features.Unknown {
		// Check ID_LIKE value
		// For some derivatives they ID will be their own but the ID_LIKE will be the parent
		// So we may be able to detect the parent and use the same family and such
		switch features.Family(values["ID_LIKE"]) {
		case features.DebianFamily:
			s.Distro = features.Debian
			s.Family = features.DebianFamily
			s.Installer = features.APTInstaller
		case features.RedHatFamily, features.Family(features.Fedora):
			s.Distro = features.Fedora
			s.Family = features.RedHatFamily
			s.Installer = features.DNFInstaller
		case features.ArchFamily:
			s.Distro = features.Arch
			s.Family = features.ArchFamily
			s.Installer = features.PacmanInstaller
		case features.SUSEFamily:
			s.Distro = features.OpenSUSELeap
			s.Family = features.SUSEFamily
			s.Installer = features.SUSEInstaller
		}
	}

	// Store the version
	s.Version = values["VERSION_ID"]

	// Store the name
	s.Name = values["PRETTY_NAME"]
	// Fallback to normal name
	if s.Name == "" {
		s.Name = values["NAME"]
	}

	s.Features = []features.Feature{}
	l.Logger.Debug().Interface("system", s).Msg("Detected system")

	return s
}
