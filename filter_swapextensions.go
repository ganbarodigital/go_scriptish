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
	"path/filepath"
	"strings"
)

// SwapExtensions treats every line in the pipeline as a filepath.
//
// It replaces every old extension with the corresponding new one.
func SwapExtensions(old []string, new []string) Command {
	// build our Scriptish command
	return func(p *Pipe) (int, error) {
		// special case - we want to replace *every extension* in old with
		// whatever is in new
		if len(old) > 1 && len(new) == 1 {
			for i := 1; i < len(old); i++ {
				new = append(new, new[0])
			}
		} else if len(old) != len(new) {
			// we don't know what to do
			return StatusNotOkay, ErrMismatchedInputs{"old", len(old), "new", len(new)}
		}

		for line := range p.Stdin.ReadLines() {
			// what extension does this filepath have?
			fileExt := filepath.Ext(line)
			swapped := false

			for i := range old {
				if fileExt == old[i] {
					p.Stdout.WriteString(strings.TrimSuffix(line, fileExt))
					p.Stdout.WriteString(new[i])
					p.Stdout.WriteRune('\n')

					// this helps us make sure that we do not output
					// the same filepath twice
					swapped = true

					// once an extension has been swapped, we do not want
					// to try and swap it a second time
					break
				}
			}

			// did we swap anything over?
			if !swapped {
				p.Stdout.WriteString(line)
				p.Stdout.WriteRune('\n')
			}
		}

		// all done
		return StatusOkay, nil
	}
}
