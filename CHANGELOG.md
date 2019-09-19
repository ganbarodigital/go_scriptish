# CHANGELOG

## develop

### New

* Added `NewPipeline()` to create a new Pipeline to execute
* Added `ExecPipeline()` to create and run a pipeline in a single step
* Added `PipelineFunc()` to create a pipeline and turn it into a function
* Sources:
  - added `CatFile()`
  - added `CatStdin()`
  - added `CutFields()`
  - added `Echo()`
  - added `EchoArgs()`
  - added `EchoSlice()`
  - added `FilepathExists()`
  - added `ListFiles()`
  - added `MkTempDir()`
  - added `MkTempFile()`
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
  - added `Tail()`
  - added `Tr()`
  - added `TrimSuffix()`
  - added `TrimWhitespace()`
  - added `Uniq()`
  - added `XargsCat()`
  - added `XargsTruncateFiles()`
* Sinks:
  - added `AppendToFile()`
  - added `RmFile()`
  - added `ToStderr()`
  - added `ToStdout()`
  - added `TruncateFile()`
  - added `WriteToFile()`
* Errors:
  - added `ErrMismatchedInputs`
* Utilities:
  - added `ParseRangeSpec()`