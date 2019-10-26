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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAlphaCharReturnsTrueForLowercaseChars(t *testing.T) {
	t.Parallel()

	for _, testData := range "abcdefghijklmnopqrstuvwxyz" {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := true

		// ----------------------------------------------------------------
		// perform the change

		actualResult := isAlphaChar(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestIsAlphaCharReturnsTrueForUppercaseChars(t *testing.T) {
	t.Parallel()

	for _, testData := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := true

		// ----------------------------------------------------------------
		// perform the change

		actualResult := isAlphaChar(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestIsAlphaCharReturnsFalseOtherwise(t *testing.T) {
	t.Parallel()

	for _, testData := range "1234567890!\"£$%^&*()[]{};:'@#~,<.>/?`" {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := false

		// ----------------------------------------------------------------
		// perform the change

		actualResult := isAlphaChar(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestIsSignedNumericStringReturnsTrueForZero(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "0"
	expectedResult := true

	// ----------------------------------------------------------------
	// perform the change

	actualResult := isSignedNumericString(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestIsSignedNumericStringReturnsTrueForPositiveNumbers(t *testing.T) {
	t.Parallel()

	for _, testData := range []string{
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"10",
		"11",
		"12",
		"23",
		"14",
		"1576",
		"15157",
	} {
		// ----------------------------------------------------------------
		// setup your test

		expectedResult := true

		// ----------------------------------------------------------------
		// perform the change

		actualResult := isSignedNumericString(testData)

		// ----------------------------------------------------------------
		// test the results

		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestIsSignedNumericStringReturnsFalseIfNumberStartsWithLeadingZero(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "0123"
	expectedResult := false

	// ----------------------------------------------------------------
	// perform the change

	actualResult := isSignedNumericString(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestIsNumericStringWithoutLeadingZeroReturnsFalseForEmptyString(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := ""
	expectedResult := false

	// ----------------------------------------------------------------
	// perform the change

	actualResult := isNumericStringWithoutLeadingZero(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestIsSignedNumericStringReturnsFalseIfNumberContainsLetters(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "123abc"
	expectedResult := false

	// ----------------------------------------------------------------
	// perform the change

	actualResult := isSignedNumericString(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
