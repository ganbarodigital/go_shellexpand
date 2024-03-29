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

func isAlphaChar(char rune) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z'
}

func isAlphaNumericChar(char rune) bool {
	return isNumericChar(char) || isAlphaChar(char)
}

func isNumericChar(char rune) bool {
	return '0' <= char && char <= '9'
}

func isNumericStartChar(char rune) bool {
	return '1' <= char && char <= '9'
}

func isNumericString(input string) bool {
	for _, c := range input {
		if !isNumericChar(c) {
			return false
		}
	}

	return true
}

func isSignedNumericString(input string) bool {
	// special case - zero
	if input == "0" {
		return true
	}

	// no leading zeros allowed ... but can be a negative value
	if input[0] != '-' && !isNumericStartChar(rune(input[0])) {
		return false
	}

	for i := 1; i < len(input); i++ {
		if !isNumericChar(rune(input[i])) {
			return false
		}
	}

	return true
}

func isNumericStringWithoutLeadingZero(input string) bool {
	if len(input) == 0 {
		return false
	}

	if !isNumericStartChar(rune(input[0])) {
		return false
	}

	for i := 1; i < len(input); i++ {
		if !isNumericChar(rune(input[i])) {
			return false
		}
	}

	return true
}

func isNameBodyChar(char rune) bool {
	return isAlphaNumericChar(char) || char == '_'
}

func isNameStartChar(char rune) bool {
	return isAlphaChar(char) || char == '_'
}

func isShellSpecialChar(char rune) bool {
	return char == '#' || char == '*' || char == '?' || char == '!' || char == '$' || char == '-' || char == '@' || char == '0'
}
