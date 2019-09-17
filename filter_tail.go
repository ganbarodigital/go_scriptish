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
	"container/ring"

	pipe "github.com/ganbarodigital/go_pipe"
)

// Tail copies the last N lines from the pipeline's Stdin to its Stdout
func Tail(n int) Command {
	// special case - deal with horrible values of n
	if n < 1 {
		return func(p *pipe.Pipe) (int, error) {
			// do nothing
			return pipe.OK, nil
		}
	}

	// general case
	return func(p *pipe.Pipe) (int, error) {
		// we'll use the ring buffer for this
		buf := ring.New(n)

		// we write everything to the ring buffer, and let it throw away
		// everything but the requested number of lines :)
		for line := range p.Stdin.ReadLines() {
			buf.Value = line
			buf = buf.Next()
		}

		// at this point, the ring buffer contains (up to) the right number
		// of lines
		buf.Do(func(line interface{}) {
			// did we fill the buffer up?
			if line == nil {
				// no, we did not
				return
			}

			// if we get here, we have a line to preserve
			p.Stdout.WriteString(line.(string))
			p.Stdout.WriteRune('\n')
		})

		// all done
		return pipe.OK, nil
	}
}
