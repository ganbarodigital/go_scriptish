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

func TestPipelineControllerWritesNothingToTheTraceOutputIfStatusCodeIsZero(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	op1 := NewSequenceStep(
		func(p *Pipe) (int, error) {
			return StatusOkay, nil
		},
	)

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline(
		op1,
	)
	pipeline.Exec()
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineControllerWritesNonZeroStatusCodesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	op1 := NewSequenceStep(
		func(p *Pipe) (int, error) {
			return StatusNotOkay, nil
		},
	)

	expectedResult := `+ status code: 1
+ error: command exited with non-zero status code 1
`

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline(
		op1,
	)
	pipeline.Exec()
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineControllerWritesErrorsToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	op1 := NewSequenceStep(
		func(p *Pipe) (int, error) {
			return StatusNotOkay, errors.New("this is a test error")
		},
	)

	expectedResult := `+ status code: 1
+ error: this is a test error
`

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline(
		op1,
	)
	pipeline.Exec()
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
