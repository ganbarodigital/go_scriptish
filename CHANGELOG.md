# CHANGELOG

## develop

### New

* Added `NewPipeline()` to create a new Pipeline to execute
* Sources:
  - added `CatFile()`
  - added `CatStdin()`
  - added `CutFields()`
  - added `Echo()`
  - added `EchoArgs()`
  - added `EchoSlice()`
  - added `ListFiles()`
  - added `MkTempFile()`
* Filters:
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
* Sinks:
  - added `ToStderr()`
  - added `ToStdout()`
* Outputs:
  - added `CountLines()`
  - added `CountWords()`
* Errors:
  - added `ErrMismatchedInputs`
* Utilities:
  - added `ParseRangeSpec()`