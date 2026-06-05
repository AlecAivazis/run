package main

import (
	"fmt"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/spf13/afero"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Sorry, there was a problem: %s.\n", err.Error())
		os.Exit(1)
	}

	config, err := resolveConfig(afero.NewOsFs(), cwd)
	if err != nil {
		fmt.Printf("Sorry, there was a problem: %s.\n", err.Error())
		os.Exit(1)
	}

	if config.rootDir != "" {
		_ = godotenv.Load(path.Join(config.rootDir, ".env"))
	}

	cmd, err := config.Cmd()
	if err != nil {
		fmt.Printf("Sorry, there was a problem: %s.\n", err.Error())
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// resolveConfig loads the project config, with special handling for built-in
// commands that don't require a project (completion, help). For dynamic
// completion queries (__complete), config errors are silently ignored so the
// shell receives an empty completion list rather than an error.
func resolveConfig(fs afero.Fs, cwd string) (*Config, error) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "completion", "help", "-h", "--help":
			return &Config{}, nil
		case "__complete", "__completeNoDesc":
			config, err := LoadConfig(fs, cwd)
			if err != nil {
				return &Config{}, nil
			}
			return config, nil
		}
	}
	return LoadConfig(fs, cwd)
}
