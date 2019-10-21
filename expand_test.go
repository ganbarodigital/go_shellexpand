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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type expandTestData struct {
	homedirs             map[string]string
	positionalVars       map[string]string
	specialVars          map[string]string
	vars                 map[string]string
	input                string
	shellExtra           []string
	expectedResult       string
	resultSubstringMatch bool
	actualResult         func(expandTestData) string
}

func TestExpand(t *testing.T) {

	// our master set of all UNIX shell string expansion permutations
	testDataSets := []expandTestData{
		// simple param, no braces
		{
			vars: map[string]string{
				"PARAM1": "foo",
			},
			input:          "$PARAM1",
			expectedResult: "foo",
		},
		// simple param inside longer string, no braces
		{
			vars: map[string]string{
				"PARAM1": "foo",
			},
			input:          "this is all $PARAM1 bar",
			expectedResult: "this is all foo bar",
		},
		// simple param, positional var $1
		{
			positionalVars: map[string]string{
				"$1": "foo",
			},
			input:          "$1",
			expectedResult: "foo",
		},
		// simple param, positional var $2
		{
			positionalVars: map[string]string{
				"$1": "foo",
				"$2": "bar",
			},
			input:          "$2",
			expectedResult: "bar",
		},
		// simple param, positional var $3
		{
			positionalVars: map[string]string{
				"$1": "foo",
				"$2": "bar",
				"$3": "alfred",
			},
			input:          "$3",
			expectedResult: "alfred",
		},
		// simple param, positional var $4
		{
			positionalVars: map[string]string{
				"$1": "foo",
				"$2": "bar",
				"$3": "alfred",
				"$4": "trout",
			},
			input:          "$4",
			expectedResult: "trout",
		},
		// simple param, positional var $5
		{
			positionalVars: map[string]string{
				"$1": "foo",
				"$2": "bar",
				"$3": "alfred",
				"$4": "trout",
				"$5": "haddock",
			},
			input:          "$5",
			expectedResult: "haddock",
		},
		// simple param, positional var $6
		{
			positionalVars: map[string]string{
				"$1": "foo",
				"$2": "bar",
				"$3": "alfred",
				"$4": "trout",
				"$5": "haddock",
				"$6": "cod",
			},
			input:          "$6",
			expectedResult: "cod",
		},
		// simple param, positional var $7
		{
			positionalVars: map[string]string{
				"$1": "foo",
				"$2": "bar",
				"$3": "alfred",
				"$4": "trout",
				"$5": "haddock",
				"$6": "cod",
				"$7": "plaice",
			},
			input:          "$7",
			expectedResult: "plaice",
		},
		// simple param, positional var $8
		{
			positionalVars: map[string]string{
				"$1": "foo",
				"$2": "bar",
				"$3": "alfred",
				"$4": "trout",
				"$5": "haddock",
				"$6": "cod",
				"$7": "plaice",
				"$8": "pollock",
			},
			input:          "$8",
			expectedResult: "pollock",
		},
		// simple param, positional var $9
		{
			positionalVars: map[string]string{
				"$1": "foo",
				"$2": "bar",
				"$3": "alfred",
				"$4": "trout",
				"$5": "haddock",
				"$6": "cod",
				"$7": "plaice",
				"$8": "pollock",
				"$9": "whitebait",
			},
			input:          "$9",
			expectedResult: "whitebait",
		},
		// simple param, positional var $10
		//
		// bash stops parsing after matching '$1'
		{
			positionalVars: map[string]string{
				"$1":  "foo",
				"$2":  "bar",
				"$3":  "alfred",
				"$4":  "trout",
				"$5":  "haddock",
				"$6":  "cod",
				"$7":  "plaice",
				"$8":  "pollock",
				"$9":  "whitebait",
				"$10": "bank",
			},
			input:          "$10",
			expectedResult: "foo0",
		},
		// simple param, positional var $10 wrapped in brackets
		//
		// this IS supported by bash
		{
			positionalVars: map[string]string{
				"$1":  "foo",
				"$2":  "bar",
				"$3":  "alfred",
				"$4":  "trout",
				"$5":  "haddock",
				"$6":  "cod",
				"$7":  "plaice",
				"$8":  "pollock",
				"$9":  "whitebait",
				"$10": "bank",
			},
			input:          "${10}",
			expectedResult: "bank",
		},
		// simple param, expand all positional vars via $*
		{
			positionalVars: map[string]string{
				"$1":  "foo",
				"$2":  "bar",
				"$3":  "alfred",
				"$4":  "trout",
				"$5":  "haddock",
				"$6":  "cod",
				"$7":  "plaice",
				"$8":  "pollock",
				"$9":  "whitebait",
				"$10": "bank",
				"$#":  "10",
			},
			input:          "$*",
			expectedResult: "foo bar alfred trout haddock cod plaice pollock whitebait bank",
		},
		// simple param, expand all positional vars via $@
		{
			positionalVars: map[string]string{
				"$1":  "foo",
				"$2":  "bar",
				"$3":  "alfred",
				"$4":  "trout",
				"$5":  "haddock",
				"$6":  "cod",
				"$7":  "plaice",
				"$8":  "pollock",
				"$9":  "whitebait",
				"$10": "bank",
				"$#":  "10",
			},
			input:          "$@",
			expectedResult: "foo bar alfred trout haddock cod plaice pollock whitebait bank",
		},
		// simple param, braces
		{
			vars: map[string]string{
				"PARAM1": "foo",
			},
			input:          "${PARAM1}",
			expectedResult: "foo",
		},
		// simple param inside longer string, braces
		{
			vars: map[string]string{
				"PARAM1": "foo",
			},
			input:          "this is all ${PARAM1}bar",
			expectedResult: "this is all foobar",
		},
		// simple param, braces, indirection
		{
			vars: map[string]string{
				"PARAM1": "PARAM2",
				"PARAM2": "foo",
			},
			input:          "${!PARAM1}",
			expectedResult: "foo",
		},
		// simple param, default value triggered
		{
			vars: map[string]string{
				"PARAM1": "",
			},
			input:          "${PARAM1:-foo}",
			expectedResult: "foo",
		},
		// simple param, default value NOT triggered
		{
			vars: map[string]string{
				"PARAM1": "foo",
			},
			input:          "${PARAM1:-bar}",
			expectedResult: "foo",
		},
		// simple param, default value triggered, indirection
		{
			vars: map[string]string{
				"PARAM1": "PARAM2",
				"PARAM2": "",
			},
			input:          "${!PARAM1:-foo}",
			expectedResult: "foo",
		},
		// simple param, default value NOT triggered, indirection
		{
			vars: map[string]string{
				"PARAM1": "PARAM2",
				"PARAM2": "foo",
			},
			input:          "${!PARAM1:-bar}",
			expectedResult: "foo",
		},
		// simple param, default value triggered AND expanded
		{
			vars: map[string]string{
				"PARAM1": "",
				"PARAM2": "bar",
			},
			input:          "${PARAM1:-${PARAM2}}",
			expectedResult: "bar",
		},
		// positional param, default value triggered
		{
			input:          "${1:-foo}",
			expectedResult: "foo",
		},
		// simple param, default value set
		{
			input: "${PARAM1:=foo}",
			shellExtra: []string{
				"dummy=${PARAM1:=foo}",
				"echo $PARAM1",
			},
			expectedResult: "foo",
			actualResult: func(testData expandTestData) string {
				return testData.vars["PARAM1"]
			},
		},
		// simple param, default value set to word expansion
		{
			vars: map[string]string{
				"PARAM2": "bar",
			},
			input: "${PARAM1:=${PARAM2}}",
			shellExtra: []string{
				"dummy=${PARAM1:=${PARAM2}}",
				"echo $PARAM1",
			},
			expectedResult: "bar",
			actualResult: func(testData expandTestData) string {
				return testData.vars["PARAM1"]
			},
		},
		// indirect param, default value set to word expansion
		{
			vars: map[string]string{
				"PARAM1": "PARAM2",
			},
			input: "${!PARAM1:=foo}",
			shellExtra: []string{
				"dummy=${!PARAM1:=foo}",
				"echo $PARAM2",
			},
			expectedResult: "foo",
			actualResult: func(testData expandTestData) string {
				return testData.vars["PARAM2"]
			},
		},
		// simple param, error written
		{
			vars: map[string]string{
				"foo": "",
			},
			input:                "${foo:?not set}",
			expectedResult:       "foo: not set",
			resultSubstringMatch: true,
		},
		// simple param, use alternative var
		{
			vars: map[string]string{
				"foo": "bar",
			},
			input:          "${foo:+alternative}",
			expectedResult: "alternative",
		},
		// simple param, expand substring to end of value
		{
			vars: map[string]string{
				"foo": "1234567890",
			},
			input:          "${foo:3}",
			expectedResult: "4567890",
		},
		// simple param, expand substring to end of value with indirection
		{
			vars: map[string]string{
				"foo": "bar",
				"bar": "1234567890",
			},
			input:          "${!foo:3}",
			expectedResult: "4567890",
		},
		// simple param, expand substring to given length
		{
			vars: map[string]string{
				"foo": "1234567890",
			},
			input:          "${foo:3:4}",
			expectedResult: "4567",
		},
		// simple param, expand substring to given length with indirection
		{
			vars: map[string]string{
				"foo": "bar",
				"bar": "1234567890",
			},
			input:          "${!foo:3:4}",
			expectedResult: "4567",
		},
		// expand param names by prefix
		{
			vars: map[string]string{
				"foo1": "bar",
				"foo2": "humbug",
			},
			input:          "${!foo*}",
			expectedResult: "foo1 foo2",
		},
		// expand param names by prefix with @ suffix
		{
			vars: map[string]string{
				"foo1": "bar",
				"foo2": "humbug",
			},
			input:          "${!foo@}",
			expectedResult: "foo1 foo2",
		},
		// length of simple param
		{
			vars: map[string]string{
				"foo": "bar humbug",
			},
			input:          "${#foo}",
			expectedResult: "10",
		},
		// number of positional parameters via $*
		{
			specialVars: map[string]string{
				"$#": "10",
			},
			positionalVars: map[string]string{
				"$1":  "foo",
				"$2":  "bar",
				"$3":  "alfred",
				"$4":  "trout",
				"$5":  "haddock",
				"$6":  "cod",
				"$7":  "plaice",
				"$8":  "pollock",
				"$9":  "whitebait",
				"$10": "bank",
			},
			input:          "${#*}",
			expectedResult: "10",
		},
		// number of positional parameters via $@
		{
			specialVars: map[string]string{
				"$#": "10",
			},
			positionalVars: map[string]string{
				"$1":  "foo",
				"$2":  "bar",
				"$3":  "alfred",
				"$4":  "trout",
				"$5":  "haddock",
				"$6":  "cod",
				"$7":  "plaice",
				"$8":  "pollock",
				"$9":  "whitebait",
				"$10": "bank",
			},
			input:          "${#@}",
			expectedResult: "10",
		},
	}

	for _, testData := range testDataSets {

		// ----------------------------------------------------------------
		// setup your test

		var buf strings.Builder

		buf.WriteString("#!/usr/bin/env bash\n\n")
		for key, value := range testData.vars {
			buf.WriteString(fmt.Sprintf("%s='%s'\n", key, value))
		}
		if len(testData.positionalVars) > 0 {
			buf.WriteString("set -- ")
			for i := 1; i <= len(testData.positionalVars); i++ {
				buf.WriteString(testData.positionalVars["$"+strconv.Itoa(i)] + " ")
			}
			buf.WriteString("\n")
		}

		// do we need to write any extra steps to get the shell to tell us
		// what the outcome was?
		if len(testData.shellExtra) > 0 {
			for _, line := range testData.shellExtra {
				buf.WriteString(line)
				buf.WriteRune('\n')
			}
		} else {
			// no, we can simply echo the string we are expanding
			buf.WriteString("echo ")
			buf.WriteString(testData.input)
			buf.WriteString("\n")
		}

		// create the shell script we'll use to verify that internal behaviour
		// matches actual shell script behaviour
		tmpFile, _ := ioutil.TempFile("", "shellexpand-expandParams-")
		cleanup := func() {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}
		defer cleanup()

		// tmpFile.Truncate(0)
		tmpFile.WriteString(buf.String())
		tmpFile.Sync()
		tmpFile.Close()

		// now, setup everything we need to test this internally
		varFuncs := VarFuncs{
			AssignToVar: func(key string, value string) error {
				if len(testData.vars) == 0 {
					testData.vars = make(map[string]string)
				}
				testData.vars[key] = value

				return nil
			},

			LookupVar: func(key string) (string, bool) {
				// special case - special parameter
				retval, ok := testData.specialVars[key]
				if ok {
					return retval, true
				}
				// special case - positional parameter
				retval, ok = testData.positionalVars[key]
				if ok {
					return retval, true
				}
				// general case
				retval, ok = testData.vars[key]
				if ok {
					return retval, true
				}
				return "", false
			},

			LookupHomeDir: func(key string) (string, bool) {
				retval, ok := testData.homedirs[key]
				if ok {
					return retval, true
				}
				return "", false
			},

			MatchVarNames: func(prefix string) []string {
				retval := []string{}

				for key := range testData.vars {
					if strings.HasPrefix(key, prefix) {
						retval = append(retval, key)
					}
				}

				return retval
			},
		}

		// shorthand
		input := testData.input
		expectedResult := testData.expectedResult

		// ----------------------------------------------------------------
		// perform the change

		cmd := exec.Command("/usr/bin/env", "bash", tmpFile.Name())
		shellRawResult, _ := cmd.CombinedOutput()
		shellActualResult := strings.TrimSpace(string(shellRawResult))

		internalActualResult := Expand(input, varFuncs)
		// special case - the result is a side effect, not a direct string
		// expansion
		if testData.actualResult != nil {
			internalActualResult = testData.actualResult(testData)
		}

		// ----------------------------------------------------------------
		// test the results

		// assert.Nil(t, shellErr)
		if testData.resultSubstringMatch {
			assert.Contains(t, shellActualResult, expectedResult, buf.String())
			assert.Contains(t, internalActualResult, expectedResult, testData)
		} else {
			assert.Equal(t, expectedResult, shellActualResult, buf.String())
			assert.Equal(t, expectedResult, internalActualResult, testData)
		}
	}
}

