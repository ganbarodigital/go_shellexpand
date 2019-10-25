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

func matchVar(input string, start int) (int, bool) {
	// have we started on a dollar?
	if input[start] != '$' {
		return 0, false
	}

	// is the dollar escaped?
	if start > 0 && input[start-1] == '\\' {
		return 0, false
	}

	// no, it is not
	//
	// special case: positional parameters are not subject to normal
	// matching rules (sigh)
	if isNumericChar(input[start+1]) {
		return start + 2, true
	}

	// general case - a non-positional parameter that may be wrapped
	// in braces
	braceDepth := 0
	for i := start + 1; i < len(input); i++ {
		if input[i] == '\\' {
			// skip escaped chars
			i++
		} else if input[i] == '{' {
			braceDepth++
		} else if input[i] == '}' {
			braceDepth--

			if braceDepth == 0 {
				return i + 1, true
			}
		} else if input[i] == ' ' {
			if braceDepth == 0 {
				// we must be looking at a var that was not surrounded
				// by braces
				return i, true
			}

			// no spaces allowed inside a var name
			// return 0, false
		}
	}

	// end of the string
	if braceDepth == 0 {
		return len(input), true
	}

	// we did not find a matching closing brace
	return 0, false
}
