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

// expandParams will expand any ${VAR} or $VAR
//
// $var -> value of var
// ${var} -> value of var
// ${var:-word} -> value of var (if set); expansion of word otherwise
// ${var:=word} -> value of var (if set); otherwise var is set to the expansion of word
// ${var:?word} -> value of var (if set); otherwise error written to stderr
// ${var:+word} -> empty string if var empty/unset; otherwise expansion of word
// ${var:offset} -> substring of var (if set), starting from offset; otherwise empty string
// ${var:offset:length} -> same as both, except also controlling length of substring
// ${!prefix*} -> return a list of names matching the given prefix
// ${#var} -> length of value of var
// ${#*} -> number of positional parameters
// ${var#word} -> value of var, with shortest match of word removed
// ${var##word} -> value of var, with longest match of word removed
// ${var%suffix} -> value of var, with shortest matching suffix removed
// ${var%%suffix} -> value of var, with longest matching suffix removed
// ${*%suffix} -> all positional params, with shorted matching suffix removed
// ${*%%suffix} -> all positional params, with longest matching suffix removed
// ${var/old/new} -> value of var, with occurances of old replaced with new
// ${*/old/new} -> all positional params, with occurances of old replaced with new
// ${var^pattern} -> value of var, with first char set to uppercase if they are in pattern
// ${var^^pattern} -> value of var, with any char set to uppercase if they are in pattern
// ${var,pattern} -> value of var, with first char set to lowercase if they are in pattern
// ${var,,pattern} -> value of var, with any char set to lowercase if they are in pattern
// ${var@a} -> a set of flags describing var
// ${var@A} -> not supported?
// ${var@E} -> escaped value of var - probably too dangerous to support
// ${var@P} -> expanded prompt string - not supported
// ${var@Q} -> quoted value of var - probably too dangerous to support
//
// traditional shell special parameters are treated as a special case:
//
// - normally, the '$' prefix is removed before calling the lookupVar
//   (e.g. "$HOME" becomes lookupVar("HOME"))
// - shell special params keep their '$' prefix when we call the lookupVar
//   (e.g) "$*" becomes lookupVar("$*")
//
// supported traditional shell params are:
//
// $# - number of positional parameters
// $* - all positional parameters as a single string
// $1, $2 ... - individual positional parameters
// $? - exit value of last command run
// $$ - PID of current process
// $0 - argv[0] of current process
// $! - PID of last created background process
// $- - flags passed into current process
// $@ - all positional params as an array
//
// it's up to the caller to ensure lookupVar() can provide a value for any
// of these params
func expandParameter(input string, lookupVar LookupVar) string {
	// we expand in a strictly left-to-right manner
	for i := 0; i < len(input); i++ {
		if input[i] == '\\' {
			// skip over escaped characters
			i++
		} else if input[i] == '$' {
			varEnd, ok := matchVar(input, i)
			if ok {
				_, ok := parseParameter(input[i : varEnd+1])
				if ok {
					i = varEnd
				}
			}
		}
	}

	return input
}

