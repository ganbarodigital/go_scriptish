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

// PipelineController executes a sequence of commands as if they were
// a UNIX shell pipeline
func PipelineController(sq *Sequence) SequenceController {
	return func() {
		// do we have a pipeline to play with?
		if sq == nil {
			return
		}

		// execute everything in our pipeline
		for _, step := range sq.Steps {
			// at this point, stdout needs to become the next
			// stdin
			preparePipeForNextCommand(sq.Pipe)

			// run the command
			step.RunStep(sq.Pipe)

			// we stop executing the moment something goes wrong
			err := sq.Pipe.Error()
			if err != nil {
				// debugging support
				statusCode := sq.StatusCode()
				Tracef("status code: %d", statusCode)
				Tracef("error: %s", err.Error())

				// we cannot continue
				return
			}
		}
	}
}

// helper function
func preparePipeForNextCommand(p *Pipe) {
	// the output from our previous command becomes the input to the next
	p.SetStdinFromString(p.Stdout.String())

	// the next command starts with no output
	p.SetNewStdout()

	// we throw away any errors that have been written here
	p.SetNewStderr()
}
