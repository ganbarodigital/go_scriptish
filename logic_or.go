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

import "io"

// Or executes the given sequence only if the previous command has returned
// some kind of error.
//
// The sequence starts with an empty Stdin. The sequence's output is written
// back to the Stdout and Stderr of the calling list or pipeline - along
// with the StatusCode() and Error().
//
// It is an emulation of UNIX shell scripting's `list1 || command`
func Or(sq *Sequence) Command {
	// we're going to wrap our sequences up as a Scriptish Command
	return func(p *Pipe) (int, error) {
		// do we need to do anything?
		statusCode, err := p.StatusError()
		if err == nil {
			// debugging support
			Tracef("Or(): not executing the given sequence")

			// make sure we do not lose the output of the sequence so far
			p.DrainStdinToStdout()

			// all done
			return statusCode, err
		}

		// debugging support
		Tracef("Or(): executing the given sequence")

		// get our parameters
		params := getParamsFromEnv(p.Env)

		// run it
		sq.Exec(params...)

		// copy the results into our pipe
		io.Copy(p.Stdout, sq.Pipe.Stdout.NewReader())
		io.Copy(p.Stderr, sq.Pipe.Stderr.NewReader())

		// all done
		return sq.StatusCode(), sq.Error()
	}
}