const (
	// we want '0' to mean something went wrong
	paramExpandNotSupported = iota
	// $var -> value of var (if set); empty string otherwise
	// ${var} -> value of var (if set); empty string otherwise
	paramExpandToValue
	// ${var:-word} -> value of var (if set); expansion of word otherwise
	paramExpandWithDefaultValue
	// ${var:=word} -> value of var (if set); otherwise var is set to the expansion of word
	paramExpandSetDefaultValue
	// ${var:?word} -> value of var (if set); otherwise error written to stderr
	paramExpandWriteError
	// ${var:+word} -> empty string if var empty/unset; otherwise expansion of word
	paramExpandAlternativeValue
	// ${var:offset} -> substring of var (if set), starting from offset; otherwise empty string
	paramExpandSubstring
	// ${var:offset:length} -> same as both, except also controlling length of substring
	paramExpandSubstringLength
	// ${!prefix*} -> return a list of names matching the given prefix
	paramExpandPrefixNames
	// ${!prefix@} -> return a list of names matching the given prefix
	paramExpandPrefixNamesDoubleQuoted
	// ${#var} -> length of value of var
	paramExpandLen
	// ${#*} -> number of positional parameters
	paramExpandNoOfPositionalParams
	// ${var#word} -> value of var, with shortest match of word removed
	paramExpandRemoveWordShortestMatch
	// ${var##word} -> value of var, with longest match of word removed
	paramExpandRemoveWordLongestMatch
	// ${var%suffix} -> value of var, with shortest matching suffix removed
	paramExpandRemoveSuffixShortestMatch
	// ${var%%suffix} -> value of var, with longest matching suffix removed
	paramExpandRemoveSuffixLongestMatch
	// ${*%suffix} -> all positional params, with shorted matching suffix removed
	paramExpandRemoveSuffixAllPositionalParamsShortestMatch
	// ${*%%suffix} -> all positional params, with longest matching suffix removed
	paramExpandRemoveSuffixAllPositionalParamsLongestMatch
	// ${var/old/new} -> value of var, with first occurance of old replaced with new
	paramExpandSearchReplaceFirstMatch
	// ${var//old/new} -> value of var, with all occurances of old replaced with new
	paramExpandSearchReplaceAllMatches
	// ${var/#old/new} -> value of var, with old replaced with new if the string starts with old
	paramExpandSearchReplacePrefix
	// ${var/%old/new} -> value of var, with old replaced with new if the string ends with old
	paramExpandSearchReplaceSuffix
	// ${*/old/new} -> all positional params, with occurances of old replaced with new
	paramExpandAllPositionalParamsSearchReplace
	// ${var^pattern} -> value of var, with first char set to uppercase if they are in pattern
	paramExpandUppercaseFirstChar
	// ${var^^pattern} -> value of var, with any char set to uppercase if they are in pattern
	paramExpandUppercaseAllChars
	// ${var,pattern} -> value of var, with first char set to lowercase if they are in pattern
	paramExpandLowercaseFirstChar
	// ${var,,pattern} -> value of var, with any char set to lowercase if they are in pattern
	paramExpandLowercaseAllChars
	// ${var@a} -> a set of flags describing var
	paramExpandDescribeFlags
	// ${var@A} -> exapnded value of var as declare statement - not supported?
	paramExpandAsDeclare
	// ${var@E} -> escaped value of var - escaped how, exactly?
	paramExpandEscaped
	// ${var@P} -> expanded prompt string - not supported
	paramExpandAsPrompt
	// ${var@Q} -> single quoted value of var
	paramExpandSingleQuoted
)

type paramDesc struct {
	kind     int
	parts    []string
	indirect bool
}

