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
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	glob "github.com/ganbarodigital/go_glob"
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
func expandParameters(input string, varFuncs VarFuncs) (string, error) {
	// keep track of whether we're dealing with an escaped character
	// or not
	inEscape := false

	// keep track of the end of the last param we matched
	varEnd := -1

	// and this will be where we build up our return value
	var buf strings.Builder

	// we expand in a strictly left-to-right manner
	var c rune
	w := 0
	for i := 0; i < len(input); {
		c, w = utf8.DecodeRuneInString(input[i:])
		if inEscape {
			// skip over escaped characters
			inEscape = false
			buf.WriteRune(c)
			i += w
		} else if c == '\\' {
			// skip over escaped characters
			inEscape = true
			i += w
		} else if c == '$' {
			var ok bool
			varEnd, ok = matchVar(input[i:])
			if ok {
				varEnd += i
				paramDesc, ok := parseParameter(input[i:varEnd])
				if !ok {
					buf.WriteRune(c)
					i += w
					continue
				}

				replacement, err := expandParameter(input[i:varEnd], paramDesc, varFuncs)
				if err != nil {
					return input, err
				}

				buf.WriteString(replacement)

				i = varEnd
			} else {
				buf.WriteRune(c)
				i += w
			}
		} else {
			buf.WriteRune(c)
			i += w
		}
	}

	return buf.String(), nil
}

type paramExpandFunc func(string, string, paramDesc, VarFuncs) (string, bool, error)

func expandParameter(original string, paramDesc paramDesc, varFuncs VarFuncs) (string, error) {
	paramExpandFuncs := map[int]paramExpandFunc{
		paramExpandToValue:                   expandParamToValue,
		paramExpandWithDefaultValue:          expandParamWithDefaultValue,
		paramExpandSetDefaultValue:           expandParamSetDefaultValue,
		paramExpandWriteError:                expandParamWriteError,
		paramExpandAlternativeValue:          expandParamAlternativeValue,
		paramExpandSubstring:                 expandParamSubstring,
		paramExpandSubstringLength:           expandParamSubstringLength,
		paramExpandPrefixNames:               expandParamPrefixNames,
		paramExpandPrefixNamesDoubleQuoted:   expandParamPrefixNames,
		paramExpandParamLength:               expandParamLength,
		paramExpandRemovePrefixShortestMatch: expandParamRemovePrefixShortestMatch,
		paramExpandRemovePrefixLongestMatch:  expandParamRemovePrefixLongestMatch,
		paramExpandRemoveSuffixShortestMatch: expandParamRemoveSuffixShortestMatch,
		paramExpandRemoveSuffixLongestMatch:  expandParamRemoveSuffixLongestMatch,
		paramExpandUppercaseFirstChar:        expandParamUppercaseFirstChar,
		paramExpandUppercaseAllChars:         expandParamUppercaseAllChars,
		paramExpandLowercaseFirstChar:        expandParamLowercaseFirstChar,
		paramExpandLowercaseAllChars:         expandParamLowercaseAllChars,
	}

	// what we will (eventually) send back
	var retval []string

	// and, because we may be building it up bit by bit, we need somewhere
	// to store it temporarily
	var buf string

	// step 1: we need to expand the paramName first, to support any
	// possible use of indirection
	paramName, ok := expandParamName(paramDesc, varFuncs.LookupVar)
	if !ok {
		return "", nil
	}

	// special case
	if paramDesc.kind == paramExpandNoOfPositionalParams {
		buf, ok = varFuncs.LookupVar("$#")
		return buf, nil
	}

	// step 2: we need to feed that into all the different ways that
	// parameters can be expanded in strings
	//
	// this is complicated by some parameters ($*, $@, and arrays if we
	// ever add support for them in the future) having the expansion applied
	// to each part of their value
	for paramValue := range expandParamValue(paramName, varFuncs.LookupVar) {
		expandFunc, ok := paramExpandFuncs[paramDesc.kind]
		if !ok {
			return "", nil
		}

		var err error
		buf, ok, err = expandFunc(paramName, paramValue, paramDesc, varFuncs)
		if err != nil {
			return "", err
		}

		if len(buf) > 0 {
			retval = append(retval, buf)
		}
	}

	// if we get here, then yes, we are happy
	return strings.Join(retval, " "), nil
}

