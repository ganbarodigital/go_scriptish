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

	pipe "github.com/ganbarodigital/go_pipe"
	"github.com/stretchr/testify/assert"
)

func TestRmFileRemovesAGivenFile(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	tmpFilename, err := ExecPipeline(
		// create the temporary file
		MkTempFile(os.TempDir(), "scriptify-"),
	).TrimmedString()
	assert.Nil(t, err)
	assert.NotEmpty(t, tmpFilename)

	// clean up after ourselves
	defer ExecPipeline(RmFile(tmpFilename))

	// this pipeline will prove if the temporary file does/does not exist
	tmpFileExists := PipelineFunc(FilepathExists(tmpFilename))
	fileExists, err := tmpFileExists().Okay()
	assert.Nil(t, err)
	assert.True(t, fileExists)

	pipeline := NewPipeline(
		RmFile(tmpFilename),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Empty(t, actualResult)

	// the file should be gone
	fileExists, err = tmpFileExists().Okay()
	assert.Error(t, err)
	assert.False(t, fileExists)
}

func TestRmFileSetsErrorIfFileDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewPipeline(
		RmFile("/does/not/exist/and/never/will"),
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

func TestRmFileDoesNotRespectFilePermissions(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we need a valid filename to try to write to
	tmpFilename, err := ExecPipeline(MkTempFile(os.TempDir(), "scriptify-")).TrimmedString()
	assert.Nil(t, err)

	// make sure the file cannot be deleted
	err = ExecPipeline(Chmod(tmpFilename, 0)).Error()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(
		Chmod(tmpFilename, 0644),
		RmFile(tmpFilename),
	)

	pipeline := NewPipeline(
		RmFile(tmpFilename),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, pipe.OK, pipeline.StatusCode)
	assert.Empty(t, actualResult)
}