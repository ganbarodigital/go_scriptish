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

func TestRedirectStderrToDevNullEmptiesThePipesStderr(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	pipeline := NewPipeline(
		EchoToStderr("this is a test", RedirectStderrToDevNull()),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, pipeline.Error())
	assert.Equal(t, "", pipeline.Pipe.Stderr.String())
	assert.Equal(t, "", pipeline.Pipe.Stdout.String())
}

func TestRedirectStderrToDevNullPreservesStdoutInAPipeline(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	expectedResult := "this is a test\n"
	pipeline := NewPipeline(
		Echo(expectedResult, RedirectStderrToDevNull()),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, pipeline.Error())
	assert.Equal(t, "", pipeline.Pipe.Stderr.String())
	assert.Equal(t, expectedResult, pipeline.Pipe.Stdout.String())
}

func TestRedirectStderrToDevNullWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := `+ RedirectStderrToDevNull()
+ EchoToStderr("this is a test")
+ => EchoToStderr("this is a test")
+ p.Stderr> this is a test
`
	traceBuf := NewTextBuffer()
	GetShellOptions().EnableTrace(traceBuf)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	pipeline := NewPipeline(
		EchoToStderr("this is a test", RedirectStderrToDevNull()),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()
	actualResult := traceBuf.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
