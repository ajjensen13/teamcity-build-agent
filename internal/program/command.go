package program

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	flagVar *string
)

func init() {
	rootCmd.AddCommand(commandCmd)
	flagVar = commandCmd.Flags().String("flag", "", "command flag")
}

var (
	commandCmd = &cobra.Command{
		Use:   "command",
		Short: "A brief description of your application",
		Long: "" +
			`A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

		Run: runCommand,
	}
)

func runCommand(_ *cobra.Command, _ []string) {
	fmt.Printf("command called with %s", *flagVar)
}
