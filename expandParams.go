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
	// ${@:offset} -> expansion of positional params, starting at offset
	paramExpandPositionalParamsFromOffset
	// ${@:offset:length} -> expansion of 'length' number of positional params, starting at offset
	paramExpandPositionalParamsFromOffsetLength
	// ${!prefix*} -> return a list of names matching the given prefix
	paramExpandPrefixNames
	// ${!prefix@} -> return a list of names matching the given prefix
	paramExpandPrefixNamesDoubleQuoted
	// ${#var} -> length of value of var
	paramExpandParamLength
	// ${#*} -> number of positional parameters
	paramExpandNoOfPositionalParams
	// ${var#word} -> value of var, with shortest matching prefix of word removed
	paramExpandRemovePrefixShortestMatch
	// ${var##word} -> value of var, with longest matching prefix of word removed
	paramExpandRemovePrefixLongestMatch
	// ${var%suffix} -> value of var, with shortest matching suffix removed
	paramExpandRemoveSuffixShortestMatch
	// ${var%%suffix} -> value of var, with longest matching suffix removed
	paramExpandRemoveSuffixLongestMatch
	// ${*%suffix} -> all positional params, with shorted matching suffix removed
	paramExpandRemoveSuffixAllPositionalParamsShortestMatch
	// ${*%%suffix} -> all positional params, with longest matching suffix removed
	paramExpandRemoveSuffixAllPositionalParamsLongestMatch
	// ${var/old/new} -> value of var, with first occurance of old replaced with new
	paramExpandSearchReplaceLongestFirstMatch
	// ${var//old/new} -> value of var, with all occurances of old replaced with new
	paramExpandSearchReplaceLongestAllMatches
	// ${var/#old/new} -> value of var, with old replaced with new if the string starts with old
	paramExpandSearchReplaceLongestPrefix
	// ${var/%old/new} -> value of var, with old replaced with new if the string ends with old
	paramExpandSearchReplaceLongestSuffix
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
	var opType int
	var opEnd int
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
		if input[maxInput:] == "*}" {
			return paramDesc{
				kind:  paramExpandPrefixNames,
				parts: []string{input[3:maxInput]},
			}, true
		} else if input[maxInput:] == "@}" {
			return paramDesc{
				kind:  paramExpandPrefixNamesDoubleQuoted,
				parts: []string{input[3:maxInput]},
			}, true
		}
	}

	// special case - handle ${#parameter} here
	if input[0:3] == "${#" && (isNameStartChar(input[3]) || isNumericStartChar(input[3]) || isShellSpecialChar(input[3])) {
		// we don't check the boolean return value, because we're 100%
		// guaranteed to match the 1st char
		paramType, paramEnd, _ = matchParam(input, 3)

		// there can't be anything else in the input string
		if paramEnd == maxInput {
			switch paramType {
			case paramTypeName:
				return paramDesc{
					kind:  paramExpandParamLength,
					parts: []string{input[3:inputLen]},
				}, true
			default:
				return paramDesc{
					kind:  paramExpandParamLength,
					parts: []string{"$" + input[3:inputLen]},
				}, true
			}
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
	//
	// this is not the easy question it should be
	if input[2] == '!' {
		// special case - indirect expansion is not supported for '$!'
		// according to my testing
		if input[3] == '!' {
			return paramDesc{}, false
		}

		// according to my testing, '${!' is *always* interpreted
		// as indirection by POSIX shells
		//
		// if you come up with test cases that prove otherwise,
		// I want to know!
		retval.indirect = true
		start++
	}

	// this helps us get out of the indirection check
	// afterIndirectionCheck:

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
	opType, opEnd, ok = matchParamOp(input, opStart)
	if !ok {
		return paramDesc{}, false
	}

	switch opType {
	case paramOpUseDefaultValue:
		retval.kind = paramExpandWithDefaultValue
		if opEnd < maxInput {
			retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		}
		return retval, true
	case paramOpAssignDefaultValue:
		retval.kind = paramExpandSetDefaultValue
		if opEnd < maxInput {
			retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		}
		return retval, true
	case paramOpWriteError:
		retval.kind = paramExpandWriteError
		if opEnd < maxInput {
			retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		}
		return retval, true
	case paramOpUseAlternativeValue:
		retval.kind = paramExpandAlternativeValue
		if opEnd < maxInput {
			retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		}
		return retval, true
	case paramOpSubstring:
		// there must be *something* after the op
		if opEnd == maxInput {
			return paramDesc{}, false
		}

		// must be a substring operation ... but which one?
		parts := strings.Split(input[opEnd+1:inputLen], ":")
		if len(parts) > 2 {
			return paramDesc{}, false
		}
		for _, part := range parts {
			// offset and length can both be negative
			// although until we have arithmetic expansion, there's no
			// way to pass a negative offset into this function
			if !isSignedNumericString(part) {
				return paramDesc{}, false
			}
		}

		// special case - positional parameter expansion
		if retval.parts[0] == "$@" {
			if len(parts) == 1 {
				retval.kind = paramExpandPositionalParamsFromOffset
			} else {
				retval.kind = paramExpandPositionalParamsFromOffsetLength
			}
		} else if len(parts) == 1 {
			retval.kind = paramExpandSubstring
		} else {
			retval.kind = paramExpandSubstringLength
		}
		retval.parts = append(retval.parts, parts...)
		return retval, true
	case paramOpRemoveShortestSuffix:
		// assume shortest match variant for now
		retval.kind = paramExpandRemoveSuffixShortestMatch
		if opEnd < maxInput {
			retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		} else {
			retval.parts = append(retval.parts, "")
		}
		return retval, true

	case paramOpRemoveLongestSuffix:
		retval.kind = paramExpandRemoveSuffixLongestMatch
		if opEnd < maxInput {
			retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		} else {
			retval.parts = append(retval.parts, "")
		}
		return retval, true

	case paramOpRemoveShortestPrefix:
		retval.kind = paramExpandRemovePrefixShortestMatch
		if opEnd < maxInput {
			retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		} else {
			retval.parts = append(retval.parts, "")
		}
		return retval, true

	case paramOpRemoveLongestPrefix:
		retval.kind = paramExpandRemovePrefixLongestMatch
		if opEnd < maxInput {
			retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		} else {
			retval.parts = append(retval.parts, "")
		}
		return retval, true

	case paramOpSearchReplace:
		// according to my testing, if there's nothing after the operand,
		// UNIX shells simply do an expand-to-value
		if opEnd == maxInput {
			retval.kind = paramExpandToValue
			return retval, true
		}

		// things get messy here, because the first char of `pattern`
		// changes the behaviour ... and can be an unescaped '/'
		switch input[opEnd+1] {
		case '/':
			// according to my testing, if there's nothing after the
			// 'all matches' /, UNIX shells effectively do an expand-to-value
			if opEnd+1 == maxInput {
				retval.kind = paramExpandToValue
				return retval, true
			}

			retval.kind = paramExpandSearchReplaceLongestAllMatches
			retval.parts = append(retval.parts, strings.Split(input[opEnd+2:inputLen], "/")...)

			// if the replace string is missing, default is an empty string
			if len(retval.parts) < 3 {
				retval.parts = append(retval.parts, "")
			}

			// all done
			return retval, true
		case '%':
			// according to my testing, if there's nothing after the
			// 'all matches' /, UNIX shells effectively do an expand-to-value
			if opEnd+1 == maxInput {
				retval.kind = paramExpandToValue
				return retval, true
			}

			retval.kind = paramExpandSearchReplaceLongestSuffix
			retval.parts = append(retval.parts, strings.Split(input[opEnd+2:inputLen], "/")...)

			// if the replace string is missing, default is an empty string
			if len(retval.parts) < 3 {
				retval.parts = append(retval.parts, "")
			}
			return retval, true
		case '#':
			// according to my testing, if there's nothing after the
			// 'all matches' /, UNIX shells effectively do an expand-to-value
			if opEnd+1 == maxInput {
				retval.kind = paramExpandToValue
				return retval, true
			}

			retval.kind = paramExpandSearchReplaceLongestPrefix
			retval.parts = append(retval.parts, strings.Split(input[opEnd+2:inputLen], "/")...)
			// if the replace string is missing, default is an empty string
			if len(retval.parts) < 3 {
				retval.parts = append(retval.parts, "")
			}
			return retval, true
		default:
			// this is the easy bit!
			retval.kind = paramExpandSearchReplaceLongestFirstMatch
			retval.parts = append(retval.parts, strings.Split(input[opEnd+1:inputLen], "/")...)
			// if the replace string is missing, default is an empty string
			if len(retval.parts) < 3 {
				retval.parts = append(retval.parts, "")
			}
			return retval, true
		}

	case paramOpUppercaseFirstChar:
		// according to my testing, if there's nothing after the operand,
		// UNIX shells simply do an expand-to-value
		if opStart == maxInput {
			retval.kind = paramExpandToValue
			return retval, true
		}
		retval.kind = paramExpandUppercaseFirstChar
		retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		return retval, true

	case paramOpUppercaseAllMatches:
		// according to my testing, if there's nothing after the operand,
		// UNIX shells simply do an expand-to-value
		if opEnd == maxInput {
			retval.kind = paramExpandToValue
			return retval, true
		}

		retval.kind = paramExpandUppercaseAllChars
		retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		return retval, true

	case paramOpLowercaseFirstChar:
		// according to my testing, if there's nothing after the operand,
		// UNIX shells simply do an expand-to-value
		if opEnd == maxInput {
			retval.kind = paramExpandToValue
			return retval, true
		}
		retval.kind = paramExpandLowercaseFirstChar
		retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		return retval, true

	case paramOpLowercaseAllMatches:
		// according to my testing, if there's nothing after the operand,
		// UNIX shells simply do an expand-to-value
		if opEnd == maxInput {
			retval.kind = paramExpandToValue
			return retval, true
		}

		retval.kind = paramExpandLowercaseAllChars
		retval.parts = append(retval.parts, input[opEnd+1:inputLen])
		return retval, true

	case paramOpDescribeFlags:
		retval.kind = paramExpandDescribeFlags
		return retval, true
	case paramOpDeclare:
		retval.kind = paramExpandAsDeclare
		return retval, true
	case paramOpEscape:
		retval.kind = paramExpandEscaped
		return retval, true
	case paramOpExpandAsPrompt:
		retval.kind = paramExpandAsPrompt
		return retval, true
	case paramOpExpandDoubleQuotes:
		retval.kind = paramExpandSingleQuoted
		return retval, true

	default:
		// unknown or unsupported operand
		return paramDesc{}, false
	}
}
