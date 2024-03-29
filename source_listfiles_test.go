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
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListFilesReturnsPathIfFile(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we're going to use our own source code as the test data
	_, filename, _, _ := runtime.Caller(0)

	expectedResult := filename + "\n"
	pipeline := NewPipeline(
		ListFiles(filename),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, _ := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestListFilesReturnsContentsIfPathIsFolder(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := []string{
		"testdata/listfiles/one.txt",
		"testdata/listfiles/three.yaml",
		"testdata/listfiles/two.txt",
	}

	pipeline := NewPipeline(
		ListFiles("./testdata/listfiles/"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, _ := pipeline.Exec().Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestListFilesReturnsMatchesIfPathContainsWildcards(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	expectedResult := []string{
		"testdata/listfiles/one.txt",
		"testdata/listfiles/two.txt",
	}

	pipeline := NewPipeline(
		ListFiles("./testdata/listfiles/*.txt"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, _ := pipeline.Exec().Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestListFilesReturnsErrIfPathDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewPipeline(
		ListFiles("./testdata/does/not/exist/"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "", actualResult)
	assert.NotNil(t, err)
}

func TestListFilesReturnsErrIfPathIsFolderThatCannotBeRead(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// make the folder unreadable
	os.Chmod("./testdata/listfiles_err", 0)

	// best put the folder permissions back again afterwards
	defer os.Chmod("./testdata/listfiles_err", 0755)

	pipeline := NewPipeline(
		ListFiles("./testdata/listfiles_err/"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "", actualResult)
	assert.NotNil(t, err)
}

func TestListFilesReturnsErrorIfPathWildcardsPatternInvalid(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewPipeline(
		ListFiles("./testdata/listfiles/[]"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "", actualResult)
	assert.NotNil(t, err)
}

func TestListFilesReturnsNothingIfPathWildcardsDoNotMatch(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewPipeline(
		ListFiles("./testdata/listfiles/*.does_not_exist"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "", actualResult)
	assert.Nil(t, err)
}

func TestListFilesWritesToTheTraceOutputIfPathToFolder(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `+ ListFiles("./testdata/listfiles/")
+ p.Stdout> testdata/listfiles/one.txt
+ p.Stdout> testdata/listfiles/three.yaml
+ p.Stdout> testdata/listfiles/two.txt
`
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	pipeline := NewPipeline(
		ListFiles("./testdata/listfiles/"),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestListFilesWritesToTheTraceOutputIfGlobbing(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `+ ListFiles("./testdata/listfiles/*")
+ p.Stdout> testdata/listfiles/one.txt
+ p.Stdout> testdata/listfiles/three.yaml
+ p.Stdout> testdata/listfiles/two.txt
`
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	pipeline := NewPipeline(
		ListFiles("./testdata/listfiles/*"),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
