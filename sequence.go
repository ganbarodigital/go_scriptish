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
	"io/ioutil"
	"os"

	pipe "github.com/ganbarodigital/go_pipe/v4"
)

// Sequence is a set of commands to be executed.
//
// Provide your own logic to do the actual command execution.
type Sequence struct {
	// our commands read from / write to this pipe
	Pipe *pipe.Pipe

	// keep track of the steps that belong to this sequence
	Steps []Command

	// How we will run the sequence
	Controller SequenceController
}

// NewSequence creates a sequence that's ready to run
func NewSequence(steps ...Command) *Sequence {
	sequence := Sequence{
		Pipe:  pipe.NewPipe(),
		Steps: steps,
	}

	return &sequence
}

// Bytes returns the contents of the sequence's stdout as a byte slice
func (sq *Sequence) Bytes() ([]byte, error) {
	// do we have a sequence?
	if sq == nil {
		return []byte{}, nil
	}

	// was the sequence initialised correctly?
	if sq.Pipe == nil {
		return []byte{}, nil
	}

	// return what we have
	retval, _ := ioutil.ReadAll(sq.Pipe.Stdout.NewReader())
	return retval, sq.Pipe.Error()
}

// Error returns the sequence's error status.
func (sq *Sequence) Error() error {
	// do we have a sequence to play with?
	if sq == nil {
		return nil
	}

	// if we get here, then all is well
	return sq.Pipe.Error()
}

// Exec executes a sequence
//
// If you embed the sequence in another struct, make sure to override this
// to return your own return type!
func (sq *Sequence) Exec() *Sequence {
	// do we have a sequence to work with?
	if sq == nil {
		return sq
	}

	// do we have a controller?
	if sq.Controller == nil {
		return sq
	}

	// we start with a new Pipe
	sq.Pipe = pipe.NewPipe()

	// use the embedded controller to animate the sequence
	sq.Controller()

	// all done
	return sq
}

// Expand replaces ${var} or $var in the input string.
//
// It uses the sequence's private environment (if the sequence has one),
// or the program's environment otherwise.
func (sq *Sequence) Expand(fmt string) string {
	// do we have a sequence to work with?
	if sq == nil {
		return os.Expand(fmt, os.Getenv)
	}

	// yes we do
	return sq.Pipe.Expand(fmt)
}

// Okay returns false if a sequence operation set the StatusCode to
// anything other than StatusOkay. It returns true otherwise.
func (sq *Sequence) Okay() bool {
	// do we have a sequence to play with?
	if sq == nil || sq.Pipe == nil {
		return true
	}

	return sq.Pipe.Okay()
}

// ParseInt returns the pipe's stdout as an integer
//
// If the integer conversion fails, error will be the conversion error.
// If the integer conversion succeeds, error will be the pipe's error
// (which may be nil)
func (sq *Sequence) ParseInt() (int, error) {
	// do we have a sequence to play with?
	if sq == nil {
		return 0, nil
	}

	// was the sequence correctly initialised?
	if sq.Pipe == nil || sq.Pipe.Stdout == nil {
		return 0, sq.Error()
	}

	// do we have an integer to return?
	retval, err := sq.Pipe.Stdout.ParseInt()
	if err != nil {
		return retval, err
	}

	// all done
	return retval, sq.Error()
}

// StatusCode returns the UNIX-like status code from the last step to execute
func (sq *Sequence) StatusCode() int {
	// do we have a sequence to work with?
	if sq == nil || sq.Pipe == nil {
		return StatusOkay
	}

	// yes we do
	return sq.Pipe.StatusCode()
}

// StatusError is a shorthand to save having to call Sequence.StatusCode()
// followed by Sequence.Error() from your code
func (sq *Sequence) StatusError() (int, error) {
	// do we have a sequence to work with?
	if sq == nil || sq.Pipe == nil {
		return StatusOkay, nil
	}

	// yes we do
	return sq.Pipe.StatusError()
}

// String returns the pipe's stdout as a single string
func (sq *Sequence) String() (string, error) {
	// do we have a sequence to play with?
	if sq == nil {
		return "", nil
	}

	// was the sequence correctly initialised?
	if sq.Pipe == nil {
		return "", sq.Error()
	}

	// return what we have
	return sq.Pipe.Stdout.String(), sq.Error()
}

// Strings returns the sequence's stdout, one string per line
func (sq *Sequence) Strings() ([]string, error) {
	// do we have a sequence to play with?
	if sq == nil {
		return []string{}, nil
	}

	// was the sequence correctly initialised?
	if sq.Pipe == nil {
		return []string{}, sq.Error()
	}

	// return what we have
	return sq.Pipe.Stdout.Strings(), sq.Error()
}

// TrimmedString returns the pipe's stdout as a single string.
// Any leading or trailing whitespace is removed.
func (sq *Sequence) TrimmedString() (string, error) {
	// do we have a sequence to play with?
	if sq == nil {
		return "", nil
	}

	// was the sequence correctly initialised?
	if sq.Pipe == nil {
		return "", sq.Error()
	}

	// return what we have
	return sq.Pipe.Stdout.TrimmedString(), sq.Error()
}