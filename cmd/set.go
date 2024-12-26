/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set information on your LibOps environment.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("You must specify the type of resource to set")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
