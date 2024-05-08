package main

import (
	"fmt"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/cfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	CFEPackageURLFlag  = "cfe-package-url"
	CFEPackageSHAFlag  = "cfe-package-checksum"
	CFEScopingTagsFlag = "cfe-scoping-tags"
)

func newCFECmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:  "cfe [scenario-name]",
		Args: bigipNameValidator,
		RunE: cfeMain,
	}
	addCFEFlags(cmd)
	if err := viper.BindPFlags(cmd.LocalFlags()); err != nil {
		return nil, fmt.Errorf("failed to bind local pflags to CFE: %w", err)
	}
	return cmd, nil
}

func addCFEFlags(cmd *cobra.Command) {
	cmd.LocalFlags().String(CFEPackageURLFlag, cfe.DefaultPackageURL, "The URL to download CFE package")
	cmd.LocalFlags().String(CFEPackageSHAFlag, cfe.DefaultPackageSHA, "The checksum of the CFE package to verify the download")
	cmd.LocalFlags().StringToString(CFEScopingTagsFlag, map[string]string{}, "Defines the scoping tags (labels) to use in CFE")
}

func cfeMain(_ *cobra.Command, args []string) error {
	name := args[0]
	context := cfeContextFromFlags(name)
	context.Header = newHeader(name, cfe.DefaultSchemaURL)
	contextPreparer, err := context.NewContextPreparer()
	if err != nil {
		return fmt.Errorf("failed to create an CFE ContextPreparer: %w", err)
	}
	generator, err := generators.NewGenerator(
		generators.WithTemplatePath(cfe.TemplatePath),
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

// Returns a CFE context updated from viper.
func cfeContextFromFlags(_ string) *cfe.Context {
	context := cfe.NewDefaultContext()
	context.PackageURL = viper.GetString(CFEPackageURLFlag)
	context.PackageSHA = viper.GetString(CFEPackageSHAFlag)
	return context
}
