// scriptish is a library to help you port bash scripts to Golang
//
// inspired by:
//
// - http://labix.org/pipe
// - https://github.com/bitfield/script
//
// Copyright 2021-present Ganbaro Digital Ltd
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

func TestRedirectStdoutToTextReaderWriterSendsOutputToGivenDestination(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data to send
	testData := "this is a test"

	// we need somewhere to send the output to
	dest := NewTextBuffer()

	// this is the pipeline that we will test
	unit := NewPipeline(
		Echo(testData, RedirectStdoutToTextReaderWriter(dest)),
	)

	expectedResult := testData + "\n"

	// ----------------------------------------------------------------
	// perform the change

	pipelineOutput, err := unit.Exec().String()
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	// the pipeline should have executed without error
	assert.Nil(t, err)

	// the pipeline's original stdout should be empty
	assert.Empty(t, pipelineOutput)

	// our testData should have been written to our preferred output
	// destination
	assert.Equal(t, expectedResult, actualResult)
}

func TestRedirectStdoutToTextReaderWriterWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data to send
	testData := "this is a test"

	// we need somewhere to send the output to
	dest := NewTextBuffer()

	// this is the pipeline that we will test
	unit := NewPipeline(
		Echo(testData, RedirectStdoutToTextReaderWriter(dest)),
	)

	// we need to enable tracing output
	traceBuf := NewTextBuffer()
	GetShellOptions().EnableTrace(traceBuf)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	expectedResult := `+ RedirectStdoutToTextReaderWriter()
+ Echo("this is a test")
+ => Echo("this is a test")
+ p.Stdout> this is a test
`

	// ----------------------------------------------------------------
	// perform the change

	unit.Exec()
	actualResult := traceBuf.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
