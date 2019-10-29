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
	"unicode/utf8"
)

// expandBraces performs UNIX shell brace expansion on the input string
func expandBraces(input string) string {
	// this is what we're assessing
	var r rune

	// this is always how many bytes 'r' is
	w := 0

	// this is true when we are skipping over escaped characters
	inEscape := false

	// we expand in a strictly left-to-right manner
	for i := 0; i < len(input); {
		r, w = utf8.DecodeRuneInString(input[i:])

		// what are we looking at?
		if inEscape {
			// skip over escaped character
			inEscape = false
			i += w
		} else if r == '\\' {
			// next character is escaped
			inEscape = true
			i += w
		} else if r == '$' {
			// possible variable?
			//
			// variables are immune to brace expansion
			varEnd, ok := matchVar(input[i:])
			if ok {
				i += varEnd - 1
			} else {
				i += w
			}
		} else if r == '{' {
			// probably the start of something we can expand
			var ok bool
			input, ok = matchAndExpandBraceSequence(input, i)
			if !ok {
				input, ok = matchAndExpandBracePattern(input, i)
			}
			i += w
		} else {
			// just another character, nothing for us to do with it
			i += w
		}
	}

	// all done
	return input
}

func expandBracePattern(preamble, part, postscript string) string {
	// we'll build our substitution here
	var buf strings.Builder

	// may be empty
	if len(preamble) > 0 {
		buf.WriteString(preamble)
	}

	// we always have a pattern part to add
	buf.WriteString(part)

	// may also be empty
	if len(postscript) > 0 {
		buf.WriteString(postscript)
	}

	// all done
	return buf.String()
}

