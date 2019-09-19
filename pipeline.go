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
	pipe "github.com/ganbarodigital/go_pipe"
)

// Command is a wrapper around the underlying pipe library's
// PipelineOperation.
//
// It's just easier to describe them as Commands
// in conversation and in documentaton.
type Command = pipe.PipelineOperation

// Pipe is an alias for the underlying pipe library's Pipe
//
// It just saves us having to import the pipe library into every single
// file in the project.
type Pipe = pipe.Pipe

// OK is an alias for the underlying pipe library's OK
//
// It just saves us having to import the pipe library into every single
// file in the project.
const OK = pipe.OK

// NOT_OK is an alias for the underlying pipe library's NOT_OK
//
// It just saves us having to import the pipe library into every single
// file in the project.
const NOT_OK = pipe.NOT_OK

// Pipeline is a wrapper around our underlying pipe library's pipeline,
// so that we can extend it with extra functionality
type Pipeline struct {
	pipe.Pipeline
}

// NewPipeline creates a pipeline ready to be executed
func NewPipeline(steps ...Command) *Pipeline {
	newPipe := pipe.NewPipeline(steps...)

	return &Pipeline{*newPipe}
}

// Exec executes the current pipeline
func (pl *Pipeline) Exec() *Pipeline {
	pl.Exec_()
	return pl
}

// ExecPipeline creates and runs a pipeline. Use this for short, throwaway
// actions.
func ExecPipeline(steps ...Command) *Pipeline {
	pipeline := NewPipeline(steps...).Exec()
	return pipeline
}

// PipelineFunc creates a pipeline, and wraps it in a function to make
// it easier to call.
func PipelineFunc(steps ...Command) func() *Pipeline {
	newPipe := NewPipeline(steps...)
	return func() *Pipeline {
		return newPipe.Exec()
	}
}
