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

func TestIfDoesNotExecuteTheSecondSequenceIfFirstSequenceReturnedAnError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\n"
	expectedResult := 100

	list := NewList(
		If(
			NewList(
				Return(expectedResult),
			),
			NewPipeline(
				Echo(testData),
			),
		),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := list.Exec().StatusError()

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
	assert.Equal(t, expectedResult, actualResult)
	actualStdout, _ := list.String()
	assert.Empty(t, actualStdout)
}

func TestListExecutesTheSecondSequenceIfFirstSequenceExitedSuccessfully(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"
	list := NewList(
		If(
			NewList(
				Return(0),
			),
			NewPipeline(
				Echo(expectedResult),
			),
		),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := list.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, StatusOkay, list.StatusCode())
}

func TestIfReturnsStatusCodeFromLastSequenceExecuted(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewList(
		If(
			NewList(
				Return(100),
			),
			NewPipeline(
				Echo("hello world"),
				Return(5),
			),
		),
	)

	// ----------------------------------------------------------------
	// perform the change

	err := pipeline.Exec().Error()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, 100, pipeline.StatusCode())
}

func TestIfReturnsStdoutFromLastSequenceExecuted(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "everything is fine\n"
	list := NewList(
		If(
			NewList(
				Return(0),
			),
			NewPipeline(
				Echo(expectedResult),
				Return(5),
			),
		),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := list.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, 5, list.StatusCode())
	assert.Equal(t, expectedResult, actualResult)
}
