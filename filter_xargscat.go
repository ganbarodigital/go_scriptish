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
	"strings"
)

// XargsCat treats each line in the pipeline's stdin as a filepath.
// It reads each file, and writes them to the pipeline's stdout.
func XargsCat(opts ...*StepOption) *SequenceStep {
	// build our Scriptish command
	return NewSequenceStep(
		func(p *Pipe) (int, error) {
			// debugging support
			Tracef("XargsCat()")

			// treat each line as a valid filepath
			for line := range p.Stdin.ReadLines() {
				Tracef("reading from file %#v", line)

				// can we read the file?
				contents, err := ioutil.ReadFile(line)
				if err != nil {
					return StatusNotOkay, err
				}

				// add the file contents to the pipeline
				fileContents := string(contents)
				TracePipeStdout("%s", fileContents)
				p.Stdout.WriteString(fileContents)

				// we don't want content from two files ending up on
				// the very same line
				if !strings.HasSuffix(fileContents, "\n") {
					p.Stdout.WriteRune('\n')
				}
			}

			// all done
			return StatusOkay, nil
		},
		opts...,
	)
}
