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

func TestParseParamSetDefaultValue(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR:=FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandSetDefaultValue,
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

func TestParseParamSetDefaultValueWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR:=FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandSetDefaultValue,
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

func TestParseParamSetDefaultValueSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V:=FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandSetDefaultValue,
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

func TestParseParamSetDefaultValueSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V:=FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandSetDefaultValue,
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

func TestParseParamPositionalParamSetDefaultValue(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + ":=FOO}"
		expectedResult := paramDesc{
			kind:  paramExpandSetDefaultValue,
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

func TestParseParamPositionalParamSetDefaultValueWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + ":=FOO}"
		expectedResult := paramDesc{
			kind:     paramExpandSetDefaultValue,
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

func TestParseParamShellSpecialSetDefaultValue(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$:=foo}",
		"${*:=foo}",
		"${@:=foo}",
		"${#:=foo}",
		"${?:=foo}",
		"${-:=foo}",
		"${0:=foo}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandSetDefaultValue,
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

func TestParseParamShellSpecialSetDefaultValueWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$:=foo}",
		"${!*:=foo}",
		"${!@:=foo}",
		"${!#:=foo}",
		"${!?:=foo}",
		"${!-:=foo}",
		"${!0:=foo}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandSetDefaultValue,
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

func TestParseParamPlingSetDefaultValueDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!:=foo}"
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

func TestParseParamWriteError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR:?FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandWriteError,
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

func TestParseParamWriteErrorWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR:?FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandWriteError,
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

func TestParseParamWriteErrorSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V:?FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandWriteError,
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

func TestParseParamWriteErrorSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V:?FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandWriteError,
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

func TestParseParamPositionalParamWriteError(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + ":?FOO}"
		expectedResult := paramDesc{
			kind:  paramExpandWriteError,
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

func TestParseParamPositionalParamWriteErrorWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + ":?FOO}"
		expectedResult := paramDesc{
			kind:     paramExpandWriteError,
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

func TestParseParamShellSpecialWriteError(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$:?foo}",
		"${*:?foo}",
		"${@:?foo}",
		"${#:?foo}",
		"${?:?foo}",
		"${-:?foo}",
		"${0:?foo}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandWriteError,
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

func TestParseParamShellSpecialWriteErrorWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$:?foo}",
		"${!*:?foo}",
		"${!@:?foo}",
		"${!#:?foo}",
		"${!?:?foo}",
		"${!-:?foo}",
		"${!0:?foo}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandWriteError,
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

func TestParseParamPlingWriteErrorDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!:?foo}"
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

func TestParseParamAlternativeValue(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR:+FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandAlternativeValue,
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

func TestParseParamAlternativeValueWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR:+FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandAlternativeValue,
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

func TestParseParamAlternativeValueSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V:+FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandAlternativeValue,
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

func TestParseParamAlternativeValueSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V:+FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandAlternativeValue,
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

func TestParseParamPositionalParamAlternativeValue(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + ":+FOO}"
		expectedResult := paramDesc{
			kind:  paramExpandAlternativeValue,
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

func TestParseParamPositionalParamAlternativeValueWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + ":+FOO}"
		expectedResult := paramDesc{
			kind:     paramExpandAlternativeValue,
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

func TestParseParamShellSpecialAlternativeValue(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$:+foo}",
		"${*:+foo}",
		"${@:+foo}",
		"${#:+foo}",
		"${?:+foo}",
		"${-:+foo}",
		"${0:+foo}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandAlternativeValue,
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

func TestParseParamShellSpecialAlternativeValueWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$:+foo}",
		"${!*:+foo}",
		"${!@:+foo}",
		"${!#:+foo}",
		"${!?:+foo}",
		"${!-:+foo}",
		"${!0:+foo}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandAlternativeValue,
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

func TestParseParamPlingAlternativeValueDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!:+foo}"
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

func TestParseParamSubstring(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR:500}"
	expectedResult := paramDesc{
		kind:  paramExpandSubstring,
		parts: []string{"VAR", "500"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSubstringMustHaveOffset(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR:}"
	expectedResult := paramDesc{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSubstringOffsetMustBeNumeric(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"abcdef",
		"1hundred",
		"500 ",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		testData := "${VAR:" + testData + "}"
		expectedResult := paramDesc{}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.False(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamSubstringWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR:500}"
	expectedResult := paramDesc{
		kind:     paramExpandSubstring,
		parts:    []string{"VAR", "500"},
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

func TestParseParamSubstringSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V:500}"
	expectedResult := paramDesc{
		kind:  paramExpandSubstring,
		parts: []string{"V", "500"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSubstringSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V:500}"
	expectedResult := paramDesc{
		kind:     paramExpandSubstring,
		parts:    []string{"V", "500"},
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

func TestParseParamPositionalParamSubstring(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + ":500}"
		expectedResult := paramDesc{
			kind:  paramExpandSubstring,
			parts: []string{"$" + testValue, "500"},
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

func TestParseParamPositionalParamSubstringWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + ":500}"
		expectedResult := paramDesc{
			kind:     paramExpandSubstring,
			parts:    []string{"$" + testValue, "500"},
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

func TestParseParamShellSpecialSubstring(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$:500}",
		"${*:500}",
		"${#:500}",
		"${?:500}",
		"${-:500}",
		"${0:500}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandSubstring,
			parts: []string{"$" + testData[2:3], "500"},
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

func TestParseParamShellSpecialSubstringWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$:500}",
		"${!*:500}",
		"${!#:500}",
		"${!?:500}",
		"${!-:500}",
		"${!0:500}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandSubstring,
			parts:    []string{"$" + testData[3:4], "500"},
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

func TestParseParamSubstringListOfPositionalParameters(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${@:1}"
	expectedResult := paramDesc{
		kind:  paramExpandPositionalParamsFromOffset,
		parts: []string{"$@", "1"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamPlingSubstringDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!:500}"
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

func TestParseParamSubstringLength(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR:500:1000}"
	expectedResult := paramDesc{
		kind:  paramExpandSubstringLength,
		parts: []string{"VAR", "500", "1000"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

// resurrect this test when we have arithmetic expansion supported!
//
// func TestParseParamSubstringLengthSupportsNegativeOffset(t *testing.T) {
// 	t.Parallel()

// 	// ----------------------------------------------------------------
// 	// setup your test

// 	testData := "${VAR:((-3)):1}"
// 	expectedResult := paramDesc{
// 		kind:  paramExpandSubstringLength,
// 		parts: []string{"VAR", "((-3))", "1"},
// 	}

// 	// ----------------------------------------------------------------
// 	// perform the change

// 	actualResult, ok := parseParameter(testData)

// 	// ----------------------------------------------------------------
// 	// test the results

// 	assert.True(t, ok)
// 	assert.Equal(t, expectedResult, actualResult)
// }

func TestParseParamSubstringLengthSupportsNegativeLength(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR:500:-1}"
	expectedResult := paramDesc{
		kind:  paramExpandSubstringLength,
		parts: []string{"VAR", "500", "-1"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSubstringLengthMustBeNumeric(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"100:abcdef",
		"100:1hundred",
		"100:500 ",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		testData := "${VAR:" + testData + "}"
		expectedResult := paramDesc{}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.False(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamSubstringLengthCannotHaveExtraParts(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"100:500:600",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		testData := "${VAR:" + testData + "}"
		expectedResult := paramDesc{}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.False(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamSubstringLengthWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR:500:1000}"
	expectedResult := paramDesc{
		kind:     paramExpandSubstringLength,
		parts:    []string{"VAR", "500", "1000"},
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

func TestParseParamSubstringLengthSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V:500:1000}"
	expectedResult := paramDesc{
		kind:  paramExpandSubstringLength,
		parts: []string{"V", "500", "1000"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSubstringLengthSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V:500:1000}"
	expectedResult := paramDesc{
		kind:     paramExpandSubstringLength,
		parts:    []string{"V", "500", "1000"},
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

func TestParseParamPositionalParamSubstringLength(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + ":500:1000}"
		expectedResult := paramDesc{
			kind:  paramExpandSubstringLength,
			parts: []string{"$" + testValue, "500", "1000"},
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

func TestParseParamPositionalParamSubstringLengthWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + ":500:1000}"
		expectedResult := paramDesc{
			kind:     paramExpandSubstringLength,
			parts:    []string{"$" + testValue, "500", "1000"},
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

func TestParseParamShellSpecialSubstringLength(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$:500:1000}",
		"${*:500:1000}",
		"${#:500:1000}",
		"${?:500:1000}",
		"${-:500:1000}",
		"${0:500:1000}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandSubstringLength,
			parts: []string{"$" + testData[2:3], "500", "1000"},
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

func TestParseParamSubstringLengthListOfPositionalParameters(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${@:1:3}"
	expectedResult := paramDesc{
		kind:  paramExpandPositionalParamsFromOffsetLength,
		parts: []string{"$@", "1", "3"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamShellSpecialSubstringLengthWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$:500:1000}",
		"${!*:500:1000}",
		"${!#:500:1000}",
		"${!?:500:1000}",
		"${!-:500:1000}",
		"${!0:500:1000}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandSubstringLength,
			parts:    []string{"$" + testData[3:4], "500", "1000"},
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

func TestParseParamPlingSubstringLengthDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!:500:1000}"
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

func TestParseParamNameMatchPrefix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR*}"
	expectedResult := paramDesc{
		kind:  paramExpandPrefixNames,
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

func TestParseParamNameMatchPrefixDoubleQuoted(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR@}"
	expectedResult := paramDesc{
		kind:  paramExpandPrefixNamesDoubleQuoted,
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

func TestParseParamParamLength(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${#VAR}"
	expectedResult := paramDesc{
		kind:  paramExpandParamLength,
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

func TestParseParamLengthMustHaveValidParamName(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		" ",
		"£",
		"0129",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		testData := "${#" + testData + "}"
		expectedResult := paramDesc{}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.False(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamLengthCanHaveNothingAfterParamName(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		" ",
		":0",
		":100:500",
		":-WORD",
		":=WORD",
		":?WORD",
		":+WORD",
		"*",
		"@",
		"#WORD",
		"##WORD",
		"%WORD",
		"%%WORD",
		"/old/new",
		"^pattern",
		"^^pattern",
		",pattern",
		",,pattern",
		"@a",
		"@A",
		"@E",
		"@P",
		"@Q",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		testData := "${#VAR" + testData + "}"
		expectedResult := paramDesc{}

		// ----------------------------------------------------------------
		// perform the change

		actualResult, ok := parseParameter(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.False(t, ok)
		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestParseParamParamLengthSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${#V}"
	expectedResult := paramDesc{
		kind:  paramExpandParamLength,
		parts: []string{"V"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamPositionalParamParamLength(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${#" + testValue + "}"
		expectedResult := paramDesc{
			kind:  paramExpandParamLength,
			parts: []string{"$" + testValue},
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

func TestParseParamShellSpecialParamLength(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${#!}",
		"${#$}",
		"${##}",
		"${#?}",
		"${#-}",
		"${#0}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandParamLength,
			parts: []string{"$" + testData[3:4]},
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

func TestParseParamNoOfPositionalParameters(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${#*}",
		"${#@}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandNoOfPositionalParams,
			parts: []string{"$" + testData[3:4]},
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

func TestParseParamRemoveShortestPrefix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR#FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandRemovePrefixShortestMatch,
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

func TestParseParamRemoveShortestPrefixWithNothingToRemove(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR#}"
	expectedResult := paramDesc{
		kind:  paramExpandRemovePrefixShortestMatch,
		parts: []string{"VAR", ""},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamRemoveShortestPrefixWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR#FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandRemovePrefixShortestMatch,
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

func TestParseParamRemoveShortestPrefixSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V#FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandRemovePrefixShortestMatch,
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

func TestParseParamRemoveShortestPrefixSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V#FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandRemovePrefixShortestMatch,
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

func TestParseParamPositionalParamRemoveShortestPrefix(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "#FOO}"
		expectedResult := paramDesc{
			kind:  paramExpandRemovePrefixShortestMatch,
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

func TestParseParamPositionalParamRemoveShortestPrefixWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "#FOO}"
		expectedResult := paramDesc{
			kind:     paramExpandRemovePrefixShortestMatch,
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

func TestParseParamShellSpecialRemoveShortestPrefix(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$#FOO}",
		"${##FOO}",
		"${?#FOO}",
		"${-#FOO}",
		"${0#FOO}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandRemovePrefixShortestMatch,
			parts: []string{"$" + testData[2:3], "FOO"},
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

func TestParseParamShellSpecialShortestPrefixWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$#FOO}",
		"${!?#FOO}",
		"${!-#FOO}",
		"${!0#FOO}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandRemovePrefixShortestMatch,
			parts:    []string{"$" + testData[3:4], "FOO"},
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

func TestParseParamRemoveShortestPrefixAllPositionalParams(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testDataSet := []string{
		"${*#FOO}",
		"${@#FOO}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandAllPositionalParamsRemovePrefixShortestMatch,
			parts:    []string{"$" + testData[2:3], "FOO"},
			indirect: false,
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

func TestParseParamPlingRemoveShortestPrefixDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!#FOO}"
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

func TestParseParamRemoveLongestPrefix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR##FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandRemovePrefixLongestMatch,
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

func TestParseParamRemoveLongestPrefixWithNothingToRemove(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR##}"
	expectedResult := paramDesc{
		kind:  paramExpandRemovePrefixLongestMatch,
		parts: []string{"VAR", ""},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamRemoveLongestPrefixWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR##FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandRemovePrefixLongestMatch,
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

func TestParseParamRemoveLongestPrefixSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V##FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandRemovePrefixLongestMatch,
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

func TestParseParamRemoveLongestPrefixSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V##FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandRemovePrefixLongestMatch,
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

func TestParseParamPositionalParamRemoveLongestPrefix(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "##FOO}"
		expectedResult := paramDesc{
			kind:  paramExpandRemovePrefixLongestMatch,
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

func TestParseParamPositionalParamRemoveLongestPrefixWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "##FOO}"
		expectedResult := paramDesc{
			kind:     paramExpandRemovePrefixLongestMatch,
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

func TestParseParamShellSpecialRemoveLongestPrefix(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$##FOO}",
		"${*##FOO}",
		"${@##FOO}",
		"${###FOO}",
		"${?##FOO}",
		"${-##FOO}",
		"${0##FOO}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandRemovePrefixLongestMatch,
			parts: []string{"$" + testData[2:3], "FOO"},
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

func TestParseParamShellSpecialLongestPrefixWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$##FOO}",
		"${!*##FOO}",
		"${!@##FOO}",
		"${!###FOO}",
		"${!?##FOO}",
		"${!-##FOO}",
		"${!0##FOO}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandRemovePrefixLongestMatch,
			parts:    []string{"$" + testData[3:4], "FOO"},
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

func TestParseParamPlingRemoveLongestPrefixDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!##FOO}"
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

func TestParseParamRemoveShortestSuffix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR%FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandRemoveSuffixShortestMatch,
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

func TestParseParamRemoveShortestSuffixWithNothingToRemove(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR%}"
	expectedResult := paramDesc{
		kind:  paramExpandRemoveSuffixShortestMatch,
		parts: []string{"VAR", ""},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamRemoveShortestSuffixWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR%FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandRemoveSuffixShortestMatch,
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

func TestParseParamRemoveShortestSuffixSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V%FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandRemoveSuffixShortestMatch,
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

func TestParseParamRemoveShortestSuffixSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V%FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandRemoveSuffixShortestMatch,
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

func TestParseParamPositionalParamRemoveShortestSuffix(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "%FOO}"
		expectedResult := paramDesc{
			kind:  paramExpandRemoveSuffixShortestMatch,
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

func TestParseParamPositionalParamRemoveShortestSuffixWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "%FOO}"
		expectedResult := paramDesc{
			kind:     paramExpandRemoveSuffixShortestMatch,
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

func TestParseParamShellSpecialRemoveShortestSuffix(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$%FOO}",
		"${#%FOO}",
		"${?%FOO}",
		"${-%FOO}",
		"${0%FOO}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandRemoveSuffixShortestMatch,
			parts: []string{"$" + testData[2:3], "FOO"},
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

func TestParseParamShellSpecialShortestSuffixWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$%FOO}",
		"${!?%FOO}",
		"${!-%FOO}",
		"${!0%FOO}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandRemoveSuffixShortestMatch,
			parts:    []string{"$" + testData[3:4], "FOO"},
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

func TestParseParamRemoveShortestSuffixAllPositionalParams(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testDataSet := []string{
		"${*%FOO}",
		"${@%FOO}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandAllPositionalParamsRemoveSuffixShortestMatch,
			parts:    []string{"$" + testData[2:3], "FOO"},
			indirect: false,
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

func TestParseParamPlingRemoveShortestSuffixDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!%FOO}"
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

func TestParseParamRemoveLongestSuffix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR%%FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandRemoveSuffixLongestMatch,
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

func TestParseParamRemoveLongestSuffixWithNothingToRemove(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR%%}"
	expectedResult := paramDesc{
		kind:  paramExpandRemoveSuffixLongestMatch,
		parts: []string{"VAR", ""},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamRemoveLongestSuffixWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR%%FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandRemoveSuffixLongestMatch,
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

func TestParseParamRemoveLongestSuffixSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V%%FOO}"
	expectedResult := paramDesc{
		kind:  paramExpandRemoveSuffixLongestMatch,
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

func TestParseParamRemoveLongestSuffixSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V%%FOO}"
	expectedResult := paramDesc{
		kind:     paramExpandRemoveSuffixLongestMatch,
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

func TestParseParamPositionalParamRemoveLongestSuffix(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "%%FOO}"
		expectedResult := paramDesc{
			kind:  paramExpandRemoveSuffixLongestMatch,
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

func TestParseParamPositionalParamRemoveLongestSuffixWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "%%FOO}"
		expectedResult := paramDesc{
			kind:     paramExpandRemoveSuffixLongestMatch,
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

func TestParseParamShellSpecialRemoveLongestSuffix(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$%%FOO}",
		"${#%%FOO}",
		"${?%%FOO}",
		"${-%%FOO}",
		"${0%%FOO}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandRemoveSuffixLongestMatch,
			parts: []string{"$" + testData[2:3], "FOO"},
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

func TestParseParamShellSpecialLongestSuffixWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$%%FOO}",
		"${!#%%FOO}",
		"${!?%%FOO}",
		"${!-%%FOO}",
		"${!0%%FOO}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandRemoveSuffixLongestMatch,
			parts:    []string{"$" + testData[3:4], "FOO"},
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

func TestParseParamRemoveLongestSuffixAllPositionalParams(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testDataSet := []string{
		"${*%%FOO}",
		"${@%%FOO}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandAllPositionalParamsRemoveSuffixLongestMatch,
			parts:    []string{"$" + testData[2:3], "FOO"},
			indirect: false,
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

func TestParseParamPlingRemoveLongestSuffixDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!%%FOO}"
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

func TestParseParamSearchReplaceLongestFirstMatch(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR/FOO/BAR}"
	expectedResult := paramDesc{
		kind:  paramExpandSearchReplaceLongestFirstMatch,
		parts: []string{"VAR", "FOO", "BAR"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSearchReplaceLongestFirstMatchWithNoReplacement(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR/FOO/}"
	expectedResult := paramDesc{
		kind:  paramExpandSearchReplaceLongestFirstMatch,
		parts: []string{"VAR", "FOO", ""},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSearchReplaceLongestFirstMatchWithNoSearchOrReplacement(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR/}"
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

func TestParseParamSearchReplaceLongestFirstMatchWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR/FOO/BAR}"
	expectedResult := paramDesc{
		kind:     paramExpandSearchReplaceLongestFirstMatch,
		parts:    []string{"VAR", "FOO", "BAR"},
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

func TestParseParamSearchReplaceLongestFirstMatchSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V/FOO/BAR}"
	expectedResult := paramDesc{
		kind:  paramExpandSearchReplaceLongestFirstMatch,
		parts: []string{"V", "FOO", "BAR"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSearchReplaceLongestFirstMatchSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V/FOO/BAR}"
	expectedResult := paramDesc{
		kind:     paramExpandSearchReplaceLongestFirstMatch,
		parts:    []string{"V", "FOO", "BAR"},
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

func TestParseParamPositionalParamSearchReplaceLongestFirstMatch(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "/FOO/BAR}"
		expectedResult := paramDesc{
			kind:  paramExpandSearchReplaceLongestFirstMatch,
			parts: []string{"$" + testValue, "FOO", "BAR"},
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

func TestParseParamPositionalParamSearchReplaceLongestFirstMatchWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "/FOO/BAR}"
		expectedResult := paramDesc{
			kind:     paramExpandSearchReplaceLongestFirstMatch,
			parts:    []string{"$" + testValue, "FOO", "BAR"},
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

func TestParseParamShellSpecialSearchReplaceLongestFirstMatch(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$/FOO/BAR}",
		"${*/FOO/BAR}",
		"${@/FOO/BAR}",
		"${#/FOO/BAR}",
		"${?/FOO/BAR}",
		"${-/FOO/BAR}",
		"${0/FOO/BAR}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandSearchReplaceLongestFirstMatch,
			parts: []string{"$" + testData[2:3], "FOO", "BAR"},
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

func TestParseParamShellSpecialSearchReplaceLongestFirstMatchWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$/FOO/BAR}",
		"${!*/FOO/BAR}",
		"${!@/FOO/BAR}",
		"${!#/FOO/BAR}",
		"${!?/FOO/BAR}",
		"${!-/FOO/BAR}",
		"${!0/FOO/BAR}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandSearchReplaceLongestFirstMatch,
			parts:    []string{"$" + testData[3:4], "FOO", "BAR"},
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

func TestParseParamPlingSearchReplaceLongestFirstMatchDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!/FOO/BAR}"
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

func TestParseParamSearchReplaceLongestAllMatches(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR//FOO/BAR}"
	expectedResult := paramDesc{
		kind:  paramExpandSearchReplaceLongestAllMatches,
		parts: []string{"VAR", "FOO", "BAR"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSearchReplaceLongestAllMatchesWithNoReplacement(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${VAR//FOO}",
		"${VAR//FOO/}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandSearchReplaceLongestAllMatches,
			parts: []string{"VAR", "FOO", ""},
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

func TestParseParamSearchReplaceLongestAllMatchesWithNoSearchOrReplacement(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR//}"
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

func TestParseParamSearchReplaceLongestAllMatchesWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR//FOO/BAR}"
	expectedResult := paramDesc{
		kind:     paramExpandSearchReplaceLongestAllMatches,
		parts:    []string{"VAR", "FOO", "BAR"},
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

func TestParseParamSearchReplaceLongestAllMatchesSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V//FOO/BAR}"
	expectedResult := paramDesc{
		kind:  paramExpandSearchReplaceLongestAllMatches,
		parts: []string{"V", "FOO", "BAR"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSearchReplaceLongestAllMatchesSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V//FOO/BAR}"
	expectedResult := paramDesc{
		kind:     paramExpandSearchReplaceLongestAllMatches,
		parts:    []string{"V", "FOO", "BAR"},
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

func TestParseParamPositionalParamSearchReplaceLongestAllMatches(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "//FOO/BAR}"
		expectedResult := paramDesc{
			kind:  paramExpandSearchReplaceLongestAllMatches,
			parts: []string{"$" + testValue, "FOO", "BAR"},
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

func TestParseParamPositionalParamSearchReplaceLongestAllMatchesWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "//FOO/BAR}"
		expectedResult := paramDesc{
			kind:     paramExpandSearchReplaceLongestAllMatches,
			parts:    []string{"$" + testValue, "FOO", "BAR"},
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

func TestParseParamShellSpecialSearchReplaceLongestAllMatches(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$//FOO/BAR}",
		"${*//FOO/BAR}",
		"${@//FOO/BAR}",
		"${#//FOO/BAR}",
		"${?//FOO/BAR}",
		"${-//FOO/BAR}",
		"${0//FOO/BAR}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandSearchReplaceLongestAllMatches,
			parts: []string{"$" + testData[2:3], "FOO", "BAR"},
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

func TestParseParamShellSpecialSearchReplaceLongestAllMatchesWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$//FOO/BAR}",
		"${!*//FOO/BAR}",
		"${!@//FOO/BAR}",
		"${!#//FOO/BAR}",
		"${!?//FOO/BAR}",
		"${!-//FOO/BAR}",
		"${!0//FOO/BAR}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandSearchReplaceLongestAllMatches,
			parts:    []string{"$" + testData[3:4], "FOO", "BAR"},
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

func TestParseParamPlingSearchReplaceLongestAllMatchesDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!//FOO/BAR}"
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

func TestParseParamSearchReplaceLongestPrefix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR/#FOO/BAR}"
	expectedResult := paramDesc{
		kind:  paramExpandSearchReplaceLongestPrefix,
		parts: []string{"VAR", "FOO", "BAR"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSearchReplaceLongestPrefixWithNoReplacement(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${VAR/#FOO}",
		"${VAR/#FOO/}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandSearchReplaceLongestPrefix,
			parts: []string{"VAR", "FOO", ""},
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

func TestParseParamSearchReplaceLongestPrefixWithNoSearchOrReplacement(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR/#}"
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

func TestParseParamSearchReplaceLongestPrefixWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR/#FOO/BAR}"
	expectedResult := paramDesc{
		kind:     paramExpandSearchReplaceLongestPrefix,
		parts:    []string{"VAR", "FOO", "BAR"},
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

func TestParseParamSearchReplaceLongestPrefixSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V/#FOO/BAR}"
	expectedResult := paramDesc{
		kind:  paramExpandSearchReplaceLongestPrefix,
		parts: []string{"V", "FOO", "BAR"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSearchReplaceLongestPrefixSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V/#FOO/BAR}"
	expectedResult := paramDesc{
		kind:     paramExpandSearchReplaceLongestPrefix,
		parts:    []string{"V", "FOO", "BAR"},
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

func TestParseParamPositionalParamSearchReplaceLongestPrefix(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "/#FOO/BAR}"
		expectedResult := paramDesc{
			kind:  paramExpandSearchReplaceLongestPrefix,
			parts: []string{"$" + testValue, "FOO", "BAR"},
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

func TestParseParamPositionalParamSearchReplaceLongestPrefixWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "/#FOO/BAR}"
		expectedResult := paramDesc{
			kind:     paramExpandSearchReplaceLongestPrefix,
			parts:    []string{"$" + testValue, "FOO", "BAR"},
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

func TestParseParamShellSpecialSearchReplaceLongestPrefix(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$/#FOO/BAR}",
		"${*/#FOO/BAR}",
		"${@/#FOO/BAR}",
		"${#/#FOO/BAR}",
		"${?/#FOO/BAR}",
		"${-/#FOO/BAR}",
		"${0/#FOO/BAR}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandSearchReplaceLongestPrefix,
			parts: []string{"$" + testData[2:3], "FOO", "BAR"},
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

func TestParseParamShellSpecialSearchReplaceLongestPrefixWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$/#FOO/BAR}",
		"${!*/#FOO/BAR}",
		"${!@/#FOO/BAR}",
		"${!#/#FOO/BAR}",
		"${!?/#FOO/BAR}",
		"${!-/#FOO/BAR}",
		"${!0/#FOO/BAR}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandSearchReplaceLongestPrefix,
			parts:    []string{"$" + testData[3:4], "FOO", "BAR"},
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

func TestParseParamPlingSearchReplaceLongestPrefixDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!/#FOO/BAR}"
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

func TestParseParamSearchReplaceLongestSuffix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR/%FOO/BAR}"
	expectedResult := paramDesc{
		kind:  paramExpandSearchReplaceLongestSuffix,
		parts: []string{"VAR", "FOO", "BAR"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSearchReplaceLongestSuffixWithNoReplacement(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${VAR/%FOO}",
		"${VAR/%FOO/}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandSearchReplaceLongestSuffix,
			parts: []string{"VAR", "FOO", ""},
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

func TestParseParamSearchReplaceLongestSuffixWithNoSearchOrReplacement(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR/%}"
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

func TestParseParamSearchReplaceLongestSuffixWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR/%FOO/BAR}"
	expectedResult := paramDesc{
		kind:     paramExpandSearchReplaceLongestSuffix,
		parts:    []string{"VAR", "FOO", "BAR"},
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

func TestParseParamSearchReplaceLongestSuffixSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V/%FOO/BAR}"
	expectedResult := paramDesc{
		kind:  paramExpandSearchReplaceLongestSuffix,
		parts: []string{"V", "FOO", "BAR"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSearchReplaceLongestSuffixSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V/%FOO/BAR}"
	expectedResult := paramDesc{
		kind:     paramExpandSearchReplaceLongestSuffix,
		parts:    []string{"V", "FOO", "BAR"},
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

func TestParseParamPositionalParamSearchReplaceLongestSuffix(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "/%FOO/BAR}"
		expectedResult := paramDesc{
			kind:  paramExpandSearchReplaceLongestSuffix,
			parts: []string{"$" + testValue, "FOO", "BAR"},
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

func TestParseParamPositionalParamSearchReplaceLongestSuffixWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "/%FOO/BAR}"
		expectedResult := paramDesc{
			kind:     paramExpandSearchReplaceLongestSuffix,
			parts:    []string{"$" + testValue, "FOO", "BAR"},
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

func TestParseParamShellSpecialSearchReplaceLongestSuffix(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$/%FOO/BAR}",
		"${*/%FOO/BAR}",
		"${@/%FOO/BAR}",
		"${#/%FOO/BAR}",
		"${?/%FOO/BAR}",
		"${-/%FOO/BAR}",
		"${0/%FOO/BAR}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandSearchReplaceLongestSuffix,
			parts: []string{"$" + testData[2:3], "FOO", "BAR"},
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

func TestParseParamShellSpecialSearchReplaceLongestSuffixWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$/%FOO/BAR}",
		"${!*/%FOO/BAR}",
		"${!@/%FOO/BAR}",
		"${!#/%FOO/BAR}",
		"${!?/%FOO/BAR}",
		"${!-/%FOO/BAR}",
		"${!0/%FOO/BAR}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandSearchReplaceLongestSuffix,
			parts:    []string{"$" + testData[3:4], "FOO", "BAR"},
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

func TestParseParamPlingSearchReplaceLongestSuffixDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!/%FOO/BAR}"
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

func TestParseParamUppercaseFirstChar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR^abcde}"
	expectedResult := paramDesc{
		kind:  paramExpandUppercaseFirstChar,
		parts: []string{"VAR", "abcde"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamUppercaseFirstCharWithNoReplacement(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${VAR^}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

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
}

func TestParseParamUppercaseFirstCharWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR^abcde}"
	expectedResult := paramDesc{
		kind:     paramExpandUppercaseFirstChar,
		parts:    []string{"VAR", "abcde"},
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

func TestParseParamUppercaseFirstCharSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V^abcde}"
	expectedResult := paramDesc{
		kind:  paramExpandUppercaseFirstChar,
		parts: []string{"V", "abcde"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamUppercaseFirstCharSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V^abcde}"
	expectedResult := paramDesc{
		kind:     paramExpandUppercaseFirstChar,
		parts:    []string{"V", "abcde"},
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

func TestParseParamPositionalParamUppercaseFirstChar(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "^abcde}"
		expectedResult := paramDesc{
			kind:  paramExpandUppercaseFirstChar,
			parts: []string{"$" + testValue, "abcde"},
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

func TestParseParamPositionalParamUppercaseFirstCharWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "^abcde}"
		expectedResult := paramDesc{
			kind:     paramExpandUppercaseFirstChar,
			parts:    []string{"$" + testValue, "abcde"},
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

func TestParseParamShellSpecialUppercaseFirstChar(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$^abcde}",
		"${*^abcde}",
		"${@^abcde}",
		"${#^abcde}",
		"${?^abcde}",
		"${-^abcde}",
		"${0^abcde}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandUppercaseFirstChar,
			parts: []string{"$" + testData[2:3], "abcde"},
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

func TestParseParamShellSpecialUppercaseFirstCharWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$^abcde}",
		"${!*^abcde}",
		"${!@^abcde}",
		"${!#^abcde}",
		"${!?^abcde}",
		"${!-^abcde}",
		"${!0^abcde}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandUppercaseFirstChar,
			parts:    []string{"$" + testData[3:4], "abcde"},
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

func TestParseParamPlingUppercaseFirstCharDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!^abcde}"
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

func TestParseParamUppercaseAllChars(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR^^abcde}"
	expectedResult := paramDesc{
		kind:  paramExpandUppercaseAllChars,
		parts: []string{"VAR", "abcde"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamUppercaseAllCharsWithNoReplacement(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${VAR^^}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

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
}

func TestParseParamUppercaseAllCharsWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR^^abcde}"
	expectedResult := paramDesc{
		kind:     paramExpandUppercaseAllChars,
		parts:    []string{"VAR", "abcde"},
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

func TestParseParamUppercaseAllCharsSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V^^abcde}"
	expectedResult := paramDesc{
		kind:  paramExpandUppercaseAllChars,
		parts: []string{"V", "abcde"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamUppercaseAllCharsSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V^^abcde}"
	expectedResult := paramDesc{
		kind:     paramExpandUppercaseAllChars,
		parts:    []string{"V", "abcde"},
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

func TestParseParamPositionalParamUppercaseAllChars(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "^^abcde}"
		expectedResult := paramDesc{
			kind:  paramExpandUppercaseAllChars,
			parts: []string{"$" + testValue, "abcde"},
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

func TestParseParamPositionalParamUppercaseAllCharsWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "^^abcde}"
		expectedResult := paramDesc{
			kind:     paramExpandUppercaseAllChars,
			parts:    []string{"$" + testValue, "abcde"},
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

func TestParseParamShellSpecialUppercaseAllChars(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$^^abcde}",
		"${*^^abcde}",
		"${@^^abcde}",
		"${#^^abcde}",
		"${?^^abcde}",
		"${-^^abcde}",
		"${0^^abcde}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandUppercaseAllChars,
			parts: []string{"$" + testData[2:3], "abcde"},
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

func TestParseParamShellSpecialUppercaseAllCharsWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$^^abcde}",
		"${!*^^abcde}",
		"${!@^^abcde}",
		"${!#^^abcde}",
		"${!?^^abcde}",
		"${!-^^abcde}",
		"${!0^^abcde}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandUppercaseAllChars,
			parts:    []string{"$" + testData[3:4], "abcde"},
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

func TestParseParamPlingUppercaseAllCharsDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!^^abcde}"
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

func TestParseParamLowercaseFirstChar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR,abcde}"
	expectedResult := paramDesc{
		kind:  paramExpandLowercaseFirstChar,
		parts: []string{"VAR", "abcde"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamLowercaseFirstCharWithNoReplacement(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${VAR,}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

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
}

func TestParseParamLowercaseFirstCharWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR,abcde}"
	expectedResult := paramDesc{
		kind:     paramExpandLowercaseFirstChar,
		parts:    []string{"VAR", "abcde"},
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

func TestParseParamLowercaseFirstCharSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V,abcde}"
	expectedResult := paramDesc{
		kind:  paramExpandLowercaseFirstChar,
		parts: []string{"V", "abcde"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamLowercaseFirstCharSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V,abcde}"
	expectedResult := paramDesc{
		kind:     paramExpandLowercaseFirstChar,
		parts:    []string{"V", "abcde"},
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

func TestParseParamPositionalParamLowercaseFirstChar(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + ",abcde}"
		expectedResult := paramDesc{
			kind:  paramExpandLowercaseFirstChar,
			parts: []string{"$" + testValue, "abcde"},
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

func TestParseParamPositionalParamLowercaseFirstCharWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + ",abcde}"
		expectedResult := paramDesc{
			kind:     paramExpandLowercaseFirstChar,
			parts:    []string{"$" + testValue, "abcde"},
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

func TestParseParamShellSpecialLowercaseFirstChar(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$,abcde}",
		"${*,abcde}",
		"${@,abcde}",
		"${#,abcde}",
		"${?,abcde}",
		"${-,abcde}",
		"${0,abcde}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandLowercaseFirstChar,
			parts: []string{"$" + testData[2:3], "abcde"},
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

func TestParseParamShellSpecialLowercaseFirstCharWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$,abcde}",
		"${!*,abcde}",
		"${!@,abcde}",
		"${!#,abcde}",
		"${!?,abcde}",
		"${!-,abcde}",
		"${!0,abcde}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandLowercaseFirstChar,
			parts:    []string{"$" + testData[3:4], "abcde"},
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

func TestParseParamPlingLowercaseFirstCharDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!,abcde}"
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

func TestParseParamLowercaseAllChars(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR,,abcde}"
	expectedResult := paramDesc{
		kind:  paramExpandLowercaseAllChars,
		parts: []string{"VAR", "abcde"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamLowercaseAllCharsWithNoReplacement(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${VAR,,}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

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
}

func TestParseParamLowercaseAllCharsWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR,,abcde}"
	expectedResult := paramDesc{
		kind:     paramExpandLowercaseAllChars,
		parts:    []string{"VAR", "abcde"},
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

func TestParseParamLowercaseAllCharsSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V,,abcde}"
	expectedResult := paramDesc{
		kind:  paramExpandLowercaseAllChars,
		parts: []string{"V", "abcde"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamLowercaseAllCharsSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V,,abcde}"
	expectedResult := paramDesc{
		kind:     paramExpandLowercaseAllChars,
		parts:    []string{"V", "abcde"},
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

func TestParseParamPositionalParamLowercaseAllChars(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + ",,abcde}"
		expectedResult := paramDesc{
			kind:  paramExpandLowercaseAllChars,
			parts: []string{"$" + testValue, "abcde"},
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

func TestParseParamPositionalParamLowercaseAllCharsWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + ",,abcde}"
		expectedResult := paramDesc{
			kind:     paramExpandLowercaseAllChars,
			parts:    []string{"$" + testValue, "abcde"},
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

func TestParseParamShellSpecialLowercaseAllChars(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$,,abcde}",
		"${*,,abcde}",
		"${@,,abcde}",
		"${#,,abcde}",
		"${?,,abcde}",
		"${-,,abcde}",
		"${0,,abcde}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandLowercaseAllChars,
			parts: []string{"$" + testData[2:3], "abcde"},
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

func TestParseParamShellSpecialLowercaseAllCharsWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$,,abcde}",
		"${!*,,abcde}",
		"${!@,,abcde}",
		"${!#,,abcde}",
		"${!?,,abcde}",
		"${!-,,abcde}",
		"${!0,,abcde}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandLowercaseAllChars,
			parts:    []string{"$" + testData[3:4], "abcde"},
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

func TestParseParamPlingLowercaseAllCharsDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!,,abcde}"
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

func TestParseParamDescribeFlags(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR@a}"
	expectedResult := paramDesc{
		kind:  paramExpandDescribeFlags,
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

func TestParseParamDescribeFlagsWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR@a}"
	expectedResult := paramDesc{
		kind:     paramExpandDescribeFlags,
		parts:    []string{"VAR"},
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

func TestParseParamDescribeFlagsSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V@a}"
	expectedResult := paramDesc{
		kind:  paramExpandDescribeFlags,
		parts: []string{"V"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamDescribeFlagsSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V@a}"
	expectedResult := paramDesc{
		kind:     paramExpandDescribeFlags,
		parts:    []string{"V"},
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

func TestParseParamPositionalParamDescribeFlags(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "@a}"
		expectedResult := paramDesc{
			kind:  paramExpandDescribeFlags,
			parts: []string{"$" + testValue},
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

func TestParseParamPositionalParamDescribeFlagsWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "@a}"
		expectedResult := paramDesc{
			kind:     paramExpandDescribeFlags,
			parts:    []string{"$" + testValue},
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

func TestParseParamShellSpecialDescribeFlags(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$@a}",
		"${*@a}",
		"${@@a}",
		"${#@a}",
		"${?@a}",
		"${-@a}",
		"${0@a}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandDescribeFlags,
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

func TestParseParamShellSpecialDescribeFlagsWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$@a}",
		"${!*@a}",
		"${!@@a}",
		"${!#@a}",
		"${!?@a}",
		"${!-@a}",
		"${!0@a}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandDescribeFlags,
			parts:    []string{"$" + testData[3:4]},
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

func TestParseParamPlingDescribeFlagsDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!@a}"
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

func TestParseParamAsDeclare(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR@A}"
	expectedResult := paramDesc{
		kind:  paramExpandAsDeclare,
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

func TestParseParamAsDeclareWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR@A}"
	expectedResult := paramDesc{
		kind:     paramExpandAsDeclare,
		parts:    []string{"VAR"},
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

func TestParseParamAsDeclareSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V@A}"
	expectedResult := paramDesc{
		kind:  paramExpandAsDeclare,
		parts: []string{"V"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamAsDeclareSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V@A}"
	expectedResult := paramDesc{
		kind:     paramExpandAsDeclare,
		parts:    []string{"V"},
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

func TestParseParamPositionalParamAsDeclare(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "@A}"
		expectedResult := paramDesc{
			kind:  paramExpandAsDeclare,
			parts: []string{"$" + testValue},
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

func TestParseParamPositionalParamAsDeclareWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "@A}"
		expectedResult := paramDesc{
			kind:     paramExpandAsDeclare,
			parts:    []string{"$" + testValue},
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

func TestParseParamShellSpecialAsDeclare(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$@A}",
		"${*@A}",
		"${@@A}",
		"${#@A}",
		"${?@A}",
		"${-@A}",
		"${0@A}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandAsDeclare,
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

func TestParseParamShellSpecialAsDeclareWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$@A}",
		"${!*@A}",
		"${!@@A}",
		"${!#@A}",
		"${!?@A}",
		"${!-@A}",
		"${!0@A}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandAsDeclare,
			parts:    []string{"$" + testData[3:4]},
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

func TestParseParamPlingAsDeclareDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!@A}"
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

func TestParseParamEscaped(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR@E}"
	expectedResult := paramDesc{
		kind:  paramExpandEscaped,
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

func TestParseParamEscapedWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR@E}"
	expectedResult := paramDesc{
		kind:     paramExpandEscaped,
		parts:    []string{"VAR"},
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

func TestParseParamEscapedSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V@E}"
	expectedResult := paramDesc{
		kind:  paramExpandEscaped,
		parts: []string{"V"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamEscapedSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V@E}"
	expectedResult := paramDesc{
		kind:     paramExpandEscaped,
		parts:    []string{"V"},
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

func TestParseParamPositionalParamEscaped(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "@E}"
		expectedResult := paramDesc{
			kind:  paramExpandEscaped,
			parts: []string{"$" + testValue},
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

func TestParseParamPositionalParamEscapedWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "@E}"
		expectedResult := paramDesc{
			kind:     paramExpandEscaped,
			parts:    []string{"$" + testValue},
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

func TestParseParamShellSpecialEscaped(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$@E}",
		"${*@E}",
		"${@@E}",
		"${#@E}",
		"${?@E}",
		"${-@E}",
		"${0@E}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandEscaped,
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

func TestParseParamShellSpecialEscapedWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$@E}",
		"${!*@E}",
		"${!@@E}",
		"${!#@E}",
		"${!?@E}",
		"${!-@E}",
		"${!0@E}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandEscaped,
			parts:    []string{"$" + testData[3:4]},
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

func TestParseParamPlingEscapedDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!@E}"
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

func TestParseParamAsPrompt(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR@P}"
	expectedResult := paramDesc{
		kind:  paramExpandAsPrompt,
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

func TestParseParamAsPromptWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR@P}"
	expectedResult := paramDesc{
		kind:     paramExpandAsPrompt,
		parts:    []string{"VAR"},
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

func TestParseParamAsPromptSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V@P}"
	expectedResult := paramDesc{
		kind:  paramExpandAsPrompt,
		parts: []string{"V"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamAsPromptSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V@P}"
	expectedResult := paramDesc{
		kind:     paramExpandAsPrompt,
		parts:    []string{"V"},
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

func TestParseParamPositionalParamAsPrompt(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "@P}"
		expectedResult := paramDesc{
			kind:  paramExpandAsPrompt,
			parts: []string{"$" + testValue},
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

func TestParseParamPositionalParamAsPromptWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "@P}"
		expectedResult := paramDesc{
			kind:     paramExpandAsPrompt,
			parts:    []string{"$" + testValue},
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

func TestParseParamShellSpecialAsPrompt(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$@P}",
		"${*@P}",
		"${@@P}",
		"${#@P}",
		"${?@P}",
		"${-@P}",
		"${0@P}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandAsPrompt,
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

func TestParseParamShellSpecialAsPromptWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$@P}",
		"${!*@P}",
		"${!@@P}",
		"${!#@P}",
		"${!?@P}",
		"${!-@P}",
		"${!0@P}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandAsPrompt,
			parts:    []string{"$" + testData[3:4]},
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

func TestParseParamPlingAsPromptDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!@P}"
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

func TestParseParamSingleQuoted(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${VAR@Q}"
	expectedResult := paramDesc{
		kind:  paramExpandSingleQuoted,
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

func TestParseParamSingleQuotedWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!VAR@Q}"
	expectedResult := paramDesc{
		kind:     paramExpandSingleQuoted,
		parts:    []string{"VAR"},
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

func TestParseParamSingleQuotedSingleLetterVar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${V@Q}"
	expectedResult := paramDesc{
		kind:  paramExpandSingleQuoted,
		parts: []string{"V"},
	}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseParameter(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseParamSingleQuotedSingleLetterVarWithIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!V@Q}"
	expectedResult := paramDesc{
		kind:     paramExpandSingleQuoted,
		parts:    []string{"V"},
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

func TestParseParamPositionalParamSingleQuoted(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${" + testValue + "@Q}"
		expectedResult := paramDesc{
			kind:  paramExpandSingleQuoted,
			parts: []string{"$" + testValue},
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

func TestParseParamPositionalParamSingleQuotedWithIndirection(t *testing.T) {
	t.Parallel()

	for i := 1; i < 20; i++ {
		testValue := strconv.Itoa(i)
		// ----------------------------------------------------------------
		// setup your test

		testData := "${!" + testValue + "@Q}"
		expectedResult := paramDesc{
			kind:     paramExpandSingleQuoted,
			parts:    []string{"$" + testValue},
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

func TestParseParamShellSpecialSingleQuoted(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${$@Q}",
		"${*@Q}",
		"${@@Q}",
		"${#@Q}",
		"${?@Q}",
		"${-@Q}",
		"${0@Q}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:  paramExpandSingleQuoted,
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

func TestParseParamShellSpecialSingleQuotedWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$@Q}",
		"${!*@Q}",
		"${!@@Q}",
		"${!#@Q}",
		"${!?@Q}",
		"${!-@Q}",
		"${!0@Q}",
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := paramDesc{
			kind:     paramExpandSingleQuoted,
			parts:    []string{"$" + testData[3:4]},
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

func TestParseParamPlingSingleQuotedDoesNotSupportIndirection(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${!!@Q}"
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
