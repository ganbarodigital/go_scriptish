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

	"github.com/stretchr/testify/assert"
)

func TestRmDirRemovesAGivenEmptyFolder(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	tmpDir, err := ExecPipeline(
		// create the temporary folder
		MkTempDir(os.TempDir(), "scriptify-rmdir-"),
	).TrimmedString()
	assert.Nil(t, err)
	assert.NotEmpty(t, tmpDir)

	// this pipeline will prove if the temporary file does/does not exist
	tmpDirExists := PipelineFunc(FilepathExists(tmpDir))
	dirExists := tmpDirExists().Okay()
	assert.True(t, dirExists)

	pipeline := NewPipeline(
		RmDir(tmpDir),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Empty(t, actualResult)

	// the file should be gone
	dirExists = tmpDirExists().Okay()
	assert.False(t, dirExists)
}

func TestRmDirSetsErrorIfDirDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewPipeline(
		RmDir("./does/not/exist/and/never/will"),
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

func TestRmDirDoesNotRespectFilePermissions(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we need a valid folder to remove
	tmpDir, err := ExecPipeline(MkTempDir(os.TempDir(), "scriptify-rmfile-")).TrimmedString()
	assert.Nil(t, err)

	// make sure the dir cannot be deleted
	err = ExecPipeline(Chmod(tmpDir, 0)).Error()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(
		Chmod(tmpDir, 0644),
		RmFile(tmpDir),
	)

	pipeline := NewPipeline(
		RmFile(tmpDir),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, StatusOkay, pipeline.StatusCode())
	assert.Empty(t, actualResult)
}
