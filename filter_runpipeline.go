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
	"io"
)

// RunPipeline allows you to call one pipeline from another.
//
// Use this to create reusable pipelines.
func RunPipeline(pl *Pipeline, opts ...*StepOption) *SequenceStep {
	// build our Scriptish command
	return NewSequenceStep(
		func(p *Pipe) (int, error) {
			// make sure our sub pipeline starts nice and empty
			pl.NewPipe()

			// copy the pipeline's content into our sub pipeline
			for line := range p.Stdin.ReadLines() {
				pl.Pipe.Stdout.WriteString(line)
				pl.Pipe.Stdout.WriteRune('\n')
			}

			// get our parameters
			params := getParamsFromEnv(p.Env)

			// set them in the target pipeline
			//
			// we're not calling pl.Exec() (see below) so we have to set them
			// this way instead
			pl.SetParams(params...)

			// run our sub-pipeline
			//
			// NOTE that we cannot call pl.Exec(), as that (by design) starts
			// the pipeline with a brand-new pipe
			pl.Controller()

			// copy our pipeline's stdout to become the pipe's next stdin
			io.Copy(p.Stdout, pl.Pipe.Stdout)

			// all done
			return pl.StatusError()
		},
		opts...,
	)
}
