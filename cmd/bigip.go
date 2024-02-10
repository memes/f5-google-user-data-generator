package main

import (
	"fmt"
	"log/slog"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/as3"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/cloudinit"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/runtimeinit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Implements the bigip sub-command which generates classic BIG-IP declarations
// for use on Google Cloud.
func newBigIPCmd() (*cobra.Command, error) {
	slog.Debug("Creating BigIP command")
	bigipCmd := &cobra.Command{
		Use:   "bigip [scenario-name]",
		Short: "Run Continuity client to make requests to a Continuity Service endpoint",
		Long:  "Launches a gRPC client that will connect to Pi Service target and request the fractional digits of pi.",
		Args:  bigipNameValidator,
		RunE:  bigIPMain,
	}
	addAS3Flags(bigipCmd)
	addCFEFlags(bigipCmd)
	addDOFlags(bigipCmd)
	addRuntimeInitFlags(bigipCmd)
	addTSFlags(bigipCmd)
	addAppFlags(bigipCmd)
	as3Cmd, err := newAS3Cmd()
	if err != nil {
		return nil, err
	}
	cfeCmd, err := newCFECmd()
	if err != nil {
		return nil, err
	}
	doCmd, err := newDOCmd()
	if err != nil {
		return nil, err
	}
	riCmd, err := newRuntimeInitCmd()
	if err != nil {
		return nil, err
	}
	tsCmd, err := newTSCmd()
	if err != nil {
		return nil, err
	}
	appCmd, err := newAppCmd()
	if err != nil {
		return nil, err
	}
	bigipCmd.AddCommand(as3Cmd, cfeCmd, doCmd, riCmd, tsCmd, appCmd)
	return bigipCmd, nil
}

func bigipNameValidator(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}
	if err := as3.ValidateLabel(args[0]); err != nil {
		return fmt.Errorf("invalid name %q: %w", args[0], err)
	}
	return nil
}

func bigIPMain(_ *cobra.Command, args []string) error {
	name := args[0]
	context := cloudConfigContextFromFlags(name)
	context.Header = newHeader(name, cloudinit.DefaultSchemaURL)
	contextPreparer, err := context.NewContextPreparer()
	if err != nil {
		return fmt.Errorf("failed to create a cloud-config ContextPreparer: %w", err)
	}
	generator, err := generators.NewGenerator(
		generators.WithTemplatePath(cloudinit.TemplatePath),
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

// Returns a cloud-config context updated from viper.
func cloudConfigContextFromFlags(name string) *cloudinit.Context {
	context := cloudinit.NewDefaultContext()
	context.RuntimeInit = runtimeInitContextFromFlags(name)
	return context
}
