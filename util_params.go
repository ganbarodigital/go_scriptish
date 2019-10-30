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
	"strconv"
	"strings"

	envish "github.com/ganbarodigital/go_envish/v3"
)

func getParamsFromEnv(e envish.Reader) []string {
	// our return value
	retval := []string{}

	// how many parameters are there?
	paramCount := getParameterCountFromEnv(e)

	// go get them
	for i := 1; i <= paramCount; i++ {
		retval = append(retval, e.Getenv(fmt.Sprintf("$%d", i)))
	}

	// all done
	return retval
}

func getParameterCountFromEnv(e envish.Reader) int {
	paramCount, err := strconv.Atoi(e.Getenv("$#"))
	if err != nil {
		return 0
	}
	return paramCount
}

func setParamsInEnv(e envish.ReaderWriter, params []string) {
	// step one - remove any existing positional params
	for i := 1; i < getParameterCountFromEnv(e); i++ {
		e.Unsetenv(fmt.Sprintf("$%d", i))
	}

	// step two - set the new param count
	//
	// params start at $1
	newParamCount := len(params)
	e.Setenv("$#", strconv.Itoa(newParamCount))

	for i := 1; i <= newParamCount; i++ {
		e.Setenv(fmt.Sprintf("$%d", i), params[i-1])
	}

	// step three - set $* and $@
	e.Setenv("$*", strings.Join(params, " "))
	e.Setenv("$@", strings.Join(params, " "))

	// all done
}
