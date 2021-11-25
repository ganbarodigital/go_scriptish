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
	"os/exec"
)

// Exec runs an operating system command, and posts the results to
// the pipeline's Stdout and Stderr.
//
// The command's status code is stored in the pipeline.StatusCode.
func Exec(args []string, opts ...*StepOption) *SequenceStep {
	// build our Scriptish command
	return NewSequenceStep(
		func(p *Pipe) (int, error) {
			// expand our input
			expArgs := make([]string, len(args))
			for i := 0; i < len(args); i++ {
				expArgs[i] = p.Env.Expand(args[i])
			}

			// debugging support
			Tracef("Exec(%#v)", args)
			Tracef("=> Exec(%#v)", expArgs)

			// build our command
			cmd := exec.Command(expArgs[0], expArgs[1:]...)

			// attach all of our inputs and outputs
			stdout := NewTextBuffer()
			stderr := NewTextBuffer()
			cmd.Stdin = p.Stdin
			cmd.Stdout = stdout
			cmd.Stderr = stderr

			// let's do it
			err := cmd.Start()
			if err != nil {
				return StatusNotOkay, err
			}

			// wait for it to finish
			err = cmd.Wait()

			// copy the output to our pipe
			//
			// it's not ideal, because we can't preserve the original mixed
			// order of the command's output atm
			//
			// at some point, we'll need a new version of pipe that does
			// support preserving mixed order output!
			for line := range stdout.ReadLines() {
				TracePipeStdout("%s", line)
				p.Stdout.WriteString(line)
				p.Stdout.WriteRune('\n')
			}
			for line := range stderr.ReadLines() {
				TracePipeStderr("%s", line)
				p.Stderr.WriteString(line)
				p.Stderr.WriteRune('\n')
			}

			// we want the process's status code
			statusCode := cmd.ProcessState.ExitCode()

			// all done
			return statusCode, err
		},
		opts...,
	)
}
