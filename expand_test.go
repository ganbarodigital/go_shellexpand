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
	expectedError        string
	resultSubstringMatch bool
	actualResult         func(expandTestData) string
}

func TestExpandBraceExpansion(t *testing.T) {
	// simple string, w/ brace expansion
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "ab{c,d,e}fg",
		expectedResult: "abcfg abdfg abefg",
	}
	testExpandTestCase(t, testData)
}

func TestExpandUnterminatedBraceExpansion(t *testing.T) {
	// simple string, w/ mismatched braces
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "ab{c,d,efg",
		expectedResult: "ab{c,d,efg",
	}
	testExpandTestCase(t, testData)
}

func TestExpandBraceExpansionSinglePattern(t *testing.T) {
	// simple string, w/ single pattern
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "ab{c}de",
		expectedResult: "ab{c}de",
	}
	testExpandTestCase(t, testData)
}

func TestExpandEscapedBraces(t *testing.T) {
	// simple string, w/ escaped braces
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "\\{PARAM1\\}",
		expectedResult: "{PARAM1}",
	}
	testExpandTestCase(t, testData)
}

func TestExpandBraceSequenceAlphasNoIterator(t *testing.T) {
	// simple string, w/ alpha sequence
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "ab{a..g}de",
		expectedResult: "abade abbde abcde abdde abede abfde abgde",
	}
	testExpandTestCase(t, testData)
}

func TestExpandBraceSequenceAlphasAndIterator(t *testing.T) {
	// simple string, w/ alpha sequence
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "ab{a..g..2}de",
		expectedResult: "abade abcde abede abgde",
	}
	testExpandTestCase(t, testData)
}

func TestExpandBraceSequenceReverseAlphasNoIterator(t *testing.T) {
	// simple string, w/ alpha sequence descending
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "ab{g..a}de",
		expectedResult: "abgde abfde abede abdde abcde abbde abade",
	}
	testExpandTestCase(t, testData)
}

func TestExpandBraceSequenceReverseAlphasAndIterator(t *testing.T) {
	// simple string, w/ alpha sequence descending
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "ab{g..a..2}de",
		expectedResult: "abgde abede abcde abade",
	}
	testExpandTestCase(t, testData)
}

func TestExpandBraceSequenceWithNoIterator(t *testing.T) {
	// simple string, w/ numerical sequence
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "ab{1..4}de",
		expectedResult: "ab1de ab2de ab3de ab4de",
	}
	testExpandTestCase(t, testData)
}

func TestExpandBraceSequenceWithIterator(t *testing.T) {
	// simple string, w/ numerical sequence
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "ab{1..11..3}de",
		expectedResult: "ab1de ab4de ab7de ab10de",
	}
	testExpandTestCase(t, testData)
}

func TestExpandBraceSequenceWithIteratorHasWrongSign(t *testing.T) {
	// simple string, w/ numerical sequence and incr in wrong direction
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "ab{1..11..-3}de",
		expectedResult: "ab1de ab4de ab7de ab10de",
	}
	testExpandTestCase(t, testData)
}

func TestExpandSimpleParam(t *testing.T) {
	// simple param, no braces
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "$PARAM1",
		expectedResult: "foo",
	}
	testExpandTestCase(t, testData)
}

