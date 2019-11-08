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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTouchWillCreateAFileIfItDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	// we need a file to touch
	tmpFilename, err := NewPipeline(
		MkTempFilename(os.TempDir(), "scriptify-truncatefile-*"),
	).Exec().TrimmedString()
	assert.Nil(t, err)

	// this is the pipeline we'll use to test Touch()
	pipeline := NewPipeline(
		Touch(tmpFilename),
	)

	// prove that the file currently does not exist
	_, err = os.Stat(tmpFilename)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))

	// ----------------------------------------------------------------
	// perform the change

	// this will create the file
	err = pipeline.Exec().Error()
	assert.Nil(t, err)

	// ----------------------------------------------------------------
	// test the results

	fi, err := os.Stat(tmpFilename)
	assert.Nil(t, err)
	modTime := fi.ModTime().Unix()
	now := time.Now().Unix()
	assert.True(t, now-modTime < 2)
}

func TestTouchUpdatesTheModifiedTimeOfAFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// we need a file to touch
	tmpFilename, err := NewPipeline(
		MkTempFile(os.TempDir(), "scriptify-truncatefile-*"),
	).Exec().TrimmedString()
	assert.Nil(t, err)

	now := time.Now()
	then := now.AddDate(0, 0, -1)

	// let's force the mtime on this file
	os.Chtimes(tmpFilename, then, then)

	// prove that worked
	fi, err := os.Stat(tmpFilename)
	assert.Equal(t, then, fi.ModTime())

	// this is the pipeline we are going to use
	pipeline := NewPipeline(
		Touch(tmpFilename),
	)

	// ----------------------------------------------------------------
	// perform the change

	// this will create the file
	err = pipeline.Exec().Error()
	assert.Nil(t, err)

	// ----------------------------------------------------------------
	// test the results

	fi, err = os.Stat(tmpFilename)
	assert.Nil(t, err)
	modTime := fi.ModTime().Unix()
	assert.True(t, now.Unix()-modTime < 2)
}
