package main

import (
	"fmt"
	"github.com/kairos-io/kairos-init/pkg/features"
	. "github.com/kairos-io/kairos-init/pkg/log"
	"github.com/kairos-io/kairos-init/pkg/system"
	"github.com/kairos-io/kairos-init/pkg/validator"
	"github.com/spf13/cobra"
	"os"
	"strings"
)
import "github.com/spf13/viper"

func main() {
	var err error

	Log.Info("Initializing system as a Kairos system.")
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
			// Override logger if the level has changed
			Log.SetLevel(viper.GetString("loglevel"))

			s := system.DetectSystem(Log)

			if len(viper.GetStringSlice("features")) == 1 && viper.GetStringSlice("features")[0] == "all" {
				Log.Logger.Info().Msg("Adding all features to queue")
				for _, f := range features.FeatSupported() {
					s.AddFeature(features.GetFeature(f))
				}
			} else {
				for _, f := range viper.GetStringSlice("features") {
					Log.Logger.Info().Str("feature", f).Msg("Adding feature to queue")
					s.AddFeature(features.GetFeature(f))
				}
			}

			err = s.ApplyFeatures(Log)
			if err != nil {
				Log.Logger.Err(err).Msg("Error applying features")
				return err
			}
			err = validator.Validate()
			return err
		},
	}

	c.Flags().StringArrayP("features", "f", []string{}, fmt.Sprintf("Features to install. Available features: %s", strings.Join(features.FeatSupported(), ", ")))
	err = viper.BindEnv("features", "KAIROS_INIT_FEATURES")
	if err != nil {
		Log.Logger.Err(err).Msg("Error binding environment variable")
		return
	}
	c.Flags().BoolP("dry-run", "d", false, "Dry run")
	err = viper.BindEnv("features", "KAIROS_INIT_DRY_RUN")
	if err != nil {
		Log.Logger.Err(err).Msg("Error binding environment variable")
		return
	}
	// Global flag
	c.PersistentFlags().StringP("loglevel", "l", "info", "Log level")
	err = viper.BindEnv("loglevel", "KAIROS_INIT_LOGLEVEL")
	if err != nil {
		Log.Logger.Err(err).Msg("Error binding environment variable")
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
			s := system.DetectSystem(Log)
			Log.Logger.Info().Str("feature", args[0]).Msg("Getting feature")
			f := s.GetFeature(args[0], Log)
			if f == nil {
				Log.Logger.Err(fmt.Errorf("feature %s not found", args[0])).Msg("Error")
				return fmt.Errorf("feature %s not found", args[0])
			}
			f.Info(s, Log)
			return nil
		},
	}

	// Add the subcommand to the parent command
	c.AddCommand(subCmd)

	validatorCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the system",
		Args:  cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return validator.Validate()
		},
	}
	c.AddCommand(validatorCmd)

	// Bind persistent flag especifically
	_ = viper.BindPFlag("loglevel", c.PersistentFlags().Lookup("loglevel"))
	err = viper.BindPFlags(c.Flags())

	if err != nil {
		Log.Logger.Err(err).Msg("Error binding flags")
		return
	}
	err = c.Execute()
	if err != nil {
		os.Exit(1)
	}
	Log.Logger.Info().Msg("Done")
}