func TestExpandSimpleParamInLongerString(t *testing.T) {
	// simple param inside longer string, no braces
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "this is all $PARAM1 bar",
		expectedResult: "this is all foo bar",
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalVar1(t *testing.T) {
	// simple param, positional var $1
	testData := expandTestData{
		positionalVars: map[string]string{
			"$1": "foo",
		},
		input:          "$1",
		expectedResult: "foo",
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalVar2(t *testing.T) {
	// simple param, positional var $2
	testData := expandTestData{
		positionalVars: map[string]string{
			"$1": "foo",
			"$2": "bar",
		},
		input:          "$2",
		expectedResult: "bar",
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalVar3(t *testing.T) {
	// simple param, positional var $3
	testData := expandTestData{
		positionalVars: map[string]string{
			"$1": "foo",
			"$2": "bar",
			"$3": "alfred",
		},
		input:          "$3",
		expectedResult: "alfred",
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalVar4(t *testing.T) {
	// simple param, positional var $4
	testData := expandTestData{
		positionalVars: map[string]string{
			"$1": "foo",
			"$2": "bar",
			"$3": "alfred",
			"$4": "trout",
		},
		input:          "$4",
		expectedResult: "trout",
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalVar5(t *testing.T) {
	// simple param, positional var $5
	testData := expandTestData{
		positionalVars: map[string]string{
			"$1": "foo",
			"$2": "bar",
			"$3": "alfred",
			"$4": "trout",
			"$5": "haddock",
		},
		input:          "$5",
		expectedResult: "haddock",
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalVar6(t *testing.T) {
	// simple param, positional var $6
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalVar7(t *testing.T) {
	// simple param, positional var $7
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalVar8(t *testing.T) {
	// simple param, positional var $8
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalVar9(t *testing.T) {
	// simple param, positional var $9
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalVar10WithoutBraces(t *testing.T) {
	// simple param, positional var $10
	//
	// bash stops parsing after matching '$1'
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalVar10WithBraces(t *testing.T) {
	// simple param, positional var $10 wrapped in brackets
	//
	// this IS supported by bash
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandSpecialVarDollarStar(t *testing.T) {
	// simple param, expand all positional vars via $*
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandSpecialVarDollarAt(t *testing.T) {
	// simple param, expand all positional vars via $@
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandSimpleParamInBraces(t *testing.T) {
	// simple param, braces
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "${PARAM1}",
		expectedResult: "foo",
	}
	testExpandTestCase(t, testData)
}

func TestExpandSimpleParamInBracesAndInLongerString(t *testing.T) {
	// simple param inside longer string, braces
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "this is all ${PARAM1}bar",
		expectedResult: "this is all foobar",
	}
	testExpandTestCase(t, testData)
}

// func TestExpandInvalidUnterminatedParamInLongerString(t *testing.T) {
// 	// invalid (unterminated) param inside longer string, braces
// 	testData := expandTestData{
// 		vars: map[string]string{
// 			"PARAM1": "foo",
// 		},
// 		input:          "this is all ${++bar",
// 		expectedResult: "this is all ${++bar",
// 	}
// 	testExpandTestCase(t, testData)
// }

// func TestExpandInvalidParamNameInLongerString(t *testing.T) {
// 	// invalid param inside longer string, braces
// 	testData := expandTestData{
// 		vars: map[string]string{
// 			"PARAM1": "foo",
// 		},
// 		input:          "this is all ${++}bar",
// 		expectedResult: "this is all ${++}bar",
// 	}
// 	testExpandTestCase(t, testData)
// }

func TestExpandParamWithIndirection(t *testing.T) {
	// simple param, braces, indirection
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "PARAM2",
			"PARAM2": "foo",
		},
		input:          "${!PARAM1}",
		expectedResult: "foo",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamToDefaultValue(t *testing.T) {
	// simple param, default value triggered
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "",
		},
		input:          "${PARAM1:-foo}",
		expectedResult: "foo",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamNotToDefaultValue(t *testing.T) {
	// simple param, default value NOT triggered
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "foo",
		},
		input:          "${PARAM1:-bar}",
		expectedResult: "foo",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamToDefaultValueWithIndirection(t *testing.T) {
	// simple param, default value triggered, indirection
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "PARAM2",
			"PARAM2": "",
		},
		input:          "${!PARAM1:-foo}",
		expectedResult: "foo",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamNotToDefaultValueWithIndirection(t *testing.T) {
	// simple param, default value NOT triggered, indirection
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "PARAM2",
			"PARAM2": "foo",
		},
		input:          "${!PARAM1:-bar}",
		expectedResult: "foo",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamToDefaultValueWithWordExpansion(t *testing.T) {
	// simple param, default value triggered AND expanded
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "",
			"PARAM2": "bar",
		},
		input:          "${PARAM1:-${PARAM2}}",
		expectedResult: "bar",
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalParamToDefaultValue(t *testing.T) {
	// positional param, default value triggered
	testData := expandTestData{
		input:          "${1:-foo}",
		expectedResult: "foo",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamSetToDefaultValue(t *testing.T) {
	// simple param, default value set
	testData := expandTestData{
		input: "${PARAM1:=foo}",
		shellExtra: []string{
			"dummy=${PARAM1:=foo}",
			"echo $PARAM1",
		},
		expectedResult: "foo",
		actualResult: func(testData expandTestData) string {
			return testData.vars["PARAM1"]
		},
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamSetToDefaultValueWithWordExpansion(t *testing.T) {
	// simple param, default value set to word expansion
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamSetToDefaultValueWithIndirection(t *testing.T) {
	// indirect param, default value set to word expansion
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamErrorWritten(t *testing.T) {
	// simple param, error written
	testData := expandTestData{
		vars: map[string]string{
			"foo": "",
		},
		input:                "${foo:?not set}",
		expectedResult:       "foo: not set",
		resultSubstringMatch: true,
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamToAlternativeValue(t *testing.T) {
	// simple param, use alternative value
	testData := expandTestData{
		vars: map[string]string{
			"foo": "bar",
		},
		input:          "${foo:+alternative}",
		expectedResult: "alternative",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamSubstring(t *testing.T) {
	// simple param, expand substring to end of value
	testData := expandTestData{
		vars: map[string]string{
			"foo": "1234567890",
		},
		input:          "${foo:3}",
		expectedResult: "4567890",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamSubstringWithIndirection(t *testing.T) {
	// simple param, expand substring to end of value with indirection
	testData := expandTestData{
		vars: map[string]string{
			"foo": "bar",
			"bar": "1234567890",
		},
		input:          "${!foo:3}",
		expectedResult: "4567890",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamToSubstringAndLength(t *testing.T) {
	// simple param, expand substring to given length
	testData := expandTestData{
		vars: map[string]string{
			"foo": "1234567890",
		},
		input:          "${foo:3:4}",
		expectedResult: "4567",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamToSubstringAndLengthWithIndirection(t *testing.T) {
	// simple param, expand substring to given length with indirection
	testData := expandTestData{
		vars: map[string]string{
			"foo": "bar",
			"bar": "1234567890",
		},
		input:          "${!foo:3:4}",
		expectedResult: "4567",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamNamesByPrefixStar(t *testing.T) {
	// expand param names by prefix with * suffix
	testData := expandTestData{
		vars: map[string]string{
			"foo1": "bar",
			"foo2": "humbug",
		},
		input:          "${!foo*}",
		expectedResult: "foo1 foo2",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamNamesByPrefixAt(t *testing.T) {
	// expand param names by prefix with @ suffix
	testData := expandTestData{
		vars: map[string]string{
			"foo1": "bar",
			"foo2": "humbug",
		},
		input:          "${!foo@}",
		expectedResult: "foo1 foo2",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamLength(t *testing.T) {
	// length of simple param
	testData := expandTestData{
		vars: map[string]string{
			"foo": "bar humbug",
		},
		input:          "${#foo}",
		expectedResult: "10",
	}
	testExpandTestCase(t, testData)
}

func TestExpandNumberOfPositionalParamsDollarStar(t *testing.T) {
	// number of positional parameters via $*
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandNumberOfPositionalParamsDollarAt(t *testing.T) {
	// number of positional parameters via $@
	testData := expandTestData{
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
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamRemoveShortestPrefix(t *testing.T) {
	// remove prefix shortest match
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "docdoc",
		},
		input:          "${PARAM1#d*c}",
		expectedResult: "doc",
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalParamsRemoveShortestPrefix(t *testing.T) {
	// remove prefix, shortest match, applied to $*
	testData := expandTestData{
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
		input:          "${*#[a-z]}",
		expectedResult: "oo ar lfred rout addock od laice ollock hitebait ank",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamRemoveLongestPrefix(t *testing.T) {
	// remove prefix longest match
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "docdocgo",
		},
		input:          "${PARAM1##d*c}",
		expectedResult: "go",
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalParamsRemoveLongestPrefix(t *testing.T) {
	// remove prefix, longest match, applied to $*
	testData := expandTestData{
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
		input:          "${*##*o}",
		expectedResult: "bar alfred ut ck d plaice ck whitebait bank",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamRemoveShortestSuffix(t *testing.T) {
	// remove suffix shortest match
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "godocdoc",
		},
		input:          "${PARAM1%d*c}",
		expectedResult: "godoc",
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalParamsRemoveShortestSuffix(t *testing.T) {
	// remove suffix, shortest match, applied to $*
	testData := expandTestData{
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
		input:          "${*%o*}",
		expectedResult: "fo bar alfred tr hadd c plaice poll whitebait bank",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamRemoveLongestSuffix(t *testing.T) {
	// remove suffix longest match
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "godocdoc",
		},
		input:          "${PARAM1%%d*c}",
		expectedResult: "go",
	}
	testExpandTestCase(t, testData)
}

func TestExpandPositionalParamsRemoveLongestSuffix(t *testing.T) {
	// remove suffix, longest match, applied to $*
	testData := expandTestData{
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
		input:          "${*%%o*}",
		expectedResult: "f bar alfred tr hadd c plaice p whitebait bank",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamUppercaseFirstLetterNoPattern(t *testing.T) {
	// uppercase first letter, no replacement pattern
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "alfred",
		},
		input:          "${PARAM1^}",
		expectedResult: "Alfred",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamUppercaseFirstLetterWithPattern(t *testing.T) {
	// uppercase first letter, replacement pattern matches
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "alfred",
		},
		input:          "${PARAM1^[a-z]}",
		expectedResult: "Alfred",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamUppercaseFirstLetterDoesNotMatchPattern(t *testing.T) {
	// uppercase first letter, replacement pattern does not match
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "alfred",
		},
		input:          "${PARAM1^[A-Z]}",
		expectedResult: "alfred",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamUppercaseAllCharsNoPattern(t *testing.T) {
	// uppercase all chars, no replacement pattern
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "alfred",
		},
		input:          "${PARAM1^^}",
		expectedResult: "ALFRED",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamUppercaseAllCharsMatchesPattern(t *testing.T) {
	// uppercase all chars, replacement pattern matches
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "alfred",
		},
		input:          "${PARAM1^^[a-z]}",
		expectedResult: "ALFRED",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamUppercaseAllCharsPartiallyMatchesPattern(t *testing.T) {
	// uppercase all chars, replacement pattern partial matches
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "alfred",
		},
		input:          "${PARAM1^^[a-m]}",
		expectedResult: "ALFrED",
	}
	testExpandTestCase(t, testData)
}

func TestExpandUppercaseAllCharsDoNotMatchPattern(t *testing.T) {
	// uppercase all chars, replacement pattern does not match
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "alfred",
		},
		input:          "${PARAM1^^[0-9]}",
		expectedResult: "alfred",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamLowercaseFirstCharNoPattern(t *testing.T) {
	// lowercase first letter, no replacement pattern
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "ALFRED",
		},
		input:          "${PARAM1,}",
		expectedResult: "aLFRED",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamLowercaseFirstCharMatchesPattern(t *testing.T) {
	// lowercase first letter, replacement pattern matches
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "ALFRED",
		},
		input:          "${PARAM1,[A-Z]}",
		expectedResult: "aLFRED",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamLowercaseFirstCharDoesNotMatchPattern(t *testing.T) {
	// lowercase first letter, replacement pattern does not match
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "ALFRED",
		},
		input:          "${PARAM1,[0-9]}",
		expectedResult: "ALFRED",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamLowercaseAllCharsNoPattern(t *testing.T) {
	// lowercase all chars, no replacement pattern
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "ALFRED",
		},
		input:          "${PARAM1,,}",
		expectedResult: "alfred",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamLowercaseAllCharsMatchesPattern(t *testing.T) {
	// lowercase all chars, replacement pattern matches
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "ALFRED",
		},
		input:          "${PARAM1,,[A-Z]}",
		expectedResult: "alfred",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamLowercaseAllCharsPartiallyMatchesPattern(t *testing.T) {
	// lowercase all chars, replacement pattern partial matches
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "ALFRED",
		},
		input:          "${PARAM1,,[A-M]}",
		expectedResult: "alfRed",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamLowercaseAllCharsDoesNotMatchPattern(t *testing.T) {
	// lowercase all chars, replacement pattern does not match
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "ALFRED",
		},
		input:          "${PARAM1,,[0-9]}",
		expectedResult: "ALFRED",
	}
	testExpandTestCase(t, testData)
}

func TestExpandParamLowercaseAllCharsInvalidPattern(t *testing.T) {
	testData := expandTestData{
		vars: map[string]string{
			"PARAM1": "ALFRED",
		},
		input:          "${PARAM1,,[0-9}",
		expectedResult: "",
		expectedError:  "bad or unsupported glob pattern '[0-9': error parsing regexp: missing closing ]: `[0-9$`",
	}
	testExpandTestCase(t, testData)
}

func testExpandTestCase(t *testing.T, testData expandTestData) {
	// ----------------------------------------------------------------
	// create the shell script we'll run

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

	// export the shell script we'll use to verify that internal behaviour
	// matches actual shell script behaviour
	tmpFile, _ := ioutil.TempFile("", "shellexpand-expandParams-")
	cleanup := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}
	defer cleanup()

	tmpFile.WriteString(buf.String())
	tmpFile.Sync()
	tmpFile.Close()

	// ----------------------------------------------------------------
	// to run the test, we need to create some helper methods

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
	expectedError := testData.expectedError

	// ----------------------------------------------------------------
	// perform the change

	cmd := exec.Command("/usr/bin/env", "bash", tmpFile.Name())
	shellRawResult, _ := cmd.CombinedOutput()
	shellActualResult := strings.TrimSpace(string(shellRawResult))

	internalActualResult, internalActualError := Expand(input, varFuncs)
	// special case - the result is a side effect, not a direct string
	// expansion
	if testData.actualResult != nil {
		internalActualResult = testData.actualResult(testData)
	}

	// ----------------------------------------------------------------
	// test the results

	// assert.Nil(t, shellErr)
	if len(expectedError) > 0 {
		assert.Error(t, internalActualError)
		assert.Equal(t, expectedError, internalActualError.Error())

		// if we are in an error situation, we don't care what
		// the shell did
		if testData.resultSubstringMatch {
			assert.Contains(t, internalActualResult, expectedResult, testData)
		} else {
			assert.Equal(t, expectedResult, internalActualResult, testData)
		}
	} else {
		assert.Nil(t, internalActualError)

		if testData.resultSubstringMatch {
			assert.Contains(t, shellActualResult, expectedResult, buf.String())
			assert.Contains(t, internalActualResult, expectedResult, testData)
		} else {
			assert.Equal(t, expectedResult, shellActualResult, buf.String())
			assert.Equal(t, expectedResult, internalActualResult, testData)
		}
	}
}
