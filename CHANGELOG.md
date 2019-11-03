# CHANGELOG

## develop

### New

* Added support for shell options
  - Added `GetShellOptions()`
* Added debugging output support
  - Added `ShellOptions.EnableTrace()`
  - Added `ShellOptions.DisableTrace()`
  - Added `ShellOptions.IsTraceEnabled()`
  - Added `IsTraceEnabled()`
  - Added `Tracef()`
  - Added `TraceOutput()`
  - Added `TraceOsStderr()`
  - Added `TraceOsStdout()`
  - Added `TracePipeStderr()`
  - Added `TracePipeStdout()`
  - Lists now write status code and error message when tracing enabled
  - Pipelines now write status code and error message when tracing enabled
  - `And()` now supports tracing
  - `AppendToTempFile()` now supports tracing
  - `Chmod()` now supports tracing
  - `CountLines()` now supports tracing
  - `CountWords()` now supports tracing
  - `CutFields()` now supports tracing
  - `DropEmptyLines()` now supports tracing
  - `Echo()` now supports tracing
  - `EchoArgs()` now supports tracing
  - `EchoRawSlice()` now supports tracing
  - `EchoSlice()` now supports tracing
  - `Grep()` now supports tracing
  - `GrepV()` now supports tracing
  - `Head()` now supports tracing
  - `If()` now supports tracing
  - `IfElse()` now supports tracing
  - `ListFiles()` now supports tracing
  - `Or()` now supports tracing
  - `Return()` now supports tracing
  - `Rsort()` now supports tracing
  - `Sort()` now supports tracing
  - `StripExtension()` now supports tracing
  - `SwapExtensions()` now supports tracing
  - `Tail()` now supports tracing
  - `TestEmpty()` now supports tracing
  - `TestFilepathExists()` now supports tracing
  - `TestNotEmpty()` now supports tracing
  - `Tr()` now supports tracing
  - `TrimSuffix()` now supports tracing
  - `TrimWhitespace()` now supports tracing
  - `Uniq()` now supports tracing
  - `XargsBasename()` now supports tracing
  - `XargsCat()` now supports tracing
  - `XargsDirname()` now supports tracing
  - `XargsRmFile()` now supports tracing
  - `XargsTestFilepathExists()` now supports tracing
  - `XargsTruncateFiles()` now supports tracing
* Pipelines now set a context flag in their Pipe
  - this is used to tell sinks where to read from!

### Fixes

* `AppendToFile()` now supports lists
* `AppendToTempFile()` now supports lists
* `EchoArgs()` no longer includes `os.Args[0]` in the output (compatibility fix)

### Dependencies

* Update to `go_pipe` v5.1.0

## v1.1.1

Released Saturday, 2nd November 2019.

### Fixes

* Logic constructs now pass on the positional parameters

## v1.1.0

Released Saturday, 2nd November 2019.

### New

* Added `TestEmpty() builtin (emulates `[[ -z $VAR ]]`)
* Added `TestNotEmpty() builtin (emulates `[[ -n $VAR ]]`)
* Added `XargsBasename()`
* Added `XargsDirname()`

### Fixes

* `Basename()` now takes an input (compatibility fix)
* `Chmod()` no longer writes to the pipeline on success (compatibility fix)
* `Dirname()` now takes an input (compatibility fix)
* `Exit()` is no longer conditionally-built
  - dents our code coverage figures, but more important that it's available when you build your tests
* `TestFilepathExists()` no longer writes to the pipeline on success
  - now more accurate for how a UNIX shell `[[ -e filepath ]]` behaves

## v1.0.0

Released Wednesday, 30th October 2019.

### New

* Added `NewPipeline()` to create a new Pipeline to execute
* Added `NewPipelineFunc()` to create a pipeline and turn it into a function
* Added `ExecPipeline()` to create and run a pipeline in a single step
* Added `NewList()` to create a new List to execute
* Added `NewListFunc()` to create a List and turn it into a function
* Added `ExecList()` to create and run a List in a single step
* Sources:
  - added `CatFile()`
  - added `CatStdin()`
  - added `Chmod()`
  - added `CutFields()`
  - added `Echo()`
  - added `EchoArgs()`
  - added `EchoSlice()`
  - added `EchoToStderr()`
  - added `Exec()`
  - added `ListFiles()`
  - added `Lsmod()`
  - added `MkTempDir()`
  - added `MkTempFile()`
  - added `Which()`
* Filters:
  - added `AppendToTempFile()`
  - added `Basename()`
  - added `CountLines()`
  - added `CountWords()`
  - added `Dirname()`
  - added `DropEmptyLines()`
  - added `Grep()`
  - added `GrepV()`
  - added `Head()`
  - added `Rsort()`
  - added `RunList()`
  - added `RunPipeline()`
  - added `Sort()`
  - added `StripExtension()`
  - added `SwapExtensions()`
  - added `Tail()`
  - added `Tr()`
  - added `TrimSuffix()`
  - added `TrimWhitespace()`
  - added `Uniq()`
  - added `XargsCat()`
  - added `XargsTestFilepathExists()`
  - added `XargsRmFile()`
  - added `XargsTruncateFiles()`
* Sinks:
  - added `AppendToFile()`
  - added `Exit()`
  - added `Return()`
  - added `RmDir()`
  - added `RmFile()`
  - added `ToStderr()`
  - added `ToStdout()`
  - added `TruncateFile()`
  - added `WriteToFile()`
* Logic:
  - added `And()`
  - added `If()`
  - added `IfElse()`
  - added `Or()`
* Expressions:
  - added `TestFilepathExists()`
* Errors:
  - added `ErrMismatchedInputs`
* Utilities:
  - added `ParseRangeSpec()`