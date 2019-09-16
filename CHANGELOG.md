# CHANGELOG

## develop

### New

* Added `NewPipeline()` to create a new Pipeline to execute
* Added `ExecPipeline()` to create and run a pipeline in a single step
* Sources:
  - added `CatFile()`
  - added `CatStdin()`
  - added `CutFields()`
  - added `Echo()`
  - added `EchoArgs()`
  - added `EchoSlice()`
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
  - added `ToStderr()`
  - added `ToStdout()`
  - added `TruncateFile()`
  - added `WriteToFile()`
* Outputs:
  - added `CountLines()`
  - added `CountWords()`
* Errors:
  - added `ErrMismatchedInputs`
* Utilities:
  - added `ParseRangeSpec()`