// scriptish is a library to help you port bash scripts to Golang
//
// inspired by:
//
// - http://labix.org/pipe
// - https://github.com/bitfield/script
//
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

package scriptish

import (
	"os"
	"testing"

	pipe "github.com/ganbarodigital/go_pipe/v6"
	"github.com/stretchr/testify/assert"
)

func TestAppendToFileWritesPipelineToGivenFile(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	tmpFilename, err := ExecPipeline(MkTempFile(os.TempDir(), "scriptish-appendtofile-*")).TrimmedString()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(RmFile(tmpFilename))

	// we will use this to prove that the file append worked
	lineCountPL := NewPipeline(
		CatFile(tmpFilename),
		CountLines(),
	)

	// we need to put some content into the temp file to start with
	ExecPipeline(
		Echo("this is a test line"),
		WriteToFile(tmpFilename),
	)

	pipeline := NewPipeline(
		CatFile("./testdata/truncatefile/content.txt"),
		AppendToFile(tmpFilename),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Empty(t, actualResult)

	// the file should have content
	lineCount, err := lineCountPL.Exec().ParseInt()
	assert.Nil(t, err)
	assert.Equal(t, 4, lineCount)
}

func TestAppendToFileSetsErrorIfFileCannotBeCreated(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewPipeline(
		CatFile("./testdata/truncatefile/content.txt"),
		AppendToFile("/does/not/exist/and/never/will"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Empty(t, actualResult)
}

func TestAppendToFileDoesNothingWhenReadFromPipelineStdinFails(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we need a valid filename to try to write to
	tmpFilename, err := ExecPipeline(MkTempFile(os.TempDir(), "scriptish-appendtofile-*")).TrimmedString()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(RmFile(tmpFilename))

	// we can't use a full pipeline to trigger this branch of the code
	singlePipe := pipe.NewPipe()

	// we need to replace the pipeline's normal Stdin with our own
	// broken version
	singlePipe.Stdin = &brokenReader{}

	// ----------------------------------------------------------------
	// perform the change

	op := AppendToFile(tmpFilename)
	statusCode, err := op.RunStep(singlePipe)
	actualResult := singlePipe.Stdout.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, StatusOkay, statusCode)

	// AppendToFile() is a sink, Stdout should be empty
	assert.Empty(t, actualResult)
}

func TestAppendToFileWritesToTheTraceOutputWhenInList(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	tmpFilename, err := ExecPipeline(MkTempFile(os.TempDir(), "scriptish-appendtofile-*")).TrimmedString()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(RmFile(tmpFilename))

	// we need to put some content into the temp file to start with
	ExecPipeline(
		Echo("this is a test line"),
		WriteToFile(tmpFilename),
	)

	expectedResult := `+ CatFile("./testdata/truncatefile/content.txt")
+ => CatFile("./testdata/truncatefile/content.txt")
+ p.Stdout> This is a file of test data.
+ p.Stdout> ` + "" + `
+ p.Stdout> We copy the contents of this file to other files, as part of our testing.
+ AppendToFile("$1")
+ => AppendToFile("` + tmpFilename + `")
+ file> This is a file of test data.
+ file> ` + "" + `
+ file> We copy the contents of this file to other files, as part of our testing.
`
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	list := NewList(
		CatFile("./testdata/truncatefile/content.txt"),
		AppendToFile("$1"),
	)

	// ----------------------------------------------------------------
	// perform the change

	list.Exec(tmpFilename)
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestAppendToFileWritesToTheTraceOutputWhenInPipeline(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	tmpFilename, err := ExecPipeline(MkTempFile(os.TempDir(), "scriptish-appendtofile-*")).TrimmedString()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(RmFile(tmpFilename))

	// we need to put some content into the temp file to start with
	ExecPipeline(
		Echo("this is a test line"),
		WriteToFile(tmpFilename),
	)

	expectedResult := `+ CatFile("./testdata/truncatefile/content.txt")
+ => CatFile("./testdata/truncatefile/content.txt")
+ p.Stdout> This is a file of test data.
+ p.Stdout> ` + "" + `
+ p.Stdout> We copy the contents of this file to other files, as part of our testing.
+ AppendToFile("$1")
+ => AppendToFile("` + tmpFilename + `")
+ file> This is a file of test data.
+ file> ` + "" + `
+ file> We copy the contents of this file to other files, as part of our testing.
`
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	pipeline := NewPipeline(
		CatFile("./testdata/truncatefile/content.txt"),
		AppendToFile("$1"),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec(tmpFilename)
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestAppendToFileEmptiesThePipeWhenInList(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	tmpFilename, err := ExecPipeline(MkTempFile(os.TempDir(), "scriptish-appendtofile-*")).TrimmedString()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(RmFile(tmpFilename))

	// we need to put some content into the temp file to start with
	ExecPipeline(
		Echo("this is a test line"),
		WriteToFile(tmpFilename),
	)

	list := NewList(
		CatFile("./testdata/truncatefile/content.txt"),
		AppendToFile("$1"),
	)

	// ----------------------------------------------------------------
	// perform the change

	list.Exec(tmpFilename)

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, list.Pipe.Stdin.String())
	assert.Empty(t, list.Pipe.Stdout.String())
	assert.Empty(t, list.Pipe.Stderr.String())
}

func TestAppendToFileEmptiesThePipeWhenInPipeline(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	tmpFilename, err := ExecPipeline(MkTempFile(os.TempDir(), "scriptish-appendtofile-*")).TrimmedString()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(RmFile(tmpFilename))

	// we need to put some content into the temp file to start with
	ExecPipeline(
		Echo("this is a test line"),
		WriteToFile(tmpFilename),
	)

	pipeline := NewPipeline(
		CatFile("./testdata/truncatefile/content.txt"),
		AppendToFile("$1"),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec(tmpFilename)

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, pipeline.Pipe.Stdin.String())
	assert.Empty(t, pipeline.Pipe.Stdout.String())
	assert.Empty(t, pipeline.Pipe.Stderr.String())
}
