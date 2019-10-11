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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParamNoBraces(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "$VAR"
	expectedResult := paramDesc{
		kind:  paramExpandToValue,
		parts: []string{"VAR"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSimpleBraces(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR}"
	expectedResult := paramDesc{
		kind:  paramExpandToValue,
		parts: []string{"VAR"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamShellSpecialNoBraces(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"$!",
		"$$",
		"$*",
		"$@",
		"$#",
		"$?",
		"$-",
		"$0",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandToValue,
			parts: []string{testData},
		}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.True(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamShellSpecialWithBraces(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!}",
		"${$}",
		"${*}",
		"${@}",
		"${#}",
		"${?}",
		"${-}",
		"${0}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandToValue,
			parts: []string{"$" + testData[2:3]},
		}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.True(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamPositionalParamsNoBraces(t *testing.T) {
	t.Parallel()

	var testDataSet []string

	for i := 1; i < 20; i++ {
		testDataSet = append(testDataSet, "$"+strconv.Itoa(i))
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandToValue,
			parts: []string{testData},
		}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.True(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamPositionalWithBraces(t *testing.T) {
	t.Parallel()

	var testDataSet []string

	for i := 1; i < 20; i++ {
		testDataSet = append(testDataSet, "${"+strconv.Itoa(i)+"}")
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandToValue,
			parts: []string{"$" + testData[2:len(testData)-1]},
		}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.True(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamDefaultValue(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR:-FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandWithDefaultValue,
		parts: []string{"VAR", "FOO"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamDefaultValueWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR:-FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandWithDefaultValue,
		parts:    []string{"VAR", "FOO"},
		indirect: true,
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamDefaultValueSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V:-FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandWithDefaultValue,
		parts: []string{"V", "FOO"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamDefaultValueSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V:-FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandWithDefaultValue,
		parts:    []string{"V", "FOO"},
		indirect: true,
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamPositionalParamDefaultValue(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + ":-FOO}"
		expectedResult := paramDesc{
			kind:  paramExpandWithDefaultValue,
			parts: []string{"$" + testValue, "FOO"},
		}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.True(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamPositionalParamDefaultValueWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + ":-FOO}"
		expectedResult := paramDesc{
			kind:     paramExpandWithDefaultValue,
			parts:    []string{"$" + testValue, "FOO"},
			indirect: true,
		}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.True(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamShellSpecialDefaultValue(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!:-foo}",
		"${$:-foo}",
		"${*:-foo}",
		"${@:-foo}",
		"${#:-foo}",
		"${?:-foo}",
		"${-:-foo}",
		"${0:-foo}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandWithDefaultValue,
			parts: []string{"$" + testData[2:3], "foo"},
		}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.True(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamShellSpecialDefaultValueWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$:-foo}",
		"${!*:-foo}",
		"${!@:-foo}",
		"${!#:-foo}",
		"${!?:-foo}",
		"${!-:-foo}",
		"${!0:-foo}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandWithDefaultValue,
			parts:    []string{"$" + testData[3:4], "foo"},
			indirect: true,
		}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.True(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamPlingDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!:-foo}"
	expectedResult := paramDesc{
		kind: paramExpandNotSupported,
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}
