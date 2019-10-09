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
	"strconv"
	"strings"
)

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

func matchSequence(input string, start int) (int, bool) {
	// are we looking at the start of a sequence?
	if input[start] != '{' {
		return 0, false
	}

	// a sequence has the format:
	//
	// {[:alphanum:]..[:alphanum:]} or
	// {[:alphanum:]..[:alphanum:]..[:num:]}
	//
	// no escape chars, no vars to worry about, no nesting either

	braceDepth := 0
	for i := start; i < len(input); i++ {
		if input[i] == '{' {
			braceDepth++

			// no nesting allowed!
			if braceDepth > 1 {
				return 0, false
			}
		} else if input[i] == '}' {
			braceDepth--

			if braceDepth == 0 {
				return i, true
			}
		} else if isSequenceChar(input[i]) {
			continue
		} else {
			return 0, false
		}
	}

	return 0, false
}

func isSequenceChar(c byte) bool {
	return c == '.' || c == '-' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
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

	// all done
	return parts, true
}

type braceSequence struct {
	// are we rendering characters?
	//
	// if false, we are rendering integers
	chars bool

	// what number are we starting on?
	start int

	// what number are we ending on?
	end int

	// are we going up or down, and by how much?
	incr int
}

func parseSequence(pattern string) (braceSequence, bool) {
	var retval braceSequence

	// sequences are (relatively!) simple ... we can use strings.Split()
	// here to get started
	parts := strings.Split(pattern[1:len(pattern)-1], "..")

	// did we get enough parts?
	if len(parts) < 2 || len(parts) > 3 {
		return retval, false
	}

	// the first two parts are the start and end of the sequence
	//
	// they can be single chars or integers, as long as both are the same
	isNumericStart := isNumericString(parts[0])
	isNumericEnd := isNumericString(parts[1])

	if len(parts[0]) == 1 && len(parts[1]) == 1 {
		// we have chars or all-numbers
		if isNumericStart && isNumericEnd {
			// all numbers
			retval.start, _ = strconv.Atoi(parts[0])
			retval.end, _ = strconv.Atoi(parts[1])
		} else {
			// must be chars
			retval.chars = true
			retval.start = int(parts[0][0])
			retval.end = int(parts[1][0])
		}
	} else {
		// if we get here, both parts must be numbers
		if !isNumericStart || !isNumericEnd {
			return braceSequence{}, false
		}

		retval.start, _ = strconv.Atoi(parts[0])
		retval.end, _ = strconv.Atoi(parts[1])
	}

	// do we have an incr element?
	if len(parts) == 3 {
		incr, err := strconv.Atoi(parts[2])
		if err != nil {
			return braceSequence{}, false
		}
		retval.incr = incr
	} else {
		// no we do not, so we must set it ourselves
		if retval.start < retval.end {
			// low to high
			retval.incr = 1
		} else {
			// high to low
			retval.incr = -1
		}
	}

	// all done
	return retval, true
}

func isNumericChar(char byte) bool {
	return '0' <= char && char <= '9'
}

func isNumericString(input string) bool {
	for i := 0; i < len(input); i++ {
		if !isNumericChar(input[i]) {
			return false
		}
	}

	return true
}
