package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/hashicorp/hcl"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// TaskFileName is the name of the file that defines run's tasks
const TaskFileName = "_tasks.hcl"

// Executable is the internal interface for something that can act as the main entry point for run
type Executable interface {
	Execute() error
}

// Config holds the entire configuration for a project
type Config struct {
	Tasks     map[string]*Task `hcl:"task"`
	Variables map[string]interface{}
	Settings  struct {
		TemplateDelimiters []string `hcl:"delimiters"`
	} `hcl:"config"`
}

// Cmd turns the config into a
func (c *Config) Cmd() (Executable, error) {
	// the delimiters to use for our templates
	if len(c.Settings.TemplateDelimiters) == 0 {
		c.Settings.TemplateDelimiters = []string{"{{", "}}"}
	}

	// at the moment, run wraps over afero/cobra so let's create a root command
	cmd := &cobra.Command{
		Use: "run",
	}

	// each task in the config represents a command to cobra
	for taskName, task := range c.Tasks {
		// save the name of the task
		task.Name = taskName

		// the description of the task can have variables
		tmpl, err := template.New("task-description").Delims(c.Settings.TemplateDelimiters[0], c.Settings.TemplateDelimiters[1]).Parse(task.Description)
		if err != nil {
			return nil, err
		}
		var description bytes.Buffer
		tmpl.Execute(&description, c.Variables)

		// create a sub command to pass to cobra
		subCmd := func(task *Task) *cobra.Command {
			return &cobra.Command{
				Use:   taskName,
				Short: string(description.Bytes()),
				Run: func(cmd *cobra.Command, args []string) {
					if err := task.Run(args, c); err != nil {
						fmt.Printf("Sorry something went wrong: %s\n", err.Error())
						os.Exit(1)
						return
					}
				},
			}
		}(task)

		// add the sub command to the root
		cmd.AddCommand(subCmd)
	}

	// return the result
	return cmd, nil
}

// LoadConfig takes a filesystem and loads the appropriate configuration file by walking up from
// the current working directory until it finds a taskfile
func LoadConfig(fs afero.Fs, dir string) (*Config, error) {
	// the path of the task file relative to this location
	taskFilePath := path.Join(dir, TaskFileName)
	// if the current directory has the config file
	if _, err := fs.Stat(taskFilePath); err == nil {
		// read its contents
		contents, err := afero.ReadFile(fs, taskFilePath)
		if err != nil {
			return nil, err
		}

		// a place to hold the result
		result := &Config{}

		// parse the contents as hcl
		err = hcl.Decode(result, string(contents))
		if err != nil {
			return nil, err
		}

		// we're done here
		return result, nil
	}

	// if we get to the top of the filesystem
	if dir == "/" {
		return nil, errors.New("could not find task file")
	}

	// keep walking up
	return LoadConfig(fs, path.Join(dir, ".."))
}
