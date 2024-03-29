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

func TestShoptTracingIsDisabledByDefault(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := false

	// ----------------------------------------------------------------
	// perform the change

	actualResult1 := IsTraceEnabled()
	actualResult2 := GetShellOptions().IsTraceEnabled()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult1)
	assert.Equal(t, expectedResult, actualResult2)
}

func TestShoptTracingCanBeEnabled(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := true

	// ----------------------------------------------------------------
	// perform the change

	GetShellOptions().EnableTrace(NewTextBuffer())

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	actualResult1 := IsTraceEnabled()
	actualResult2 := GetShellOptions().IsTraceEnabled()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult1)
	assert.Equal(t, expectedResult, actualResult2)
}

func TestShoptTracingCanBeDisabled(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := false
	GetShellOptions().EnableTrace(NewTextBuffer())

	// ----------------------------------------------------------------
	// perform the change

	GetShellOptions().DisableTrace()

	actualResult1 := IsTraceEnabled()
	actualResult2 := GetShellOptions().IsTraceEnabled()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult1)
	assert.Equal(t, expectedResult, actualResult2)
}

func TestShoptTracefWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := "this is my test output"
	expectedResult := "+ " + testData + "\n"
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	// ----------------------------------------------------------------
	// perform the change

	Tracef(testData)
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestShoptTraceOutputWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := "this is my test output"
	expectedResult := "+ test> " + testData + "\n"
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	// ----------------------------------------------------------------
	// perform the change

	TraceOutput("test", "%s", testData)
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestShoptTraceOsStderrWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := "this is my test output"
	expectedResult := "+ os.Stderr> " + testData + "\n"
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	// ----------------------------------------------------------------
	// perform the change

	TraceOsStderr("%s", testData)
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestShoptTraceOsStdoutWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := "this is my test output"
	expectedResult := "+ os.Stdout> " + testData + "\n"
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	// ----------------------------------------------------------------
	// perform the change

	TraceOsStdout("%s", testData)
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestShoptTracePipeStderrWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := "this is my test output"
	expectedResult := "+ p.Stderr> " + testData + "\n"
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	// ----------------------------------------------------------------
	// perform the change

	TracePipeStderr("%s", testData)
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestShoptTracePipeStdoutWritesToTheTraceOutput(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	testData := "this is my test output"
	expectedResult := "+ p.Stdout> " + testData + "\n"
	dest := NewTextBuffer()
	GetShellOptions().EnableTrace(dest)

	// clean up after ourselves
	defer GetShellOptions().DisableTrace()

	// ----------------------------------------------------------------
	// perform the change

	TracePipeStdout("%s", testData)
	actualResult := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
