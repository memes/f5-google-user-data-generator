// F5-google-declaration-generator is a utility to build a full-formed and self-referential onboarding declaration for
// various F5 products that can be deployed to Google Cloud as virtual machines. Each first-level subcommand encapsulates
// the generation of a declaration for an F5 product.
package main

import (
	"log/slog"
	"os"
)

func main() {
	rootCmd, err := NewRootCmd()
	if err != nil {
		slog.Error("Error building commands", "error", err)
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Error executing command", "error", err)
		os.Exit(1)
	}
}
