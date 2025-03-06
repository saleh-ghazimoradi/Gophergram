package cmd

import (
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Launching the app via http",
	Run: func(cmd *cobra.Command, args []string) {
		if err := gateway.Server(); err != nil {
			sLogger.SLogger.Error("Failed to start server", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
