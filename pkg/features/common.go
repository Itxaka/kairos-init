package features

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/kairos-io/kairos-init/pkg/values"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var Features = map[string]values.Feature{
	"immutability": Immutability{Order: 1},
	"release":      KairosRelease{Order: 2},
	"kernel":       Kernel{Order: 3},
	"initrd":       Initrd{Order: 5},
	"clean":        Cleanup{Order: 10},
}

var FeatSupported = func() []string {
	var f []string
	for k := range Features {
		f = append(f, k)
	}
	return f
}

func GetFeature(name string) values.Feature {
	return Features[strings.ToLower(name)]
}

func FeatureSupported(name string) bool {
	for _, f := range FeatSupported() {
		if strings.ToLower(f) == strings.ToLower(name) {
			return true
		}
	}
	return false
}

// GetOrderedFeatures Returns the features in order
func GetOrderedFeatures() []values.Feature {
	var orderedFeatures []values.Feature
	for _, feature := range Features {
		orderedFeatures = append(orderedFeatures, feature)
	}
	sort.Slice(orderedFeatures, func(i, j int) bool {
		return orderedFeatures[i].GetOrder() < orderedFeatures[j].GetOrder()
	})
	return orderedFeatures
}

func CommandToLogger(cmd string, args []string, l sdkTypes.KairosLogger) (err error) {
	command := exec.Command(cmd, args...)
	stdout, _ := command.StdoutPipe()
	stderr, _ := command.StderrPipe()
	var stderrBuffer bytes.Buffer

	if err = command.Start(); err != nil {
		l.Logger.Err(err).Str("command", cmd).Strs("args", args).Msg("Error running command")
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			l.Logger.Debug().Msg(scanner.Text())
		}
	}()

	// Store the error output
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			stderrBuffer.WriteString(scanner.Text() + "\n")
		}
	}()

	if err := command.Wait(); err != nil {
		l.Logger.Err(err).Str("command", cmd).Err(fmt.Errorf(stderrBuffer.String())).Strs("args", args).Msg("Error running command")
		return err
	}
	return err
}

func GetLatestKernel(l sdkTypes.KairosLogger) (string, error) {
	var kernelVersion string
	modulesPath := "/lib/modules"
	// Read the directories under /lib/modules
	dirs, err := os.ReadDir(modulesPath)
	if err != nil {
		l.Logger.Error().Msgf("Failed to read the directory %s: %s", modulesPath, err)
		return kernelVersion, err

	}

	var versions []*semver.Version

	for _, dir := range dirs {
		if dir.IsDir() {
			// Parse the directory name as a semver version
			version, err := semver.NewVersion(dir.Name())
			if err != nil {
				l.Logger.Error().Msgf("Failed to parse the version %s: %s", dir.Name(), err)
				continue
			}
			versions = append(versions, version)
		}
	}

	sort.Sort(semver.Collection(versions))
	kernelVersion = versions[0].String()
	return kernelVersion, nil
}
