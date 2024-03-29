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

// Expand replaces ${var} and $var in the input string. Variable values
// are found by calling the supplied mapping function.
//
// This is a replacement for Golang's `os.Expand()` that supports full
// UNIX shell string expansion. It is not a drop-in replacement, but it
// should be straight-forward to migrate from `os.Expand()`
func Expand(input string, cb ExpansionCallbacks) (string, error) {
	// step 1: brace expansion
	input = expandBraces(input)

	// step 2: tilde expansion
	input = ExpandTilde(input, cb)

	// step 3: parameter & variable expansion
	var err error
	input, err = expandParameters(input, cb)
	if err != nil {
		return "", err
	}

	// step 4: arithmetic expansion
	// step 5: quote removal
	input = expandQuoteRemoval(input)

	// all done
	return input, nil
}
