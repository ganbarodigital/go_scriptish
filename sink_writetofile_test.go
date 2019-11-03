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

	pipe "github.com/ganbarodigital/go_pipe/v5"
	"github.com/stretchr/testify/assert"
)

func TestWriteToFileWritesPipelineToGivenFile(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	tmpFilename, err := ExecPipeline(MkTempFile(os.TempDir(), "scriptify-writetofile-")).TrimmedString()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(RmFile(tmpFilename))

	// we will use this to prove that the file write worked
	lineCountPL := NewPipeline(
		CatFile(tmpFilename),
		CountLines(),
	)

	pipeline := NewPipeline(
		CatFile("./testdata/truncatefile/content.txt"),
		WriteToFile(tmpFilename),
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
	assert.Equal(t, 3, lineCount)
}

func TestWriteToFileOverwritesExistingFileContents(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	tmpFilename, err := ExecPipeline(MkTempFile(os.TempDir(), "scriptify-writetofile-")).TrimmedString()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(RmFile(tmpFilename))

	// we will use this to prove that the file write worked
	lineCountPL := NewPipeline(
		CatFile(tmpFilename),
		CountLines(),
	)

	// we need to put some content into the temp file to start with
	err = ExecPipeline(
		Echo("this is a test line"),
		WriteToFile(tmpFilename),
	).Error()
	assert.Nil(t, err)

	pipeline := NewPipeline(
		CatFile("./testdata/truncatefile/content.txt"),
		WriteToFile(tmpFilename),
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
	assert.Equal(t, 3, lineCount)
}

func TestWriteToFileSetsErrorWhenCreateFileFails(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	pipeline := NewPipeline(
		CatFile("./testdata/truncatefile/content.txt"),
		WriteToFile("/does/not/exist/invalid/path"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestWriteToFileDoesNothingReadFromPipelineStdinFails(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	tmpFilename, err := ExecPipeline(
		MkTempFile(os.TempDir(), "scriptify-writetofile-"),
	).TrimmedString()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(RmFile(tmpFilename))

	// we can't use a full pipeline to trigger this branch of the code
	singlePipe := pipe.NewPipe()

	// we need to replace the pipeline's normal Stdin with our own
	// broken version
	singlePipe.Stdin = pipe.NewSourceFromReader(brokenReader{})

	// ----------------------------------------------------------------
	// perform the change

	op := WriteToFile(tmpFilename)
	statusCode, err := op(singlePipe)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, StatusOkay, statusCode)
}
