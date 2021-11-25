// Copyright 2019-present Ganbaro Digital Ltd
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//   * Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//
//   * Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in
//     the documentation and/or other materials provided with the
//     distribution.
//
//   * Neither the names of the copyright holders nor the names of his
//     contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
// FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
// COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
// ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"fmt"
	"os"

	scriptish "github.com/ganbarodigital/go_scriptish"
)

// the best way to use Scriptish is like this:
//
// don't try to translate all of your original shell script into a
// single Scriptish pipeline
//
// create a Golang function per Scriptish pipeline, and keep your
// pipelines focused on doing a single task at a time
//
// use the capture functions like TrimmedString() to get the results
// back into your Golang code
//
// write your logic in native Golang code
//
// and take advantage of the tracing support to help you debug things
// like escaping characters in regexes :)

// ErrNotAGitRepo is the error we return when we are not in a Git repo
type ErrNotAGitRepo struct{}

func (e ErrNotAGitRepo) Error() string {
	return "fatal: not a git repository (or any of the parent directories): .git"
}

func findGitBinary() (string, error) {
	pipeline := scriptish.NewPipeline(
		scriptish.Which("git"),
	)

	return pipeline.Exec().TrimmedString()
}

func gitCurrentBranch(gitPath string) (string, error) {
	pipeline := scriptish.NewPipeline(
		scriptish.Exec([]string{gitPath, "branch", "--no-color"}),
		scriptish.Grep("^\\\\* "),
		scriptish.GrepV("no branch"),
		scriptish.CutFields("2"),
	)

	err := pipeline.Exec().Error()

	// are we in a Git rep?
	if err != nil {
		return "", ErrNotAGitRepo{}
	}

	// all done
	return pipeline.TrimmedString()
}

func main() {
	// do we want to trace?
	if len(os.Args) > 1 && os.Args[1] == "-x" {
		scriptish.GetShellOptions().EnableTrace(os.Stderr)
	}

	// do we have a git binary anywhere?
	//
	// we can print a friendlier error message if we make an
	// explicit check for it
	gitPath, err := findGitBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "*** error: unable to find 'git' in your PATH")
		os.Exit(1)
	}

	// what is the current branch?
	branch, err := gitCurrentBranch(gitPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("%s\n", branch)
}