// func TestExpandParamsDebugCase(t *testing.T) {

// 	// if you add a test here, you must also add it to the main
// 	// Expand test suite
// 	testDataSets := []expandTestData{
// 		// simple param, error written
// 		{
// 			vars: map[string]string{
// 				"foo": "",
// 			},
// 			input:          "${foo:?not set}",
// 			expectedResult: "foo",
// 		},
// 	}

// 	for _, testData := range testDataSets {

// 		// ----------------------------------------------------------------
// 		// setup your test

// 		var buf strings.Builder

// 		buf.WriteString("#!/usr/bin/env bash\n\n")
// 		for key, value := range testData.vars {
// 			buf.WriteString(fmt.Sprintf("%s='%s'\n", key, value))
// 		}
// 		if len(testData.positionalVars) > 0 {
// 			buf.WriteString("set -- ")
// 			for i := 1; i <= len(testData.positionalVars); i++ {
// 				buf.WriteString(testData.positionalVars["$"+strconv.Itoa(i)] + " ")
// 			}
// 			buf.WriteString("\n")
// 		}

// 		// do we need to write any extra steps to get the shell to tell us
// 		// what the outcome was?
// 		if len(testData.shellExtra) > 0 {
// 			for _, line := range testData.shellExtra {
// 				buf.WriteString(line)
// 				buf.WriteRune('\n')
// 			}
// 		} else {
// 			// no, we can simply echo the string we are expanding
// 			buf.WriteString("echo ")
// 			buf.WriteString(testData.input)
// 			buf.WriteString("\n")
// 		}

