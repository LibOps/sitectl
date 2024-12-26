package gcloud

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func AccessToken() (string, error) {
	args := []string{
		"auth",
		"print-identity-token",
	}
	cmd := exec.Command("gcloud", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	token := stdout.String()
	return strings.TrimSpace(token), nil
}

func GetEmail() (string, error) {

	args := []string{
		"auth",
		"list",
		"--filter",
		"status:ACTIVE",
		"--format",
		"value(account)",
	}
	cmd := exec.Command("gcloud", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Unable to get email from gcloud auth list. Are you authenticated to gcloud?")
	}

	return strings.TrimSpace(string(output)), nil
}
