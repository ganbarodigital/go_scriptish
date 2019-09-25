# CHANGELOG

## develop

### New

* Added `NewPipeline()` to create a new Pipeline to execute
* Added `ExecPipeline()` to create and run a pipeline in a single step
* Added `PipelineFunc()` to create a pipeline and turn it into a function
* Sources:
  - added `CatFile()`
  - added `CatStdin()`
  - added `Chmod()`
  - added `CutFields()`
  - added `Echo()`
  - added `EchoArgs()`
  - added `EchoSlice()`
  - added `Exec()`
  - added `FilepathExists()`
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
  - added `XargsFilepathExists()`
  - added `XargsRmFile()`
  - added `XargsTruncateFiles()`
* Sinks:
  - added `AppendToFile()`
  - added `Return()`
  - added `RmDir()`
  - added `RmFile()`
  - added `ToStderr()`
  - added `ToStdout()`
  - added `TruncateFile()`
  - added `WriteToFile()`
* Errors:
  - added `ErrMismatchedInputs`
* Utilities:
  - added `ParseRangeSpec()`