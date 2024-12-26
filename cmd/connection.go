/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/libops/homebrew-cli/pkg/gcloud"
	"github.com/libops/homebrew-cli/pkg/libops"
	"github.com/spf13/cobra"
)

// getConnectionInfoCmd represents the drupal command
var getConnectionInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get connection information, including secret credentials.",
	Long: `
Info:
	Get connection information for services deployed to your LibOps environment.

    Database and SSH connection information will have a host, name, and port. Along with relevant credentials.

    Services exposed over HTTPS will have a URL and relevant credentials.

    Examples:
	# get all the production connection information
    libops get info -e production

	# print the database connection information for your development environment
    libops get info | jq .database

	# print the URL to your development environment Drupal URL
    libops get info | jq -r .drupal.url
`,
	Run: func(cmd *cobra.Command, args []string) {
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
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/info", url), nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode > 299 {
			log.Fatalf("Unable to get environment info: %v", resp.StatusCode)
		}

		r := libops.ConnectionInfo{}
		json.NewDecoder(resp.Body).Decode(&r)
		b, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	},
}

func init() {
	getCmd.AddCommand(getConnectionInfoCmd)
	getConnectionInfoCmd.Flags().StringP("token", "t", "", "(optional/machines-only) The gcloud identity token to access your LibOps environment")
}
