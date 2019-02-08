# run

A task runner

## Installing

Installing `run` is easily done with `go get`:

```bash
$ go get github.com/alecaivazis/run
```

## Task File

`Run` uses a hcl file called `_task.hcl` to define the valid tasks for a give project. As long as the task
file is in the current directory or its parents, `run` will find it. For an example of a valid taskfile
as well as various configuration values, see [_task]
