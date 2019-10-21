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

const (
	paramOpInvalid = iota
	paramOpUseDefaultValue
	paramOpAssignDefaultValue
	paramOpWriteError
	paramOpUseAlternativeValue
	paramOpSubstring
	paramOpRemoveShortestPrefix
	paramOpRemoveLongestPrefix
	paramOpRemoveShortestSuffix
	paramOpRemoveLongestSuffix
	paramOpSearchReplace
	paramOpUppercaseFirstChar
	paramOpUppercaseAllMatches
	paramOpLowercaseFirstChar
	paramOpLowercaseAllMatches
	paramOpDescribeFlags
	paramOpDeclare
	paramOpEscape
	paramOpExpandAsPrompt
	paramOpExpandDoubleQuotes
	// this has been added to help us test unsupported operand rejection
	// in the parameter parser
	paramOpEmptyObject
)

func matchParamOp(input string, start int) (int, int, bool) {
	// shorthand
	inputLen := len(input)
	maxInput := inputLen - 1
	startPlus1 := start + 1

	// what are we looking at?
	switch input[start] {
	case ':':
		if start == maxInput {
			return paramOpInvalid, 0, false
		}
		switch input[startPlus1] {
		case '-':
			return paramOpUseDefaultValue, startPlus1, true
		case '=':
			return paramOpAssignDefaultValue, startPlus1, true
		case '?':
			return paramOpWriteError, startPlus1, true
		case '+':
			return paramOpUseAlternativeValue, startPlus1, true
		default:
			return paramOpSubstring, start, true
		}
	case '#':
		if start < maxInput && input[startPlus1] == '#' {
			return paramOpRemoveLongestPrefix, startPlus1, true
		}

		return paramOpRemoveShortestPrefix, start, true
	case '%':
		if start < maxInput && input[startPlus1] == '%' {
			return paramOpRemoveLongestSuffix, startPlus1, true
		}
		return paramOpRemoveShortestSuffix, start, true
	case '/':
		return paramOpSearchReplace, start, true
	case '^':
		if start < maxInput && input[startPlus1] == '^' {
			return paramOpUppercaseAllMatches, startPlus1, true
		}
		return paramOpUppercaseFirstChar, start, true
	case ',':
		if start < maxInput && input[startPlus1] == ',' {
			return paramOpLowercaseAllMatches, startPlus1, true
		}
		return paramOpLowercaseFirstChar, start, true
	case '@':
		if start == maxInput {
			return paramOpInvalid, 0, false
		}
		switch input[startPlus1] {
		case 'a':
			return paramOpDescribeFlags, startPlus1, true
		case 'A':
			return paramOpDeclare, startPlus1, true
		case 'E':
			return paramOpEscape, startPlus1, true
		case 'P':
			return paramOpExpandAsPrompt, startPlus1, true
		case 'Q':
			return paramOpExpandDoubleQuotes, startPlus1, true
		default:
			return paramOpInvalid, 0, false
		}
	case '{':
		if start == maxInput {
			return paramOpInvalid, 0, false
		}
		switch input[startPlus1] {
		case '}':
			return paramOpEmptyObject, startPlus1, true
		default:
			return paramOpInvalid, 0, false
		}
	default:
		return paramOpInvalid, 0, false
	}
}
