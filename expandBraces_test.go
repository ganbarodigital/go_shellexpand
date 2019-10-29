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

func TestExpandBracesSingleSequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "so is {1..9}"
	expectedResult := "so is 1 2 3 4 5 6 7 8 9"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := expandBraces(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandBracesMalformedVariable(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "so is ${++"
	expectedResult := "so is ${++"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := expandBraces(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandBracesMalformedVariableInsidePattern(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "so is {${++"
	expectedResult := "so is {${++"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := expandBraces(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandBracesPatternAndSequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "this is a te{st,ab}{1..3}ing"
	expectedResult := "this is a test1ing test2ing test3ing teab1ing teab2ing teab3ing"

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
	expectedResult := len(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBracePattern(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult])
}

func TestMatchPatternNestedSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{ucb/{ex,edit}/tmp1,lib/{ex?.?*,how_ex}/tmp2}"
	expectedResult := len(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBracePattern(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult])
}

func TestMatchPatternNoOpeningBrace(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "this is not a pattern}"
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBracePattern(testData)

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
	expectedResult := len(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBracePattern(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult])
}

func TestMatchPatternSkipDollarVars(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{this is ${a} pattern}"
	expectedResult := len(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBracePattern(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult])
}

func TestMatchPatternIgnoresUnterminatedPatterns(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{this is ${a} pattern"
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBracePattern(testData)

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
	expectedResult := len(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBraceSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult])
}

func TestMatchSequenceSingleSetWithUpperCaseChars(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{A..Z}"
	expectedResult := len(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBraceSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult])
}

func TestMatchSequenceSingleSetWithNumbers(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{1..99}"
	expectedResult := len(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBraceSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult])
}

func TestMatchSequenceSingleSetWithIterator(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{1..99..3}"
	expectedResult := len(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBraceSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult])
}

func TestMatchSequenceSingleSetWithNegativeIterator(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{1..99..-3}"
	expectedResult := len(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBraceSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, testData, testData[:actualResult])
}

func TestMatchSequenceMustStartWithBrace(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "a{..z}"
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBraceSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchSequenceRejectsNestedBraces(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{a..z{a..z}}"
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBraceSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchSequenceRejectsMismatchedBraces(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{1..99"
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchBraceSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParsePatternSingleSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{b,c,d}"
	expectedResult := []string{"b", "c", "d"}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseBracePattern(testData)

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

	actualResult, ok := parseBracePattern(testData)

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

	actualResult, ok := parseBracePattern(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParsePatternWithEscapedChars(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{\\{b,c,d\\}}"
	expectedResult := []string{"\\{b", "c", "d\\}"}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseBracePattern(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParsePatternWithMismatchedBraces(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{b,c,d"
	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseBracePattern(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParsePatternWithSinglePattern(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{b}"
	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseBracePattern(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseSequenceSingleSetWithLowerCaseChars(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{a..z}"
	expectedResult := braceSequence{true, 97, 122, 1}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseSequenceSingleSetWithUpperCaseChars(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{A..Z}"
	expectedResult := braceSequence{true, 65, 90, 1}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseSequenceSingleSetWithNumbers(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{1..99}"
	expectedResult := braceSequence{false, 1, 99, 1}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseSequenceSingleSetWithIterator(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{1..99..3}"
	expectedResult := braceSequence{false, 1, 99, 3}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseSequenceSingleSetWithNegativeIterator(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{99..1..-3}"
	expectedResult := braceSequence{false, 99, 1, -3}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseSequenceSingleSetHighToLow(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{99..1}"
	expectedResult := braceSequence{false, 99, 1, -1}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseSequenceSingleSetHighToLowWithIncrement(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{99..1..2}"
	expectedResult := braceSequence{false, 99, 1, -2}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseSequenceRejectsMismatchedSequenceCharNum(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{a..1}"
	expectedResult := braceSequence{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseSequenceRejectsMismatchedSequenceNumChar(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{1..a}"
	expectedResult := braceSequence{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseSequenceRejectsNonIntegerIncrement(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "{1..5..a}"
	expectedResult := braceSequence{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseSequence(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}
