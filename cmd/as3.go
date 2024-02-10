package main

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/as3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	AS3PackageURLFlag           = "as3-package-url"
	AS3PackageSHAFlag           = "as3-package-checksum"
	AS3LivezHealthCheckFlag     = "as3-livez-health-check"
	AS3LivezHealthCheckPortFlag = "as3-livez-health-check-port"
)

func newAS3Cmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use: "as3 [scenario-name]",
		Aliases: []string{
			"as",
			"application-services",
			"application-services3",
			"application-services-3",
		},
		Short: "Generate an AS3 declaration with static entries.",
		Args:  bigipNameValidator,
		RunE:  as3Main,
	}
	addAS3Flags(cmd)
	if err := viper.BindPFlags(cmd.LocalFlags()); err != nil {
		return nil, fmt.Errorf("failed to bind local pflags to AS3: %w", err)
	}
	return cmd, nil
}

func addAS3Flags(cmd *cobra.Command) {
	cmd.LocalFlags().String(AS3PackageURLFlag, as3.DefaultPackageURL, "The URL to download AS3 package")
	cmd.LocalFlags().String(AS3PackageSHAFlag, as3.DefaultPackageSHA, "The checksum of the AS3 package to verify the download")
	cmd.LocalFlags().Bool(AS3LivezHealthCheckFlag, true, "Add a virtual service that responds to GCP Health Checks for liveness")
	cmd.LocalFlags().Int(AS3LivezHealthCheckPortFlag, 26000, "The TCP Port to bind for liveness probes")
}

func as3Main(_ *cobra.Command, args []string) error {
	name := args[0]
	context := as3ContextFromFlags(name)
	context.Header = newHeader(name, as3.DefaultSchemaURL)
	contextPreparer, err := context.NewContextPreparer()
	if err != nil {
		return fmt.Errorf("failed to create an AS3 ContextPreparer: %w", err)
	}
	generator, err := generators.NewGenerator(
		generators.WithTemplatePath(as3.TemplatePath),
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

// Returns an AS3 context updated from viper.
func as3ContextFromFlags(name string) *as3.Context {
	context := as3.NewDefaultContext()
	if viper.GetString(URLFlagName) != "" {
		context.ID = uuid.NewSHA1(uuid.NameSpaceURL, []byte(viper.GetString(URLFlagName))).String()
	}
	context.Label = name
	context.PackageURL = viper.GetString(AS3PackageURLFlag)
	context.PackageSHA = viper.GetString(AS3PackageSHAFlag)
	context.LivenessHealthCheck = viper.GetBool(AS3LivezHealthCheckFlag)
	context.LivenessHealthCheckPort = strconv.Itoa(viper.GetInt(context.LivenessHealthCheckPort))
	return context
}