func parseParameter(input string) (paramDesc, bool) {
	// shorthand
	inputLen := len(input)
	maxInput := inputLen - 1

	// we'll use these throughout the function
	var paramType int
	var paramEnd int
	var ok bool
	var retval paramDesc

	// make sure we're looking at something that has the shape of a parameter
	if input[0] != '$' {
		return paramDesc{}, false
	}
	if input[1] == '{' && input[maxInput] != '}' {
		return paramDesc{}, false
	}
	if input[1] != '{' && input[maxInput] == '}' {
		return paramDesc{}, false
	}

	// is the string wrapped in braces?
	if input[1] != '{' && input[maxInput] != '}' {
		// no
		paramType, paramEnd, ok = matchParam(input, 1)
		if !ok {
			return paramDesc{}, false
		}
		if paramEnd != maxInput {
			return paramDesc{}, false
		}

		switch paramType {
		case paramTypeName:
			return paramDesc{
				kind:  paramExpandToValue,
				parts: []string{input[1:inputLen]},
			}, true
		default:
			return paramDesc{
				kind:  paramExpandToValue,
				parts: []string{input},
			}, true
		}
	}

	// at this point, we know we're looking at an input string wrapped
	// in braces
	maxInput--
	inputLen--

	// special case - handle *all* single-char names here
	//
	// this greatly simplifies the code later on
	if len(input) == 4 {
		paramType, paramEnd, ok = matchParam(input, 2)
		if !ok {
			return paramDesc{}, false
		}
		if paramEnd != maxInput {
			return paramDesc{}, false
		}

		switch paramType {
		case paramTypeName:
			return paramDesc{
				kind:  paramExpandToValue,
				parts: []string{input[2:inputLen]},
			}, true
		default:
			return paramDesc{
				kind:  paramExpandToValue,
				parts: []string{"$" + input[2:inputLen]},
			}, true
		}
	}

	// special case - handle positional params
	if isNumericStringWithoutLeadingZero(input[2:inputLen]) {
		return paramDesc{
			kind:  paramExpandToValue,
			parts: []string{"$" + input[2:inputLen]},
		}, true
	}

	// special case - handle ${!prefix*} and ${prefix@} here
	if input[0:3] == "${!" {
		if input[len(input)-2:] == "*}" {
			return paramDesc{
				kind:  paramExpandPrefixNames,
				parts: []string{input[3:maxInput]},
			}, true
		} else if input[len(input)-2:] == "@}" {
			return paramDesc{
				kind:  paramExpandPrefixNames,
				parts: []string{input[3:maxInput]},
			}, true
		}
	}

	// special case - handle ${#parameter} here
	if input[0:3] == "${#" && (isNameStartChar(input[3]) || isNumericStartChar(input[3]) || isShellSpecialChar(input[3])) {
		paramType, paramEnd, ok = matchParam(input, 3)
		if !ok {
			return paramDesc{}, false
		}
		if paramEnd != maxInput {
			return paramDesc{}, false
		}

		switch paramType {
		case paramTypeName:
			return paramDesc{
				kind:  paramExpandToValue,
				parts: []string{input[3 : paramEnd+1]},
			}, true
		default:
			return paramDesc{
				kind:  paramExpandToValue,
				parts: []string{"$" + input[3:paramEnd+1]},
			}, true
		}
	}

	// at this point, what's left is everything of the form:
	//
	// ${[!]parameter<op>[<op-specific parts>]}
	//
	// we just have to work through them

	// where are we going to start looking for the name of the param?
	start := 2

	// do we have indirect expansion?
	if input[2] == '!' && (isNameStartChar(input[3]) || isNumericStartChar(input[3]) || isShellSpecialChar(input[3])) {
		// special case - indirect expansion is not supported for '$!'
		// according to my testing
		if input[3] == '!' {
			return paramDesc{}, false
		}

		retval.indirect = true
		start++
	}

	// the param name must be valid
	paramType, paramEnd, ok = matchParam(input, start)
	if !ok {
		return paramDesc{}, false
	}
	switch paramType {
	case paramTypeName:
		retval.parts = append(retval.parts, input[start:paramEnd+1])
	default:
		retval.parts = append(retval.parts, "$"+input[start:paramEnd+1])
	}

	// special case - is that it?
	if paramEnd == maxInput {
		retval.kind = paramExpandToValue
		return retval, true
	}

	// what kind of operator do we have?
	//
	// remember that it may be the last part of the parameter expansion
	opStart := paramEnd + 1
	switch input[opStart] {
	case ':':
		if opStart == maxInput {
			// cannot have this as the last character for parameter expansion
			return paramDesc{}, false
		}

		switch input[opStart+1] {
		case '-':
			retval.kind = paramExpandWithDefaultValue
			if opStart < maxInput {
				retval.parts = append(retval.parts, input[opStart+2:inputLen])
			}
			return retval, true
		case '=':
			retval.kind = paramExpandSetDefaultValue
			if opStart < maxInput {
				retval.parts = append(retval.parts, input[opStart+2:inputLen])
			}
			return retval, true
		case '?':
			retval.kind = paramExpandWriteError
			if opStart < maxInput {
				retval.parts = append(retval.parts, input[opStart+2:inputLen])
			}
			return retval, true
		case '+':
			retval.kind = paramExpandAlternativeValue
			if opStart < maxInput {
				retval.parts = append(retval.parts, input[opStart+2:inputLen])
			}
			return retval, true
		default:
			// must be a substring operation ... but which one?
			parts := strings.Split(input[opStart+1:inputLen], ":")
			if len(parts) > 2 {
				return paramDesc{}, false
			}
			for _, part := range parts {
				if !isNumericString(part) {
					return paramDesc{}, false
				}
			}
			if len(parts) == 1 {
				retval.kind = paramExpandSubstring
			} else {
				retval.kind = paramExpandSubstringLength
			}
			retval.parts = append(retval.parts, parts...)
			return retval, true
		}
	case '%':
		// assume shortest match variant for now
		retval.kind = paramExpandRemoveSuffixShortestMatch
		if opStart == maxInput {
			retval.parts = append(retval.parts, "")
			return retval, true
		}

		// is it actually longest-match variant?
		if input[opStart+1] == '%' {
			retval.kind = paramExpandRemoveSuffixLongestMatch
			if opStart < maxInput {
				retval.parts = append(retval.parts, input[opStart+2:inputLen])
			}
			return retval, true
		}

		retval.parts = append(retval.parts, input[opStart+1:inputLen], "")
		return retval, true
	case '#':
		// assume shortest match variant for now
		retval.kind = paramExpandRemoveWordShortestMatch
		if opStart == maxInput {
			retval.parts = append(retval.parts, "")
			return retval, true
		}

		// is it actually longest-match variant?
		if input[opStart+1] == '#' {
			retval.kind = paramExpandRemoveWordLongestMatch
			if opStart < maxInput {
				retval.parts = append(retval.parts, input[opStart+2:inputLen])
			}
			return retval, true
		}

		retval.parts = append(retval.parts, input[opStart+1:inputLen], "")
		return retval, true
	case '/':
		// according to my testing, if there's nothing after the operand,
		// UNIX shells simply do an expand-to-value
		if opStart == maxInput {
			retval.kind = paramExpandToValue
			return retval, true
		}

		// things get messy here, because the first char of `pattern`
		// changes the behaviour ... and can be an unescaped '/'
		switch input[opStart+1] {
		case '/':
			// are we looking at a pattern that starts with '/'?
			if strings.ContainsRune(input[opStart+2:inputLen], '/') {
				// yes, we are
				retval.kind = paramExpandSearchReplaceAllMatches
				retval.parts = append(retval.parts, strings.Split(input[opStart+2:inputLen], "/")...)
			} else {
				retval.kind = paramExpandSearchReplaceFirstMatch
				retval.parts = append(retval.parts, strings.Split(input[opStart+1:inputLen], "/")...)
			}

			// if the replace string is missing, default is an empty string
			if len(retval.parts) < 3 {
				retval.parts = append(retval.parts, "")
			}

			// all done
			return retval, true
		case '%':
			retval.kind = paramExpandSearchReplacePrefix
			retval.parts = append(retval.parts, strings.Split(input[opStart+2:inputLen], "/")...)
			// if the replace string is missing, default is an empty string
			if len(retval.parts) < 3 {
				retval.parts = append(retval.parts, "")
			}
			return retval, true
		case '#':
			retval.kind = paramExpandSearchReplaceSuffix
			retval.parts = append(retval.parts, strings.Split(input[opStart+2:inputLen], "/")...)
			// if the replace string is missing, default is an empty string
			if len(retval.parts) < 3 {
				retval.parts = append(retval.parts, "")
			}
			return retval, true
		default:
			// this is the easy bit!
			retval.kind = paramExpandSearchReplaceFirstMatch
			retval.parts = append(retval.parts, strings.Split(input[opStart+1:inputLen], "/")...)
			// if the replace string is missing, default is an empty string
			if len(retval.parts) < 3 {
				retval.parts = append(retval.parts, "")
			}
			return retval, true
		}
	case '^':
		// according to my testing, if there's nothing after the operand,
		// UNIX shells simply do an expand-to-value
		if opStart == maxInput {
			retval.kind = paramExpandToValue
			return retval, true
		}

		switch input[opStart+1] {
		case '^':
			if opStart+2 < maxInput {
				retval.kind = paramExpandUppercaseAllChars
				retval.parts = append(retval.parts, input[opStart+2:inputLen])
			} else {
				// nothing after the operand, so once again we default to
				// expand-to-value
				retval.kind = paramExpandToValue
				retval.parts = append(retval.parts, "")
			}
			return retval, true
		default:
			retval.kind = paramExpandUppercaseFirstChar
			retval.parts = append(retval.parts, input[opStart+2:inputLen])
			return retval, true
		}
	case ',':
		// according to my testing, if there's nothing after the operand,
		// UNIX shells simply do an expand-to-value
		if opStart == maxInput {
			retval.kind = paramExpandToValue
			return retval, true
		}

		switch input[opStart+1] {
		case ',':
			if opStart+2 < len(input)-1 {
				retval.kind = paramExpandLowercaseAllChars
				retval.parts = append(retval.parts, input[opStart+2:inputLen])
			} else {
				// nothing after the operand, so once again we default to
				// expand-to-value
				retval.kind = paramExpandToValue
				retval.parts = append(retval.parts, "")
			}
			return retval, true
		default:
			retval.kind = paramExpandLowercaseFirstChar
			retval.parts = append(retval.parts, input[opStart+2:inputLen])
			return retval, true
		}
	case '@':
		if opStart == maxInput {
			return paramDesc{}, false
		}

		// using a string comparison here, for future expansion and to
		// catch any ops that are too long
		switch input[opStart+1 : inputLen] {
		case "a":
			retval.kind = paramExpandDescribeFlags
			return retval, true
		case "A":
			retval.kind = paramExpandAsDeclare
			return retval, true
		case "E":
			retval.kind = paramExpandSingleQuoted
			return retval, true
		case "P":
			retval.kind = paramExpandAsPrompt
			return retval, true
		case "Q":
			retval.kind = paramExpandSingleQuoted
			return retval, true
		default:
			// unknown or unsupported operand
			return paramDesc{}, false
		}
	}

	// if we somehow end up here, we don't understand the parameter
	return paramDesc{}, false
}
