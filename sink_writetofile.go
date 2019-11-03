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
)

// WriteToFile writes the contents of the pipe's stdin to the given file.
// The existing contents of the file are replaced.
//
// If the file does not exist, it is created.
func WriteToFile(filename string) Command {
	// build our Scriptish command
	return func(p *Pipe) (int, error) {
		// expand our input
		expFilename := p.Env.Expand(filename)

		// debugging support
		Tracef("WriteToFile(%#v)", filename)
		Tracef("=> WriteToFile(%#v)", expFilename)

		// open / create the file
		fh, err := os.OpenFile(expFilename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return StatusNotOkay, err
		}
		defer fh.Close()

		// write to the file
		for line := range getSinkReader(p) {
			TraceOutput("file", "%s", line)
			_, err = fh.WriteString(line)
			if err != nil {
				return StatusNotOkay, err
			}
			_, err = fh.WriteString("\n")
			if err != nil {
				return StatusNotOkay, err
			}
		}

		// all done
		return StatusOkay, err
	}
}
