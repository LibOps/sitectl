package gcloud

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func GetCloudRunUrl(p, env string) (string, error) {
	if err := godotenv.Load(); err != nil {
		if err = os.Setenv("LIBOPS_REGION", "us-central1"); err != nil {
                    log.Fatal("Error loading .env file")
                }
	}
	region := os.Getenv("LIBOPS_REGION")

	args := []string{
		"run",
		"services",
		"describe",
		fmt.Sprintf("remote-%s", env),
		fmt.Sprintf("--region=%s", region),
		fmt.Sprintf("--project=%s", p),
		"--format=value(status.url)",
	}
	cmd := exec.Command("gcloud", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		log.Println("Failed to location remote service")
		return "", err
	}

	url := stdout.String()
	return strings.TrimSpace(url), nil
}
