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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMkTempfilename(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// this is the pipeline we'll use to test MkTempFilename()
	pipeline := NewPipeline(
		MkTempFilename(os.TempDir(), ""),
	)

	// ----------------------------------------------------------------
	// perform the change

	// this will create the filename
	actualResult, err := pipeline.Exec().TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(actualResult, os.TempDir()), actualResult)
}

func TestMkTempfilenameGeneratesFilenameThatDoesNotExist(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// this is the pipeline we'll use to test MkTempFilename()
	pipeline := NewPipeline(
		MkTempFilename(os.TempDir(), ""),
	)

	// ----------------------------------------------------------------
	// perform the change

	// this will truncate the file
	tmpFilename, err := pipeline.Exec().TrimmedString()
	assert.Nil(t, err)

	fileExists, err := ExecPipeline(TestFilepathExists(tmpFilename)).StatusError()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, "stat "+tmpFilename+": no such file or directory", err.Error())
	assert.Equal(t, 1, fileExists)
}
