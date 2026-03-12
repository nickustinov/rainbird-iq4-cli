package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".iq4")
}

func tokenPath() string {
	return filepath.Join(configDir(), "token")
}

// SaveToken stores the JWT token to ~/.iq4/token.
func SaveToken(token string) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	return os.WriteFile(tokenPath(), []byte(token), 0600)
}

// LoadToken reads the stored JWT token, or returns empty string if none.
func LoadToken() string {
	data, err := os.ReadFile(tokenPath())
	if err != nil {
		return ""
	}
	return string(data)
}

// ClearToken removes the stored token.
func ClearToken() error {
	return os.Remove(tokenPath())
}
