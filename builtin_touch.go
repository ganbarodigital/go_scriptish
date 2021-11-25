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
	"time"
)

// Touch creates the named file (if it doesn't exist), or updates its
// atime and mtime (if it does exist)
//
// It ignores the contents of the pipeline.
//
// On success, it returns the status code `StatusOkay`. On failure,
// it returns the status code `StatusNotOkay`.
func Touch(filepath string, opts ...*StepOption) *SequenceStep {
	// build our Scriptish command
	return NewSequenceStep(
		func(p *Pipe) (int, error) {
			// expand our input
			expFilepath := p.Env.Expand(filepath)

			// debugging support
			Tracef("Touch(%#v)", filepath)
			Tracef("=> Touch(%#v)", expFilepath)

			var fh *os.File

			// does the file exist?
			_, err := os.Stat(expFilepath)
			if err != nil {
				if os.IsNotExist(err) {
					fh, err = os.OpenFile(expFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						return StatusNotOkay, err
					}
					defer fh.Close()
				}
			} else {
				// if we get here, then the file does exist
				//
				// we need to modify its inode data
				now := time.Now()
				err = os.Chtimes(expFilepath, now, now)
			}

			if err != nil {
				return StatusNotOkay, err
			}

			// all done
			return StatusOkay, nil
		},
		opts...,
	)
}
