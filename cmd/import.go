/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import resources to your LibOps environment.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("You must specify the type of resource to import")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
