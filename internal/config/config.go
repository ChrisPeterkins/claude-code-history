package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Config holds user preferences for claude-history.
// Stored at ~/.claude-history.json.
type Config struct {
	// ProjectRoots are directory names that contain projects.
	// Paths after these segments become the project display name.
	// Example: ["Projects", "code", "work/clients"]
	ProjectRoots []string `json:"projectRoots,omitempty"`

	// Theme is the name of the color theme to use on startup.
	// One of: nord, dracula, catppuccin, light
	Theme string `json:"theme,omitempty"`

	// DefaultFilter is the session filter applied on startup.
	// One of: all, code, long, recent
	DefaultFilter string `json:"defaultFilter,omitempty"`
}

// DefaultProjectRoots are used when no config file exists or projectRoots is empty.
var DefaultProjectRoots = []string{
	"Projects", "projects",
	"Code", "code",
	"Dev", "dev",
	"src", "repos",
	"Workspace", "workspace",
}

var configPath string
var current Config

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	configPath = filepath.Join(home, ".claude-history.json")
	current = load()
}

// Get returns the current configuration.
func Get() Config {
	return current
}

// ProjectRoots returns the effective project root names, expanding any
// multi-segment roots (like "work/clients") into individual segments to match.
func ProjectRoots() map[string]bool {
	roots := current.ProjectRoots
	if len(roots) == 0 {
		roots = DefaultProjectRoots
	}

	result := make(map[string]bool)
	for _, root := range roots {
		// Support multi-segment roots like "work/clients"
		// by matching on the last segment
		parts := strings.Split(root, "/")
		result[parts[len(parts)-1]] = true
	}
	return result
}

// Save writes the current config to disk.
func Save(c Config) error {
	current = c
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// DefaultFilterName returns the configured default filter, or "all".
func DefaultFilterName() string {
	if current.DefaultFilter != "" {
		return current.DefaultFilter
	}
	return "all"
}

func load() Config {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{} // no config file is fine
	}
	var c Config
	json.Unmarshal(data, &c) // ignore parse errors, use defaults
	return c
}
