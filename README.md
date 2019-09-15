# Welcome to scriptish

## From Bash To Scriptish

Here's a handy table to help you quickly translate an action from a Bash shell script to the equivalent `scriptish` function.

Bash       | Scriptish
-----------|----------
cat "..."  | `Cat(...)`
echo "..." | `Echo(...)`
echo "$@"  | `EchoArgs()`

## Sources

### CatFile()

`CatFile()` writes the contents of a file to the pipeline's stdout.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("a/file.txt"),
).Exec().String()
```

### CatStdin()

`CatStdin()` copies the program's stdin to the pipeline's stdout.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatStdin(),
).Exec().String()
```

### Echo()

`Echo()` writes a string to the pipeline's stdout.

```go
result, err := scriptish.NewPipeline(
    scriptish.Echo("hello world"),
).Exec().String()
```

### EchoArgs()

`EchoArgs()` writes the program's arguments to the pipeline's stdout, one line per argument.

```go
result, err := scriptish.NewPipeline(
    scriptish.EchoArgs(),
).Exec().String()
```

### EchoSlice()

`EchoSlice()` writes an array of strings to the pipeline's stdout, one line per array entry.

```go
myStrings := []string{"hello world", "have a nice day"}

result, err := scriptish.NewPipeline(
    scriptish.EchoSlice(myStrings),
).Exec().String()
```

## Sinks()

### ToStderr()

`ToStdout()` writes the pipeline's stdin to the program's stderr.

```go
result, err := scriptish.NewPipeline(
    scriptish.Echo("usage: simpleca <command>"),
    scriptish.ToStderr(),
).Exec().String()
```

### ToStdout()

`ToStdout()` writes the pipeline's stdin to the program's stdout.

```go
result, err := scriptish.NewPipeline(
    scriptish.Echo("usage: simpleca <command>"),
    scriptish.ToStdout(),
).Exec().String()
```