package main

import (
	"fmt"
	"github.com/kairos-io/kairos-init/pkg/system"
	"github.com/kairos-io/kairos-init/pkg/system/features"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
	"github.com/spf13/cobra"
	"os"
	"strings"
)
import "github.com/spf13/viper"

func main() {
	var l sdkTypes.KairosLogger
	var err error

	l = sdkTypes.NewKairosLogger("kairos-init", "info", false)
	l.Info("Initializing system as a Kairos system.")
	c := cobra.Command{
		Use:   "kairos-init",
		Short: "Initialize the system as a Kairos system",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(viper.GetStringSlice("features")) == 0 {
				return fmt.Errorf("no features specified")
			}
			for _, feature := range viper.GetStringSlice("features") {
				if feature == "all" {
					continue
				}
				if !features.FeatureSupported(feature) {
					return fmt.Errorf("feature %s not supported. Available features: %s", feature, strings.Join(features.FeatSupported(), ", "))
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Override logger
			l = sdkTypes.NewKairosLogger("kairos-init", viper.GetString("loglevel"), false)
			s := system.DetectSystem(l)

			if len(viper.GetStringSlice("features")) == 1 && viper.GetStringSlice("features")[0] == "all" {
				l.Logger.Info().Msg("Adding all features to queue")
				for _, f := range features.FeatSupported() {
					s.AddFeature(features.GetFeature(f))
				}
			} else {
				for _, f := range viper.GetStringSlice("features") {
					l.Logger.Info().Str("feature", f).Msg("Adding feature to queue")
					s.AddFeature(features.GetFeature(f))
				}
			}

			return s.ApplyFeatures(l)
		},
	}

	c.Flags().StringArrayP("features", "f", []string{}, fmt.Sprintf("Features to install. Available features: %s", strings.Join(features.FeatSupported(), ", ")))
	err = viper.BindEnv("features", "KAIROS_INIT_FEATURES")
	if err != nil {
		l.Logger.Err(err).Msg("Error binding environment variable")
		return
	}
	c.Flags().BoolP("dry-run", "d", false, "Dry run")
	err = viper.BindEnv("features", "KAIROS_INIT_DRY_RUN")
	if err != nil {
		l.Logger.Err(err).Msg("Error binding environment variable")
		return
	}
	// Global flag
	c.PersistentFlags().StringP("loglevel", "l", "info", "Log level")
	err = viper.BindEnv("loglevel", "KAIROS_INIT_LOGLEVEL")
	if err != nil {
		l.Logger.Err(err).Msg("Error binding environment variable")
		return
	}

	// Define the subcommand
	subCmd := &cobra.Command{
		Use:                   "show FEATURE",
		Short:                 "Show info about a feature",
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no feature specified. Available features: %s", strings.Join(features.FeatSupported(), ", "))
			}

			if !features.FeatureSupported(args[0]) {
				return fmt.Errorf("feature %s not supported. Available features: %s", args[0], strings.Join(features.FeatSupported(), ", "))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			l = sdkTypes.NewKairosLogger("kairos-init", viper.GetString("loglevel"), false)
			s := system.DetectSystem(l)
			l.Logger.Info().Str("feature", args[0]).Msg("Getting feature")
			f := s.GetFeature(args[0], l)
			if f == nil {
				l.Logger.Err(fmt.Errorf("feature %s not found", args[0])).Msg("Error")
				return fmt.Errorf("feature %s not found", args[0])
			}
			f.Info(s, l)
			return nil
		},
	}

	// Add the subcommand to the parent command
	c.AddCommand(subCmd)

	// Bind persistent flag especifically
	_ = viper.BindPFlag("loglevel", c.PersistentFlags().Lookup("loglevel"))
	err = viper.BindPFlags(c.Flags())

	if err != nil {
		l.Logger.Err(err).Msg("Error binding flags")
		return
	}
	err = c.Execute()
	if err != nil {
		l.Logger.Err(err).Msg("Error executing command")
		os.Exit(1)
	}
	l.Logger.Info().Msg("Done")
}
