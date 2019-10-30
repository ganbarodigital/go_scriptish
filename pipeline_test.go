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

	pipe "github.com/ganbarodigital/go_pipe/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewPipelineCreatesEmptyPipeline(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "", pipeline.Pipe.Stdin.String())
	assert.Equal(t, "", pipeline.Pipe.Stdout.String())
	assert.Equal(t, "", pipeline.Pipe.Stderr.String())
}

func ExecPipelineCreatesAndRunsAPipeline(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"

	// ----------------------------------------------------------------
	// perform the change

	pipeline := ExecPipeline(Echo("hello world!"))

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "", pipeline.Pipe.Stdin.String())
	assert.Equal(t, expectedResult, pipeline.Pipe.Stdout.String())
	assert.Equal(t, "", pipeline.Pipe.Stderr.String())

	assert.Nil(t, pipeline.Error())
}

func PipelineFuncCreatesAPipelineAsAFunction(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipelineFunc(Echo("hello world!"))()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "", pipeline.Pipe.Stdin.String())
	assert.Equal(t, expectedResult, pipeline.Pipe.Stdout.String())
	assert.Equal(t, "", pipeline.Pipe.Stderr.String())

	assert.Nil(t, pipeline.Error())
}

func ExecCopesWithANilPipelinePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline *Pipeline
	pipeline = nil

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	// as long as the test doesn't crash, it has passed
}

func ExecCopesWithAnEmptyPipelineStruct(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline Pipeline

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	// as long as the test doesn't crash, it has passed
}

func TestNewPipelineCreatesPipeWithEmptyStdin(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline()
	actualResult := pipeline.Pipe.Stdin.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipelineCreatesPipeWithEmptyStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline()
	actualResult := pipeline.Pipe.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipelineCreatesPipeWithEmptyStderr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline()
	actualResult := pipeline.Pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipelineCreatesPipelineWithNilErrSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, pipeline.Error())
}

func TestNewPipelineCreatesPipelineWithZeroStatusCode(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, 0, pipeline.StatusCode())
}

func TestPipelineControllerCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence
	controller := PipelineController

	// ----------------------------------------------------------------
	// perform the change

	controller(sequence)

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestPipelineControllerCopesWithEmptySequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence Sequence
	sequence.Controller = PipelineController(&sequence)

	// ----------------------------------------------------------------
	// perform the change

	sequence.Exec()

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestPipelineExecRunsAllStepsInOrder(t *testing.T) {
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
		// copy what op1 did first
		p.DrainStdinToStdout()

		// add our own content
		p.Stdout.WriteString("have a nice day")
		p.Stdout.WriteRune('\n')

		// all done
		return 0, nil
	}

	pipeline := NewPipeline(op1, op2)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineExecStopsWhenAStepReportsAnError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedStdout := "hello world\n"
	expectedStderr := "alfred the great\n"
	op1 := func(p *Pipe) (int, error) {
		p.Stdout.WriteString(expectedStdout)
		p.Stderr.WriteString(expectedStderr)

		// all done
		return 0, errors.New("stop at step 1")
	}
	op2 := func(p *Pipe) (int, error) {
		// copy what op1 did first
		p.DrainStdinToStdout()

		// add our own content
		p.Stdout.WriteString("have a nice day")
		p.Stdout.WriteRune('\n')

		// all done
		return 0, nil
	}

	pipeline := NewPipeline(op1, op2)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	finalOutput, err := pipeline.String()
	actualStderr := pipeline.Pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	// pipeline.String() should have returned an error
	assert.NotNil(t, err)

	// pipeline.String() should have returned the contents of our
	// Pipe.Stdout buffer
	assert.Equal(t, expectedStdout, finalOutput)

	// our pipeline's Stderr should still contain what the first step
	// did ... and only the first step
	assert.Equal(t, expectedStderr, actualStderr)
}

func TestPipelineExecSetsErrWhenOpReturnsNonZeroStatusCodeAndNilErr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *pipe.Pipe) (int, error) {
		// fail, but without an error to say why
		return StatusNotOkay, nil
	}

	pipeline := NewPipeline(op1)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	// pipeline.Err should have been set by Exec()
	assert.NotNil(t, pipeline.Error())
	_, ok := pipeline.Error().(pipe.ErrNonZeroStatusCode)
	assert.True(t, ok)
}
