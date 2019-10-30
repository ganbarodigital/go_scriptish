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
	"os"
	"path/filepath"
	"strings"
)

// ListFiles is the equivalent of `ls -1 <path>`.
//
// If `path` is a file, ListFiles writes the file to the pipeline's stdout
// If `path` is a folder, ListFiles writes the contents of the folder
// to the pipeline's stdout. The path to the folder is included.
//
// If `path` contains wildcards, ListFiles writes any files that matches
// to the pipeline's stdout.
func ListFiles(path string) Command {
	// build our Scriptish command
	return func(p *Pipe) (int, error) {
		// special case: user wants a list of files that match a wildcard
		if strings.ContainsAny(path, "[]^*?\\{}!") {
			return globFiles(p, path)
		}

		// general case: user has given us a path with no wildcards
		info, err := os.Stat(path)
		if err != nil {
			return StatusNotOkay, err
		}

		// what are we looking at?
		if info.IsDir() {
			return listFolder(p, path)
		}

		// if we get here, then `path` maps onto a single file
		return listFile(p, path)
	}
}

func globFiles(p *Pipe, path string) (int, error) {
	// can we find any files?
	filenames, err := filepath.Glob(path)
	if err != nil {
		return StatusNotOkay, err
	}

	// we have something to pass on
	for _, filename := range filenames {
		p.Stdout.WriteString(filename)
		p.Stdout.WriteRune('\n')
	}

	// all done
	return StatusOkay, nil
}

func listFile(p *Pipe, path string) (int, error) {
	p.Stdout.WriteString(path)
	p.Stdout.WriteRune('\n')

	return StatusOkay, nil
}

func listFolder(p *Pipe, path string) (int, error) {
	// can we read what's in the folder?
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return StatusNotOkay, err
	}

	// we have some filenames to pass on
	for _, entry := range files {
		p.Stdout.WriteString(filepath.Join(path, entry.Name()))
		p.Stdout.WriteRune('\n')
	}

	// all done
	return StatusOkay, nil
}
