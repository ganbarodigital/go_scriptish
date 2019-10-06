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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewListCreatesEmptyList(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	list := NewList()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "", list.Pipe.Stdin.String())
	assert.Equal(t, "", list.Pipe.Stdout.String())
	assert.Equal(t, "", list.Pipe.Stderr.String())
}

func ExecListCreatesAndRunsAList(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"

	// ----------------------------------------------------------------
	// perform the change

	list := ExecList(Echo("hello world!"))

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "", list.Pipe.Stdin.String())
	assert.Equal(t, expectedResult, list.Pipe.Stdout.String())
	assert.Equal(t, "", list.Pipe.Stderr.String())

	assert.Nil(t, list.Error())
}

func ListFuncCreatesAListAsAFunction(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"

	// ----------------------------------------------------------------
	// perform the change

	list := NewListFunc(Echo("hello world!"))()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "", list.Pipe.Stdin.String())
	assert.Equal(t, expectedResult, list.Pipe.Stdout.String())
	assert.Equal(t, "", list.Pipe.Stderr.String())

	assert.Nil(t, list.Error())
}

func ListExecCopesWithANilListPointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var list *List = nil

	// ----------------------------------------------------------------
	// perform the change

	list.Exec()

	// ----------------------------------------------------------------
	// test the results

	// as long as the test doesn't crash, it has passed
}

func ListExecCopesWithAnEmptyListStruct(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var list List

	// ----------------------------------------------------------------
	// perform the change

	list.Exec()

	// ----------------------------------------------------------------
	// test the results

	// as long as the test doesn't crash, it has passed
}

func TestNewListCreatesPipeWithEmptyStdin(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	list := NewList()
	actualResult := list.Pipe.Stdin.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewListCreatesPipeWithEmptyStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	list := NewList()
	actualResult := list.Pipe.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewListCreatesPipeWithEmptyStderr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	list := NewList()
	actualResult := list.Pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewListCreatesPipelineWithNilErrSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	list := NewList()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, list.Error())
}

func TestNewListCreatesPipeWithZeroStatusCode(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	list := NewList()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, 0, list.StatusCode())
}

func TestListExecRunsAllStepsInOrder(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"
	op1 := func(p *Pipe) (int, error) {
		p.Stdout.WriteString("hello world")
		p.Stdout.WriteRune('\n')

		// all done
		return 0, nil
	}
	op2 := func(p *Pipe) (int, error) {
		// add our own content
		p.Stdout.WriteString("have a nice day")
		p.Stdout.WriteRune('\n')

		// all done
		return 0, nil
	}

	list := NewList(op1, op2)
	list.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := list.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestListExecContinuesWhenAStepReportsAnError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedStdout := "hello world\nhave a nice day\n"
	expectedStderr := "alfred the great\n"
	op1 := func(p *Pipe) (int, error) {
		p.Stdout.WriteString("hello world")
		p.Stdout.WriteRune('\n')
		p.Stderr.WriteString(expectedStderr)

		// all done
		return 0, errors.New("stop at step 1")
	}
	op2 := func(p *Pipe) (int, error) {
		// add our own content
		p.Stdout.WriteString("have a nice day")
		p.Stdout.WriteRune('\n')

		// all done
		return 0, nil
	}

	list := NewList(op1, op2)
	list.Exec()

	// ----------------------------------------------------------------
	// perform the change

	finalOutput, err := list.String()
	actualStderr := list.Pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)

	// pipeline.String() should have returned the contents of our
	// Pipe.Stdout buffer
	assert.Equal(t, expectedStdout, finalOutput)

	// our pipeline's Stderr should still contain what the first step
	// did ... and only the first step
	assert.Equal(t, expectedStderr, actualStderr)
}

func TestListExecReturnsLastCommandsStatusCodeAndError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedStatus := 100
	expectedError := errors.New("this is an error")

	op1 := func(p *Pipe) (int, error) {
		// fail, but without an error to say why
		return StatusNotOkay, nil
	}
	op2 := func(p *Pipe) (int, error) {
		return expectedStatus, expectedError
	}

	list := NewList(op1, op2)

	// ----------------------------------------------------------------
	// perform the change

	list.Exec()
	actualStatus, actualError := list.StatusError()

	// ----------------------------------------------------------------
	// test the results

	// pipeline.Err should have been set by Exec()
	assert.Equal(t, expectedStatus, actualStatus)
	assert.Equal(t, expectedError, actualError)
}

func TestExecListRunsAListOfSteps(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedStatus := 100
	expectedError := errors.New("this is an error")

	op1 := func(p *Pipe) (int, error) {
		// fail, but without an error to say why
		return StatusNotOkay, nil
	}
	op2 := func(p *Pipe) (int, error) {
		return expectedStatus, expectedError
	}

	list := ExecList(op1, op2)

	// ----------------------------------------------------------------
	// perform the change

	actualStatus, actualError := list.StatusError()

	// ----------------------------------------------------------------
	// test the results

	// pipeline.Err should have been set by Exec()
	assert.Equal(t, expectedStatus, actualStatus)
	assert.Equal(t, expectedError, actualError)
}

func TestNewListFuncReturnsAListAsAFunction(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedStatus := 100
	expectedError := errors.New("this is an error")

	op1 := func(p *Pipe) (int, error) {
		// fail, but without an error to say why
		return StatusNotOkay, nil
	}
	op2 := func(p *Pipe) (int, error) {
		return expectedStatus, expectedError
	}

	listFunc := NewListFunc(op1, op2)
	list := listFunc()

	// ----------------------------------------------------------------
	// perform the change

	actualStatus, actualError := list.StatusError()

	// ----------------------------------------------------------------
	// test the results

	// pipeline.Err should have been set by Exec()
	assert.Equal(t, expectedStatus, actualStatus)
	assert.Equal(t, expectedError, actualError)
}
