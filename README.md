# Welcome to scriptish

## Introduction

Scriptish is a Golang library. It helps you port UNIX shell scripts to Golang.

It is released under the 3-clause New BSD license. See [LICENSE.md](LICENSE.md) for details.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.CountLines(),
).Exec().ParseInt()
```

(We're going to create Scriptish for other languages too, and we'll update this README when those are available!)

## Table of Contents <!-- omit in toc -->

- [Introduction](#introduction)
- [Why Use Scriptish?](#why-use-scriptish)
  - [Who Is Scriptish For?](#who-is-scriptish-for)
  - [Why UNIX Shell Scripts?](#why-unix-shell-scripts)
  - [Why Not Use UNIX Shell Scripts?](#why-not-use-unix-shell-scripts)
  - [Enter Scriptish](#enter-scriptish)
- [How Does It Work?](#how-does-it-work)
  - [Getting Started](#getting-started)
  - [What Is A Pipeline?](#what-is-a-pipeline)
  - [What Happens When A Pipeline Runs?](#what-happens-when-a-pipeline-runs)
  - [How Are Errors Handled?](#how-are-errors-handled)
  - [Sources, Filters, Sinks, Logic and Capture Methods](#sources-filters-sinks-logic-and-capture-methods)
- [Creating A Pipeline](#creating-a-pipeline)
  - [NewPipeline()](#newpipeline)
  - [NewPipelineFunc()](#newpipelinefunc)
  - [ExecPipeline()](#execpipeline)
- [Running An Existing Pipeline](#running-an-existing-pipeline)
- [Passing Parameters Into Pipelines](#passing-parameters-into-pipelines)
- [Calling A Pipeline From Another Pipeline](#calling-a-pipeline-from-another-pipeline)
- [Capturing The Output](#capturing-the-output)
- [Pipelines vs Lists](#pipelines-vs-lists)
- [Creating A List](#creating-a-list)
  - [NewList()](#newlist)
  - [NewListFunc()](#newlistfunc)
  - [ExecList()](#execlist)
- [Running An Existing List](#running-an-existing-list)
- [Passing Parameters Into Lists](#passing-parameters-into-lists)
- [Calling A List From Another List Or Pipeline](#calling-a-list-from-another-list-or-pipeline)
- [Pipelines, Lists and Sequences](#pipelines-lists-and-sequences)
- [UNIX Shell String Expansion](#unix-shell-string-expansion)
  - [What Is String Expansion?](#what-is-string-expansion)
  - [Setting Positional Parameters](#setting-positional-parameters)
  - [Setting Local Variables](#setting-local-variables)
  - [Escaping Strings](#escaping-strings)
  - [Filename Globbing / Pathname Expansion](#filename-globbing--pathname-expansion)
- [From Bash To Scriptish](#from-bash-to-scriptish)
- [Sources](#sources)
  - [CatFile()](#catfile)
  - [CatStdin()](#catstdin)
  - [Chmod()](#chmod)
  - [Echo()](#echo)
  - [EchoArgs()](#echoargs)
  - [EchoSlice()](#echoslice)
  - [EchoToStderr()](#echotostderr)
  - [Exec()](#exec)
  - [ListFiles()](#listfiles)
  - [Lsmod()](#lsmod)
  - [MkTempDir()](#mktempdir)
  - [MkTempFile()](#mktempfile)
  - [Which()](#which)
- [Filters](#filters)
  - [AppendToTempFile()](#appendtotempfile)
  - [Basename()](#basename)
  - [CountLines()](#countlines)
  - [CountWords()](#countwords)
  - [CutFields()](#cutfields)
  - [DropEmptyLines()](#dropemptylines)
  - [Grep()](#grep)
  - [GrepV()](#grepv)
  - [Head()](#head)
  - [Rsort()](#rsort)
  - [RunPipeline()](#runpipeline)
  - [Sort()](#sort)
  - [StripExtension()](#stripextension)
  - [SwapExtensions()](#swapextensions)
  - [Tail()](#tail)
  - [Tr()](#tr)
  - [TrimSuffix()](#trimsuffix)
  - [TrimWhitespace()](#trimwhitespace)
  - [Uniq()](#uniq)
  - [XargsCat()](#xargscat)
  - [XargsDirname()](#xargsdirname)
  - [XargsRmFile()](#xargsrmfile)
  - [XargsTestFilepathExists()](#xargstestfilepathexists)
  - [XargsTruncateFiles()](#xargstruncatefiles)
- [Sinks](#sinks)
  - [AppendToFile()](#appendtofile)
  - [Exit()](#exit)
  - [Return()](#return)
  - [RmDir()](#rmdir)
  - [RmFile()](#rmfile)
  - [ToStderr()](#tostderr)
  - [ToStdout()](#tostdout)
  - [TruncateFile()](#truncatefile)
  - [WriteToFile()](#writetofile)
- [Builtins](#builtins)
  - [TestFilepathExists()](#testfilepathexists)
  - [TestNotEmpty()](#testnotempty)
- [Capture Methods](#capture-methods)
  - [Bytes()](#bytes)
  - [Error()](#error)
  - [Okay()](#okay)
  - [ParseInt()](#parseint)
  - [String()](#string)
  - [Strings()](#strings)
  - [TrimmedString()](#trimmedstring)
- [Logic Calls](#logic-calls)
  - [And()](#and)
  - [If()](#if)
  - [IfElse()](#ifelse)
  - [Or()](#or)
- [Errors](#errors)
  - [ErrMismatchedInputs](#errmismatchedinputs)
- [Inspirations](#inspirations)
  - [Compared To Labix's Pipe](#compared-to-labixs-pipe)
  - [Compared To Bitfield's Script](#compared-to-bitfields-script)

## Why Use Scriptish?

### Who Is Scriptish For?

We've built Scriptish for anyone who needs to replace UNIX shell scripts with compiled Golang binaries.

We're going to be doing that ourselves for some of our projects:

* Dockhand - Docker management utility
* HubFlow - the GitFlow extension for Git
* SimpleCA - local TLS certificate authority for internal infrastructure

We'll add links to those projects when they're available.

### Why UNIX Shell Scripts?

UNIX shell scripts are one of the most practical inventions in computer programming.

* They're very quick to write.
* They're very powerful.
* They treat everything as text.

They're fantastic for knocking up utilities in a matter of minutes, for automating things, and for gluing things together. Our hard drives are littered with little shell scripts - and some large ones too! - and we create new ones all the time.

If you're using any sort of UNIX system (Linux, or MacOS), **shell scripting is a must-have skill** - whether you're a developer or a sysadmin.

### Why Not Use UNIX Shell Scripts?

UNIX shell scripts are great until you want to share them with other people. They're just not a great choice if you want to distribute your work outside your team, organisation or community.

* If someone else is going to run your shell scripts, they need to make sure that they've installed all the commands that your shell scripts call. This can end up being a trial-and-error process. And what happens if they can't install those commands for any reason?

* Creating portable shell scripts (e.g. scripts that run on both Linux and MacOS) isn't always easy, and is very difficult (if not impossible) to test via a CI process.

* What about your Windows users? UNIX shell scripts don't work on a vanilla Windows box.

### Enter Scriptish

If you want to distribute shell scripts, it's best not to write them as shell scripts. Use Scriptish to quickly do the same thing in Golang:

* There's one binary to ship to your users.
* Scriptish is self-contained. No need to worry about installing additional commands (unless you call [scriptish.Exec()](#exec) ...)
* Use Golang's `go test` to create tests for your tools.
* Use the power of Golang to cross-compile binaries for Linux, MacOS and Windows.

## How Does It Work?

### Getting Started

Import Scriptish into your Golang code:

```go
import scriptish "github.com/ganbarodigital/go_scriptish"
```

Create a pipeline, and provide it with a list of commands:

```go
pipeline := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.CountWords()
)
```

Once you have your pipeline, run it:

```go
pipeline.Exec()
```

Once you've run your pipeline, call one of the [capture methods](#capture-methods) to find out what happened:

```go
result, err := pipeline.ParseInt()
```

### What Is A Pipeline?

UNIX shell scripts compose UNIX commands into a pipeline:

```bash
cat /path/to/file.txt | wc -w
```

The UNIX commands execute from left to right. The output (known as `stdout`) of each command becomes the input (known as `stdin`) of the next command.

The output of the final command can be captured by your shell script to become the value of a variable:

```bash
current_branch=$(git branch --no-color | grep "^[*] " | sed -e 's/^[*] //')
```

Scriptish works the same way. You create a pipeline of Scriptish commands:

```go
pipeline := scriptish.NewPipeline(
    scriptish.Exec("git", "branch", "--no-color"),
    scriptish.Grep("^[* ]"),
    scriptish.Tr([]string{"* "}, []string{""}),
)
```

and then you run it:

```go
pipeline.Exec()
```

The output of the final command can be captured by your Golang code to become the value of a variable, using [capture methods](#capture-methods):

```
current_branch, err := pipeline.TrimmedString()
```

### What Happens When A Pipeline Runs?

UNIX commands in a pipeline:

* read text input from `stdin`
* write their results (as text!) to `stdout`
* write any errors (as text!) out to `stderr`
* return a status code to indicate what happened

Each Scriptish command works the same way:

* they read text input from the pipeline's `Stdin`
* they write their results to the pipeline's `Stdout`
* they write any error messages out to the pipeline's `Stderr`
* they return a status code and a Golang error to indicate what happened

When a single command has finished, its `Stdout` becomes the `Stdin` for the next command in the pipeline.

### How Are Errors Handled?

One difference between UNIX commands and Golang is error handling. Scriptish combines the best of both.

* UNIX commands return a status code to indicate what happened. A status code of 0 (zero) means success.
* Scriptish commands return the UNIX-like status code, and any Golang error that has occurred. We store these in the pipeline.

If you're calling external commands using `scriptish.Exec()`, you've still got access to the UNIX status code exactly like a shell script does. And you've always got access to any Golang errors that have occurred too.

**Unlike** UNIX shell scripts, a Scriptish pipeline stops executing if any command returns an error.

You might not be aware of it, but by default, a pipeline in a UNIX shell script continues to run even if one of the commands returns an error. This causes error values to propagate - and error propagation is a major cause of robustness issues in software.

Philosophically, we believe that good software engineering practices are more important than UNIX shell compatibility.

### Sources, Filters, Sinks, Logic and Capture Methods

Scriptish commands fall into one of four categories:

* [Sources](#sources) create content in the pipeline, e.g. `scriptish.CatFile()`. They ignore whatever's already in the pipeline.
* [Filters](#filters) do something with (or to) the pipeline's content, and they write the results back into the pipeline. These results form the input content for the next pipeline command.
* [Sinks](#sinks) do something with (or two) the pipeline's content, and don't write any new content back into the pipeline.
* [Logic](#logic-calls) implement support for `if`-like statements directly in Scriptish.

A pipeline normally:

* starts with a _source_
* applies one or more _filters_
* finishes with a _sink_ to send the results somewhere

But what if we want to get the results back into our Golang code, to reuse in some way? Instead of using a [sink](#sinks), use a [capture method](#capture-methods) instead.

A capture method isn't a Scriptish command. It's a method on the `Pipeline` struct:

```go
fileExists = scriptish.ExecPipeline(
    scriptish.TestFilepathExists("/path/to/file.txt")
).Okay()
```

## Creating A Pipeline

You can create a pipeline in several ways.

Pipeline                      | Produces                             | Best For
------------------------------|--------------------------------------|---------------------------------
`scriptish.NewPipeline()`     | Pipeline that's ready to run         | Reusable pipelines
`scriptish.NewPipelineFunc()` | Function that will run your pipeline | Getting results back into Golang
`scriptish.ExecPipeline()`    | Pipeline that has been run once      | Throwaway pipelines

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

Most of the examples in this README (and most of the unit tests) use `scriptish.NewPipeline()`.

### NewPipelineFunc()

`NewPipelineFunc()` builds the pipeline and turns it into a function.

```go
fileExistsFunc := scriptish.NewPipelineFunc(
    scriptish.FileExists("/path/to/file")
)
```

Whenever you call the function, the pipeline executes. The function returns a `*Pipeline`. Use any of the [capture methods](#capture-methods) to find out what happened when the pipeline executed.

```go
fileExists := fileExistsFunc().Okay()
```

You can re-use the function as often as you want.

`NewPipelineFunc()` is great for pipelines where you want to get the results back into your Golang code:

```go
getCurrentBranch := scriptish.NewPipelineFunc(
    scriptish.Exec("git", "branch", "--no-color"),
    scriptish.Grep("^[* ]"),
    scriptish.Tr([]string{"* "}, []string{""}),
)

