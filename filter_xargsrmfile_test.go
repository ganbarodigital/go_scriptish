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

func TestXargsRmFile(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we need somewhere to create our temporary files
	tmpDir, err := ExecPipeline(MkTempDir(os.TempDir(), "scriptify-")).TrimmedString()
	assert.Nil(t, err)

	// clean up after ourselves
	defer ExecPipeline(RmDir(tmpDir))

	// we need some files to remove
	var testData []string
	for i := 0; i < 5; i++ {
		tmpFile, err := ExecPipeline(MkTempFile(tmpDir, "rmtest-")).TrimmedString()
		assert.Nil(t, err)
		testData = append(testData, tmpFile)
	}

	// prove that the files exist
	listingFunc := PipelineFunc(ListFiles(tmpDir))
	expectedResult, err := listingFunc().Strings()
	assert.Nil(t, err)
	assert.Equal(t, 5, len(expectedResult))

	// this is the pipeline we'll use to test XargsRmFile()
	pipeline := NewPipeline(
		ListFiles(tmpDir),
		XargsRmFile(),
	)

	// ----------------------------------------------------------------
	// perform the change

	// this will delete the files
	actualResult, err := pipeline.Exec().Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, pipeline.Err)
	assert.Equal(t, expectedResult, actualResult)

	// the folder should now be empty
	listing, err := listingFunc().Strings()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(listing))
}

func TestXargsRmFileSetsErrorWhenSomethingGoesWrong(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	pipeline := NewPipeline(
		Echo("/does/not/exist"),
		XargsRmFile(),
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
