package helpers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func Contains(s string, slice []string) bool {
	found := false
	for _, e := range slice {
		if e == s {
			found = true
			break
		}
	}

	return found
}

func GetIp() (string, error) {
	ipServiceURL := "https://ifconfig.me"
	resp, err := http.Get(ipServiceURL)
	if err != nil {
		return "", fmt.Errorf("Error making HTTP request: %v\n", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Unable to get IP: %v\n", err)
	}

	return string(body), nil
}

func FindHostBlocks(config []byte, targetHost string) ([]string, error) {
	pattern := `(Host .*libops.*\n[\s\S]*?)(?:\n{2,}|\z)`
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(string(config), -1)
	var blocks []string

	for _, match := range matches {
		if !containsHostLine(match[1], targetHost) {
			blocks = append(blocks, match[1])
		}
	}

	return blocks, nil
}

func containsHostLine(input string, targetHost string) bool {
	lines := strings.Split(input, "\n")
	if len(lines) > 0 && lines[0] == fmt.Sprintf("Host %s", targetHost) {
		return true
	}

	return false
}
