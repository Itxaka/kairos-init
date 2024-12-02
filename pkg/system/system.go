package system

import (
	"github.com/joho/godotenv"
	"github.com/kairos-io/kairos-init/pkg/features"
	"github.com/kairos-io/kairos-init/pkg/values"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"github.com/sanity-io/litter"
	"os"
	"runtime"
)

func DetectSystem(l sdkTypes.KairosLogger) values.System {
	// Detects the system
	s := values.System{
		Distro: values.Unknown,
		Family: values.UnknownFamily,
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
		s.Family = values.DebianFamily
		s.Installer = features.APTInstaller
	case values.Ubuntu:
		s.Distro = values.Ubuntu
		s.Family = values.DebianFamily
		s.Installer = features.APTInstaller
	case values.Fedora:
		s.Distro = values.Fedora
		s.Family = values.RedHatFamily
		s.Installer = features.DNFInstaller
	case values.RockyLinux:
		s.Distro = values.RockyLinux
		s.Family = values.RedHatFamily
		s.Installer = features.DNFInstaller
	case values.AlmaLinux:
		s.Distro = values.AlmaLinux
		s.Family = values.RedHatFamily
		s.Installer = features.DNFInstaller
	case values.RedHat:
		s.Distro = values.RedHat
		s.Family = values.RedHatFamily
		s.Installer = features.DNFInstaller
	case values.Arch:
		s.Distro = values.Arch
		s.Family = values.ArchFamily
		s.Installer = features.PacmanInstaller
	case values.Alpine:
		s.Distro = values.Alpine
		s.Family = values.AlpineFamily
		s.Installer = features.AlpineInstaller
	case values.OpenSUSELeap:
		s.Distro = values.OpenSUSELeap
		s.Family = values.SUSEFamily
		s.Installer = features.SUSEInstaller
	case values.OpenSUSETumbleweed:
		s.Distro = values.OpenSUSETumbleweed
		s.Family = values.SUSEFamily
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
		switch values.Family(val["ID_LIKE"]) {
		case values.DebianFamily:
			s.Distro = values.Debian
			s.Family = values.DebianFamily
			s.Installer = features.APTInstaller
		case values.RedHatFamily, values.Family(values.Fedora):
			s.Distro = values.Fedora
			s.Family = values.RedHatFamily
			s.Installer = features.DNFInstaller
		case values.ArchFamily:
			s.Distro = values.Arch
			s.Family = values.ArchFamily
			s.Installer = features.PacmanInstaller
		case values.SUSEFamily:
			s.Distro = values.OpenSUSELeap
			s.Family = values.SUSEFamily
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

	s.Features = []values.Feature{}

	// Check if we have any workarounds for the system
	if s.Distro != values.Unknown {
		if workarounds, ok := values.WorkaroundsMap[s.Distro][s.Arch][s.Version]; ok {
			for _, w := range workarounds {
				l.Logger.Debug().Str("workaround", litter.Sdump(w)).Str("version", s.Version).Str("distro", s.Distro.String()).Str("arch", s.Arch.String()).Msg("Adding workaround")
				s.Workarounds = append(s.Workarounds, w)
			}
		}
	}

	return s
}
