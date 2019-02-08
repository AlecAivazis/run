package main

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
)

func main() {
	// look for the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Sorry, there was a problem: %s.\n", err.Error())
		os.Exit(1)
	}

	// load the current configuration from the filesystem
	config, err := LoadConfig(afero.NewOsFs(), cwd)
	if err != nil {
		fmt.Printf("Sorry, there was a problem: %s.\n", err.Error())
		os.Exit(1)
	}

	// convert the configuration into a command we can execute
	cmd, err := config.Cmd()
	if err != nil {
		fmt.Printf("Sorry, there was a problem: %s.\n", err.Error())
		os.Exit(1)
	}

	// execute the command
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
