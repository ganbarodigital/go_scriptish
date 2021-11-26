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

func TestRedirectStdoutToStderrSendsWritesToStderr(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// this pipeline contains the code that we are testing
	unit := NewPipeline(
		CatFile(
			"./testdata/truncatefile/content.txt",
			RedirectStdoutToStderr(),
		),
	)

	// let's grab the contents of the file, to compare later
	expectedResult, _ := ExecPipeline(
		CatFile("./testdata/truncatefile/content.txt"),
	).String()
	assert.NotEmpty(t, expectedResult)

	// ----------------------------------------------------------------
	// perform the change

	actualStatus, actualErr := unit.Exec().StatusError()
	actualResult := unit.Pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	// make sure the pipeline ran with no errors
	assert.Nil(t, actualErr)
	assert.Equal(t, StatusOkay, actualStatus)

	// stdout should be empty
	assert.Empty(t, unit.Pipe.Stdout.String())

	// stderr should contain the data we expect
	assert.Equal(t, expectedResult, actualResult)
}

func TestRedirectStdoutToStderrWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `+ RedirectStdoutToStderr()
+ Echo("this is a test")
+ => Echo("this is a test")
+ p.Stdout> this is a test
`
	traceBuf := NewTextBuffer()
	GetShellOptions().EnableTrace(traceBuf)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	pipeline := NewPipeline(
		Echo("this is a test", RedirectStdoutToStderr()),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()
	actualResult := traceBuf.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestRedirectStdoutToStderrAllowsMixedWrites(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `this is data written to stdout
this is data written to stderr
and this is more data written to stdout
`
	list := NewList(
		Echo("this is data written to stdout", RedirectStdoutToStderr()),
		EchoToStderr("this is data written to stderr"),
		Echo("and this is more data written to stdout", RedirectStdoutToStderr()),
	)

	// ----------------------------------------------------------------
	// perform the change

	statusCode, err := list.Exec().StatusError()

	actualStdoutResult := list.Pipe.Stdout.String()
	actualStderrResult := list.Pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Zero(t, statusCode)
	assert.Empty(t, actualStdoutResult)
	assert.Equal(t, expectedResult, actualStderrResult)
}
