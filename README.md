# Welcome to scriptish

## From Bash To Scriptish

Here's a handy table to help you quickly translate an action from a Bash shell script to the equivalent `scriptish` function.

Bash           | Scriptish
---------------|----------
`${x%.*}`      | `StripExtension()`
`${x%$y}`      | `TrimSuffix()`
`basename ...` | `Basename()`
`cat "..."`    | `Cat(...)`
`cut -f`       | `CutFields()`
`dirname ...`  | `Dirname()`
`echo "..."`   | `Echo(...)`
`echo "$@"`    | `EchoArgs()`
`function`     | `RunPipeline()`
`head -n X`    | `Head(X)`
`ls -1 ...`    | `ListFiles(...)`
`mktemp`       | `MkTempFile()`
`mktemp -d`    | `MkTempDir()`
`sort`         | `Sort()`
`sort -r`      | `Rsort()`
`tail -n X`    | `Tail(X)`
`tr old new`   | `Tr(old, new)`
`uniq`         | `Uniq()`
`wc -l`        | `CountLines()`
`wc -w`        | `CountWords()`
`xargs cat`    | `XargsCat()`

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
).Exec().Int()
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

## Sinks

Sinks take the contents of the pipeline's stdin, and write it to somewhere outside the pipeline.

A sink should be the last operation in your pipeline. You can add more operations afterwards if you really want to. Just be aware that the first operation after any sink will be starting with an empty stdin.

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