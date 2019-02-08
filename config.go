package run

import (
	"errors"
	"path"

	"github.com/hashicorp/hcl"
	"github.com/spf13/afero"
)

// TaskFileName is the name of the file that defines run's tasks
const TaskFileName = "_tasks.hcl"

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
