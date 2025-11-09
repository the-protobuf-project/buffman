// package cmd implements the root command for the Buffman CLI application.
// It sets up the main command structure using Cobra, defines global flags,
// and adds subcommands for specific functionalities like 'convert' and 'generate'.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/machanirobotics/buffman/cmd/completion"
	"github.com/machanirobotics/buffman/cmd/convert"
	"github.com/machanirobotics/buffman/cmd/generate"
	"github.com/machanirobotics/buffman/internal/install"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "buffman",
	Short: "A powerful CLI tool for managing and converting buffer schemas.",
	Long: `Buffman is a versatile command-line tool designed to streamline
working with different buffer schemas. It simplifies converting Protocol Buffer
(.proto) files into other formats like FlatBuffers (.fbs), with more formats
planned for the future.`,
}

// Execute is the primary entry point for the CLI application. It is called by
// main.main() and is responsible for checking dependencies and executing the
// root command.
func Execute() {
	// Ensure critical dependencies like the FlatBuffers compiler are installed.
	ensureFlatcInstalled()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// init adds all child commands to the root command.
func init() {
	rootCmd.AddCommand(convert.ConvertCmd, generate.GenerateCmd, completion.CompletionCmd)
}

// ensureFlatcInstalled checks if the FlatBuffers compiler (flatc) is available.
// If flatc is not found, it prompts the user for permission to install it. The
// application will exit if the user declines or if the installation fails, as
// flatc is a required dependency.
func ensureFlatcInstalled() {
	installer := install.NewInstaller(install.FlatbuffersInstaller)
	if installer.Exists() {
		return // Dependency is already satisfied.
	}

	fmt.Print("Required dependency 'flatc' is missing. Would you like to install it? (y/n): ")

	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Normalize the response to handle variations like "Y", "yes", or " Yes ".
	response = strings.ToLower(strings.TrimSpace(response))

	if response == "y" || response == "yes" {
		fmt.Println("Installing flatc...")
		if err := installer.Install(); err != nil {
			fmt.Printf("Failed to install flatc. Please try installing it manually: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("flatc installed successfully!")
	} else {
		fmt.Println("Installation cancelled. 'flatc' is required for Buffman to work. Exiting.")
		os.Exit(1)
	}
}
