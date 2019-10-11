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
		"${!:=foo}",
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
		"${!:?foo}",
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
		"${!:+foo}",
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
		"${!:500}",
		"${$:500}",
		"${*:500}",
		"${@:500}",
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
		"${!@:500}",
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
		"${!:500:1000}",
		"${$:500:1000}",
		"${*:500:1000}",
		"${@:500:1000}",
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

func TestParseParamShellSpecialSubstringLengthWithIndirection(t *testing.T) {
	t.Parallel()

	testDataSet := []string{
		"${!$:500:1000}",
		"${!*:500:1000}",
		"${!@:500:1000}",
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
		"${#*}",
		"${#@}",
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
		"${!#FOO}",
		"${$#FOO}",
		"${*#FOO}",
		"${@#FOO}",
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
		"${!*#FOO}",
		"${!@#FOO}",
		"${!##FOO}",
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
