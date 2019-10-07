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
	"os"
	"testing"

	envish "github.com/ganbarodigital/go_envish/v2"
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

func TestSequenceExpandCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence

	// ----------------------------------------------------------------
	// perform the change

	sequence.Expand("hello ${HOME}")

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestSequenceExpandCopesWithEmptySequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence Sequence

	// ----------------------------------------------------------------
	// perform the change

	sequence.Expand("hello ${HOME}")

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestSequenceExpandUsesTheProgramEnvironmentByDefault(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestSequenceKey"
	testValue := "this is a test"
	os.Setenv(testKey, testValue)

	expectedResult := "hello this is a test"
	sequence := NewSequence()

	// ----------------------------------------------------------------
	// perform the change

	actualResult := sequence.Expand("hello ${TestSequenceKey}")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestSequenceExpandUsesTheSequenceEnvironmentIfAvailable(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestSequenceKey"
	testValue1 := "this is a test"
	testValue2 := "this is another test"
	os.Setenv(testKey, testValue1)

	expectedResult := "hello this is another test"

	sequence := NewSequence()
	sequence.Pipe.Env = envish.NewEnv()
	sequence.Pipe.Env.Setenv(testKey, testValue2)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := sequence.Expand("hello ${TestSequenceKey}")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
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
