package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	AppName            = "f5-google-declaration-generator"
	InterfacesFlagName = "interfaces"
	URLFlagName        = "url"
	VerboseFlagName    = "verbose"
)

// Version is updated from git tags during build.
var version = "snapshot"

// Determine the outcome of command line flags, environment variables, and an
// optional configuration file to perform initialization of the application.
func initConfig() {
	logLevel := slog.LevelVar{}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     &logLevel,
	})))
	slog.Debug("Initializing config")
	viper.AddConfigPath(".")
	if home, err := homedir.Dir(); err == nil {
		viper.AddConfigPath(home)
	}
	viper.SetConfigName("." + AppName)
	viper.SetEnvPrefix(AppName)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		slog.Warn("Error reading configuration file", "error", err)
	}
	verbosity := viper.GetInt(VerboseFlagName)
	switch {
	case verbosity == 1:
		logLevel.Set(slog.LevelInfo)
	case verbosity > 1:
		logLevel.Set(slog.LevelDebug)
	default:
		logLevel.Set(slog.LevelWarn)
	}
}

func NewRootCmd() (*cobra.Command, error) {
	slog.Debug("Building new root command")
	cobra.OnInitialize(initConfig)
	rootCmd := &cobra.Command{
		Use:     AppName,
		Version: version,
		Short:   "Calculate and retrieve a fractional digit of pi at an arbitrary index",
		Long:    `Provides a gRPC client/server demo for distributed calculation of fractional digits of pi.`,
	}
	rootCmd.PersistentFlags().Count(VerboseFlagName, "Increase the verbosity of logging, can be specified multiple times")
	rootCmd.PersistentFlags().Int(InterfacesFlagName, 3, "The number of interfaces attached to VM")
	rootCmd.PersistentFlags().String(URLFlagName, "", "The URL to use for generating UUID values")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		return nil, fmt.Errorf("failed to bind root pflags: %w", err)
	}
	bigipCmd, err := newBigIPCmd()
	if err != nil {
		return nil, err
	}
	rootCmd.AddCommand(bigipCmd)
	return rootCmd, nil
}

type headerImpl struct {
	name        func() string
	description func() string
	version     func() string
	timestamp   func() time.Time
	schemaURL   func() string
}

func (h headerImpl) Name() string {
	return h.name()
}

func (h headerImpl) Description() string {
	return h.description()
}

func (h headerImpl) Version() string {
	return h.version()
}

func (h headerImpl) Timestamp() time.Time {
	return h.timestamp()
}

func (h headerImpl) SchemaURL() string {
	return h.schemaURL()
}

// Helper function to create an implementation of generator.Header interface that
// can be referenced in the the rendered output.
func newHeader(name /*description, */, schemaURL string) generators.Header {
	ts := time.Now()
	return headerImpl{
		name:        func() string { return name },
		description: func() string { return "Foo" },
		version:     func() string { return AppName + " " + version },
		timestamp:   func() time.Time { return ts },
		schemaURL:   func() string { return schemaURL },
	}
}
