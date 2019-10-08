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

// bracePair tracks the location of opening and closing braces
// in a string
type bracePair struct {
	start int
	end   int
}

// bracePairs is a list of bracePair structs
type bracePairs []bracePair

// Len is used for sorting the brace pairs so that nested pairs come first
func (e bracePairs) Len() int {
	return len(e)
}

// Less is used for sorting the brace pairs so that nested pairs come first
func (e bracePairs) Less(i, j int) bool {
	return (e[i].end <= e[j].end)
}

// Swap is used for sorting the brace pairs so that nested pairs come first
func (e bracePairs) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// matchBraces searches a string to find all the braces that match
//
// nested braces are returned earlier than outer braces
// quotes are completely ignored
func matchBraces(input string) ([]bracePair, error) {
	// the list of braces that we will return
	var retval []bracePair

	// we'll build a stack of nested braces as we go
	var braceStack []bracePair

	// keep track of where we are in the list
	pairIndex := -1

	// search the string
	for i := 0; i < len(input); i++ {
		// find our '{'
		if input[i] == '{' {
			pairIndex++
			braceStack = append(braceStack, bracePair{i, -1})
		} else if input[i] == '}' {
			if pairIndex < 0 {
				return []bracePair{}, ErrMismatchedClosingBrace{i + 1}
			}

			braceStack[pairIndex].end = i
			retval = append(retval, braceStack[pairIndex])

			pairIndex--
			braceStack = braceStack[:len(braceStack)-1]
		}
	}

	// did we run into mismatched braces?
	if len(braceStack) > 0 {
		return []bracePair{}, ErrMismatchedBrace{braceStack[0].start}
	}

	// all done
	return retval, nil
}
