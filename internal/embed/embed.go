// Package embed provides embedded assets for TAW.
package embed

import (
	"embed"
	"io/fs"
)

//go:embed assets/*
var Assets embed.FS

// GetPrompt returns the appropriate prompt content based on git mode.
func GetPrompt(isGitRepo bool) (string, error) {
	var filename string
	if isGitRepo {
		filename = "assets/PROMPT.md"
	} else {
		filename = "assets/PROMPT-nogit.md"
	}

	data, err := Assets.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetHelp returns the help content.
func GetHelp() (string, error) {
	data, err := Assets.ReadFile("assets/HELP.md")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetCommand returns the content of a slash command.
func GetCommand(name string) (string, error) {
	data, err := Assets.ReadFile("assets/commands/" + name + ".md")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ListCommands returns all available slash commands.
func ListCommands() ([]string, error) {
	entries, err := Assets.ReadDir("assets/commands")
	if err != nil {
		return nil, err
	}

	var commands []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Remove .md extension
		if len(name) > 3 && name[len(name)-3:] == ".md" {
			commands = append(commands, name[:len(name)-3])
		}
	}
	return commands, nil
}

// GetAsset returns the content of any embedded asset.
func GetAsset(path string) ([]byte, error) {
	return Assets.ReadFile(path)
}

// WalkAssets walks all embedded assets.
func WalkAssets(fn func(path string, d fs.DirEntry, err error) error) error {
	return fs.WalkDir(Assets, "assets", fn)
}
