package main

import (
	"fmt"
	"log/slog"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/ts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	TSPackageURLFlag = "ts-package-url"
	TSPackageSHAFlag = "ts-package-checksum"
)

func newTSCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:  "ts [scenario-name]",
		Args: bigipNameValidator,
		RunE: tsMain,
	}
	addTSFlags(cmd)
	if err := viper.BindPFlags(cmd.LocalFlags()); err != nil {
		return nil, fmt.Errorf("failed to bind local pflags to TS: %w", err)
	}
	return cmd, nil
}

func addTSFlags(cmd *cobra.Command) {
	cmd.LocalFlags().String(TSPackageURLFlag, ts.DefaultPackageURL, "The URL to download TS package")
	cmd.LocalFlags().String(TSPackageSHAFlag, ts.DefaultPackageSHA, "The checksum of the TS package to verify the download")
}

func tsMain(_ *cobra.Command, args []string) error {
	name := args[0]
	logger := slog.With("name", name)
	logger.Info("Preparing TS context")
	context := tsContextFromFlags(name)
	context.Header = newHeader(name, ts.DefaultSchemaURL)
	contextPreparer, err := context.NewContextPreparer()
	if err != nil {
		return fmt.Errorf("failed to create an TS ContextPreparer: %w", err)
	}
	logger.Info("Creating generator")
	generator, err := generators.NewGenerator(
		generators.WithTemplatePath(ts.TemplatePath),
		generators.WithContextPreparer(contextPreparer),
	)
	if err != nil {
		return fmt.Errorf("failed to create a generator: %w", err)
	}
	logger.Info("Executing generator")
	if err := generator.Execute(viper.GetInt(InterfacesFlagName)); err != nil {
		return fmt.Errorf("generator raised an error: %w", err)
	}
	return nil
}

// Returns a TS context updated from viper.
func tsContextFromFlags(_ string) *ts.Context {
	context := ts.NewDefaultContext()
	context.PackageURL = viper.GetString(TSPackageURLFlag)
	context.PackageSHA = viper.GetString(TSPackageSHAFlag)
	return context
}
