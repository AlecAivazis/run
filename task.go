package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/spf13/cobra"
)

// Task represents a single task that can be executed by run
type Task struct {
	Name        string
	Description string
	Script      string `hcl:"command"`
	Pipeline    []string
}

// Run executes the task
func (t *Task) Run(args []string, c *Config) error {
	// make sure there is only one way to run the given task
	waysToRun := 0
	if t.Script != "" {
		waysToRun++
	}
	if len(t.Pipeline) > 0 {
		waysToRun++
	}
	if waysToRun != 1 {
		return fmt.Errorf("encountered invalid number of ways to run %s", t.Name)
	}

	// we are safe to run the command
	if t.Script != "" {
		return t.runScript(args, c)
	}
	return t.runPipeline(args, c)
}

func (t *Task) runScript(args []string, c *Config) error {
	// create a com
	return t.execute(args, c, t.Script)
}

func (t *Task) runPipeline(args []string, c *Config) error {
	return t.execute(args, c, t.Pipeline...)
}

func (t *Task) execute(arguments []string, c *Config, cmds ...string) error {
	// for each command we have to run
	for _, command := range cmds {
		// the command could be a template string
		tmpl, err := template.New("task-command").Delims(c.Settings.TemplateDelimiters[0], c.Settings.TemplateDelimiters[1]).Parse(command)
		if err != nil {
			return err
		}
		var cmdStr bytes.Buffer
		tmpl.Execute(&cmdStr, c.Variables)

		// build up the command
		args := append([]string{"-c", string(cmdStr.Bytes()), "sh"}, arguments...)
		cmd := exec.Command("bash", args...)

		// make sure the command prints to the right spots
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// if this command fails
		if err := cmd.Run(); err != nil {
			// abort
			return err
		}
	}
	return nil
}

func (t *Task) CobraCommand(config *Config) *cobra.Command {
	return &cobra.Command{
		Use:   t.Name,
		Short: t.Description,
		Run: func(cmd *cobra.Command, args []string) {
			if err := t.Run(args, config); err != nil {
				fmt.Printf("Sorry something went wrong: %s\n", err.Error())
				os.Exit(1)
				return
			}
		},
	}
}
