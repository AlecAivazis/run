package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	rootDir string
	fs      afero.Fs
}

// Cmd turns the config into an Executable
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
			return task.CobraCommand(c)
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
	// we also need to recnogize node projects so look for a package.json
	packageJSONPath := path.Join(dir, "package.json")
	// if we get to the top of the filesystem
	if dir == "/" {
		return nil, errors.New("could not find task file")
	}

	// if we don't see either in this directory. an error indicates it doesn't exist
	_, taskFileStatErr := fs.Stat(taskFilePath)
	taskFileExists := taskFileStatErr == nil
	_, packageStatErr := fs.Stat(packageJSONPath)
	packageJSONExists := packageStatErr == nil
	// if neither exist
	if !taskFileExists && !packageJSONExists {
		// we have to keep looking up
		return LoadConfig(fs, path.Join(dir, ".."))
	}

	// if we got this far, we are at the root of a run project.
	// there is at least one of:
	//    - _tasks.hcl
	//    - package.json

	// a place to hold the result
	result := &Config{
		rootDir: dir,
		fs:      fs,
		Tasks:   map[string]*Task{},
	}

	// if we have a task file
	if taskFileExists {
		// read its contents
		contents, err := afero.ReadFile(fs, taskFilePath)
		if err != nil {
			return nil, err
		}

		// parse the contents as hcl and use that as the starting point for
		// language-specific configuration
		err = hcl.Unmarshal(contents, result)
		if err != nil {
			return nil, err
		}
	}

	// if we have a package.json file
	if packageJSONExists {
		// each script in the package.json is a run task
		packageJSON := struct{ Scripts map[string]string }{}

		// read the contents of the file
		contents, err := afero.ReadFile(fs, packageJSONPath)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(contents, &packageJSON)

		for scriptName := range packageJSON.Scripts {
			// add a task to the config
			result.Tasks[scriptName] = &Task{
				Name:        scriptName,
				Description: "<none>",
				Script:      fmt.Sprintf("npm run %s", scriptName),
			}
		}
	}

	// nothing went wrong
	return result, nil
}