func expandParamName(paramDesc paramDesc, lookupVar LookupVar) (string, bool) {
	varName := paramDesc.parts[0]
	ok := true
	if paramDesc.indirect {
		varName, ok = lookupVar(varName)
	}

	return varName, ok
}

func expandParamToValue(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	// nothing else to do
	return paramValue, true, nil
}

func expandParamWithDefaultValue(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	// do we need to return the default value?
	if paramValue != "" {
		return paramValue, true, nil
	}

	retval, err := expandWord(paramDesc.parts[1], varFuncs)
	return retval, true, err
}

func expandParamSetDefaultValue(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	// do we need to do anything?
	if paramValue != "" {
		return paramValue, true, nil
	}

	// at this point, we need to assign a new value
	word, err := expandWord(paramDesc.parts[1], varFuncs)
	if err != nil {
		return "", false, err
	}
	err = varFuncs.AssignToVar(paramName, word)
	if err != nil {
		return "", false, err
	}

	// all done
	retval, success := varFuncs.LookupVar(paramName)
	return retval, success, nil
}

func expandParamWriteError(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	// do we have a value?
	if paramValue != "" {
		return paramValue, true, nil
	}

	word, err := expandWord(paramDesc.parts[1], varFuncs)
	if err != nil {
		return "", false, err
	}

	return paramName + ": " + word, true, nil
}

func expandParamAlternativeValue(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	// do we need to return the alternative value?
	if paramValue == "" {
		return paramValue, true, nil
	}

	word, err := expandWord(paramDesc.parts[1], varFuncs)
	if err != nil {
		return "", false, err
	}

	return word, true, nil
}

func expandParamSubstring(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	start, err := strconv.Atoi(paramDesc.parts[1])
	if err != nil {
		return paramValue, true, nil
	}

	// range overflow?
	if start > len(paramValue) {
		return "", true, nil
	}

	return paramValue[start:], true, nil
}

func expandParamSubstringLength(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	// where do we start from?
	start, err := strconv.Atoi(paramDesc.parts[1])
	if err != nil {
		return paramValue, true, nil
	}
	// range overflow?
	if start > len(paramValue) {
		return "", true, nil
	}

	// and where do we end?
	amount, err := strconv.Atoi(paramDesc.parts[2])
	if err != nil {
		return "", false, nil
	}
	end := start + amount

	// watch out for this range overflowing too!
	if end > len(paramValue) {
		end = len(paramValue)
	}

	return paramValue[start:end], true, nil
}

func expandParamPrefixNames(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	varNames := varFuncs.MatchVarNames(paramName)
	sort.Strings(varNames)
	return strings.Join(varNames, " "), true, nil
}

func expandParamLength(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	return strconv.Itoa(len(paramValue)), true, nil
}

func expandParamRemovePrefixShortestMatch(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	g := glob.NewGlob(paramDesc.parts[1])

	pos, success, err := g.MatchShortestPrefix(paramValue)
	if err != nil {
		return "", false, err
	}
	if success {
		return paramValue[pos:], true, nil
	}

	return paramValue, true, nil
}

func expandParamRemovePrefixLongestMatch(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	g := glob.NewGlob(paramDesc.parts[1])

	pos, success, err := g.MatchLongestPrefix(paramValue)
	if err != nil {
		return "", false, err
	}
	if success {
		return paramValue[pos:], true, nil
	}

	return paramValue, true, nil
}

func expandParamRemoveSuffixShortestMatch(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	g := glob.NewGlob(paramDesc.parts[1])

	pos, success, err := g.MatchShortestSuffix(paramValue)
	if err != nil {
		return "", false, err
	}
	if success {
		if pos < len(paramValue) {
			return paramValue[:pos], true, nil
		}
		return paramValue, true, nil
	}

	return paramValue, true, nil
}

