package features

import (
	"github.com/kairos-io/kairos-init/pkg/values"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"io"
	"os"
	"path/filepath"
)

type Overrides struct {
	Order int
}

func (o Overrides) Install(system values.System, logger sdkTypes.KairosLogger) error {
	// If there is a dir called overrides, copy whatever is on to / and then remove the dir

	_, err := os.Stat("/overrides")
	if err == nil {
		// Copy the contents of the dir to /
		filepath.WalkDir("/overrides", func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				// Create the dir
				err := os.MkdirAll(path, os.ModePerm)
				if err != nil {
					return err
				}
			} else {
				// Copy the file
				// Open the source file
				src, err := os.Open(path)
				if err != nil {
					return err
				}
				defer src.Close()

				// Create the destination file
				dst, err := os.Create(path)
				if err != nil {
					return err
				}
				defer dst.Close()

				// Copy the file
				_, err = io.Copy(dst, src)
				if err != nil {
					return err
				}
			}
			return nil
		})
	} else {
		logger.Info("No overrides found.")
	}

	return nil
}

func (o Overrides) Remove(system values.System, logger sdkTypes.KairosLogger) error {
	return nil
}

func (o Overrides) Info(system values.System, logger sdkTypes.KairosLogger) {
	logger.Info("Overrides feature.")
}

func (o Overrides) Installed(system values.System, logger sdkTypes.KairosLogger) bool {
	return false
}

func (o Overrides) HasServices() bool {
	return false
}

func (o Overrides) InstallsPackages() bool {
	return false
}

func (o Overrides) Name() string {
	return "Overrides"
}

func (o Overrides) GetOrder() int {
	return o.Order
}
