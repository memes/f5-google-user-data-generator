package main

import (
	"fmt"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/do"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DOPackageURLFlag    = "do-package-url"
	DOPackageSHAFlag    = "do-package-checksum"
	LicenseRegKeyFlag   = "license-regkey"
	LicensePoolNameFlag = "license-pool-name"
)

func newDOCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:  "do [scenario-name]",
		Args: bigipNameValidator,
		RunE: doMain,
	}
	addDOFlags(cmd)
	if err := viper.BindPFlags(cmd.LocalFlags()); err != nil {
		return nil, fmt.Errorf("failed to bind local pflags to DO: %w", err)
	}
	return cmd, nil
}

func addDOFlags(cmd *cobra.Command) {
	cmd.LocalFlags().String(DOPackageURLFlag, do.DefaultPackageURL, "The URL to download DO package")
	cmd.LocalFlags().String(DOPackageSHAFlag, do.DefaultPackageSHA, "The checksum of the DO package to verify the download")
	cmd.LocalFlags().String(LicenseRegKeyFlag, do.DefaultRegKey, "Use a RegKey for licensing")
	cmd.LocalFlags().String(LicensePoolNameFlag, "", "The BIG-IQ license pool to use")
}

func doMain(_ *cobra.Command, args []string) error {
	name := args[0]
	context := doContextFromFlags(name)
	context.Header = newHeader(name, do.DefaultSchemaURL)
	contextPreparer, err := context.NewContextPreparer()
	if err != nil {
		return fmt.Errorf("failed to create an DO ContextPreparer: %w", err)
	}
	generator, err := generators.NewGenerator(
		generators.WithTemplatePath(do.TemplatePath),
		generators.WithContextPreparer(contextPreparer),
	)
	if err != nil {
		return fmt.Errorf("failed to create a generator: %w", err)
	}
	if err := generator.Execute(viper.GetInt(InterfacesFlagName)); err != nil {
		return fmt.Errorf("generator raised an error: %w", err)
	}
	return nil
}

// Returns a DO context updated from viper.
func doContextFromFlags(_ string) *do.Context {
	context := do.NewDefaultContext()
	context.PackageURL = viper.GetString(DOPackageURLFlag)
	context.PackageSHA = viper.GetString(DOPackageSHAFlag)
	return context
}