currentBranch, err := getCurrentBranch().TrimmedString()
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

You can then use any of the [capture methods](#capture-methods) to find out what happened:

```go
result, err = pipeline.ParseInt()
```

You can re-use the resulting pipeline as often as you want.

`ExecPipeline()` is great for pipelines that you want to throw away after use:

```go
result, err := scriptish.ExecPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.CountWords()
).ParseInt()
```

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
result, err := pipeline.Exec().ParseInt()
```

```go
// in this example, we don't keep a reference to the pipeline
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.CountWords()
).Exec().ParseInt()
```

## Passing Parameters Into Pipelines

Both `Pipeline.Exec()` and the function returned by `NewPipelineFunc()` accept a list of parameters.

```golang
pipeline := scriptish.NewPipeline(
    scriptish.CatFile("$1"),
    scriptish.CountWords()
)
wordCount, _ := pipeline.Exec("/path/to/file").ParseInt()
fmt.Printf("file has %d words\n", wordCount)
```

```golang
countWordsInFile := scriptish.NewPipelineFunc(
    scriptish.CatFile("$1"),
    scriptish.CountWords()
)
wordCount, _ := countWordsInFile("/path/to/file").ParseInt()
fmt.Printf("file has %d words\n", wordCount)
```

The positional variables `$1`, `$2`, `$3` et al are available inside the pipeline, just like they would be in a UNIX shell script function.

## Calling A Pipeline From Another Pipeline

UNIX shell scripts can be broken up into functions to make them easier to maintain. You can do something similar in Scriptish, by calling a pipeline from another pipeline:

```go
// this will parse the output of Git to find the selected branch
//
// the selected branch depends on the Git command called
filterSelectedBranch := scriptish.NewPipeline(
    scriptish.Grep("^[*] "),
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

## Capturing The Output

If you're familiar with UNIX shell scripting, you'll know that every shell command creates three different outputs:

* `stdout` - normal text output
* `stderr` - any error messages
* status code - an integer representing what happened. 0 (zero) means success, any other value means an error occurred.

Scriptish commands work the same way. They also track any Golang errors that occur when the commands run.

Property              | Description
----------------------| -------------------
`pipeline.Stdout`     | normal text output
`pipeline.Stderr`     | an error messages (this is normally blank, because we have Golang errors too)
`pipeline.Err`        | Golang errors
`pipeline.StatusCode` | an integer representing what happened. Normally 0 for success

When the pipeline has executed, you can call one of [the capture methods](#capture-methods) to find out what happened:

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.CountWords()
).Exec().ParseInt()

// if the pipeline worked ...
// - result now contains the number of words in the file
// - err is nil
//
// and if the pipeline didn't work ...
// - result is 0
// - err contains a Golang error
```

If you want to run a Scriptish command and you don't care about capturing the output, call `Pipeline.Okay()`:

```go
success := scriptish.NewPipeline(
    scriptish.RmFile("/path/to/file.txt")
).Exec().Okay()

// if the pipeline worked ...
// - success is `true`
//
// and if the pipeline didn't work ...
// - success is `false`
```

## Pipelines vs Lists

UNIX shell scripts support two main ways (known as sequences) to string individual commands together:

* _pipelines_ feed the output from one command into the next one
* _lists_ simply append the output from each command to `stdout` and `stderr`

Most of the time, you'll want to stick to _pipelines_, and port the rest of your shell script's behaviour over to native Golang code. That gives you the convenience of Scriptish's emulation of classic UNIX shell commands and the power of everything that Golang can do.

Sometimes, you'll find it less effort to use a few _lists_ too.

A classic example is `die()`. It's very common for UNIX shell scripts to define their own `die()` function like this:

```bash
die() {
    echo "*** error: $*"
    exit 1
}

[[ -e ./Dockerfile ]] || die "cannot find Dockerfile"
```

Here's the equivalent Scriptish:

```golang
dieFunc := scriptish.NewList(
    scriptish.Echo("*** error: $*"),
    scriptish.ToStderr(),
    scriptish.Exit(1),
)

scriptish.ExecList(
    scriptish.TestFileExists("./Dockerfile"),
    scriptish.Or(dieFunc("cannot find Dockerfile")),
)
```

## Creating A List

You can create a list in several ways.

Pipeline                  | Produces                         | Best For
--------------------------|----------------------------------|---------------------------------
`scriptish.NewList()`     | List that's ready to run         | Reusable lists
`scriptish.NewListFunc()` | Function that will run your list | Getting results back into Golang
`scriptish.ExecList()`    | List that has been run once      | Throwaway lists

### NewList()

Call `NewList()` when you want to build a list:

```go
list := scriptish.NewList(
    scriptish.Echo("*** warning: $*"),
)
```

`list` can now be executed as often as you want.

```go
list.Exec("cannot find Dockerfile")
```

### NewListFunc()

`NewListFunc()` builds the list and turns it into a function.

```go
fileExistsFunc := scriptish.NewListFunc(
    scriptish.FileExists("/path/to/file")
)
```

Whenever you call the function, the list executes. The function returns a `*List`. Use any of the [capture methods](#capture-methods) to find out what happened when the list executed.

```go
fileExists := fileExistsFunc().Okay()
```

You can re-use the function as often as you want.

`NewListFunc()` is great for lists where you want to get the results back into your Golang code.

### ExecList()

`ExecList()` builds a list and executes it in a single step.

```go
list := scriptish.ExecList(
    scriptish.CatFile("/path/to/file1.txt"),
    scriptish.CatFile("/path/to/file2.txt"),
)
```

Behind the scenes, it simply does a `scriptish.NewList(...).Exec()` for you.

You can then use any of the [capture methods](#capture-methods) to find out what happened:

```go
result, err = list.ParseInt()
```

You can re-use the resulting list as often as you want.

`ExecList()` is great for lists that you want to throw away after use.

## Running An Existing List

Once you have built a list, call the `Exec()` method to execute it:

```go
list := scriptish.NewList(
    scriptish.CatFile("/path/to/file1.txt"),
    scriptish.CatFile("/path/to/file2.txt"),
)
list.Exec()
```

`Exec()` always returns a pointer to the same list, so that you can use method chaining to create nicer-looking code.

```go
// in this example, `list` is available to be used more than once
list := scriptish.NewList(
    scriptish.CatFile("/path/to/file1.txt"),
    scriptish.CatFile("/path/to/file2.txt"),
)
result, err := pipeline.Exec().String()
```

```go
// in this example, we don't keep a reference to the list
result, err := scriptish.NewList(
    scriptish.CatFile("/path/to/file1.txt"),
    scriptish.CatFile("/path/to/file2.txt"),
).Exec().String()
```

## Passing Parameters Into Lists

Both `List.Exec()` and the function returned by `NewListFunc()` accept a list of parameters.

```golang
list := scriptish.NewPipeline(
    scriptish.CatFile("$1"),
    scriptish.CatFile("$2"),
)
fileContents := list.Exec("/path/to/file1.txt", "/path/to/file2.txt").String()
```

```golang
getTwoFileContents := scriptish.NewPipelineFunc(
    scriptish.CatFile("$1"),
    scriptish.CatFile("$2"),
)
fileContents := getTwoFileContents("/path/to/file1.txt", "/path/to/file2.txt").String()
```

The positional variables `$1`, `$2`, `$3` et al are available inside the list, just like they would be in a UNIX shell script function.

## Calling A List From Another List Or Pipeline

UNIX shell scripts can be broken up into functions to make them easier to maintain. You can do something similar in Scriptish, by calling a list from another pipeline or list:

```golang
// this will fetch the latest changes from the upstream Git repo
fetch_changes_from_origin := scriptish.NewList(
    scriptish.Exec("git", "remote", "update", "$ORIGIN"),
    scriptish.Or(dieFunc("Unable to get list of changes from '$ORIGIN'")),
    scriptish.Exec("git", "fetch", "$ORIGIN"),
    scriptish.Or(dieFunc("Unable to fetch latest changes from '$ORIGIN'")),
    scriptish.Exec("git", "fetch", "--tags"),
    scriptish.Or(dieFunc("Unable to fetch latest tags from '$ORIGIN'")),
)

// this will fetch the latest changes from upstream, and then
// merge them into local branches
merge_latest_changes = scriptish.NewList(
    scriptish.RunList(fetch_changes_from_origin),
    ...
)
```

## Pipelines, Lists and Sequences

In UNIX shell programming, pipelines and lists are both examples of a _sequence of commands_. Each one is a set of commands that are wrapped in slightly different execution logic.

In Scriptish, a `Pipeline` and a `List` are type aliases for a `Sequence`. A call to `NewPipeline()` or `NewList()` creates a `Sequence` that also has the right execution logic for a pipeline or a list. We've done it this way so that you can use both pipelines and lists in our [logic calls](#logic-calls).

All of our logic calls accept `Sequence` parameters. You can pass in a `Pipeline` or a `List` to suit, and either will work just fine.

## UNIX Shell String Expansion

### What Is String Expansion?

One of the things that makes UNIX shells so powerful is the way they expand a string, or a line of code, before executing it.

```bash
#!/usr/bin/env bash

PARAM1=hello world

# output: "Hello world"
echo ${PARAM1^H}
```

We've integrated the [ShellExpand package](https://github.com/ganbarodigital/go_shellexpand) so that string expansion is available to you.

### Setting Positional Parameters

The positional parameters are `$1`, `$2`, `$3` all the way up to `$9`, as well as `$#` and `$*`. These are exactly the same as their equivalents in shell scripts.

To set these, pass parameters [into pipeline](#passing-parameters-into-pipelines) or [into lists](#passing-parameters-into-lists).

### Setting Local Variables

Every `Pipeline` and `List` struct comes with a `LocalVars` member. You can call its `Setenv()` method to create more variables to use in string expansion:

```golang
// create a reusable List
fetch_changes_from_remote := scriptish.NewList(
    scriptish.Exec("git", "remote", "update", "$REMOTE"),
    scriptish.Or(dieFunc("Unable to get list of changes from '$REMOTE'")),
    scriptish.Exec("git", "fetch", "$REMOTE"),
    scriptish.Or(dieFunc("Unable to fetch latest changes from '$REMOTE'")),
    scriptish.Exec("git", "fetch", "--tags"),
    scriptish.Or(dieFunc("Unable to fetch latest tags from '$REMOTE'")),
)

// set the value of '$REMOTE'
fetch_changes_from_remote.LocalVars.Setenv("REMOTE", "origin")

// run it
fetch_changes_from_remote.Exec()
```

Any local variables that you set will remain set if you reuse the pipeline or list - ie, they are persistent.

### Escaping Strings

The one downside of string expansion is that you will need to escape characters in your strings, to avoid them being interpreted as instructions to the string expansion engine.

The basic rule of thumb is that if you'd need to escape it in a shell script, you'll also need to escape it in a string passed into Scriptish.

### Filename Globbing / Pathname Expansion

At the moment, the string expansion does not support globbing (properly known as _pathname expansion_). That means you can't use wildcards in filepaths anywhere.

This is something we might add in a future release.

## From Bash To Scriptish

Here's a handy table to help you quickly translate an action from a Bash shell script to the equivalent Scriptish command.

Bash                         | Scriptish
-----------------------------|------------------------------------------------
`$(...)`                     | [`scriptish.Exec()`](#exec)
`${x%.*}`                    | [`scriptish.StripExtension()`](#stripextension)
`${x%$y}%z`                  | [`scriptish.SwapExtensions()](#swapextensions)
`${x%$y}`                    | [`scriptish.TrimSuffix()`](#trimsuffix)
`[[ -e $x ]]`                | [`scriptish.TestFilepathExists()`](#testfilepathexists)
`[[ -n $x ]]`                | [`scriptish.TestNotEmpty()`](#testnotempty)
`> $file`                    | [`scriptish.WriteToFile()`](#writetofile)
`>> $file`                   | [`scriptish.AppendToFile()`](#appendtofile)
`||`                         | [`scriptish.Or()`](#or)
`&&`                         | [`scriptish.And()`](#and)
`basename ...`               | [`scriptish.Basename()`](#basename)
`cat "..."`                  | [`scriptish.CatFile(...)`](#catfile)
`cat /dev/null > $x`         | [`scriptish.TruncateFile($x)`](#truncatefile)
`chmod`                      | [`scriptish.Chmod()`](#chmod)
`cut -f`                     | [`scriptish.CutFields()`](#cutfields)
`dirname ...`                | [`scriptish.Dirname()`](#dirname)
`echo "..."`                 | [`scriptish.Echo(...)`](#echo)
`echo "$@"`                  | [`scriptish.EchoArgs()`](#echoargs)
`exit ...`                   | [`scriptish.Exit()`](#exit)
`function`                   | [`scriptish.RunPipeline()`](#runpipeline)
`grep ...`                   | [`scriptish.Grep()`](#grep)
`grep -v ..`                 | [`scriptish.GrepV()`](#grepv)
`head -n X`                  | [`scriptish.Head(X)`](#head)
`if expr ; then body ; fi`   | [`scriptish.If()`](#if)
`if expr ; then body ; else elseBlock ; fi` | [`scriptish.IfElse()`](#ifelse)
`ls -1 ...`                  | [`scriptish.ListFiles(...)`](#listfiles)
`ls -l | awk '{ print $1 }'` | [`scriptish.Lsmod()`](#lsmod)
`mktemp`                     | [`scriptish.MkTempFile()`](#mktempfile)
`mktemp -d`                  | [`scriptish.MkTempDir()`](#mktempdir)
`return`                     | [`scriptish.Return()`](#return)
`rm -f`                      | [`scriptish.RmFile()`](#rmfile)
`rm -r`                      | [`scriptish.RmDir()`](#rmdir)
`sort`                       | [`scriptish.Sort()`](#sort)
`sort -r`                    | [`scriptish.Rsort()`](#rsort)
`tail -n X`                  | [`scriptish.Tail(X)`](#tail)
`tr old new`                 | [`scriptish.Tr(old, new)`](#tr)
`uniq`                       | [`scriptish.Uniq()`](#uniq)
`wc -l`                      | [`scriptish.CountLines()`](#countlines)
`wc -w`                      | [`scriptish.CountWords()`](#countwords)
`which`                      | [`scriptish.Which()`](#which)
`xargs cat`                  | [`scriptish.XargsCat()`](#xargscat)
`xargs rm`                   | [`scriptish.XargsRmFile()`](#xargsrmfile)
`xargs test -e`              | [`scriptish.XargsTestFilepathExists()`](#xargstestfilepathexists)

## Sources

Sources get data from outside the pipeline, and write it into the pipeline's `Stdout`.

Every pipeline normally begins with a source, and is then followed by one or more [filters](#filters).

### CatFile()

`CatFile()` writes the contents of a file to the pipeline's `Stdout`.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("a/file.txt"),
).Exec().String()
```

### CatStdin()

`CatStdin()` copies the program's `stdin` (`os.Stdin` in Golang terms) to the pipeline's `Stdout`.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatStdin(),
).Exec().String()
```

### Chmod()

`Chmod()` attempts to change the permissions on the given filepath.

It ignores the contents of the pipeline.

On success, it writes the filepath to the pipeline's stdout, in case anything else in the pipeline can use it.

```go
result, err := scriptish.NewPipeline(
    scriptish.Chmod("/path/to/file", 0644)
).Exec().String()
```

### Echo()

`Echo()` writes a string to the pipeline's `Stdout`.

```go
result, err := scriptish.NewPipeline(
    scriptish.Echo("hello world"),
).Exec().String()
```

### EchoArgs()

`EchoArgs()` writes the program's arguments to the pipeline's `Stdout`, one line per argument.

```go
result, err := scriptish.NewPipeline(
    scriptish.EchoArgs(),
).Exec().String()
```

### EchoSlice()

`EchoSlice()` writes an array of strings to the pipeline's `Stdout`, one line per array entry.

```go
myStrings := []string{"hello world", "have a nice day"}

result, err := scriptish.NewPipeline(
    scriptish.EchoSlice(myStrings),
).Exec().String()
```

### EchoToStderr()

`EchoToStderr()` writes a string to the pipeline's `Stderr`.

```go
result, err := scriptish.NewPipeline(
    scriptish.EchoToStderr("*** error: file not found"),
).Exec()
```

### Exec()

`Exec()` executes an operating system command. The command's `stdin` will be the pipeline's `Stdin`, and it will write to the pipeline's `Stdout` and `Stderr`.

The command's status code will be stored in the pipeline's `StatusCode`.

```go
localBranch, err := scriptish.ExecPipeline(
    scriptish.Exec("git", "branch", "--no-color"),
    scriptish.Grep("^[* ]"),
    scriptish.Tr([]string{"* "}, []string{""}),
).TrimmedString()
```

Use the [`.Okay()`](#okay) capture method if you simply want to know if the command worked or not:

```go
success := scriptish.ExecPipeline(scriptish.Exec("git", "push")).Okay()
```

Golang will set `err` to an [`exec.ExitError`](https://golang.org/pkg/os/exec/#ExitError) if the command's status code is not 0 (zero).

Golang will set `err` to an [`os.PathError`](https://golang.org/pkg/os/#PathError) if the command could not be found in the first place.

### ListFiles()

`ListFiles()` writes a list of matching files to the pipeline's `Stdout`, one line per filename found.

* If `path` is a file, ListFiles writes the file to the pipeline's `Stdout`
* If `path` is a folder, ListFiles writes the contents of the folder to the pipeline's `Stdout`. The path to the folder is included.
* If `path` contains wildcards, ListFiles writes any files that matches to the pipeline's `Stdout`. `ListFiles()` uses Golang's `os.Glob()` to expand the wildcards.

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

### Lsmod()

`Lsmod()` writes the permissions of the given filepath to the pipe's stdout.

* Symlinks are followed.
* Permissions are in the form '-rwxrwxrwx'.

It ignores the contents of the pipeline.

```go
result, err := scriptish.NewPipeline(
    scriptish.Lsmod("/path/to/file"),
).Exec().TrimmedString()
```

### MkTempDir()

`MkTempDir()` creates a temporary folder, and writes the filename to the pipeline's `Stdout`.

```go
result, err := scriptish.NewPipeline(
    scriptish.MkTempFile(os.TempDir(), "scriptish-")
).Exec().String()
```

### MkTempFile()

`MkTempFile()` creates a temporary file, and writes the filename to the pipeline's `Stdout`.

```go
result, err := scriptish.NewPipeline(
    scriptish.MkTempFile(os.TempDir(), "scriptish-*")
).Exec().String()
```

### Which()

`Which()` searches the current PATH to find the given path. If one is found, the command's path is written to the pipeline's `Stdout`.

It ignores the contents of the pipeline.

```go
result, err := scriptish.NewPipeline(
    scriptish.Which("git"),
).Exec().String()
```

## Filters

Filters read the contents of the pipeline's `Stdin`, do something to that data, and write the results out to the pipeline's `Stdout`.

When you've finished adding filters to your pipeline, you should either add a [sink](#sinks), or call one of the [capture methods](#capture-methods) to get the results back into your Golang code.

### AppendToTempFile()

`AppendToTempFile()` writes the contents of the pipeline's `Stdin` to a
temporary file. The temporary file's filename is then written to
the pipeline's `Stdout`.

If the file does not exist, it is created.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.AppendToTempFile(os.TempDir(), "scriptish-*"),
).Exec().TrimmedString()

// result now contains the temporary filename
```

### Basename()

`Basename()` treats each line in the pipeline's `Stdin` as a filepath. Any parent elements are stripped from the line, and the results written to the pipeline's `Stdout`.

Any blank lines are preserved.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder/*.txt"),
    scriptish.Basename()
).Exec().Strings()
```

### CountLines()

`CountLines()` counts the number of lines in the pipeline's `Stdin`, and writes that to the pipeline's `Stdout`.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("path/to/folder/*.txt"),
    scriptish.CountLines(),
).Exec().ParseInt()
```

### CountWords()

`CountWords()` counts the number of words in the pipeline's `Stdin`, and writes that to the pipeline's `Stdout`.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("path/to/file.txt"),
    scriptish.CountWords(),
).Exec().ParseInt()
```

### CutFields()

`CutFields()` retrieves only the fields specified on each line of the pipeline's `Stdin`, and writes them to the pipeline's `Stdout`.

```go
result, err := scriptish.NewPipeline(
    scriptish.Echo("one two three four five"),
    scriptish.CutFields("2-3,5")
).Exec().String()
```

### DropEmptyLines()

`DropEmptyLines()` removes any lines that are blank, or that only contain whitespace.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.DropEmptyLines()
).Exec().String()
```

### Grep()

`Grep()` filters out lines that do not match the given regex.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Grep("second|third"),
).Exec().String()
```

### GrepV()

`GrepV()` filters out lines that match the given regex.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.GrepV("second|third"),
).Exec().String()
```

### Head()

`Head()` copies the first N lines of the pipeline's `Stdin` to its `Stdout`.

If N is zero or negative, `Head()` copies no lines.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Head(100),
).Exec().String()
```

### Rsort()

`Rsort()` sorts the contents of the pipeline into descending alphabetical order.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Rsort(),
).Exec().String()
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
).Exec().String()
```

### StripExtension()

`StripExtension()` treats every line in the pipeline as a filepath. It removes the extension from each filepath.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder"),
    scriptish.StripExtension(),
).Exec().Strings()
```

### SwapExtensions()

`SwapExtensions()` treats every line in the pipeline as a filepath.

It replaces every old extension with the corresponding new one.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder"),
    scriptish.SwapExtensions([]string{"txt","yml"}, []string{"text","yaml"}),
).Exec().Strings()
```

If the second parameter is a string slice of length 1, every old file extension will be replaced by that parameter.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder"),
    scriptish.SwapExtensions([]string{"yml","YAML"}, []string{"yaml"}),
).Exec().Strings()
```

If the first and second parameters are different lengths, `SwapExtensions()` will return an `scriptish.ErrMismatchedInputs`.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder"),
    scriptish.SwapExtensions([]string{"one"}, []string{"1","2"}),
).Exec().Strings()

// err is an ErrMismatchedInputs, and result is empty
```

### Tail()

`Tail()` copies the last N lines from the pipeline's `Stdin` to its `Stdout`.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Tail(50),
).Exec().String()
```

### Tr()

`Tr()` replaces all occurances of one string with another.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Tr([]string{"one","two"}, []string{"1","2"}),
).Exec().String()
```

If the second parameter is a string slice of length 1, everything from the first parameter will be replaced by that slice.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Tr([]string{"one","two"}, []string{"numberwang"}),
).Exec().String()
```

If the first and second parameters are different lengths, `Tr()` will return an `scriptish.ErrMismatchedInputs`.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.Tr([]string{"one","two","three"}, []string{"1","2"}),
).Exec().String()

// err is an ErrMismatchedInputs, and result is empty
```

### TrimSuffix()

`TrimSuffix()` removes the given suffix from each line of the pipeline.

Use it to emulate `basename(1)`'s `[suffix]` parameter.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder/"),
    scriptish.TrimSuffix(".txt")
).Exec().Strings()
```

### TrimWhitespace()

`TrimWhitespace()` removes any whitespace from the front and end of each line in the pipeline.

```go
result, err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.TrimWhitespace(),
).Exec().String()
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

`XargsCat()` treats each line in the pipeline's `Stdin` as a filepath. The contents of each file are written to the pipeline's `Stdout`.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder/*.txt"),
    scriptish.XargsCat()
).Exec().String()
```

### XargsDirname()

`XargsDirname()` treats each line in the pipeline's `Stdin` as a filepath. The last element is stripped from the line, and the results written to the pipeline's `Stdout`.

Any blank lines are turned in '.'

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder/*.txt"),
    scriptish.XargsDirname()
).Exec().Strings()
```

### XargsRmFile()

`XargsRmFile()` treats every line in the pipeline as a filename. It attempts to delete each file.

It stops at the first file that cannot be deleted.

Each successfully-deleted filepath is written to the pipeline's `Stdout`, for use by the next command in the pipeline.

```go
err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/files"),
    scriptish.XargsRmFile(),
    scriptish.ForXIn(
        scriptish.Printf("deleted file: $x")
    )
).Exec().Error()
```

### XargsTestFilepathExists()

`XargsTestFilepathExists()` treats each line in the pipeline as a filepath. It checks to see if the given filepath exists. If the filepath exists, it is written to the pipeline's stdout.

It does not care what the filepath points at (file, folder, named pipe, and so on).

```go
// example: find all RAW photo files that also have a corresponding
// JPEG file
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/folder/*.raw"),
    scriptish.SwapExtensions(".raw", ".jpeg"),
    scriptish.XargsTestFilepathExists()
).Exec().Strings()
```

### XargsTruncateFiles()

`XargsTruncatesFiles()` treats each line of the pipeline's `Stdin` as a filepath. The contents of each file are truncated. If the file does not exist, it is created.

Each filepath is written to the pipeline's `Stdout`, for use by the next command in the pipeline.

```go
result, err := scriptish.NewPipeline(
    scriptish.ListFiles("/path/to/files"),
    scriptish.XargsTruncateFiles(),
).Exec().Strings()

// result now contains a list of the files that have been truncated
```

## Sinks

Sinks take the contents of the pipeline's `Stdin`, and write it to somewhere outside the pipeline.

A sink should be the last command in your pipeline. You can add more commands afterwards if you really want to. Just be aware that the first command after any sink will be starting with an empty `Stdin`.

### AppendToFile()

`AppendToFile()` writes the contents of the pipeline's `Stdin` to the given file

If the file does not exist, it is created.

```go
err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.AppendToFile("my-app.log"),
).Exec().Error()
```

### Exit()

`Exit()` terminates your Golang program with the given status code. Use with caution.

```golang
dieFunc := scriptish.NewList(
    scriptish.Echo("*** error: $*"),
    scriptish.ToStderr(),
    scriptish.Exit(1),
)

scriptish.ExecList(
    scriptish.TestFilepathExists("./Dockerfile"),
    scriptish.Or(dieFunc("cannot find Dockerfile")),
)
```

### Return()

`Return()` terminates the pipeline with the given status code.

```go
statusCode := scriptish.NewPipeline(
    scriptish.Return(3)
).Exec().StatusCode()

// statusCode will be: 3
```

### RmDir()

`RmDir()` deletes the given folder, as long as the folder is empty.

It ignores the contents of the pipeline.

It ignores the file's file permissions, because the underlying Golang os.Remove() behaves that way.

```go
err := scriptish.NewPipeline(
    scriptish.RmDir("/path/to/folder"),
).Exec().Error()
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

`ToStdout()` writes the pipeline's `Stdin` to the program's `stderr` (`os.Stderr` in Golang terms).

```go
err := scriptish.NewPipeline(
    scriptish.Echo("usage: simpleca <command>"),
    scriptish.ToStderr(),
).Exec().Error()
```

### ToStdout()

`ToStdout()` writes the pipeline's `Stdin` to the program's `Stdout` (`os.Stdout` in Golang terms).

```go
err := scriptish.NewPipeline(
    scriptish.Echo("usage: simpleca <command>"),
    scriptish.ToStdout(),
).Exec().Error()
```

### TruncateFile()

`TruncateFile()` removes the contents of the given file.

If the file does not exist, it is created.

```go
err := scriptish.NewPipeline(
    scriptish.TruncateFile("/tmp/scriptish-test"),
).Exec().Error()
```

### WriteToFile()

`WriteToFile()` writes the contents of the pipeline's `Stdin` to the given file. The existing contents of the file are replaced.

If the file does not exist, it is created.

```go
err := scriptish.NewPipeline(
    scriptish.CatFile("/path/to/file.txt"),
    scriptish.WriteToFile("/path/to/new_file.txt"),
).Exec().Error()
```

## Builtins

Builtins are UNIX shell commands and UNIX CLI utilities that don't fall into the [sources](#sources), [sinks](#sinks) and [filters](#filters) categories:

* their input is a parameter; they ignore the pipeline
* their only output is the status code; they don't write anything new to the pipeline

### TestFilepathExists()

`TestFilepathExists()` checks to see if the given filepath exists. If it does, it returns `StatusOkay`. If not, it returns `StatusNotOkay`.

* It does not care what the filepath points at (file, folder, named pipe, and so on).
* It ignores the contents of the pipeline.
* It follows symbolic links.

```go
fileExists := scriptish.ExecPipeline(
    scriptish.TestFilepathExists("/path/to/file")
).Okay()
```

### TestNotEmpty()

`TestNotEmpty()` returns `StatusOkay` if the (expanded) input is empty; `StatusNotOkay` otherwise.

It is the equivalent to `if [[ -n $VAR ]]` in a UNIX shell script.

```bash
show_usage() {
    echo "*** error: $*"
    echo
    echo "usage: $0 <name-of-arg>"
    exit 1
}

[[ -n $1 ]] || show_usage("missing parameter <name-of-arg>")
```

Here's the equivalent Scriptish:

```golang
showUsage := scriptish.NewList(
    scriptish.Echo("*** error: $*"),
    scriptish.Echo(""),
    scriptish.Echo("usage: $0 <name-of-arg>")
    scriptish.Exit(1),
)

checkArgs := scriptish.NewList(
    scriptish.TestNotEmpty("$1"),
    scriptish.Or(showUsage("missing argument")),
)

checkArgs.Exec(os.Args...)
```

## Capture Methods

Capture methods available on each `Pipeline`. Use them to get the output from the pipeline.

Don't forget to run your pipeline first, if you haven't already!

### Bytes()

`Bytes()` is the standard interface's `io.Reader` `Bytes()` method.

* It writes the contents of the pipeline's `Stdout` into the byte slice that you provide.
* It returns the number of bytes written.
* It also returns the pipeline's current Golang error value.

Normally, you wouldn't call this yourself.

### Error()

`Error()` returns the pipeline's current Golang error status, which may be `nil`.

```go
err := scriptish.ExecPipeline(scriptish.RmFile("/path/to/file")).Error()
```

### Okay()

`Okay()` returns `true|false` depending on the pipeline's current UNIX status code.

```go
success := scriptish.ExecPipeline(scriptish.Exec("git push")).Okay()
```

`success` is a `bool`:

* `false` if the pipeline's `StatusCode` property is *not* 0
* `true` otherwise

All Scriptish commands set the pipeline's `StatusCode`, so it's safe to use `Okay()` to check any pipeline you create.

It's mostly there if you want to call a pipeline in a Golang `if` statement:

```go
if !scriptish.ExecPipeline(scriptish.Exec("git", "push")).Okay() {
    // push failed, do something about it
}
```

### ParseInt()

`ParseInt()` returns the pipeline's `Stdout` as an `int` value:

```go
lineCount, err := scriptish.ExecPipeline(
    scriptish.CatFile("/path/to/file"),
    scriptish.CountLines(),
).ParseInt()
```

If the pipeline's `Stdout` can't be turned into an integer, then it will return `0` and the parsing error from Golang's `strconv.ParseInt()`.

If the pipeline didn't execute successfully, it will return `0` and the pipeline's current Golang error status.

### String()

`String()` returns the pipeline's `Stdout` as a single string:

```go
contents, err := scriptish.ExecPipeline(
    scriptish.CatFile("/path/to/file")
).String()
```

**The string will be terminated by a linefeed `\n` character.** `String()` is a good choice if you want to get content into your Golang code. If you just want a single-line value, see [TrimmedString()](#trimmedstring) below.

If the pipeline's `Stdout` is empty, an empty string will be returned.

If the pipeline didn't execute successfully, the contents of the pipeline's `Stderr` will be returned. We might change this behaviour in the future.

### Strings()

`Strings()` returns the pipeline's `Stdout` as an array of strings (aka a string slice):

```go
files, err := scriptish.ExecPipeline(
    scriptish.ListFiles("/path/to/folder"),
    scriptish.Basename(),
).Strings()
```

Each string will *not* be terminated by a linefeed `\n` character.

If the pipeline's `Stdout` is empty, an empty string slice will be returned.

If the pipeline didn't execute successfully, the contents of the pipeline's `Stderr` will be returned. We might change this behaviour in the future.

### TrimmedString()

`TrimmedString()` returns the pipeline's `Stdout` as a single string, with leading and trailing whitespace removed:

```go
localBranch, err := scriptish.ExecPipeline(
    scriptish.Exec("git", "branch", "--no-color"),
    scriptish.Grep("^[* ]"),
    scriptish.Tr([]string{"* "}, []string{""}),
).TrimmedString()
```

`TrimmedString()` is the right choice if you're expecting a single line of text back. This is very useful for getting results back into your Golang code!

If the pipeline's `Stdout` is empty, an empty string will be returned.

If the pipeline didn't execute successfully, the contents of the pipeline's `Stderr` will be returned. We might change this behaviour in the future.

## Logic Calls

Most of the time, you will use native Golang code to write `if` statements for your code. Use the [capture methods](#capture-methods) to get the results of your pipelines back into your Golang code.

Sometimes, it will be more convenient to use Scriptish's built-in logic support. The classic example is the `die()` function commonly created in UNIX shell scripts to handle errors:

```bash
die() {
    echo "*** error: $*" >&2
    exit 1
}

[[ -e ./Dockerfile ]] || die "cannot find Dockerfile"
```

Here's the equivalent Scriptish:

```golang
dieFunc := scriptish.NewList(
    scriptish.Echo("*** error: $*"),
    scriptish.ToStderr(),
    scriptish.Exit(1),
)

scriptish.ExecList(
    scriptish.TestFilepathExists("./Dockerfile"),
    scriptish.Or(dieFunc("cannot find Dockerfile")),
)
```

(It also has to be said that implementing logic support in Scriptish was a good test case for proving Scriptish's underlying design.)

Generally, all the implemented logc takes Lists or Pipelines as arguments. (Lists and Pipelines are both aliases for the `Sequence` struct, so you can pass either in to suit.) If there are any exceptions to this rule, we'll make sure to point it out in the documentation for the individual logic call.

### And()

`And()` executes the given sequence only if the previous command did not return any kind of error.

The sequence starts with an empty Stdin. The sequence's output is written back to the Stdout and Stderr of the calling list or pipeline - along with the StatusCode() and Error().

It is an emulation of UNIX shell scripting's `command1 && command2` feature.

__NOTE that `And()` only works when run inside a List__:

```golang
scriptish.ExecList(
    scriptish.Exec("git", "fetch"),
    // if `git fetch` fails, do not attempt the merge
    scriptish.And(
        scriptish.NewList(
            scriptish.Exec("git", "merge")
        )
    )
)
```

If you call `And()` inside a Pipeline, it'll always run the given sequence. Pipelines terminate whenever a command returns an error, so `And()` will only be called if the previous command succeeded.

### If()

`If()` executes the `body` if (and only if) the given `expr` does not return any kind of error.

Both `expr` and `body` start with an empty `Stdin`. Their output is written back to the pipeline's `Stdout` and `Stderr`.

It is an emulation of UNIX shell scripting's `if expr ; then body ; fi` feature.

```golang
result, err := scriptish.ExecList(
    scriptish.If(
        // this is the `expr` or expression
        scriptish.NewPipeline(
            scriptish.TestFilepathExists("/path/to/file"),
        ),
        // this is the `body` that is executed if the `expr` succeeds
        scriptish.NewPipeline(
            scriptish.CatFile("/path/to/file"),
            scriptish.Head(3),
        ),
    )
).String()
```

You can safely use `If()` inside a pipeline, because it doesn't depend upon the result of any previous command.

### IfElse()

`IfElse()` executes the body if (and only if) the expr completes without an error. Otherwise, it executes the elseBlock instead.

`IfElse()` is an emulation of UNIX shell scripting's `if expr ; then body ; else elseBlock ; fi`.

```golang
result, err := scriptish.ExecList(
    scriptish.IfElse(
        // this is the `expr` or expression
        scriptish.NewPipeline(
            scriptish.TestFilepathExists("/path/to/file"),
        ),
        // this is the `body` that is executed if the `expr` succeeds
        scriptish.NewPipeline(
            scriptish.CatFile("/path/to/file"),
            scriptish.Head(3),
        ),
        // and this is the `elseBlock` that is executed if the `expr` fails
        scriptish.NewPipeline(
            scriptish.Echo("*** error: file not found"),
            scriptish.ToStderr(),
        )
    )
).String()
```

### Or()

`Or()` executes the given sequence only if the previous command has returned some kind of error.

The sequence starts with an empty Stdin. The sequence's output is written back to the Stdout and Stderr of the calling list or pipeline - along with the StatusCode() and Error().

It is an emulation of UNIX shell scripting's `list1 || command` feature.

__NOTE that `Or()` only works when run inside a List__:

```golang
statusCode, err := scriptish.NewList(
    scriptish.TestFilepathExists("/path/to/file"),
    scriptish.Or(dieFunc("cannot find file")),
).Exec().StatusError()
```

If you call `Or()` from inside a Pipeline, it will never work. Pipelines terminate when the first command returns an error. This means that either:

* the pipeline will terminate before `Or()` is reached, or
* `Or()` never executes the given sequence (because there's no previous error)

At the moment, we can't think of a way of detecting any attempt to call `Or()` from a pipeline.

## Errors

### ErrMismatchedInputs

`ErrMismatchedInputs` is returned whenever two input arrays aren't the same length.

## Inspirations

Scriptish is inspired by:

* [Labix's Pipe package](http://labix.org/pipe)
* [John Arundel / Bitfield's Script package](https://github.com/bitfield/script)

### Compared To Labix's Pipe

_Pipe_ is a bit more low level, and seems to be aimed predominantly at executing external commands, as a more powerful alternative to Golang's `exec` package.

If that's what you need, definitely check it out!

### Compared To Bitfield's Script

We started out using (and contributing to) Script, but ran into a few things that didn't suit what we were doing:

* Built for different purposes

    _script_ is aimed at doing minimal shell-like operations from a Golang app.

    _scriptish_ is more about providing everything necessary to recreate any size UNIX shell script - including powerful ones like the HubFlow extension for Git - to Golang, without having to port everything to Golang.

* Aimed at different people

    We want it to take as little thinking as possible to port UNIX shell scripts over to Golang - especially for casual or lapsed Golang programmers!

    That means (amongst other things) using function names that are similar to the UNIX shell command that they emulate, and emulating UNIX behaviour as closely as is practical, including UNIX features like status codes and `stderr`.

* Extensibility

    _script operations_ are methods on the `script.Pipe` struct. We found that this makes it very hard to extend `script` with your own methods, because Golang doesn't support inheritance, and the method chaining prevents Golang embedding from working.

    In contrast, _script commands_ are first-order functions that take the Pipe as a function parameter. You can create your own Scriptish commands, and they can live in your own Golang package.

* Reusability

    There's currently no way to call one pipeline from another using _script_ alone. You can achieve that by writing your own Golang boiler plate code.

    _scriptish_ builds first-order pipelines that you can run, pass around as values, and call from other _scriptish_ pipelines.

We were originally attracted to _script_ because of how incredibly easy it is to use. There's a lot to like about it, and we've definitely tried to create the same feel in _scriptish_ too. We've borrowed concepts such as [sources](#sources), [filters](#filters) and [sinks](#sinks) from _script_, because they're such a great way to describe how different Scriptish commands behave.

You should definitely check [script](https://github.com/bitfield/script) out if you think that _scriptish_ is too much for what you need.