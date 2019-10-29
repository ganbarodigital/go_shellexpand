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

func TestExpandTildeHomedir(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			if key == "HOME" {
				return "/home/stuart", true
			}

			return "invalid key", true
		},
		LookupHomeDir: func(key string) (string, bool) {
			return "should not be called", true
		},
	}
	testData := "~/path/to/folder"
	expectedResult := "/home/stuart/path/to/folder"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := ExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandTildePwd(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			if key == "PWD" {
				return "/tmp", true
			}

			return "invalid key", true
		},
		LookupHomeDir: func(key string) (string, bool) {
			return "should not be called", true
		},
	}
	testData := "~+/path/to/folder"
	expectedResult := "/tmp/path/to/folder"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := ExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandTildeOldPwd(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			if key == "OLDPWD" {
				return "/tmp", true
			}

			return "invalid key", true
		},
		LookupHomeDir: func(key string) (string, bool) {
			return "should not be called", true
		},
	}
	testData := "~-/path/to/folder"
	expectedResult := "/tmp/path/to/folder"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := ExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandTildeUserHomedir(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			return "should not be called", true
		},
		LookupHomeDir: func(key string) (string, bool) {
			if key == "stuart" {
				return "/home/stuart", true
			}
			return "invalid key", true
		},
	}
	testData := "~stuart/path/to/folder"
	expectedResult := "/home/stuart/path/to/folder"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := ExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandTildeDoesNotModifyStringsWithoutTildePrefixes(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			return "invalid key", true
		},
		LookupHomeDir: func(key string) (string, bool) {
			return "should not be called", true
		},
	}
	testData := "/path/to/folder"
	expectedResult := "/path/to/folder"

	// ----------------------------------------------------------------
	// perform the change

	actualResult := ExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandTildeIgnoresTildeInsideVars(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			return "invalid key", true
		},
		LookupHomeDir: func(key string) (string, bool) {
			return "should not be called", true
		},
	}
	testData := "${VAR1:~VAR2}"
	expectedResult := testData

	// ----------------------------------------------------------------
	// perform the change

	actualResult := ExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandTildeIgnoresEscapedTilde(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			return "invalid key", true
		},
		LookupHomeDir: func(key string) (string, bool) {
			return "should not be called", true
		},
	}
	testData := "\\~/path"
	expectedResult := testData

	// ----------------------------------------------------------------
	// perform the change

	actualResult := ExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandTildeIgnoresWhenHomedirNotSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			return "invalid key", false
		},
		LookupHomeDir: func(key string) (string, bool) {
			return "should not be called", true
		},
	}
	testData := "~/path"
	expectedResult := testData

	// ----------------------------------------------------------------
	// perform the change

	actualResult := ExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandTildeIgnoresWhenPwdNotSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			return "invalid key", false
		},
		LookupHomeDir: func(key string) (string, bool) {
			return "should not be called", true
		},
	}
	testData := "~+/path"
	expectedResult := testData

	// ----------------------------------------------------------------
	// perform the change

	actualResult := ExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandTildeIgnoresWhenOldPwdNotSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			return "invalid key", false
		},
		LookupHomeDir: func(key string) (string, bool) {
			return "should not be called", true
		},
	}
	testData := "~-/path"
	expectedResult := testData

	// ----------------------------------------------------------------
	// perform the change

	actualResult := ExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestExpandTildeIgnoresWhenUsernameNotKnown(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			return "should not be called", true
		},
		LookupHomeDir: func(key string) (string, bool) {
			return "invalid key", false
		},
	}
	testData := "~user/path"
	expectedResult := testData

	// ----------------------------------------------------------------
	// perform the change

	actualResult := ExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestParseTildePrefixWithHomedir(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "~"
	expectedResult := tildePrefix{tildePrefixHome, ""}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseTildePrefix(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseTildePrefixWithPwd(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "~+"
	expectedResult := tildePrefix{tildePrefixPwd, ""}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseTildePrefix(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseTildePrefixWithOldPwd(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "~-"
	expectedResult := tildePrefix{tildePrefixOldPwd, ""}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseTildePrefix(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseTildePrefixWithUsername(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "~stuart"
	expectedResult := tildePrefix{tildePrefixUsername, "stuart"}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseTildePrefix(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestParseTildePrefixWithoutTilde(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "/root"
	expectedResult := tildePrefix{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := parseTildePrefix(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchTildePrefixWithHomedir(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "~/path/to/folder"
	expectedResult := 1

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchTildePrefix(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchTildePrefixWithPwd(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "~+/src"
	expectedResult := 2

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchTildePrefix(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchTildePrefixWithOldPwd(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "~-/bin"
	expectedResult := 2

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchTildePrefix(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchTildePrefixWithUsername(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "~stuart"
	expectedResult := 7

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchTildePrefix(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchTildePrefixWithoutTilde(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "/root"
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchTildePrefix(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchTildePrefixIgnoresEscapedSlashes(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "~\\/path/to/folder"
	expectedResult := 7

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchTildePrefix(testData)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}

func TestMatchAndExpandTildeIgnoresNonPrefix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	cb := ExpansionCallbacks{
		LookupVar: func(key string) (string, bool) {
			return "invalid key", true
		},
		LookupHomeDir: func(key string) (string, bool) {
			return "should not be called", true
		},
	}
	testData := "/path"
	expectedResult := testData

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := matchAndExpandTilde(testData, cb)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Equal(t, expectedResult, actualResult)
}
