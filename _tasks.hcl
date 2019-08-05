task "test" {
    command = "echo $MESSAGE"
    environment {
        MESSAGE = "hello world"
    }
}