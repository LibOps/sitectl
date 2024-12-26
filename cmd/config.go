/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/libops/homebrew-cli/pkg/gcloud"
	"github.com/libops/homebrew-cli/pkg/libops"
	"github.com/spf13/cobra"
)

// getConfigCmd represents the drupal command
var getConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Export your drupal config and download locally.",
	Run: func(cmd *cobra.Command, args []string) {

		// make sure this is being ran from the root directory of the site
		if _, err := os.Stat("config"); os.IsNotExist(err) {
			log.Fatal("config directory does not exist.\nThis command needs ran from your code directory.")
			return
		}

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

		url, err := gcloud.GetCloudRunUrl(site, env)
		if err != nil {
			log.Fatal("Unable to retrieve remote URL")
		}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/cex", url), nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unable to get config: %v", resp.StatusCode)
		}

		// Remove the config directory and its contents
		if err := os.RemoveAll("config"); err != nil {
			log.Fatal("Could not remove config to overwrite with new content.\nMake sure you have permission to write to config.")
		}

		// untar the config
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Printf("Error creating gzip reader: %v\n", err)
			return
		}
		defer gzipReader.Close()
		tarReader := tar.NewReader(gzipReader)

		// Loop through the files in the tar archive and extract them to the current directory
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Printf("Error reading tar header: %v\n", err)
				return
			}

			extractPath := filepath.Join(".", header.Name)
			if header.FileInfo().IsDir() {
				err := os.MkdirAll(extractPath, os.ModePerm)
				if err != nil {
					fmt.Printf("Error creating directory: %v\n", err)
					return
				}
				continue
			}

			outFile, err := os.Create(extractPath)
			if err != nil {
				fmt.Printf("Error creating output file: %v\n", err)
				return
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tarReader)
			if err != nil {
				fmt.Printf("Error copying file content: %v\n", err)
				return
			}
		}
	},
}

func init() {
	getCmd.AddCommand(getConfigCmd)
	getConfigCmd.Flags().StringP("token", "t", "", "(optional/machines-only) The gcloud identity token to access your LibOps environment")
}
