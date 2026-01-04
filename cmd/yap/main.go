package main

import (
	"fmt"
	"os"

	"github.com/rlamalama/YAP/cmd/yap/commands"
	"github.com/spf13/cobra"
)

func main() {
	// 1. Root Command (the main binary)
	var rootCmd = &cobra.Command{
		Use:   "YAP",
		Short: "YAP is the YAML to Programming CLI",
		Long:  `YAP is the YAML to Programming CLI`,
	}

	// 2. Subcommand (e.g., 'hello')
	var runCmd = &cobra.Command{
		Use:   "run [file]",
		Short: "Runs a particular .YAP file",
		Args:  cobra.MinimumNArgs(1), // Requires at least one argument
		Run: func(cmd *cobra.Command, args []string) {
			commands.RunCmd(args)
		},
	}

	// // 3. Define Flags (e.g., '--loud' or '-l')
	// runCmd.Flags().BoolP("loud", "l", false, "Shout the greeting")

	// 4. Add subcommand to root
	rootCmd.AddCommand(runCmd)

	// 5. Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
