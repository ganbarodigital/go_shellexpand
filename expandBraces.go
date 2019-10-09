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

import "strings"

func expandBraces(input string) string {
	// we expand in a strictly left-to-right manner
	for i := 0; i < len(input); i++ {
		if input[i] == '\\' {
			i++
		} else if input[i] == '$' {
			varEnd, ok := matchVar(input, i)
			if ok {
				i = varEnd
			}
		} else if input[i] == '{' {
			patternEnd, ok := matchPattern(input, i)
			if ok {
				patternParts, ok := parsePattern(input[i : patternEnd+1])
				if ok {
					preambleStart := findPreambleStart(input, i)
					postscriptEnd := findPostscriptEnd(input, patternEnd)

					var exp []string
					for _, part := range patternParts {
						exp = append(exp, expandPattern(input, part, i, patternEnd, preambleStart, postscriptEnd))
					}

					var buf strings.Builder
					if preambleStart > 0 {
						buf.WriteString(input[:preambleStart])
					}
					buf.WriteString(strings.Join(exp, " "))
					if postscriptEnd < len(input) {
						buf.WriteRune(' ')
						buf.WriteString(input[postscriptEnd+1:])
					}
					input = buf.String()
				}
			}
		}
	}

	return input
}

func expandPattern(input, part string, i, patternEnd, preambleStart, postscriptEnd int) string {
	// we'll build our substitution here
	var buf strings.Builder

	// do we have a preamble to add?
	if preambleStart < i {
		buf.WriteString(input[preambleStart:i])
	}

	// we always have a pattern part to add
	buf.WriteString(part)

	// do we have a postscript to add?
	if postscriptEnd > patternEnd+1 {
		buf.WriteString(input[patternEnd+1 : postscriptEnd])
	}

	return buf.String()
}

func findPreambleStart(input string, preambleStart int) int {
	for ; preambleStart > 0; preambleStart-- {
		if input[preambleStart] == ' ' {
			return preambleStart + 1
		}
	}

	return 0
}

func findPostscriptEnd(input string, postscriptEnd int) int {
	for ; postscriptEnd < len(input); postscriptEnd++ {
		if input[postscriptEnd] == ' ' {
			return postscriptEnd
		}
	}

	return postscriptEnd
}

func matchPattern(input string, start int) (int, bool) {
	// are we looking at the start of a pattern?
	if input[start] != '{' {
		return 0, false
	}

	braceDepth := 0
	for i := start; i < len(input); i++ {
		if input[i] == '\\' {
			// skip over escaped character
			i++
		} else if input[i] == '$' {
			varEnd, ok := matchVar(input, i)
			if ok {
				i = varEnd
			}
		} else if input[i] == '{' {
			braceDepth++
		} else if input[i] == '}' {
			braceDepth--

			if braceDepth == 0 {
				return i, true
			}
		}
	}

	return 0, false
}

func parsePattern(pattern string) ([]string, bool) {
	var parts []string

	// we can't do a simple `strings.Split()` here, because we have to
	// take nested braces into account

	// how many braces are we in?
	braceDepth := 0

	// where are we?
	start := 1

	// find the next pattern
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '\\':
			// skip over escaped characters
			i++
		case '{':
			braceDepth++
		case '}':
			braceDepth--
			// have we reached the end of the pattern?
			if braceDepth == 0 {
				parts = append(parts, pattern[start:i])
				start = i + 1
			}
		case ',':
			// are we in a sub-pattern?
			if braceDepth == 1 {
				// no, we are not :)
				parts = append(parts, pattern[start:i])
				start = i + 1
			}
		}
	}

	// did we end up with mismatched braces?
	if braceDepth > 0 {
		return []string{}, false
	}

	if len(parts) == 1 {
		return []string{}, false
	}

	// make sure no parts are empty
	for _, part := range parts {
		if len(part) == 0 {
			return []string{}, false
		}
	}

	// all done
	return parts, true
}
