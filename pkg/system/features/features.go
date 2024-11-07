package features

import (
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
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
}
