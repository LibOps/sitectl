/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "libops",
	Short: "Interact with your libops site",
	Long:  `Interact with your libops site`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	rootCmd.PersistentFlags().StringP("site", "p", filepath.Base(path), "LibOps project/site")
	rootCmd.PersistentFlags().StringP("environment", "e", "development", "LibOps environment")
}
