# Welcome to scriptish

## Creating A Pipeline

You can create a pipeline in several ways.

Pipeline         | Produces
-----------------|--------------------------------------
`NewPipeline()`  | Pipeline that's ready to run
`ExecPipeline()` | Pipeline that has been run once
`PipelineFunc()` | Function that will run your pipeline

### NewPipeline()

Call `NewPipeline()` when you want to build a pipeline:

```go
pipeline := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.CountWords()
)
```

`pipeline` can now be executed as often as you want.

```go
result, err := pipeline.Exec().ParseInt()
```

### ExecPipeline()

`ExecPipeline()` builds a pipeline and executes it in a single step.

```go
pipeline := scriptish.ExecPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.CountWords()
)
```

Behind the scenes, it simply does a `scriptish.NewPipeline(...).Exec()` for you.

You can then use any of the output methods to find out what happened:

```go
result, err = pipeline.ParseInt()
```

You can re-use the resulting pipeline as often as you want.

### PipelineFunc()

`PipelineFunc()` builds the pipeline and turns it into a function.

```go
fileExistsFunc := scriptish.PipelineFunc(
    scriptish.FileExists("/path/to/file")
)
```

Whenever you call the function, the pipeline executes. The function returns a `*Pipeline`. Use any of the output methods to find out what happened when the pipeline executed.

```go
fileExists, err := fileExistsFunc().Okay()
```

You can re-use the function as often as you want.

## Running An Existing Pipeline

Once you have built a pipeline, call the `Exec()` method to execute it:

```go
pipeline := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.CountWords()
)
pipeline.Exec()
```

`Exec()` always returns a pointer to the same pipeline, so that you can use method chaining to create nicer-looking code.

```go
// in this example, `pipeline` is available to be used more than once
pipeline := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.CountWords()
)
result, err := pipeline.Exec().String()
```

```go
// in this example, we don't keep a reference to the pipeline
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.CountWords()
).Exec().String()
```

## Calling A Pipeline From Another Pipeline

UNIX shell scripts can be broken up into functions to make them easier to maintain. You can do something close to that in Scriptish, by calling a pipeline from another pipeline:

```go
// this will parse the output of Git to find the selected branch
//
// the selected branch depends on the Git command called
filterSelectedBranch := scriptish.NewPipeline(
    scriptish.Grep("^[*]"),
    scriptish.Tr([]string{"* "}, []string{""}),
)

// which local branch are we working on?
localBranch, err := scriptish.NewPipeline(
    scriptish.Exec("git branch --no-color"),
    scriptish.RunPipeline(filterSelectedBranch),
).Exec().TrimmedString()

// what's the tracking branch?
remoteBranch, err := scriptish.NewPipeline(
    scriptish.Exec("git branch -av --no-color"),
    scriptish.RunPipeline(filterSelectedBranch),
).Exec().TrimmedString()
```

## Getting A Result

If you're familiar with UNIX shell scripting, you'll know that every shell command creates three different outputs:

* stdout - normal text output
* stderr - any error messages
* status code - an integer representing what happened. 0 (zero) means success, any other value means an error occurred.

Scriptish commands work the same way. They also track any Golang errors that occur when the commands run.

Property              | Description
----------------------| -------------------
`pipeline.Stdout`     | normal text output
`pipeline.Stderr`     | an error messages (this is normally blank, because we have Golang errors too)
`pipeline.Err`        | Golang errors
`pipeline.StatusCode` | an integer representing what happened. Normally 0 for success

When the pipeline has executed, you can call one of [the output functions](#outputs) to find out what happened:

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.CountWords()
).Exec().String()

// if the pipeline worked ...
// - result now contains the number of words in the file
// - err is nil
//
// and if the pipeline didn't work ...
// - result is 0
// - err contains a Golang error
```

If you want to run a Scriptish command and you don't care about the output, call `Pipeline.Okay()`:

```go
success, err := scriptish.NewPipeline(
    scriptish.RmFile("/path/to/file.txt")
).Exec()

