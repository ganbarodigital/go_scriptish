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
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Range struct {
	start int
	end   int
}

var rangeRegex = regexp.MustCompile("([1-9][0-9]*){0,1}-([1-9][0-9]*){0,1}")

// ParseRangeSpec takes a string of the form `X1-Y1[,X2-Y2 ...]` and turns
// it into a list of start and end ranges
//
// It emulates the `cut -f <range>` range support.
func ParseRangeSpec(spec string) ([]Range, error) {
	// this will hold all the columns that have been requested
	var retval []Range

	// do we have several ranges to parse?
	for _, item := range strings.Split(spec, ",") {
		newRange, err := parseSingleRange(item)
		if err != nil {
			return []Range{}, err
		}
		retval = append(retval, *newRange)
	}

	// all done
	return retval, nil
}

func parseSingleRange(spec string) (*Range, error) {
	// special case - have we received a single number?
	if !strings.HasPrefix(spec, "-") {
		element, err := strconv.Atoi(spec)
		if err == nil {
			return &Range{element, element}, nil
		}
	}

	// what does our regex make of this?
	matches := rangeRegex.FindAllStringSubmatch(spec, -1)

	// did the regex work?
	if len(matches) == 0 || len(matches[0]) < 3 {
		return nil, fmt.Errorf("invalid range: %s", spec)
	}

	// at this point:
	//
	// matches[0] contains item
	// matches[1] contains start (may be empty string)
	// matches[2] contains end (may be empty string)

	// both matches[1] and matches[2] cannot be empty
	if len(matches[0][1]) == 0 && len(matches[0][2]) == 0 {
		return nil, fmt.Errorf("invalid range: start and end cannot both be empty")
	}

	var start, end int

	if len(matches[0][1]) > 0 {
		start, _ = strconv.Atoi(matches[0][1])
	} else {
		start = 1
	}

	if len(matches[0][2]) > 0 {
		end, _ = strconv.Atoi(matches[0][2])
	}
	if (start > 0) && end == 0 {
		end = math.MaxInt64
	}

	return &Range{start, end}, nil
}
