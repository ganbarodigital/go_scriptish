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
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecReadsFromPipelineStdin(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"
	pipeline := NewPipeline(
		Echo("hello world"),
		Exec("/usr/bin/env", "cat"),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, pipeline.Error())
	assert.Equal(t, "", pipeline.Pipe.Stdin.String())
	assert.Equal(t, expectedResult, pipeline.Pipe.Stdout.String())
	assert.Equal(t, "", pipeline.Pipe.Stderr.String())
}

func TestExecWritesToStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"
	pipeline := NewPipeline(
		Exec("/usr/bin/env", "bash", "-c", "echo \"hello world\""),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, pipeline.Error())
	assert.Equal(t, "", pipeline.Pipe.Stdin.String())
	assert.Equal(t, expectedResult, pipeline.Pipe.Stdout.String())
	assert.Equal(t, "", pipeline.Pipe.Stderr.String())
}

func TestExecWritesToStderr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\n"
	pipeline := NewPipeline(
		Exec("/usr/bin/env", "bash", "-c", "echo \"hello world\" 1>&2"),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, pipeline.Error())
	assert.Equal(t, "", pipeline.Pipe.Stdin.String())
	assert.Equal(t, "", pipeline.Pipe.Stdout.String())
	assert.Equal(t, expectedResult, pipeline.Pipe.Stderr.String())
}

func TestExecCapturesTheStatusCode(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := 5
	pipeline := NewPipeline(
		Exec("/usr/bin/env", "bash", "-c", "exit 5"),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	// the pipeline will have an exec.ExitError() error status, because
	// we didn't exit with a status code of 0
	err := pipeline.Error()
	assert.NotNil(t, err)
	_, ok := err.(*exec.ExitError)
	assert.True(t, ok)
	assert.Equal(t, expectedResult, pipeline.StatusCode())
}

func TestExecSetsErrorIfCommandCannotBeRun(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := 1
	pipeline := NewPipeline(
		Exec("/does/not/exist"),
	)

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	// the pipeline will have an os.PathError() error status, because
	// we gave it a file that could not be found
	err := pipeline.Error()
	assert.NotNil(t, err)
	_, ok := err.(*os.PathError)
	assert.True(t, ok)
	assert.Equal(t, expectedResult, pipeline.StatusCode())
}