// if the pipeline worked ...
// - success is `true`
// - err is `nil`
//
// and if the pipeline didn't work ...
// - success is `false`
// - err contains a Golang error
```

## From Bash To Scriptish

Here's a handy table to help you quickly translate an action from a Bash shell script to the equivalent `scriptish` function.

Bash                 | Scriptish
---------------------|----------
`${x%.*}`            | `StripExtension()`
`${x%$y}`            | `TrimSuffix()`
`> $file`            | `WriteToFile()`
`>> $file`           | `AppendToFile()`
`basename ...`       | `Basename()`
`cat "..."`          | `Cat(...)`
`cat /dev/null > $x` | `TruncateFile($x)`
`cut -f`             | `CutFields()`
`dirname ...`        | `Dirname()`
`echo "..."`         | `Echo(...)`
`echo "$@"`          | `EchoArgs()`
`function`           | `RunPipeline()`
`head -n X`          | `Head(X)`
`ls -1 ...`          | `ListFiles(...)`
`mktemp`             | `MkTempFile()`
`mktemp -d`          | `MkTempDir()`
`rm -f`              | `RmFile()`
`sort`               | `Sort()`
`sort -r`            | `Rsort()`
`tail -n X`          | `Tail(X)`
`tr old new`         | `Tr(old, new)`
`uniq`               | `Uniq()`
`wc -l`              | `CountLines()`
`wc -w`              | `CountWords()`
`xargs cat`          | `XargsCat()`

## Sources

Sources get data from outside the pipeline, and write it into the pipeline's stdout.

Every pipeline normally begins with a source, and is then followed by one or more [filters](#filters).

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

### ListFiles()

`ListFiles()` writes a list of matching files to the pipeline's stdout, one line per filename found.// TrimSuffix removes the given suffix from each line of the pipeline.
//
// Use it to emulate basename(1)'s `[suffix]` parameter.


```go
// list a single file, if it exists
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("path/to/file"),
).Exec().String()
```

```go
// list all files in a folder, if the folder exists
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("path/to/folder"),
).Exec().String()
```

```go
// list all files in a folder that match wildcards, if the folder exists
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("path/to/folder/*.txt"),
).Exec().String()
```

### MkTempDir()

`MkTempDir()` creates a temporary folder, and writes the filename to the pipeline's stdout.

```go
result, err := scriptish.NewPipeline(
    scriptish.MkTempFile(os.TempDir(), "scriptish-")
).Exec.String()
```

### MkTempFile()

`MkTempFile()` creates a temporary file, and writes the filename to the pipeline's stdout.

```go
result, err := scriptish.NewPipeline(
    scriptish.MkTempFile(os.TempDir(), "scriptish-*")
).Exec.String()
```

## Filters

Filters read the contents of the pipeline's stdin, do something to that data, and write the results out to the pipeline's stdout.

When you've finished adding filters to your pipeline, you should either add a [sink](#sinks), or call one of the [output functions](#outputs) to get the results back into your Golang code.

### AppendToTempFile()

`AppendToTempFile()` writes the contents of the pipeline's stdin to a
temporary file. The temporary file's filename is then written to
the pipeline's stdout.

If the file does not exist, it is created.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.AppendToTempFile(os.TempDir(), "scriptish-*"),
).Exec.String()

// result now contains the temporary filename
```

### Basename()

`Basename()` treats each line in the pipeline's stdin as a filepath. Any parent elements are stripped from the line, and the results written to the pipeline's stdout.

