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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXargsDirnameStripsFinalElement(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"/path/to/one/file",
		"/path/to/two/file",
	}
	expectedResult := []string{
		"/path/to/one",
		"/path/to/two",
	}

	pipeline := NewPipeline(
		EchoSlice(testData),
		XargsDirname(),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestXargsDirnamePreservesRootFolder(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// Golang's filepath behaves differently on different platforms
	// We develop and test on Linux. This prevents us testing Windows paths
	// here :(
	testData := []string{
		"/",
	}
	expectedResult := []string{
		"/",
	}

	pipeline := NewPipeline(
		EchoSlice(testData),
		XargsDirname(),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestXargsDirnameReplacesBlankLinesWithDots(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"/path/to/one/file",
		"/path/to/two/file",
		"",
	}
	expectedResult := []string{
		"/path/to/one",
		"/path/to/two",
		".",
	}

	pipeline := NewPipeline(
		EchoSlice(testData),
		XargsDirname(),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestXargsDirnameWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `+ EchoSlice([]string{"/path/to/one/file", "/path/to/two/file"})
+ p.Stdout> /path/to/one/file
+ p.Stdout> /path/to/two/file
+ XargsDirname()
+ p.Stdout> /path/to/one
+ p.Stdout> /path/to/two
`
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	testData := []string{
		"/path/to/one/file",
		"/path/to/two/file",
	}

	pipeline := NewPipeline(
		EchoSlice(testData),
		XargsDirname(),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
