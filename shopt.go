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
	"fmt"
	"io"
)

// ShellOptions holds flags and settings that change Scriptish's behaviour
type ShellOptions struct {
	// trace is where we send our debugging output to
	trace io.Writer
}

// shopt holds the parameters you can set to change Scriptish's behaviour
var shopt ShellOptions

// GetShellOptions gives you access to the package-wide behaviour flags
// and settings
func GetShellOptions() *ShellOptions {
	return &shopt
}

// DisableTrace will switch off execution tracing across Scriptish
func (s *ShellOptions) DisableTrace() {
	s.trace = nil
}

// EnableTrace will switch on execution tracing across Scriptish
func (s *ShellOptions) EnableTrace(dest io.Writer) {
	s.trace = dest
}

// IsTraceEnabled return true if execution tracing is currently switched on
func (s *ShellOptions) IsTraceEnabled() bool {
	return s.trace != nil
}

// IsTraceEnabled returns true if execution tracing is currently switched on
// across Scriptish
func IsTraceEnabled() bool {
	return shopt.trace != nil
}

// Tracef writes a trace message to os.Stderr if tracing is enabled
func Tracef(format string, args ...interface{}) {
	if IsTraceEnabled() {
		fmt.Fprintf(shopt.trace, "+ "+format+"\n", args...)
	}
}

// TraceOutput writes a trace message about content written to a file or
// a buffer
func TraceOutput(dest string, format string, args ...interface{}) {
	Tracef(dest+"> "+format, args...)
}