Any blank lines are preserved.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder/*.txt"),
    scriptish.Basename()
).Exec().Strings()
```

### CountLines()

`CountLines()` counts the number of lines in the pipeline's stdin, and writes that to the pipeline's stdout.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("path/to/folder/*.txt"),
    scriptish.CountLines(),
).Exec().String()
```

### CountWords()

`CountWords()` counts the number of words in the pipeline's stdin, and writes that to the pipeline's stdout.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("path/to/file.txt"),
    scriptish.CountWords(),
).Exec().String()
```

### CutFields()

`CutFields()` retrieves only the fields specified on each line of the pipeline's stdin, and writes them to the pipeline's stdout.

```go
result, err := scriptish.NewPipeline(
    scriptish.Echo("one two three four five"),
    scriptish.CutFields("2-3,5")
).Exec().String()
```

### Dirname

`Dirname()` treats each line in the pipeline's stdin as a filepath. The last element is stripped from the line, and the results written to the pipeline's stdout.

Any blank lines are turned in '.'

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder/*.txt"),
    scriptish.Dirname()
).Exec().Strings()
```

### DropEmptyLines()

`DropEmptyLines()` removes any lines that are blank, or that only contain whitespace.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.DropEmptyLines()
).Exec().Strings()
```

### Head()

`Head()` copies the first N lines of the pipeline's stdin to its stdout.

If N is zero or negative, `Head()` copies no lines.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Head(100),
).Exec().Strings()
```

### Rsort()

`Rsort()` sorts the contents of the pipeline into descending alphabetical order.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Rsort(),
).Exec().Strings()
```

### RunPipeline()

`RunPipeline()` allows you to call one pipeline from another. Use it to create reusable pipelines, a bit like shell script functions.

```go
getWordCount := scriptish.NewPipeline(
    scriptish.SplitWords(),
    scriptish.CountLines(),
)

result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.RunPipeline(getWordCount),
).Exec().ParseInt()
```

### Sort()

`Sort()` sorts the contents of the pipeline into ascending alphabetical order.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Sort(),
).Exec().Strings()
```

### StripExtension()

`StripExtension()` treats every line in the pipeline as a filepath. It removes the extension from each filepath.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder"),
    scriptish.StripExtension(),
).Exec().Strings()
```

### Tail()

`Tail()` copies the last N lines from the pipeline's stdin to its stdout.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Tail(50),
).Exec().Strings()
```

### Tr()

`Tr()` replaces all occurances of one string with another.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Tr([]string{"one","two"}, []string{"1","2"}),
).Exec().Strings()
```

If the second parameter is a string slice of length 1, everything from the first parameter will be replaced by that slice.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Tr([]string{"one","two"}, []string{"numberwang"}),
).Exec().Strings()
```

If the first and second parameters are different lengths, `Tr()` will return an `scriptish.ErrMismatchedInputs`.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Tr([]string{"one","two"}, []string{"1","2"}),
).Exec().Strings()

// err is an ErrMismatchedInputs, and result is empty
```

### TrimSuffix()

`TrimSuffix()` removes the given suffix from each line of the pipeline.

Use it to emulate basename(1)'s `[suffix]` parameter.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder/"),
    scriptish.TrimSuffix(".txt")
).Exec().Strings()
```

### TrimWhitespace()

`TrimWhitespace()` removes any whitespace from the front and end of the line.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.TrimWhitespace(),
).Exec().Strings()
```

### Uniq()

`Uniq()` removes duplicated lines from the pipeline.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Uniq(),
).Exec().Strings()
```

### XargsCat()

`XargsCat()` treats each line in the pipeline's stdin as a filepath. The contents of each file are written to the pipeline's stdout.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder/*.txt"),
    scriptish.XargsCat()
).Exec().String()
```

### XargsTruncateFiles()

`XargsTruncatesFiles()` treats each line of the pipeline's stdin as a filepath. The contents of each file are truncated. If the file does not exist, it is created.

Each filepath is written to the pipeline's stdout, for use by the next operation in the pipeline.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/files"),
    scriptish.XargsTruncateFiles(),
).Exec().Strings()

// result now contains a list of the files that have been truncated
```

## Sinks

Sinks take the contents of the pipeline's stdin, and write it to somewhere outside the pipeline.

A sink should be the last operation in your pipeline. You can add more operations afterwards if you really want to. Just be aware that the first operation after any sink will be starting with an empty stdin.

### AppendToFile()

`AppendToFile()` writes the contents of the pipeline's stdin to the given file

If the file does not exist, it is created.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.AppendToFile("my-app.log"),
).Exec().String()
```

### RmFile()

`RmFile()` deletes the given file.

It ignores the contents of the pipeline.

It ignores the file's file permissions, because the underlying Golang os.Remove() behaves that way.

```go
err := scriptish.NewPipeline(
    scriptish.RmFile("/path/to/file"),
).Exec().Error()
```

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

### TruncateFile()

`TruncateFile()` removes the contents of the given file.

If the file does not exist, it is created.

```go
result, err := scriptish.NewPipeline(
    scriptish.TruncateFile("/tmp/scriptish-test"),
).Exec().String()
```

### WriteToFile()

`WriteToFile()` writes the contents of the pipe's stdin to the given file. The existing contents of the file are replaced.

If the file does not exist, it is created.

```go
err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.WriteToFile("/path/to/new_file.txt"),
).Exec().Error()
```

## Outputs

Outputs are methods available on each `Pipeline`. Some are inherited from the `pipe` package, and some are defined by `scriptish`.

### CountLines()

`CountLines()` returns the number of lines in the pipeline's stdout.

```go
lineCount, err := scriptish.NewPipeline(
    scriptish.ListFiles("."),
).Exec().CountLines()
```

If the pipeline failed to complete, `lineCount` will be `0`, and `err` will be the pipeline's last error status.

### CountWords()

`CountWords()` returns the number of words in the pipeline's stdout.

```go
wordCount, err := scriptish.NewPipeline(
    scriptish.ListFiles("."),
).Exec().CountWords()
```

If the pipeline failed to complete, `wordCount` will be `0`, and `err` will be the pipeline's last error status.

## Errors

### ErrMismatchedInputs

`ErrMismatchedInputs` is returned whenever two input arrays aren't the same length.