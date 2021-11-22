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

func TestChmodChangesPermissionsOnGivenFile(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we need a file to mess with
	tmpFile, err := ExecPipeline(
		MkTempFile(os.TempDir(), "scriptish-chmod-"),
	).TrimmedString()
	assert.Nil(t, err)
	assert.NotEmpty(t, tmpFile)

	// clean up after ourselves
	defer ExecPipeline(RmFile(tmpFile))

	// grab its current permissions
	//
	// the actual value of these permissions may depend upon
	// the local umask
	origMode, err := ExecPipeline(Lsmod(tmpFile)).TrimmedString()
	assert.Nil(t, err)
	assert.NotEmpty(t, origMode)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := ExecPipeline(
		Chmod(tmpFile, 0),
	).StatusError()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, StatusOkay, actualResult)

	// the file should have different permissions now
	// grab its current permissions
	newMode, err := ExecPipeline(Lsmod(tmpFile)).TrimmedString()
	assert.Nil(t, err)
	assert.Equal(t, "----------", newMode)
	assert.NotEqual(t, origMode, newMode)
}

func TestChmodSetsErrorIfFileDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewPipeline(
		Chmod("./does/not/exist/and/never/will", 0644),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().StatusError()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, StatusNotOkay, actualResult)
}

func TestChmodWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we use string expansion here to prove that our trace includes
	// the expansion
	expectedResult := `+ Chmod("$1", 0644)
+ => Chmod("./builtin_chmod_test.go", 0644)
`
	dest := NewDest()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline(
		Chmod("$1", 0644),
	)
	pipeline.Exec("./builtin_chmod_test.go")
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestChmodErrorsAppearInTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we use string expansion here to prove that our trace includes
	// the expansion
	expectedResult := `+ Chmod("$1", 0644)
+ => Chmod("./does/not/exist", 0644)
+ status code: 1
+ error: chmod ./does/not/exist: no such file or directory
`
	dest := NewDest()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline(
		Chmod("$1", 0644),
	)
	pipeline.Exec("./does/not/exist")
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
