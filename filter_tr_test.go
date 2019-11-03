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

func TestTrReplacesStrings(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"this is the first line of test data",
		"this is the second line of test data",
		"this is the third line of test data",
		"this is the fourth line of test data",
		"this is the fifth line of test data",
		"",
		"this is the seventh line of test data",
	}
	expectedResult := []string{
		"sith is the first line of trout data",
		"sith is the second line of trout data",
		"sith is the third line of trout data",
		"sith is the fourth line of trout data",
		"sith is the fifth line of trout data",
		"",
		"sith is the seventh line of trout data",
	}

	pipeline := NewPipeline(
		EchoSlice(testData),
		Tr([]string{"this", "test"}, []string{"sith", "trout"}),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestTrSupportsSingleNewString(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"this is the first line of test data",
		"this is the second line of test data",
		"this is the third line of test data",
		"this is the fourth line of test data",
		"this is the fifth line of test data",
		"",
		"this is the seventh line of test data",
	}
	expectedResult := []string{
		"trout is the first line of trout data",
		"trout is the second line of trout data",
		"trout is the third line of trout data",
		"trout is the fourth line of trout data",
		"trout is the fifth line of trout data",
		"",
		"trout is the seventh line of trout data",
	}

	pipeline := NewPipeline(
		EchoSlice(testData),
		Tr([]string{"this", "test"}, []string{"trout"}),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestTrReturnsAnErrorIfParamsNotSameLength(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"this is the first line of test data",
		"this is the second line of test data",
		"this is the third line of test data",
		"this is the fourth line of test data",
		"this is the fifth line of test data",
		"",
		"this is the seventh line of test data",
	}
	expectedResult := []string{}

	pipeline := NewPipeline(
		EchoSlice(testData),
		Tr([]string{"this"}, []string{"sith", "trout"}),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	_, ok := err.(ErrMismatchedInputs)
	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestTrWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `+ EchoSlice([]string{"this is the first line of test data", "this is the second line of test data"})
+ p.Stdout> this is the first line of test data
+ p.Stdout> this is the second line of test data
+ Tr([]string{"this", "test"}, []string{"sith", "trout"})
+ p.Stdout> sith is the first line of trout data
+ p.Stdout> sith is the second line of trout data
`
	dest := NewDest()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	testData := []string{
		"this is the first line of test data",
		"this is the second line of test data",
	}

	pipeline := NewPipeline(
		EchoSlice(testData),
		Tr([]string{"this", "test"}, []string{"sith", "trout"}),
	)
	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
