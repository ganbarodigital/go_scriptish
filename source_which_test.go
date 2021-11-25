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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhichWritesTheMatchingPathToThePipeline(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewPipeline(
		Which("bash"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.True(t, strings.HasSuffix(actualResult, "/bash"))
}

func TestWhichReturnsAnErrorIfNoMatchFound(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewPipeline(
		Which("does-not-exist-and-never-will"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Empty(t, actualResult)
}

func TestWhichWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	pipeline := NewPipeline(
		Which("$1"),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec("bash")
	actualResult := dest.String()

	path, err := pipeline.TrimmedString()
	assert.Nil(t, err)

	expectedResult := `+ Which("$1")
+ => Which("bash")
+ p.Stdout> ` + path + `
`

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
