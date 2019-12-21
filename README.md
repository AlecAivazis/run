# run

A task runner with templates and hcl for configuration.

```hcl
// inside of _tasks.hcl

task "build" {
    description = "this is the description for the build task"
    // a single command to execute with variable expansion
    command = "echo {% .hello %}"
}

task "foo" {
    description = "description with variable: {% .hello %} "
    command = "echo {% .hello %}"
}

task "bar" {
    description = "another description"
    // these get executed in series
    pipeline = [
        "echo 1",
        "echo 2",
    ]
    // set these environment variables for the command/pipeline
    environment {
        KEY = "value"
    }
}

variables {
    hello = "hello"
}

config {
    // you can change the delimiter that the templates use. Default is ["{{", "}}"]
    delimiters = ["{%", "%}"]
}
```

## Installing

Installing `run` is easily done with `go get`:

```bash
$ go get github.com/alecaivazis/run
```

## Script Definitions

As shown in the example above, scripts and various configuration for `run` is defined in a
file called `_tasks.hcl`. This was designed with the go community in mind which lacks a built in
script-runners in its build tooling. In order to smoothen the experience when
working with different languages, `run` can get script definitions from other places too!

For example, run will look at the `package.json` file in a node project for additional script definitions.
You can still use `_tasks.hcl` aswell if you want.
