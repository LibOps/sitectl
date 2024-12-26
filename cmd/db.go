/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/libops/homebrew-cli/pkg/gcloud"
	"github.com/libops/homebrew-cli/pkg/libops"
	"github.com/spf13/cobra"
)

// importDbCmd represents the drupal command
var importDbCmd = &cobra.Command{
	Use:   "db",
	Short: "View basic information about your LibOps Drupal deployment.",
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

		err = libops.WaitUntilOnline(site, env, token)
		if err != nil {
			log.Fatal("Environment not responding")
		}
		f, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Open(f)
		if err != nil {
			log.Fatalf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)
		fileField, err := writer.CreateFormFile("sql", f)
		if err != nil {
			log.Fatalf("Error creating form field: %v\n", err)
			return
		}
		_, err = io.Copy(fileField, file)
		if err != nil {
			fmt.Printf("Error copying file content: %v\n", err)
			return
		}
		writer.Close()

		url, err := gcloud.GetCloudRunUrl(site, env)
		if err != nil {
			log.Fatal("Unable to retrieve remote URL")
		}
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/import/db", url), &requestBody)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Failed importing database: %v", resp.StatusCode)
		}

		fmt.Println("Successfully imported database!")
	},
}

func init() {
	importCmd.AddCommand(importDbCmd)
	importDbCmd.Flags().StringP("file", "f", "", "The database file to import")
	importDbCmd.Flags().StringP("token", "t", "", "(optional/machines-only) The gcloud identity token to access your LibOps environment")

	importDbCmd.MarkFlagRequired("file")
}