func expandBraceSequence(entry int, isChars bool, preamble, postscript string) string {
	// we'll build our substitution here
	var buf strings.Builder

	// may be empty
	if len(preamble) > 0 {
		buf.WriteString(preamble)
	}

	// we always have a sequence entry to add
	if isChars {
		buf.WriteString(string(entry))
	} else {
		buf.WriteString(strconv.Itoa(entry))
	}

	// may also be empty
	if len(postscript) > 0 {
		buf.WriteString(postscript)
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
	var r rune
	w := 0
	for postscriptEnd < len(input) {
		r, w = utf8.DecodeRuneInString(input[postscriptEnd:])
		if r == ' ' {
			return postscriptEnd
		}
		postscriptEnd += w
	}

	return postscriptEnd
}

func matchAndExpandBracePattern(input string, i int) (string, bool) {
	// are we looking at a pattern?
	patternEnd, ok := matchBracePattern(input[i:])
	if !ok {
		return input, false
	}

	// is it really a pattern though?
	patternParts, ok := parseBracePattern(input[i : i+patternEnd])
	if !ok {
		return input, false
	}

	// if we get here, then yes it is
	preamble := ""
	preambleStart := findPreambleStart(input, i)
	if preambleStart < i {
		preamble = input[preambleStart:i]
	}
	postscript := ""
	postscriptEnd := findPostscriptEnd(input, i+patternEnd)
	if postscriptEnd > i+patternEnd {
		postscript = input[i+patternEnd : postscriptEnd]
	}

	var exp []string
	for _, part := range patternParts {
		exp = append(exp, expandBracePattern(preamble, part, postscript))
	}

	var buf strings.Builder
	if preambleStart > 0 {
		buf.WriteString(input[:preambleStart])
	}
	buf.WriteString(strings.Join(exp, " "))
	if postscriptEnd+1 < len(input) {
		buf.WriteRune(' ')
		buf.WriteString(input[postscriptEnd+1:])
	}

	return buf.String(), true
}

func matchAndExpandBraceSequence(input string, i int) (string, bool) {
	// are we looking at a sequence?
	seqEnd, ok := matchBraceSequence(input[i:])
	if !ok {
		return input, false
	}

	// but is it really a sequence?
	braceSeq, ok := parseBraceSequence(input[i : i+seqEnd])
	if !ok {
		return input, false
	}

	// if we get here, then yes it is
	preamble := ""
	preambleStart := findPreambleStart(input, i)
	if preambleStart < i {
		preamble = input[preambleStart:i]
	}
	postscript := ""
	postscriptEnd := findPostscriptEnd(input, i+seqEnd)
	if postscriptEnd > i+seqEnd {
		postscript = input[i+seqEnd : postscriptEnd]
	}

	var exp []string
	if braceSeq.incr > 0 {
		for j := braceSeq.start; j <= braceSeq.end; j += braceSeq.incr {
			exp = append(exp, expandBraceSequence(j, braceSeq.chars, preamble, postscript))
		}
	} else {
		for j := braceSeq.start; j >= braceSeq.end; j += braceSeq.incr {
			exp = append(exp, expandBraceSequence(j, braceSeq.chars, preamble, postscript))
		}
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

	// all done
	return buf.String(), true
}

func matchBracePattern(input string) (int, bool) {
	// are we looking at the start of a pattern?
	if input[0] != '{' {
		return 0, false
	}

	var r rune
	w := 0
	inEscape := false
	braceDepth := 0

	for i := 0; i < len(input); {
		r, w = utf8.DecodeRuneInString(input[i:])

		if inEscape {
			inEscape = false
			i += w
		} else if r == '\\' {
			// skip over escaped character
			inEscape = true
			i += w
		} else if r == '$' {
			varEnd, ok := matchVar(input[i:])
			if ok {
				i += varEnd
			} else {
				i += w
			}
		} else if r == '{' {
			braceDepth++
			i += w
		} else if r == '}' {
			braceDepth--

			i += w

			if braceDepth == 0 {
				return i, true
			}
		} else {
			i += w
		}
	}

	return 0, false
}

func matchBraceSequence(input string) (int, bool) {
	// are we looking at the start of a sequence?
	if input[0] != '{' {
		return 0, false
	}

	// a sequence has the format:
	//
	// {[:alphanum:]..[:alphanum:]} or
	// {[:alphanum:]..[:alphanum:]..[:num:]}
	//
	// no escape chars, no vars to worry about, no nesting either

	var r rune
	w := 0
	braceDepth := 0
	for i := 0; i < len(input); {
		// grab the next character
		r, w = utf8.DecodeRuneInString(input[i:])

		// what are we looking at?
		if r == '{' {
			braceDepth++
			i += w

			// no nesting allowed!
			if braceDepth > 1 {
				return 0, false
			}
		} else if r == '}' {
			braceDepth--
			i += w

			if braceDepth == 0 {
				return i, true
			}
		} else if isSequenceChar(r) {
			i += w
		} else {
			return 0, false
		}
	}

	return 0, false
}

func isSequenceChar(c rune) bool {
	return c == '.' || c == '-' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func parseBracePattern(pattern string) ([]string, bool) {
	var parts []string

	// we can't do a simple `strings.Split()` here, because we have to
	// take nested braces into account

	// how many braces are we in?
	braceDepth := 0

	// where are we?
	start := 1

	var r rune
	w := 0
	inEscape := false

	// find the next pattern
	for i := 0; i < len(pattern); {
		// find the next unicode character
		r, w = utf8.DecodeRuneInString(pattern[i:])

		if inEscape {
			// skip over the escaped character
			inEscape = false
			i += w
		} else if r == '\\' {
			inEscape = true
			i += w
		} else if r == '{' {
			braceDepth++
			i += w
		} else if r == '}' {
			braceDepth--
			// have we reached the end of the pattern?
			if braceDepth == 0 {
				parts = append(parts, pattern[start:i])
				start = i + 1
			}
			i += w
		} else if r == ',' {
			// are we in a sub-pattern?
			if braceDepth == 1 {
				// no, we are not :)
				parts = append(parts, pattern[start:i])
				start = i + 1
			}
			i += w
		} else {
			i += w
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

func parseBraceSequence(pattern string) (braceSequence, bool) {
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

	if isNumericStart && isNumericEnd {
		// all numbers
		retval.start, _ = strconv.Atoi(parts[0])
		retval.end, _ = strconv.Atoi(parts[1])
	} else if isNumericStart != isNumericEnd {
		return braceSequence{}, false
	} else {
		// must be chars
		retval.chars = true
		retval.start = int(parts[0][0])
		retval.end = int(parts[1][0])
	}

	// do we have an incr element?
	if len(parts) == 3 {
		incr, err := strconv.Atoi(parts[2])
		if err != nil {
			return braceSequence{}, false
		}
		retval.incr = incr
	} else {
		retval.incr = 1
	}

	// now we just need to make sure the incr element goes in the same
	// direction as the range itself does
	if retval.start < retval.end {
		if retval.incr < 0 {
			retval.incr = 0 - retval.incr
		}
	} else if retval.incr > 0 {
		retval.incr = 0 - retval.incr
	}

	// all done
	return retval, true
}
