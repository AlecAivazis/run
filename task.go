package run

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/spf13/cobra"
)

// Executable is the internal interface for something that can act as the main entry point for run
type Executable interface {
	Execute() error
}

// Task represents a single task that can be executed by run
type Task struct {
	Description string
	Command     string
	Commands    []string
	Pipeline    []string
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
	delimiters := []string{"{{", "}}"}
	if len(c.Settings.TemplateDelimiters) != 0 {
		delimiters = c.Settings.TemplateDelimiters
	}

	// at the moment, run wraps over afero/cobra so let's create a root command
	cmd := &cobra.Command{
		Use: "run",
	}

	fmt.Println(c.Variables)
	// each task in the config represents a command to cobra
	for taskName, task := range c.Tasks {
		// the description of the task can have variables
		tmpl, err := template.New("task-description").Delims(delimiters[0], delimiters[1]).Parse(task.Description)
		if err != nil {
			return nil, err
		}
		var description bytes.Buffer
		tmpl.Execute(&description, c.Variables)

		// create a sub command to pass to cobra
		subCmd := &cobra.Command{
			Use:   taskName,
			Short: string(description.Bytes()),
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
			},
		}

		// add the sub command to the root
		cmd.AddCommand(subCmd)
	}

	// return the result
	return cmd, nil
}
