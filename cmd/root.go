package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "norfrig-monitor",
	Short: "Monitor del pipeline de BI de Norfrig",
	Long:  "CLI para monitorear el estado del stock y el pipeline de datos de Norfrig SRL.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
