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

func TestOrDoesNotExecuteTheSecondSequenceIfFirstSequenceFinishesSuccessfully(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"
	pipeline := NewList(
		Echo(expectedResult),
		Or(
			NewPipeline(
				Return(100),
			),
		),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, StatusOkay, pipeline.StatusCode())
}

func TestOrExecutesTheSecondSequenceIfFirstSequenceDoesNotExitSuccessfully(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"
	pipeline := NewList(
		Return(100),
		Or(
			NewPipeline(
				Echo(expectedResult),
			),
		),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, StatusOkay, pipeline.StatusCode())
}

func TestOrExecutesTheThirdSequenceIfEarlierSequencesDoNotExitSuccessfully(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"
	pipeline := NewList(
		Return(100),
		Or(
			NewPipeline(
				Return(1),
			),
		),
		Or(
			NewPipeline(
				Echo(expectedResult),
			),
		),
	)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Exec().String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, StatusOkay, pipeline.StatusCode())
}

func TestOrReturnsStatusCodeFromLastSequenceExecuted(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewList(
		Return(100),
		Or(
			NewPipeline(
				Return(1),
			),
		),
		Or(
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
	assert.Equal(t, 5, pipeline.StatusCode())
}

func TestOrReturnsStdoutFromLastSequenceExecuted(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "everything is fine\n"
	pipeline := NewList(
		Return(100),
		Or(
			NewPipeline(
				Return(1),
			),
		),
		Or(
			NewPipeline(
				Echo(expectedResult),
				Return(5),
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

func TestOrWritesToTheTraceOutputWhenExecutingTheOrBranch(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `+ Return(1)
+ status code: 1
+ error: command exited with non-zero status code 1
+ Or(): executing the given sequence
+ Echo("this is the 'or' option")
+ => Echo("this is the 'or' option")
+ p.Stdout> this is the 'or' option
`
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	list := NewList(
		Return(1),
		Or(
			NewList(
				Echo("this is the 'or' option"),
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

func TestOrWritesToTheTraceOutputWhenNotExecutingTheOrBranch(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `+ Return(0)
+ Or(): not executing the given sequence
`
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	list := NewList(
		Return(0),
		Or(
			NewList(
				Echo("this is the 'or' option"),
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
