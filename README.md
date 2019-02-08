# run

A task runner with templates and hcl for configuration

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
}

variables {
    hello = "hello"
}

config {
    // we can even change the delimiter that our templates use
    delimiters = ["{%", "%}"]
}
```

## Installing

Installing `run` is easily done with `go get`:

```bash
$ go get github.com/alecaivazis/run
```

