package cmd

import (
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "launching the http rest listen server",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("http rest server is starting")
		if err := gateway.Server(gateway.Routes()); err != nil {
			logger.Logger.Fatal(err)
		}
	},
}
