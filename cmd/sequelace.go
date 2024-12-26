/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/libops/homebrew-cli/pkg/gcloud"
	"github.com/libops/homebrew-cli/pkg/libops"
	"github.com/spf13/cobra"
)

// addsequelAceCmd represents the addconfigSsh command
var sequelAceCmd = &cobra.Command{
	Use:   "sequelace",
	Short: "Connect to your LibOps database using Sequel Ace (Mac OS only)",
	Long: `
Info:
	Running this command opens a connection to your LibOps environment's database.

    Database and SSH connection information will have a host, name, and port. Along with relevant credentials.

    Services exposed over HTTPS will have a URL and relevant credentials.

    Examples:
	# Connect to your production database
    libops sequelace -e production

    # Connect to your development database
    libops sequelace
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

		ci := libops.ConnectionInfo{}
		json.NewDecoder(resp.Body).Decode(&ci)

		sshKeyPath, err := cmd.Flags().GetString("ssh-priv-key")
		if err != nil {
			log.Fatal(err)
		}

		sequelAcePath, err := cmd.Flags().GetString("sequel-ace-path")
		if err != nil {
			log.Fatal(err)
		}

		mysqlInfo := fmt.Sprintf("mysql://%s:%s@%s:%d/%s", ci.Database.Credentials.Username, ci.Database.Credentials.Password, ci.Database.Host, ci.Database.Port, ci.Database.Name)
		sshInfo := fmt.Sprintf("ssh_host=%s&ssh_port=%d&ssh_user=%s&ssh_keyLocation=%s&ssh_keyLocationEnabled=1", ci.Ssh.Host, ci.Ssh.Port, ci.Ssh.Credentials.Username, sshKeyPath)
		cmdArgs := []string{
			fmt.Sprintf("%s?%s", mysqlInfo, sshInfo),
			"-a",
			sequelAcePath,
		}
		openCmd := exec.Command("open", cmdArgs...)
		if err := openCmd.Run(); err != nil {
			log.Fatal("Could not open sequelace.")
		}

	},
}

func init() {
	rootCmd.AddCommand(sequelAceCmd)

	sequelAceCmd.Flags().StringP("token", "t", "", "(optional/machines-only) The gcloud identity token to access your LibOps environment")
	sequelAceCmd.Flags().StringP("ssh-priv-key", "k", fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME")), "Full path to your SSH private key (default ~/.ssh/id_rsa)")
	sequelAceCmd.Flags().StringP("sequel-ace-path", "s", "/Applications/Sequel Ace.app/Contents/MacOS/Sequel Ace", "Full path to your Sequel Ace app (default /Applications/Sequel Ace.app/Contents/MacOS/Sequel Ace)")
}