func expandParamRemoveSuffixLongestMatch(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	g := glob.NewGlob(paramDesc.parts[1])

	pos, success, err := g.MatchLongestSuffix(paramValue)
	if err != nil {
		return "", false, err
	}
	if success {
		// it is impossible for 'pos' to be out-of-bounds
		return paramValue[:pos], true, nil
	}

	return paramValue, true, nil
}

func expandParamUppercaseFirstChar(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	for pos, firstChar := range paramValue {
		// empty pattern?
		if len(paramDesc.parts[1]) == 0 {
			return string(unicode.ToUpper(firstChar)) + paramValue[pos+1:], true, nil
		}

		g := glob.NewGlob(paramDesc.parts[1])
		success, err := g.Match(string(firstChar))
		if err != nil {
			return "", false, err
		}
		if success {
			return string(unicode.ToUpper(firstChar)) + paramValue[pos+1:], true, nil
		}

		return paramValue, true, nil
	}

	// empty value
	return "", true, nil
}

func expandParamUppercaseAllChars(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	// special case
	if len(paramDesc.parts[1]) == 0 {
		return strings.ToUpper(paramValue), true, nil
	}

	// we have to do this the old-fashioned way
	var buf strings.Builder
	g := glob.NewGlob(paramDesc.parts[1])

	for _, c := range paramValue {
		success, err := g.Match(string(c))
		if err != nil {
			return "", false, err
		}
		if success {
			buf.WriteRune(unicode.ToUpper(c))
		} else {
			buf.WriteRune(c)
		}
	}

	// all done
	return buf.String(), true, nil
}

func expandParamLowercaseFirstChar(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	for pos, firstChar := range paramValue {
		// empty pattern?
		if len(paramDesc.parts[1]) == 0 {
			return string(unicode.ToLower(firstChar)) + paramValue[pos+1:], true, nil
		}

		g := glob.NewGlob(paramDesc.parts[1])
		success, err := g.Match(string(firstChar))
		if err != nil {
			return "", false, err
		}
		if success {
			return string(unicode.ToLower(firstChar)) + paramValue[pos+1:], true, nil
		}

		return paramValue, true, nil
	}

	// empty value
	return "", true, nil
}

func expandParamLowercaseAllChars(paramName, paramValue string, paramDesc paramDesc, varFuncs VarFuncs) (string, bool, error) {
	// special case
	if len(paramDesc.parts[1]) == 0 {
		return strings.ToLower(paramValue), true, nil
	}

	// we have to do this the old-fashioned way
	var buf strings.Builder
	g := glob.NewGlob(paramDesc.parts[1])

	for _, c := range paramValue {
		success, err := g.Match(string(c))
		if err != nil {
			return "", false, err
		}
		if success {
			buf.WriteRune(unicode.ToLower(c))
		} else {
			buf.WriteRune(c)
		}
	}

	// all done
	return buf.String(), true, nil
}

func expandParamValue(key string, lookupVar LookupVar) <-chan string {
	// we'll send the results bit by bit via this channel
	chn := make(chan string)

	// are we expanding the positional parameters?
	if key == "$@" || key == "$*" {
		go func() {
			// how many positional parameters are there?
			//
			// we rely on $# being correctly set by the caller
			rawMax, ok := lookupVar("$#")
			if !ok {
				chn <- ""
			} else {
				maxI, err := strconv.Atoi(rawMax)
				if err != nil {
					chn <- ""
				} else {
					for i := 1; i <= maxI; i++ {
						retval, ok := lookupVar("$" + strconv.Itoa(i))
						if ok {
							chn <- retval
						}
					}
				}
			}
			close(chn)
		}()
	} else {
		go func() {
			retval, _ := lookupVar(key)
			chn <- retval
			close(chn)
		}()
	}

	return chn
}