// 		// create the shell script we'll use to verify that internal behaviour
// 		// matches actual shell script behaviour
// 		tmpFile, _ := ioutil.TempFile("", "shellexpand-expandParams-")
// 		cleanup := func() {
// 			tmpFile.Close()
// 			os.Remove(tmpFile.Name())
// 		}
// 		defer cleanup()

// 		// tmpFile.Truncate(0)
// 		tmpFile.WriteString(buf.String())
// 		tmpFile.Sync()
// 		tmpFile.Close()

// 		// now, setup everything we need to test this internally
// 		assignVar := func(key string, value string) error {
// 			if len(testData.vars) == 0 {
// 				testData.vars = make(map[string]string)
// 			}
// 			testData.vars[key] = value

// 			return nil
// 		}

// 		varLookup := func(key string) (string, bool) {
// 			// special case - positional parameter
// 			retval, ok := testData.positionalVars[key]
// 			if ok {
// 				return retval, true
// 			}
// 			// general case
// 			retval, ok = testData.vars[key]
// 			if ok {
// 				return retval, true
// 			}
// 			return "", false
// 		}

// 		homeDirLookup := func(key string) (string, bool) {
// 			retval, ok := testData.homedirs[key]
// 			if ok {
// 				return retval, true
// 			}
// 			return "", false
// 		}

// 		// shorthand
// 		input := testData.input
// 		expectedResult := testData.expectedResult

// 		// ----------------------------------------------------------------
// 		// perform the change

// 		cmd := exec.Command("/usr/bin/env", "bash", tmpFile.Name())
// 		shellRawResult, _ := cmd.CombinedOutput()
// 		shellActualResult := strings.TrimSpace(string(shellRawResult))

// 		internalActualResult := Expand(input, varLookup, homeDirLookup, assignVar)
// 		// special case - the result is a side effect, not a direct string
// 		// expansion
// 		if testData.actualResult != nil {
// 			internalActualResult = testData.actualResult(testData)
// 		}

// 		// ----------------------------------------------------------------
// 		// test the results

// 		// assert.Nil(t, shellErr)
// 		if testData.resultSubstringMatch {
// 			assert.Contains(t, shellActualResult, expectedResult, buf.String())
// 			assert.Contains(t, internalActualResult, expectedResult, testData)
// 		} else {
// 			assert.Equal(t, expectedResult, shellActualResult, buf.String())
// 			assert.Equal(t, expectedResult, internalActualResult, testData)
// 		}
// 	}
// }
