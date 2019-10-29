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

func TestMatchVarSingleSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${this} is a test"
	expectedEnd := 7

	// ----------------------------------------------------------------
	// perform the change

	actualEnd, ok := matchVar(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedEnd, actualEnd)
	assert.True(t, ok)
}

func TestMatchVarNestedSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${HOME:${TMPDIR:-/var/tmp}} a test"
	expectedEnd := 27

	// ----------------------------------------------------------------
	// perform the change

	actualEnd, ok := matchVar(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedEnd, actualEnd)
	assert.True(t, ok)
}

func TestMatchVarIgnoresMissingDollar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{HOME:${TMPDIR:-/var/tmp}} a test"
	expectedEnd := 0

	// ----------------------------------------------------------------
	// perform the change

	actualEnd, ok := matchVar(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedEnd, actualEnd)
	assert.False(t, ok)
}

func TestMatchVarSupportsMissingOpeningBrace(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "$HOME a test"
	expectedEnd := 5

	// ----------------------------------------------------------------
	// perform the change

	actualEnd, ok := matchVar(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedEnd, actualEnd)
	assert.Equal(t, testData[:actualEnd], "$HOME")
	assert.True(t, ok)
}

func TestMatchVarIgnoresMissingClosingBraceMidString(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${HOME a test"
	expectedEnd := 0

	// ----------------------------------------------------------------
	// perform the change

	actualEnd, ok := matchVar(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedEnd, actualEnd)
	assert.False(t, ok)
}

func TestMatchVarIgnoresMissingClosingBraceWholeString(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${HOME"
	expectedEnd := 0

	// ----------------------------------------------------------------
	// perform the change

	actualEnd, ok := matchVar(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedEnd, actualEnd)
	assert.False(t, ok)
}

func TestMatchVarIgnoresEscapedClosingBrace(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "${HOME\\}}"
	expectedEnd := 9

	// ----------------------------------------------------------------
	// perform the change

	actualEnd, ok := matchVar(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedEnd, actualEnd)
	assert.True(t, ok)
}

func TestMatchVarKnownParamOperators(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := []string{
		"$var",
		"${var}",
		"${var:-word}",
		"${var:=word}",
		"${var:?word}",
		"${var:100}",
		"${var:1:5}",
		"${!prefix*}",
		"${#var}",
		"${#*}",
		"${var#word}",
		"${var##word}",
		"${var%suffix}",
		"${var%%suffix}",
		"${*%suffix}",
		"${*%%suffix}",
		"${var/old/new}",
		"${*/old/new}",
		"${var^pattern}",
		"${var^^pattern}",
		"${var,pattern}",
		"${var,,pattern}",
		"${var@a}",
		"${var@A}",
		"${var@E}",
		"${var@P}",
		"${var@Q}",
	}

	// ----------------------------------------------------------------
	// perform the change

	for i := range testData {
		testResult, ok := matchVar(testData[i])

		assert.True(t, ok)
		assert.Equal(t, testData[i], testData[i][:testResult])
	}

	// ----------------------------------------------------------------
	// test the results

}
