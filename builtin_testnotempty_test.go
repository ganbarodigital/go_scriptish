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

func TestTestNotEmptyReturnsOneIfStringIsEmpty(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := ""
	expectedResult := 1

	pipeline := NewPipeline(
		TestNotEmpty(testData),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipeline.Exec().StatusCode()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestTestNotEmptyReturnsOneIfStringExpandsToEmpty(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := "$DOES_NOT_EXIST"
	expectedResult := 1

	pipeline := NewPipeline(
		TestNotEmpty(testData),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipeline.Exec().StatusCode()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestTestNotEmptyReturnsZeroIfStringIsNotEmpty(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := "not empty"
	expectedResult := 0

	pipeline := NewPipeline(
		TestNotEmpty(testData),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipeline.Exec().StatusCode()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestTestNotEmptyReturnsZeroIfStringExpandsToNotEmpty(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := "$DOES_NOT_EXIST"
	expectedResult := 0

	pipeline := NewPipeline(
		TestNotEmpty(testData),
	)
	pipeline.LocalVars.Setenv("DOES_NOT_EXIST", "yes it does")

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipeline.Exec().StatusCode()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestTestNotEmptyWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we use string expansion here to prove that our trace includes
	// the expansion
	expectedResult := `+ TestNotEmpty("$1")
+ => TestNotEmpty("this is an expanded string")
`
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline(
		TestNotEmpty("$1"),
	)
	pipeline.Exec("this is an expanded string")
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestTestNotEmptyErrorsAppearInTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we use string expansion here to prove that our trace includes
	// the expansion
	expectedResult := `+ TestNotEmpty("$1")
+ => TestNotEmpty("")
+ status code: 1
+ error: command exited with non-zero status code 1
`
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline(
		TestNotEmpty("$1"),
	)
	pipeline.Exec()
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
