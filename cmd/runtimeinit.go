package main

import (
	"fmt"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/runtimeinit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	RuntimeInitPackageURLFlag = "runtime-init-package-url"
	RuntimeInitPackageSHAFlag = "runtime-init-package-checksum"
	AS3EnableFlag             = "as3-enable"
	DOEnableFlag              = "do-enable"
	CFEEnableFlag             = "cfe-enable"
	FASTEnableFlag            = "fast-enable"
	TSEnableFlag              = "ts-enable"
)

func newRuntimeInitCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:  "runtime-init [scenario-name]",
		Args: bigipNameValidator,
		RunE: runtimeInitMain,
	}
	addRuntimeInitFlags(cmd)
	if err := viper.BindPFlags(cmd.LocalFlags()); err != nil {
		return nil, fmt.Errorf("failed to bind local pflags to runtime-init: %w", err)
	}
	return cmd, nil
}

func addRuntimeInitFlags(cmd *cobra.Command) {
	cmd.LocalFlags().String(RuntimeInitPackageURLFlag, runtimeinit.DefaultPackageURL, "The URL to download runtime-init package")
	cmd.LocalFlags().String(RuntimeInitPackageSHAFlag, runtimeinit.DefaultPackageSHA, "The checksum of the runtime-init package to verify the download")
	cmd.LocalFlags().Bool(AS3EnableFlag, true, "Include support for AS3 Extension")
	cmd.LocalFlags().Bool(DOEnableFlag, true, "Include support for DO Extension")
	cmd.LocalFlags().Bool(CFEEnableFlag, true, "Include support for CFE Extension")
	cmd.LocalFlags().Bool(FASTEnableFlag, true, "Include support for FAST Extension")
	cmd.LocalFlags().Bool(TSEnableFlag, true, "Include support for TS Extension")
}

func runtimeInitMain(_ *cobra.Command, args []string) error {
	name := args[0]
	context := runtimeInitContextFromFlags(name)
	context.Header = newHeader(name, runtimeinit.DefaultSchemaURL)
	contextPreparer, err := context.NewContextPreparer()
	if err != nil {
		return fmt.Errorf("failed to create an runtime-init ContextPreparer: %w", err)
	}
	generator, err := generators.NewGenerator(
		generators.WithTemplatePath(runtimeinit.TemplatePath),
		generators.WithInterfaceBuilder(runtimeinit.InterfaceBuilder),
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

// Returns a runtime-init context updated from viper.
func runtimeInitContextFromFlags(name string) *runtimeinit.Context {
	context := runtimeinit.NewDefaultContext()
	context.PackageURL = viper.GetString(RuntimeInitPackageURLFlag)
	context.PackageSHA = viper.GetString(RuntimeInitPackageSHAFlag)
	if viper.GetBool(AS3EnableFlag) {
		context.ApplicationServices3 = as3ContextFromFlags(name)
	}
	if viper.GetBool(CFEEnableFlag) {
		context.CloudFailover = cfeContextFromFlags(name)
	}
	if viper.GetBool(DOEnableFlag) {
		doContext := doContextFromFlags(name)
		doContext.AdminPassword = context.AdminPassword
		if doContext.Licensing != nil && doContext.Licensing.LicensePool != nil {
			doContext.Licensing.LicensePool.BigIPPassword = context.AdminPassword
		}
		doContext.InstanceName = context.InstanceName
		context.DeclarativeOnboarding = doContext
	}
	if viper.GetBool(TSEnableFlag) {
		tsContext := tsContextFromFlags(name)
		tsContext.ProjectID = context.ProjectID
		tsContext.ServiceAccount = context.ServiceAccount
		context.TelemetryStreaming = tsContext
	}
	return context
}
