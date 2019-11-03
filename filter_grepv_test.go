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

func TestGrepVDropsLinesThatMatchTheRegex(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"this is the first line",
		"this is the second line",
		"this is the third line",
		"and this is the fourth line",
	}
	expectedResult := []string{
		"this is the first line",
		"and this is the fourth line",
	}

	pipeline := NewPipeline(
		EchoSlice(testData),
		GrepV("second|third"),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestGrepVReturnsErrorIfRegexInvalid(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"this is the first line",
		"this is the second line",
		"this is the third line",
		"and this is the fourth line",
	}
	expectedResult := []string{}

	pipeline := NewPipeline(
		EchoSlice(testData),
		GrepV("[* "),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestGrepVWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `+ EchoSlice([]string{"this is the first line", "this is the second line", "this is the third line", "and this is the fourth line"})
+ p.Stdout> this is the first line
+ p.Stdout> this is the second line
+ p.Stdout> this is the third line
+ p.Stdout> and this is the fourth line
+ GrepV("$1")
+ => GrepV("second|third")
+ p.Stdout> this is the first line
+ p.Stdout> and this is the fourth line
`
	dest := NewDest()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	testData := []string{
		"this is the first line",
		"this is the second line",
		"this is the third line",
		"and this is the fourth line",
	}

	pipeline := NewPipeline(
		EchoSlice(testData),
		GrepV("$1"),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec("second|third")
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
