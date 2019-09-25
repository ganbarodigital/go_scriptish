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
	"errors"
	"os"
	"strings"
	"testing"

	pipe "github.com/ganbarodigital/go_pipe/v2"
	"github.com/stretchr/testify/assert"
)

func TestAppendToTempFileReturnsFilename(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// this is the pipeline we'll use to test TruncateFile()
	pipeline := NewPipeline(
		CatFile("testdata/truncatefile/content.txt"),
		AppendToTempFile(os.TempDir(), "scriptify-appendtotempfile-*"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(actualResult, os.TempDir()+"/scriptify-"))

	// clean up after ourselves
	err = ExecPipeline(RmFile(actualResult)).Error()
	assert.Nil(t, err)
}

func TestAppendToTempFileSetsErrorWhenTempFileCannotBeCreated(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	pipeline := NewPipeline(
		AppendToTempFile("/does/not/exist", "scriptify-*"),
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

type brokenReader struct {
}

func (r brokenReader) Read(buf []byte) (int, error) {
	return 0, errors.New("mock reader does not work by design")
}

func TestAppendToTempFileSetsErrorWhenReadFromPipelineStdinFails(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we can't use a full pipeline to trigger this branch of the code
	singlePipe := pipe.NewPipe()

	// we need to replace the pipeline's normal Stdin with our own
	// broken version
	singlePipe.Stdin = pipe.NewSourceFromReader(brokenReader{})

	// ----------------------------------------------------------------
	// perform the change

	op := AppendToTempFile(os.TempDir(), "scriptify-appendtotempfile-*")
	statusCode, err := op(singlePipe)
	actualResult := singlePipe.Stdout.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, NOT_OK, statusCode)

	// the name of the tempfile DOES exist in the pipe's stdout
	assert.True(t, strings.HasPrefix(actualResult, os.TempDir()+"/scriptify-"))

	// clean up after ourselves
	ExecPipeline(RmFile(actualResult))
}
