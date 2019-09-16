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
* Filters:
  - added `Basename()`
  - added `CountLines()`
  - added `CountWords()`
  - added `Dirname()`
  - added `DropEmptyLines()`
  - added `Head()`
  - added `XargsCat()`
* Sinks:
  - added `ToStderr()`
  - added `ToStdout()`
* Outputs:
  - added `CountLines()`
  - added `CountWords()`
* Utilities:
  - added `ParseRangeSpec()`