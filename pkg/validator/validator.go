// Package validator
// This package is responsible for validating that a system has all the kairos bits and blops to be considered a Kairos system.
// All the checks should pass for the system to be considered a Kairos system correctly
package validator

import (
	"github.com/hashicorp/go-multierror"
	"github.com/kairos-io/kairos-init/pkg/log"
	"github.com/kairos-io/kairos-init/pkg/values"
	"os"
	"os/exec"
	"strings"
)

func ValidateFeatures(features []values.Feature) error {
	var err *multierror.Error
	for _, f := range features {
		switch f.Name() {
		case "immutability":
			err = multierror.Append(err, validateBinaries())
		case "kernel":
			err = multierror.Append(err, validateKernel())
		case "initrd":
			err = multierror.Append(err, validateInitrd())
		}
	}
	return err.ErrorOrNil()
}

// validateKernel checks if the kernel is there and its linked from /boot/vmlinuz
func validateKernel() error {
	log.Log.Logger.Info().Msg("Validating kernel")
	_, stat := os.Stat("/boot/vmlinuz")
	return stat
}

// validateBinaries checks if the expected binaries are there
func validateBinaries() error {
	log.Log.Logger.Info().Msg("Validating binaries")
	for _, bin := range values.BinariesCheck() {
		// check if we have multiple binaries to check
		// if we do, we check if any of them are present
		// if none are present, we return an error
		// we split the binaries by | for the same binary with different names in different OSes
		split := strings.Split(bin, "|")
		if len(split) > 1 {
			found := false
			for _, b := range split {
				_, err := exec.LookPath(b)
				if err == nil {
					found = true
					break
				}
			}
			if !found {
				return &exec.Error{Name: strings.Join(split, " or "), Err: exec.ErrNotFound}
			}
		} else {
			_, err := exec.LookPath(bin)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// validateSystemd checks if systemd is the correct version
// valid only fore UKI as we need a specific or higher version of systemd
func validateSystemd() error {
	log.Log.Logger.Info().Msg("Validating systemd")
	return nil
}

// validateInitrd checks if the initrd is there and its linked from /boot/initrd
// also checks if it has immucore and agent binary in the initrd
func validateInitrd() error {
	log.Log.Logger.Info().Msg("Validating initrd")
	_, stat := os.Stat("/boot/initrd")
	return stat
}
