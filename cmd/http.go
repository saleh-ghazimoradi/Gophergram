/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Launching the app via http",

	Run: func(cmd *cobra.Command, args []string) {
		if err := gateway.Server(); err != nil {
			logger.Logger.Error(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
