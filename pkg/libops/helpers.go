package libops

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/libops/homebrew-cli/pkg/gcloud"
	"github.com/spf13/cobra"
)

func LoadEnvironment(cmd *cobra.Command) (string, string, error) {
	site, err := cmd.Flags().GetString("site")
	if err != nil {
		return "", "", err
	}
	env, err := cmd.Flags().GetString("environment")
	if err != nil {
		return site, "", err
	}

	return site, env, nil
}

func IssueCommand(site, env, cmd, args, token string) error {
	var err error
	err = WaitUntilOnline(site, env, token)
	if err != nil {
		return err
	}

	log.Printf("Running `%s %s` on %s %s\n", cmd, args, site, env)
	url, err := gcloud.GetCloudRunUrl(site, env)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", url, cmd), bytes.NewBuffer([]byte(args)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("%s on %s %s returned a non-200: %v", cmd, site, env, resp.StatusCode)
	}
	// print the output to the terminal as it streams in
	for {
		buffer := make([]byte, 1024)
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		fmt.Print(string(buffer[:n]))
	}

	return nil
}

func GetToken(cmd *cobra.Command, tokenArg string) (string, error) {
	token, err := cmd.Flags().GetString(tokenArg)
	if err != nil {
		return "", err
	}
	if token == "" {
		token, err = gcloud.AccessToken()
		if err != nil {
			log.Println("Unable to run `gcloud auth print-identity-token`. Ensure you've ran `gcloud auth login`.")
			return "", err
		}
	}

	return token, nil
}

func WaitUntilOnline(site, env, token string) error {
	var err error
	timeout := 3 * time.Minute
	url, err := gcloud.GetCloudRunUrl(site, env)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ping/", url), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(5 * time.Second) {
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			log.Println("Waiting 10 seconds before trying again.")
			time.Sleep(10 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		log.Printf("Received status code %d, retrying...\n", resp.StatusCode)
	}
	log.Println("Timeout exceeded")
	return fmt.Errorf("%s %s not ready after one minute", site, env)
}
