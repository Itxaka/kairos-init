package features

import (
	"github.com/joho/godotenv"
	"github.com/kairos-io/kairos-init/pkg/values"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"os"
)

// KairosRelease implements the Feature interface.
// it fills the /etc/kairos-release file with the release version and such
type KairosRelease struct {
	Order int
}

func (g KairosRelease) GetOrder() int {
	return g.Order
}

func (k KairosRelease) Install(system values.System, logger sdkTypes.KairosLogger) error {
	// TODO: read from the current running os-release file to fill some of this fields
	// KAIROS_FLAVOR_RELEASE
	// TODO: Add missing stuff and check which ones do we really need?
	// TODO: Some inputs may be required to fill this file like k3s VS standard or the model?
	// KAIROS_ID="kairos"
	// KAIROS_IMAGE_LABEL="24.04-standard-amd64-generic-v3.2.3-4-gae5349e"
	// KAIROS_VARIANT="standard"
	// KAIROS_RELEASE="v3.2.3-4-gae5349e"
	// KAIROS_SOFTWARE_VERSION_PREFIX="k3s"
	// KAIROS_NAME="kairos-standard-ubuntu-24.04"
	// KAIROS_VERSION="v3.2.3-4-gae5349e"
	// KAIROS_ID_LIKE="kairos-standard-ubuntu-24.04"
	// KAIROS_FLAVOR_RELEASE="24.04"
	// KAIROS_REGISTRY_AND_ORG="quay.io/kairos"
	// KAIROS_GITHUB_REPO="kairos-io/kairos"
	// KAIROS_IMAGE_REPO="quay.io/kairos/ubuntu:24.04-standard-amd64-generic-v3.2.3-4-gae5349e"
	// KAIROS_ARTIFACT="kairos-ubuntu-24.04-standard-amd64-generic-v3.2.3-4-gae5349e"
	//
	// KAIROS_TARGETARCH="amd64"
	// KAIROS_BUG_REPORT_URL="https://github.com/kairos-io/kairos/issues"
	// KAIROS_HOME_URL="https://github.com/kairos-io/kairos"
	// KAIROS_VERSION_ID="v3.2.3-4-gae5349e"
	// KAIROS_PRETTY_NAME="kairos-standard-ubuntu-24.04 v3.2.3-4-gae5349e"
	releaseInfo := map[string]string{
		"KAIROS_VERSION": system.Version,
		"KAIROS_ARCH":    system.Arch.String(),
		"KAIROS_FLAVOR":  system.Distro.String(),
		"KAIROS_FAMILY":  system.Family.String(),
		"KAIROS_MODEL":   "generic", // NEEDED or it breaks boot!
		"KAIROS_VARIANT": "core",    // Maybe needed?
		"TEST":           "HALLO",
	}
	return godotenv.Write(releaseInfo, "/etc/kairos-release")
}

func (k KairosRelease) Remove(system values.System, logger sdkTypes.KairosLogger) error {
	return os.Remove("/etc/kairos-release")
}

func (k KairosRelease) Info(system values.System, logger sdkTypes.KairosLogger) {
	logger.Info("Kernel feature.")
}

func (k KairosRelease) Installed(system values.System, logger sdkTypes.KairosLogger) bool {
	_, err := os.Stat("/etc/kairos-release")
	return err == nil
}

func (k KairosRelease) HasServices() bool {
	return false
}

func (k KairosRelease) InstallsPackages() bool {
	return false
}

func (k KairosRelease) Name() string {
	return "KairosRelease"
}
