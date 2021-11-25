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
	"fmt"
	"io"
	"io/ioutil"
	"os"

	envish "github.com/ganbarodigital/go_envish/v3"
)

// flag we set if we are executing commands in a pipeline
const contextIsPipeline = 1

// SequenceStep bundles up both a Command to run, and the options to apply
// to the pipe when that command runs
type SequenceStep struct {
	Command Command
	Opts    []*StepOption
}

// NewSequenceStep creates a new runnable step for a pipeline or list
func NewSequenceStep(command Command, opts ...*StepOption) *SequenceStep {
	return &SequenceStep{
		Command: command,
		Opts:    opts,
	}
}

// RunStep uses the given Pipe to run this step.
//
// We run any StepOption setup phases first. If any setup phase fails,
// we skip the Command, but do run the teardown phases (so that we can
// clean up after ourselves).
func (st *SequenceStep) RunStep(p *Pipe) (int, error) {
	// do any per-command setup, such as redirects
	//
	// if the setup fails, we bail, and do not attempt to run
	// the command itself ... but we will run the teardown phases,
	// so that we can clean up after ourselves
	_, err := ApplySetupPhasesToPipe(p, st.Opts...)

	if err == nil {
		// run the next step
		p.RunCommand(st.Command)
	}

	// do any post-command teardown, such as closing open files
	ApplyTeardownPhasesToPipe(p, st.Opts...)

	// for convenience
	return p.StatusError()
}

// Sequence is a set of commands to be executed.
//
// Provide your own logic to do the actual command execution.
type Sequence struct {
	// our commands read from / write to this pipe
	Pipe *Pipe

	// keep track of the steps that belong to this sequence
	Steps []*SequenceStep

	// How we will run the sequence
	Controller SequenceController

	// we store local variables here
	LocalVars *envish.LocalEnv

	// the flags we pass into new pipes
	Flags int
}

// NewSequence creates a sequence that's ready to run
func NewSequence(steps ...*SequenceStep) *Sequence {
	sq := Sequence{
		Steps:     steps,
		LocalVars: envish.NewLocalEnv(),
	}

	// make sure we have a pipe, and its environment knows about our
	// LocalVars
	sq.NewPipe()

	// set special parameters here ... well, the ones that make sense
	// for Scriptish, anyways :)
	//
	// positional parameters are set when Exec() is called
	sq.LocalVars.Setenv("$0", os.Args[0])
	sq.LocalVars.Setenv("$-", os.Args[0])
	sq.LocalVars.Setenv("$$", fmt.Sprintf("%d", os.Getpid()))

	// all done
	return &sq
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
	retval, _ := ioutil.ReadAll(sq.Pipe.Stdout)
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
func (sq *Sequence) Exec(params ...string) *Sequence {
	// do we have a sequence to work with?
	if sq == nil {
		return sq
	}

	// do we have a controller?
	if sq.Controller == nil {
		return sq
	}

	// we start with a new Pipe
	sq.NewPipe()

	// we need to set the parameters
	sq.SetParams(params...)

	// use the embedded controller to animate the sequence
	sq.Controller()

	// all done
	return sq
}

// Flush writes the output from running this sequence to the given
// stdout and stderr
func (sq *Sequence) Flush(stdout io.Writer, stderr io.Writer) {
	// do we have a sequence to play with?
	if sq == nil || sq.Pipe == nil {
		return
	}

	io.Copy(stdout, sq.Pipe.Stdout)
	io.Copy(stderr, sq.Pipe.Stderr)
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

// NewPipe replaces the Sequence's existing pipe with a brand new (and empty)
// one. This is very useful for reusing Sequences.
//
// This is called from various places right before a Sequence is run.
//
// You shouldn't need to call it yourself, but it's exported just in case
// it's useful in some way.
func (sq *Sequence) NewPipe() {
	// we start with a new Pipe
	sq.Pipe = NewPipe()

	// the new pipe needs a new environment establishing
	sq.Pipe.Env = envish.NewOverlayEnv(
		sq.LocalVars,
		envish.NewProgramEnv(),
	)

	// set the flags
	sq.Pipe.Flags = sq.Flags
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

// SetParams sets $#, $1... and $* in the pipe's Var store
func (sq *Sequence) SetParams(params ...string) {
	// do we have a sequence to work with?
	if sq == nil || sq.Pipe == nil {
		return
	}

	setParamsInEnv(sq.LocalVars, params)
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
