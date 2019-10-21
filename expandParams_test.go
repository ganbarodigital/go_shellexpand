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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type expandParamsTestData struct {
	homedirs       map[string]string
	vars           map[string]string
	input          string
	expectedResult string
}

func TestExpandParams(t *testing.T) {

	// if you add a test here, you must also add it to the main
	// Expand test suite
	testDataSets := []expandParamsTestData{
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
	}

	for _, testData := range testDataSets {

		// ----------------------------------------------------------------
		// setup your test

		// create the shell script we'll use to verify that internal behaviour
		// matches actual shell script behaviour
		tmpFile, _ := ioutil.TempFile("", "shellexpand-expandParams-")
		cleanup := func() {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}
		defer cleanup()
		tmpFile.WriteString("#!/usr/bin/env bash\n\n")
		for key, value := range testData.vars {
			tmpFile.WriteString(fmt.Sprintf("%s='%s'\n", key, value))
		}
		tmpFile.WriteString("echo ")
		tmpFile.WriteString(testData.input)
		tmpFile.WriteString("\n")
		tmpFile.Close()

		// now, setup everything we need to test this internally
		varLookup := func(key string) (string, bool) {
			retval, ok := testData.vars[key]
			if ok {
				return retval, true
			}
			return "", false
		}
		homeDirLookup := func(key string) (string, bool) {
			retval, ok := testData.homedirs[key]
			if ok {
				return retval, true
			}
			return "", false
		}

		// shorthand
		input := testData.input
		expectedResult := testData.expectedResult

		// ----------------------------------------------------------------
		// perform the change

		cmd := exec.Command("/usr/bin/env", "bash", tmpFile.Name())
		shellRawResult, shellErr := cmd.CombinedOutput()
		shellActualResult := strings.TrimSpace(string(shellRawResult))

		internalActualResult := expandParameters(input, varLookup, homeDirLookup)

		// ----------------------------------------------------------------
		// test the results

		assert.Nil(t, shellErr)
		assert.Equal(t, expectedResult, shellActualResult)
		assert.Equal(t, expectedResult, internalActualResult)
	}
}
