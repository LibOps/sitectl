/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/libops/homebrew-cli/pkg/libops"
	"github.com/spf13/cobra"
)

// backupCmd represents the gsutil command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup your libops environment",
	Long: `
Info:
    Backup your libops environment.

    Right now only database backups are performed.

    Examples:
    # Backup production environment
    libops backup -e production
`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		site, env, err := libops.LoadEnvironment(cmd)
		if err != nil {
			log.Println("Unable to load environment.")
			log.Fatal(err)
		}

		// get the gcloud id token
		token, err := libops.GetToken(cmd, "token")
		if err != nil {
			log.Fatal(err)
		}

		// export the db
		exportArgs := []string{
			"sql-dump",
			"-y",
			"--skip-tables-list=cache,cache_*",
			"--structure-tables-list=cache,cache_*",
			"--result-file=/tmp/drupal.sql",
			"--debug",
		}
		drushCmd := strings.Join(exportArgs, " ")
		err = libops.IssueCommand(site, env, "drush", drushCmd, token)
		if err != nil {
			log.Fatal(err)
		}

		now := time.Now().Format("2006/01/02")
		currentTime := time.Now().Format("15-04")
		fileName := fmt.Sprintf("drupal-%s.sql", currentTime)
		gcsObject := fmt.Sprintf("gs://%s-backups/%s/%s/%s", site, now, env, fileName)
		uploadArgs := []string{
			"cp",
			"/tmp/drupal.sql",
			gcsObject,
		}
		gsutilCmd := strings.Join(uploadArgs, " ")

		err = libops.IssueCommand(site, env, "gsutil", gsutilCmd, token)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringP("token", "t", "", "(optional/machines-only) The gcloud identity token to access your LibOps environment")
}
