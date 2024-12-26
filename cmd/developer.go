/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"

	"github.com/libops/homebrew-cli/internal/helpers"
	"github.com/libops/homebrew-cli/pkg/gcloud"
	"github.com/libops/homebrew-cli/pkg/libops"
	"github.com/spf13/cobra"
)

// setDeveloperCmd represents the drupal command
var setDeveloperCmd = &cobra.Command{
	Use:   "developer",
	Short: "Add a developer's information to libops.yml.",
	Run: func(cmd *cobra.Command, args []string) {
		// don't want to accidentally overwrite any changes
		// so make sure the git directory is clean
		a := []string{
			"diff",
			"--exit-code",
		}
		c := exec.Command("git", a...)
		if err := c.Run(); err != nil {
			log.Fatal("git directory is not clean. Check git status output.")
		}

		// unmarshal libops.yml
		yamlFile, err := os.ReadFile("libops.yml")
		if err != nil {
			log.Fatal("Error reading libops.yml")
		}
		var yml libops.Yml
		err = yaml.Unmarshal(yamlFile, &yml)
		if err != nil {
			log.Fatalf("Error loading libops.yml")
		}

		// configure the firewall with the IP(s) passed
		ips, err := cmd.Flags().GetStringSlice("ip")
		if err != nil {
			log.Fatal("Could not fetch IP")
		}
		for _, ip := range ips {
			if !helpers.Contains(ip, yml.HttpsFirewall) {
				yml.HttpsFirewall = append(yml.HttpsFirewall, ip)
			}

			if !helpers.Contains(ip, yml.SshFirewall) {
				yml.SshFirewall = append(yml.SshFirewall, ip)
			}
		}

		// get the SSH public key contents, skipping if the flag was passed
		skipPubKey, err := cmd.Flags().GetBool("skip-pub-key")
		pubKey := ""
		if skipPubKey == false {
			sshPath, err := cmd.Flags().GetString("pub-key")
			if err != nil {
				log.Fatal("Unable to read the pub-key flag.")
			}
			if _, err := os.Stat(sshPath); os.IsNotExist(err) {
				log.Println("Could not find your SSH public key. Pass the full path as --pub-key")
			}
			pubKeyContents, err := os.ReadFile(sshPath)
			if err != nil {
				log.Println("Unable to find public key. Not setting key value.")
			}
			pubKey = strings.TrimSpace(string(pubKeyContents))
		}

		// add the developer to the project
		email, err := cmd.Flags().GetString("google-account")
		if email == "" {
			log.Fatal("Unable to read email address to add. Are you authenticated to gcloud?")
		}
		found := false
		for k, _ := range yml.Developers {
			if k == email {
				found = true
				if !helpers.Contains(pubKey, yml.Developers[k]) {
					if pubKey != "" {
						yml.Developers[k] = append(yml.Developers[email], pubKey)
					}
				}
				break
			}
		}
		if !found {
			yml.Developers[email] = []string{}
			if pubKey != "" {
				yml.Developers[email] = append(yml.Developers[email], pubKey)
			}
		}

		// write the changes back to libops.yml
		yamlData, err := yaml.Marshal(&yml)
		if err != nil {
			log.Fatal("Error re-loading libpps.yml")
		}
		err = os.WriteFile("libops.yml", yamlData, 0644)
		if err != nil {
			log.Fatal("Error writing to libops.yml")
		}

		// show the changes
		a = []string{
			"diff",
			"--exit-code",
		}
		c = exec.Command("git", a...)
		output, err := c.Output()
		if err != nil {
			fmt.Println("Your libops.yml has been updated.\nCommit the changes for the new settings to take effect\n%s", string(output))
		} else {
			fmt.Println("No changes have been made to libops.yml. Is this developer already configured?")
		}
	},
}

func init() {
	setCmd.AddCommand(setDeveloperCmd)

	sshPath := ""
	currentUser, err := user.Current()
	if err == nil {
		sshPath = filepath.Join(currentUser.HomeDir, ".ssh", "id_rsa.pub")
	}
	setDeveloperCmd.Flags().StringP("pub-key", "k", sshPath, "Full path to your SSH public key file")
	setDeveloperCmd.Flags().BoolP("skip-pub-key", "s", false, "Skip loading/adding the SSH public key for the developer. Use this in the event you're adding on behalf of someone else and do not have their SSH public key on disk.")

	email, err := gcloud.GetEmail()
	if err != nil {
		email = ""
	}
	setDeveloperCmd.Flags().StringP("google-account", "g", email, "Google Account email address")

	ips := []string{}
	ip, err := helpers.GetIp()
	if err == nil {
		ips = []string{
			fmt.Sprintf("%s/32", ip),
		}
	}
	setDeveloperCmd.Flags().StringSliceP("ip", "i", ips, "IP Address(es) to add the the SSH and HTTPS firewalls to allow access")
}
