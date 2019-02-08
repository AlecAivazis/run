package main

import (
	"fmt"
	"os"
	"os/exec"

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
func (t *Task) Run(cmd *cobra.Command, args []string) error {
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
		return t.runScript(cmd, args)
	}
	return t.runPipeline(cmd, args)
}

func (t *Task) runScript(cmd *cobra.Command, args []string) error {
	// create a com
	return t.execute(args, t.Script)
}

func (t *Task) runPipeline(cmd *cobra.Command, args []string) error {
	return t.execute(args, t.Pipeline...)
}

func (t *Task) execute(arguments []string, cmds ...string) error {
	// for each command we have to run
	for _, command := range cmds {
		// build up the command
		args := append([]string{"-c", command, "sh"}, arguments...)
		cmd := exec.Command("sh", args...)

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
