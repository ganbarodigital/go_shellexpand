// shellexpand is a replacement for Golang's `os.Expand()` that supports
// UNIX shell string expansion and substituation
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

package shellexpand

import (
	"strings"
)

// ExpandTilde will expand any '~' at the start of a word as follows:
//
// ~/path/to/folder -> $HOME/path/to/folder
// ~username/path/to/folder -> <user's homedir>/path/to/folder
// ~+/path/to/folder -> $PWD/path/to/folder
// ~-/path/to/folder -> $OLDPWD/path/to/folder
//
// Directory stack (~+N / ~-N) expansion is not supported (yet).
//
// If expansion fails, the input string is left unmodified.
//
// Don't call this directly; use Expand() instead.
//
// This function is exported because (for UNIX shell compatibility), you
// should call this function when setting variables.
func ExpandTilde(input string, varFuncs VarFuncs) string {
	// we expand in a strictly left-to-right manner
	for i := 0; i < len(input); i++ {
		if input[i] == '\\' {
			// skip over escaped characters
			i++
		} else if input[i] == '$' {
			varEnd, ok := matchVar(input, i)
			if ok {
				i = varEnd
			}
		} else if input[i] == '~' {
			input, _ = matchAndExpandTilde(input, i, varFuncs)
		}
	}

	return input
}

func matchAndExpandTilde(input string, start int, varFuncs VarFuncs) (string, bool) {
	var ok bool

	// are we looking at a tilde w/ optional prefix??
	prefixEnd, ok := matchTildePrefix(input, start)
	if !ok {
		return input, false
	}

	// what kind of prefix are we looking at?
	tildePrefix, _ := parseTildePrefix(input[start : prefixEnd+1])

	// this will hold our replacement
	var repl string

	// build the replacement
	switch tildePrefix.kind {
	case tildePrefixHome:
		repl, ok = varFuncs.LookupVar("HOME")
		if !ok {
			return input, false
		}
	case tildePrefixPwd:
		repl, ok = varFuncs.LookupVar("PWD")
		if !ok {
			return input, false
		}
	case tildePrefixOldPwd:
		repl, ok = varFuncs.LookupVar("OLDPWD")
		if !ok {
			return input, false
		}
	case tildePrefixUsername:
		repl, ok = varFuncs.LookupHomeDir(tildePrefix.prefix)
		if !ok {
			return input, false
		}
	}

	var buf strings.Builder
	if start > 0 {
		buf.WriteString(input[:start])
	}
	buf.WriteString(repl)
	if prefixEnd < len(input) {
		buf.WriteString(input[prefixEnd+1:])
	}

	return buf.String(), true
}

func matchTildePrefix(input string, start int) (int, bool) {
	// are we looking at the start of a prefix?
	if input[start] != '~' {
		return start, false
	}

	// find the end of the prefix
	for i := start; i < len(input); i++ {
		if input[i] == '\\' {
			// skip over escaped character
			i++
		} else if input[i] == '/' || input[i] == ' ' {
			return i - 1, true
		}
	}

	// if we get here, the '~' prefix is the last part of the string
	return len(input) - 1, true
}

const (
	tildePrefixHome = iota
	tildePrefixUsername
	tildePrefixOldPwd
	tildePrefixPwd
)

type tildePrefix struct {
	kind   int
	prefix string
}

func parseTildePrefix(input string) (tildePrefix, bool) {
	// make sure we *are* looking at a prefix
	if input[0] != '~' {
		return tildePrefix{}, false
	}

	// what kind of prefix are we looking at?
	if len(input) == 1 {
		return tildePrefix{tildePrefixHome, ""}, true
	}
	if input == "~+" {
		return tildePrefix{tildePrefixPwd, ""}, true
	}
	if input == "~-" {
		return tildePrefix{tildePrefixOldPwd, ""}, true
	}

	// must be a username; all other options eliminated
	return tildePrefix{tildePrefixUsername, input[1:]}, true
}
