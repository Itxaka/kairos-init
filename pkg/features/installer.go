package features

import (
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"os"
	"os/exec"
)

type Installer string

const (
	APTInstaller    Installer = "apt-get"
	DNFInstaller    Installer = "dnf"
	PacmanInstaller Installer = "pacman"
	SUSEInstaller   Installer = "zypper"
	AlpineInstaller Installer = "apk"
)

func (i Installer) Install(packages []string, l sdkTypes.KairosLogger) error {
	var args []string
	var updateArgs []string
	cmd := string(i)
	l.Logger.Info().Str("installer", string(i)).Msg("Installing packages")
	switch i {
	case APTInstaller, DNFInstaller, SUSEInstaller:
		os.Setenv("DEBIAN_FRONTEND", "noninteractive")
		defer os.Unsetenv("DEBIAN_FRONTEND")
		updateArgs = []string{"-y", "update"}
		args = []string{"-y", "--no-install-recommends", "install"}
	case AlpineInstaller:
		updateArgs = []string{"update"}
		args = []string{"add", "--no-cache"}
	case PacmanInstaller:
		updateArgs = []string{"-Sy"}
		args = []string{"-S", "--noconfirm"}
	}
	// Run update
	l.Logger.Debug().Str("command", cmd).Strs("args", updateArgs).Msg("Running update")
	if err := CommandToLogger(cmd, updateArgs, l); err != nil {
		return err
	}

	// Run install
	args = append(args, packages...)
	l.Logger.Debug().Str("command", cmd).Strs("args", args).Msg("Running command")
	if err := CommandToLogger(cmd, args, l); err != nil {
		return err
	}

	return nil
}

func (i Installer) Remove(packages []string, l sdkTypes.KairosLogger) error {
	var args []string
	cmd := string(i)
	l.Logger.Info().Str("installer", string(i)).Msg("Removing packages")
	switch i {
	case APTInstaller, DNFInstaller, SUSEInstaller:
		args = []string{"-y", "remove"}
	case AlpineInstaller:
		args = []string{"remove", "--no-cache"}
	case PacmanInstaller:
		args = []string{"-R", "--noconfirm"}
	}
	args = append(args, packages...)
	l.Logger.Debug().Str("command", cmd).Strs("args", args).Msg("Running command")
	command := exec.Command(cmd, args...)
	out, err := command.CombinedOutput()
	if err != nil {
		l.Logger.Err(err).Str("output", string(out)).Str("command", cmd).Strs("args", args).Msg("Error running command")
		return err
	}
	return nil
}
