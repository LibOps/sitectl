/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/libops/homebrew-cli/pkg/libops"
	"github.com/spf13/cobra"
)

// syncDbCmd represents the syncDb command
var syncDbCmd = &cobra.Command{
	Use:   "sync-db",
	Short: "Transfer the database from one environment to another",
	Long: `
Info:
	Transfer the database from one environment to another

    Example sync the production database into development:
      libops sync-db --site libops-abcdef01 --source production --target development
`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		site, _, err := libops.LoadEnvironment(cmd)
		if err != nil {
			log.Println("Unable to load environment.")
			log.Fatal(err)
		}

		source, err := cmd.Flags().GetString("source")
		if err != nil {
			log.Fatal(err)
		}
		target, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatal(err)
		}

		sourceToken, err := libops.GetToken(cmd, "source-token")
		if err != nil {
			log.Fatal(err)
		}
		targetToken, err := libops.GetToken(cmd, "target-token")
		if err != nil {
			log.Fatal(err)
		}

		// run the drush command
		exportArgs := []string{
			"sql-dump",
			"-y",
			"--skip-tables-list=cache,cache_*,watchdog",
			"--structure-tables-list=cache,cache_*,watchdog",
			"--result-file=/tmp/drupal.sql",
			"--debug",
		}
		drushCmd := strings.Join(exportArgs, " ")
		err = libops.IssueCommand(site, source, "drush", drushCmd, sourceToken)
		if err != nil {
			log.Fatal(err)
		}

		rand := rand.Int()
		fileName := fmt.Sprintf("drupal-%s-%d.sql", source, rand)
		gcsObject := fmt.Sprintf("gs://%s-backups/%s", site, fileName)
		uploadArgs := []string{
			"cp",
			"/tmp/drupal.sql",
			gcsObject,
		}
		gsutilCmd := strings.Join(uploadArgs, " ")

		err = libops.IssueCommand(site, source, "gsutil", gsutilCmd, sourceToken)
		if err != nil {
			log.Fatal(err)
		}

		downloadArgs := []string{
			"cp",
			gcsObject,
			"/tmp/",
		}
		gsutilCmd = strings.Join(downloadArgs, " ")
		err = libops.IssueCommand(site, target, "gsutil", gsutilCmd, targetToken)
		if err != nil {
			log.Fatal(err)
		}

		importArgs := []string{
			"sql-query",
			"-y",
			"--file-delete",
			fmt.Sprintf("--file=/tmp/%s", fileName),
			"--debug",
		}
		drushCmd = strings.Join(importArgs, " ")
		err = libops.IssueCommand(site, target, "drush", drushCmd, targetToken)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncDbCmd)

	syncDbCmd.Flags().StringP("source", "s", "", "The database that will be exported from")
	syncDbCmd.Flags().StringP("target", "t", "", "The database that will be overwritten")
	syncDbCmd.Flags().StringP("source-token", "x", "", "(optional/machines-only) The gcloud identity token to access the LibOps environment passed as `source`")
	syncDbCmd.Flags().StringP("target-token", "y", "", "(optional/machines-only) The gcloud identity token to access the LibOps environmen passed as `target`")

	syncDbCmd.MarkFlagRequired("source")
	syncDbCmd.MarkFlagRequired("target")
}
