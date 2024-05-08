package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newAppCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use: "app [scenario-name]",
		Aliases: []string{
			"application",
		},
		Short: "Generate an AS3 application declaration with static entries.",
		Args:  bigipNameValidator,
		RunE:  appMain,
	}
	addAppFlags(cmd)
	if err := viper.BindPFlags(cmd.LocalFlags()); err != nil {
		return nil, fmt.Errorf("failed to bind local pflags to App: %w", err)
	}
	return cmd, nil
}

func addAppFlags(_ *cobra.Command) {
}

func appMain(_ *cobra.Command, args []string) error {
	name := args[0]
	context := appContextFromFlags(name)
	context.Header = newHeader(name, app.DefaultSchemaURL)
	contextPreparer, err := context.NewContextPreparer()
	if err != nil {
		return fmt.Errorf("failed to create an App ContextPreparer: %w", err)
	}
	generator, err := generators.NewGenerator(
		generators.WithTemplatePath(app.TemplatePath),
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

// Returns an App context updated from viper.
func appContextFromFlags(name string) *app.Context {
	context := app.NewDefaultContext()
	if viper.GetString(URLFlagName) != "" {
		context.ID = uuid.NewSHA1(uuid.NameSpaceURL, []byte(viper.GetString(URLFlagName))).String()
	}
	context.Label = name
	return context
}
