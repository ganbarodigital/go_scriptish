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

func TestAndDoesNotExecuteTheSecondSequenceIfFirstSequenceReturnedAnError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\n"
	expectedResult := 100

	list := NewList(
		Return(expectedResult),
		And(
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

func TestAndExecutesTheSecondSequenceIfFirstSequenceExitedSuccessfully(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"
	list := NewList(
		Return(0),
		And(
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

func TestAndExecutesTheThirdSequenceIfEarlierSequencesDoExitSuccessfully(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\nthe sun is shining\n"
	list := NewList(
		Echo("hello world"),
		And(
			NewPipeline(
				Echo("have a nice day"),
			),
		),
		And(
			NewPipeline(
				Echo("the sun is shining"),
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

func TestAndDoesNotExecuteTheThirdSequenceIfEarlierSequencesDoNotExitSuccessfully(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\n"
	pipeline := NewList(
		Return(0),
		And(
			NewPipeline(
				Return(100),
			),
		),
		And(
			NewPipeline(
				Echo(testData),
			),
		),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Empty(t, actualResult)
	assert.Equal(t, 100, pipeline.StatusCode())
}

func TestAndReturnsStatusCodeFromLastSequenceExecuted(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewList(
		Return(100),
		And(
			NewPipeline(
				Return(1),
			),
		),
		And(
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

func TestAndReturnsStdoutFromLastSequenceExecuted(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "everything is fine\n"
	pipeline := NewList(
		Return(0),
		And(
			NewPipeline(
				Echo(expectedResult),
				Return(5),
			),
		),
		And(
			NewPipeline(
				Return(1),
			),
		),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, 5, pipeline.StatusCode())
	assert.Equal(t, expectedResult, actualResult)
}

func TestLogicAndWritesToTheTraceOutputIfNotExecutingNextSequence(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `+ Return(1)
+ status code: 1
+ error: command exited with non-zero status code 1
+ And(): not executing the given sequence
+ status code: 1
+ error: command exited with non-zero status code 1
`
	dest := NewDest()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	list := NewList(
		Return(1),
		And(
			NewPipeline(
				Echo("hello world"),
			),
		),
	)

	// ----------------------------------------------------------------
	// perform the change

	list.Exec()
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestLogicAndWritesToTheTraceOutputIfExecutingNextSequence(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `+ Return(0)
+ And(): executing the given sequence
+ Echo("hello world")
+ => Echo("hello world")
+ p.Stdout> hello world
`
	dest := NewDest()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	list := NewList(
		Return(0),
		And(
			NewPipeline(
				Echo("hello world"),
			),
		),
	)

	// ----------------------------------------------------------------
	// perform the change

	list.Exec()
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
