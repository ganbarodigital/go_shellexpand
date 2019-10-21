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
func expandParameters(input string, lookupVar LookupVar, lookupHomeDir LookupVar, assignVar AssignVar) string {
	// we expand in a strictly left-to-right manner
	for i := 0; i < len(input); i++ {
		if input[i] == '\\' {
			// skip over escaped characters
			i++
		} else if input[i] == '$' {
			varEnd, ok := matchVar(input, i)
			if ok {
				paramDesc, ok := parseParameter(input[i : varEnd+1])
				if !ok {
					continue
				}

				replacement := expandParameter(paramDesc, lookupVar, lookupHomeDir, assignVar)
				var buf strings.Builder

				if i > 0 {
					buf.WriteString(input[0:i])
				}
				buf.WriteString(replacement)

				if i < len(input) {
					buf.WriteString(input[varEnd+1:])
				}

				input = buf.String()
			}
		}
	}

	return input
}

func expandParameter(paramDesc paramDesc, lookupVar LookupVar, lookupHomeDir LookupVar, assignVar AssignVar) string {
	// what we will (eventually) send back
	var retval []string

	// ... but only if all is well
	var ok bool

	// and, because we may be building it up bit by bit, we need somewhere
	// to store it temporarily
	var buf string

	// step 1: we need to expand the paramName first, to support any
	// possible use of indirection
	paramName, ok := expandParamName(paramDesc, lookupVar)

	// step 2: we need to feed that into all the different ways that
	// parameters can be expanded in strings
	//
	// this is complicated by some parameters ($*, $@, and arrays if we
	// ever add support for them in the future) having the expansion applied
	// to each part of their value
	for paramValue := range expandParamValue(paramName, lookupVar) {
		switch paramDesc.kind {
		case paramExpandToValue:
			buf, ok = expandParamToValue(paramName, paramValue, paramDesc, lookupVar)
		case paramExpandWithDefaultValue:
			buf, ok = expandParamWithDefaultValue(paramName, paramValue, paramDesc, lookupVar, lookupHomeDir, assignVar)
		case paramExpandSetDefaultValue:
			buf, ok = expandParamSetDefaultValue(paramName, paramValue, paramDesc, lookupVar, lookupHomeDir, assignVar)
		case paramExpandWriteError:
			buf, ok = expandParamWriteError(paramName, paramValue, paramDesc, lookupVar, lookupHomeDir, assignVar)
		}

		retval = append(retval, buf)
	}

	// are we happy with our attempted expansion?
	if !ok {
		return ""
	}

	// if we get here, then yes, we are happy
	return strings.Join(retval, " ")
}

func expandParamName(paramDesc paramDesc, lookupVar LookupVar) (string, bool) {
	varName := paramDesc.parts[0]
	ok := true
	if paramDesc.indirect {
		varName, ok = lookupVar(varName)
	}

	return varName, ok
}

func expandParamWithIndirection(paramName string, lookupVar LookupVar) string {
	retval, ok := lookupVar(paramName)
	if !ok {
		return ""
	}

	return retval
}

func expandParamToValue(paramName, paramValue string, paramDesc paramDesc, lookupVar LookupVar) (string, bool) {
	// nothing else to do
	return paramValue, true
}

func expandParamWithDefaultValue(paramName, paramValue string, paramDesc paramDesc, lookupVar LookupVar, lookupHomeDir LookupVar, assignVar AssignVar) (string, bool) {
	// do we need to return the default value?
	if paramValue == "" {
		return expandWord(paramDesc.parts[1], lookupVar, lookupHomeDir, assignVar), true
	}

	return paramValue, true
}

func expandParamSetDefaultValue(paramName, paramValue string, paramDesc paramDesc, lookupVar LookupVar, lookupHomeDir LookupVar, assignVar AssignVar) (string, bool) {
	// do we need to do anything?
	if paramValue != "" {
		return paramValue, true
	}

	// at this point, we need to assign a new value
	err := assignVar(paramName, expandWord(paramDesc.parts[1], lookupVar, lookupHomeDir, assignVar))
	if err != nil {
		return "", false
	}

	// all done
	return lookupVar(paramName)
}

func expandParamWriteError(paramName, paramValue string, paramDesc paramDesc, lookupVar LookupVar, lookupHomeDir LookupVar, assignVar AssignVar) (string, bool) {
	// do we have a value?
	if paramValue == "" {
		return paramName + ": " + expandWord(paramDesc.parts[1], lookupVar, lookupHomeDir, assignVar), true
	}

	return paramValue, true
}

func expandParamValue(key string, lookupVar LookupVar) <-chan string {
	// we'll send the results bit by bit via this channel
	chn := make(chan string)

	// are we expanding the positional parameters?
	if key == "$@" || key == "$*" {
		// how many positional parameters are there?
		//
		// we rely on $# being correctly set by the caller
		rawMax, ok := lookupVar("$#")
		if !ok {
			chn <- ""
			close(chn)
		} else {
			maxI, err := strconv.Atoi(rawMax)
			if err != nil {
				chn <- ""
				close(chn)
			} else {
				go func() {
					for i := 1; i <= maxI; i++ {
						retval, ok := lookupVar("$" + strconv.Itoa(i))
						if ok {
							chn <- retval
						}
					}
					close(chn)
				}()
			}
		}
	} else {
		go func() {
			retval, _ := lookupVar(key)
			chn <- retval
			close(chn)
		}()
	}

	return chn
}
