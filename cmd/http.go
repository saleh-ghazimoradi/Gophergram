/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Launching the app via http",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("http called")
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
