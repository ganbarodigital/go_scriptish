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
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSequenceCreatesPipeWithEmptyStdin(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	sequence := NewSequence()
	actualResult := sequence.Pipe.Stdin.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewSequenceCreatesPipeWithEmptyStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	sequence := NewSequence()
	actualResult := sequence.Pipe.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewSequenceCreatesPipeWithEmptyStderr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	sequence := NewSequence()
	actualResult := sequence.Pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewSequenceCreatesSequenceWithNilErrSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	sequence := NewSequence()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, sequence.Error())
}

func TestNewSequenceCreatesSequenceWithZeroStatusCode(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	sequence := NewSequence()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, 0, sequence.StatusCode())
}

// helper for testing our sequence behaviour
type testOpResult struct {
	StatusCode int
	Err        error
}

func TestNewSequenceCreatesSequenceWithGivenSequenceOperations(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *Pipe) (int, error) { return 0, nil }
	op2 := func(p *Pipe) (int, error) { return 1, nil }

	// we can't compare functions directly in Go, but we can execute them
	// and compare their output
	expectedResult := []testOpResult{{0, nil}, {1, nil}}

	// ----------------------------------------------------------------
	// perform the change

	var actualResult []testOpResult

	sequence := NewSequence(op1, op2)
	for _, step := range sequence.Steps {
		statusCode, err := step(sequence.Pipe)
		actualResult = append(actualResult, testOpResult{statusCode, err})
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceExecCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence

	// ----------------------------------------------------------------
	// perform the change

	sequence.Exec()

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestSequenceExecCopesWithEmptySequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence Sequence

	// ----------------------------------------------------------------
	// perform the change

	sequence.Exec()

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestSequenceFlushCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// ----------------------------------------------------------------
	// perform the change

	sequence.Flush(stdout, stderr)

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, stdout.String())
	assert.Empty(t, stderr.String())
}

func TestSequenceFlushWritesTheContents(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	expectedStdout := "hello world\nhave a nice day\n"
	expectedStderr := "this is the stderr content\nit should be different to the stdout content\n"
	op1 := func(p *Pipe) (int, error) {
		p.Stdout.WriteString(expectedStdout)
		p.Stderr.WriteString(expectedStderr)

		// all done
		return 0, nil
	}

	sequence := NewSequence(op1)
	sequence.Pipe.RunCommand(op1)

	// ----------------------------------------------------------------
	// perform the change

	sequence.Flush(stdout, stderr)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedStdout, stdout.String())
	assert.Equal(t, expectedStderr, stderr.String())
}

func TestSequenceBytesCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualBytes, err := sequence.Bytes()
	actualResult := string(actualBytes)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceBytesCopesWithEmptySequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence Sequence
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualBytes, err := sequence.Bytes()
	actualResult := string(actualBytes)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceBytesReturnsContentsOfStdoutWhenNoError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"
	op1 := func(p *Pipe) (int, error) {
		// this is the content we want
		p.Stdout.WriteString(expectedResult)

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, nil
	}

	sequence := NewSequence(op1)
	sequence.Pipe.RunCommand(op1)

	// ----------------------------------------------------------------
	// perform the change

	actualBytes, err := sequence.Bytes()
	actualResult := string(actualBytes)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceBytesReturnsContentsOfStdoutWhenError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"
	op1 := func(p *Pipe) (int, error) {
		// this is the content we want
		p.Stdout.WriteString(expectedResult)

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, errors.New("an error occurred")
	}

	sequence := NewSequence(op1)
	sequence.Pipe.RunCommand(op1)

	// ----------------------------------------------------------------
	// perform the change

	actualBytes, err := sequence.Bytes()
	actualResult := string(actualBytes)

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceErrorCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence

	// ----------------------------------------------------------------
	// perform the change

	err := sequence.Error()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
}

func TestSequenceErrorCopesWithEmptySequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence Sequence

	// ----------------------------------------------------------------
	// perform the change

	err := sequence.Error()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
}

func TestSequenceErrorReturnsErrProperty(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *Pipe) (int, error) {
		// all done
		return 0, errors.New("this is an error")
	}
	expectedResult := errors.New("this is an error")

	sequence := NewSequence(op1)
	sequence.Pipe.RunCommand(op1)

	// ----------------------------------------------------------------
	// perform the change

	err := sequence.Error()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, expectedResult, err)
}

func TestSequenceOkayCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence

	// ----------------------------------------------------------------
	// perform the change

	success := sequence.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, success)
}

func TestSequenceOkayCopesWithEmptySequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence Sequence

	// ----------------------------------------------------------------
	// perform the change

	success := sequence.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, success)
}

func TestSequenceOkayReturnsFalseWhenSequenceErrorHappens(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *Pipe) (int, error) {
		// all done
		return StatusNotOkay, errors.New("this is an error")
	}

	sequence := NewSequence(op1)
	sequence.Pipe.RunCommand(op1)

	// ----------------------------------------------------------------
	// perform the change

	success := sequence.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, success)
}

func TestSequenceParseIntCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.ParseInt()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceParseIntCopesWithEmptySequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence Sequence
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.ParseInt()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceParseIntConvertsContentsOfStdoutWhenNoError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := 100
	op1 := func(p *Pipe) (int, error) {
		p.Stdout.WriteString("100\n")

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, nil
	}

	sequence := NewSequence(op1)
	op1(sequence.Pipe)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.ParseInt()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceParseIntReturnsZeroWhenError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := 0
	op1 := func(p *Pipe) (int, error) {
		// we don't want to see this in our final output
		p.Stdout.WriteString("we do not want this")
		p.Stderr.WriteString("not a number")

		// all done
		return 0, errors.New("an error occurred")
	}

	sequence := NewSequence(op1)
	sequence.Pipe.RunCommand(op1)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.ParseInt()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestNewSequenceSetParamsUpdatesParamsCount(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "5"
	sequence := NewSequence()

	// ----------------------------------------------------------------
	// perform the change

	sequence.SetParams("one", "two", "three", "four", "five")
	actualResult := sequence.Pipe.Env.Getenv("$#")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewSequenceSetParamsUpdatesPositionalParams(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	sequence := NewSequence()

	// ----------------------------------------------------------------
	// perform the change

	sequence.SetParams("one", "two", "three", "four", "five")
	actualResult1 := sequence.Pipe.Env.Getenv("$1")
	actualResult2 := sequence.Pipe.Env.Getenv("$2")
	actualResult3 := sequence.Pipe.Env.Getenv("$3")
	actualResult4 := sequence.Pipe.Env.Getenv("$4")
	actualResult5 := sequence.Pipe.Env.Getenv("$5")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "one", actualResult1)
	assert.Equal(t, "two", actualResult2)
	assert.Equal(t, "three", actualResult3)
	assert.Equal(t, "four", actualResult4)
	assert.Equal(t, "five", actualResult5)
}

func TestNewSequenceSetParamsUpdatesParamsString(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "one two three four five"
	sequence := NewSequence()

	// ----------------------------------------------------------------
	// perform the change

	sequence.SetParams("one", "two", "three", "four", "five")
	actualResult := sequence.Pipe.Env.Getenv("$*")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewSequenceSetParamsRemovesPreviousParams(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	sequence := NewSequence()
	sequence.SetParams("one", "two", "three", "four", "five")

	// ----------------------------------------------------------------
	// perform the change

	sequence.SetParams("six", "seven", "eight")

	actualCount := sequence.Pipe.Env.Getenv("$#")
	actualParams := sequence.Pipe.Env.Getenv("$*")
	actualParam1 := sequence.Pipe.Env.Getenv("$1")
	actualParam2 := sequence.Pipe.Env.Getenv("$2")
	actualParam3 := sequence.Pipe.Env.Getenv("$3")
	actualParam4 := sequence.Pipe.Env.Getenv("$4")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, "3", actualCount)
	assert.Equal(t, "six seven eight", actualParams)
	assert.Equal(t, "six", actualParam1)
	assert.Equal(t, "seven", actualParam2)
	assert.Equal(t, "eight", actualParam3)
	assert.Equal(t, "", actualParam4)
}

func TestNewSequenceSetParamsCopesWithNilPointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence = nil

	// ----------------------------------------------------------------
	// perform the change

	sequence.SetParams("one", "two", "three", "four", "five")

	// ----------------------------------------------------------------
	// test the results

}

func TestSequenceStringCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceStringCopesWithEmptySequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence Sequence
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceStringReturnsContentsOfStdoutWhenNoError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"
	op1 := func(p *Pipe) (int, error) {
		// this is the content we want
		p.Stdout.WriteString(expectedResult)

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, nil
	}

	sequence := NewSequence(op1)
	op1(sequence.Pipe)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceStringReturnsContentsOfStdoutWhenError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"
	op1 := func(p *Pipe) (int, error) {
		// this is the content we want
		p.Stdout.WriteString(expectedResult)

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, errors.New("an eccor occurred")
	}

	sequence := NewSequence(op1)
	sequence.Pipe.RunCommand(op1)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.String()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceStringsCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence
	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceStringsCopesWithEmptySequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence Sequence
	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceStringsReturnsContentsOfStdoutWhenNoError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := []string{"hello world", "have a nice day"}
	op1 := func(p *Pipe) (int, error) {
		for _, line := range expectedResult {
			p.Stdout.WriteString(line)
			p.Stdout.WriteRune('\n')
		}

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, nil
	}

	sequence := NewSequence(op1)
	sequence.Pipe.RunCommand(op1)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceStringsReturnsContentsOfStdoutWhenError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := []string{"hello world", "have a nice day"}
	op1 := func(p *Pipe) (int, error) {
		for _, line := range expectedResult {
			p.Stdout.WriteString(line)
			p.Stdout.WriteRune('\n')
		}

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, errors.New("an error occurred")
	}

	sequence := NewSequence(op1)
	sequence.Pipe.RunCommand(op1)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceTrimmedStringCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceTrimmedStringCopesWithEmptySequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence Sequence
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceTrimmedStringReturnsContentsOfStdoutWhenNoError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "   hello world\nhave a nice day\n\n\n"
	expectedResult := "hello world\nhave a nice day"

	op1 := func(p *Pipe) (int, error) {
		// this is the content we want
		p.Stdout.WriteString(testData)

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, nil
	}

	sequence := NewSequence(op1)
	sequence.Pipe.RunCommand(op1)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceTrimmedStringReturnsContentsOfStdoutWhenError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "   hello world\nhave a nice day\n\n\n"
	expectedResult := "hello world\nhave a nice day"
	op1 := func(p *Pipe) (int, error) {
		// this is the content we want
		p.Stdout.WriteString(testData)

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, errors.New("an error occurred")
	}

	sequence := NewSequence(op1)
	sequence.Pipe.RunCommand(op1)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := sequence.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}
