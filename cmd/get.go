/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display information about your LibOps environment.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("You must specify the type of resource to get")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
