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
	pipe "github.com/ganbarodigital/go_pipe/v5"
)

// Command is an alias for the underlying pipe library's Command
//
// It just saves us having to import the pipe library into every single
// file in the project.
type Command = pipe.Command

// Dest is an alias for the underlying pipe library's Dest
type Dest = pipe.Dest

// NewDest creates a new pipe.Dest struct
func NewDest() *Dest {
	return pipe.NewDest()
}

// List is an alias for a Sequence
type List = Sequence

// Pipe is an alias for the underlying pipe library's Pipe
//
// It just saves us having to import the pipe library into every single
// file in the project.
type Pipe = pipe.Pipe

// Pipeline is an alias for a Sequence
type Pipeline = Sequence

// Source is an alias for the underlying pipe library's Source
type Source = pipe.Source

// NewSourceFromReader is an alias for the underlying pipe library's
// NewSourceFromReader()
var NewSourceFromReader = pipe.NewSourceFromReader

// NewSourceFromString is an alias for the underlying pipe library's
// NewSourceFromString()
var NewSourceFromString = pipe.NewSourceFromString

// StatusOkay is an alias for the underlying pipe library's StatusOkay
//
// It just saves us having to import the pipe library into every single
// file in the project.
const StatusOkay = pipe.StatusOkay

// StatusNotOkay is an alias for the underlying pipe library's StatusNotOkay
//
// It just saves us having to import the pipe library into every single
// file in the project.
const StatusNotOkay = pipe.StatusNotOkay
