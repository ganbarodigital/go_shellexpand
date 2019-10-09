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

func TestExpandBracesSingleSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "a{b,c,d}e"
	expectedResult := "abe ace ade"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := expandBraces(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandBracesSingleSetWithEmptyPart(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "/var/log/kern.log{,.bak}"
	expectedResult := "/var/log/kern.log /var/log/kern.log.bak"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := expandBraces(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandBracesNestedSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "/usr/{ucb/{ex,edit}/tmp1,lib/{ex?.?*,how_ex}/tmp2}"
	expectedResult := "/usr/ucb/ex/tmp1 /usr/ucb/edit/tmp1 /usr/lib/ex?.?*/tmp2 /usr/lib/how_ex/tmp2"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := expandBraces(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchPatternSingleSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{b,c,d}"
	expectedResult := len(testData) - 1

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchPattern(testData, 0)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult+1])
}

func TestMatchPatternNestedSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{ucb/{ex,edit}/tmp1,lib/{ex?.?*,how_ex}/tmp2}"
	expectedResult := len(testData) - 1

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchPattern(testData, 0)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult+1])
}

func TestMatchPatternNoOpeningBrace(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "this is not a pattern}"
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchPattern(testData, 0)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchPatternSkipEscapedBraces(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{this is \\{ a \\}pattern}"
	expectedResult := len(testData) - 1

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchPattern(testData, 0)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult+1])
}

func TestMatchPatternSkipDollarVars(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{this is ${a} pattern}"
	expectedResult := len(testData) - 1

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchPattern(testData, 0)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult+1])
}

func TestMatchPatternIgnoresUnterminatedPatterns(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{this is ${a} pattern"
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchPattern(testData, 0)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchSequenceSingleSetWithLowerCaseChars(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{a..z}"
	expectedResult := len(testData) - 1

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchSequence(testData, 0)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult+1])
}

func TestMatchSequenceSingleSetWithUpperCaseChars(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{A..Z}"
	expectedResult := len(testData) - 1

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchSequence(testData, 0)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult+1])
}

func TestMatchSequenceSingleSetWithNumbers(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{1..99}"
	expectedResult := len(testData) - 1

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchSequence(testData, 0)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult+1])
}

func TestMatchSequenceSingleSetWithIterator(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{1..99..3}"
	expectedResult := len(testData) - 1

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchSequence(testData, 0)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult+1])
}

func TestParsePatternSingleSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{b,c,d}"
	expectedResult := []string{"b", "c", "d"}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parsePattern(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParsePatternNestedSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{ucb/{ex,edit}/tmp1,lib/{ex?.?*,how_ex}/tmp2}"
	expectedResult := []string{"ucb/{ex,edit}/tmp1", "lib/{ex?.?*,how_ex}/tmp2"}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parsePattern(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParsePatternWithEmptyPart(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{,.bak}"
	expectedResult := []string{"", ".bak"}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parsePattern(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}
