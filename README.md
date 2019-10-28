# Welcome to shellexpand

## Introduction

_ShellExpand_ is a replacement for Golang's standard `os.Expand()` that supports UNIX shell string substitution and expansion.

It is released under the 3-clause New BSD license. See [LICENSE.md](LICENSE.md) for details.

```golang
import (
    "os"
    shellexpand "github.com/ganbarodigital/go_shellexpand"
)

cb := ExpansionCallbacks{
    AssignVar: os.Setenv,
    LookupVar: os.Lookupenv,
}
result, err := shellexpand.Expand(input, cb)
```

## Table of Contents <!-- omit in toc -->

- [Introduction](#introduction)
- [Why Use ShellExpand?](#why-use-shellexpand)
  - [Who Is ShellExpand For?](#who-is-shellexpand-for)
  - [Why UNIX Shell String Expansion?](#why-unix-shell-string-expansion)
- [How Does It Work?](#how-does-it-work)
  - [Getting Started](#getting-started)
  - [How Are Errors Handled?](#how-are-errors-handled)
- [Expansion Callbacks](#expansion-callbacks)
  - [ExpansionCallbacks.AssignToVar()](#expansioncallbacksassigntovar)
  - [ExpansionCallbacks.LookupVar()](#expansioncallbackslookupvar)
  - [ExpansionCallbacks.LookupHomeDir()](#expansioncallbackslookuphomedir)
  - [ExpansionCallbacks.MatchVarNames()](#expansioncallbacksmatchvarnames)
- [Supported Expansions](#supported-expansions)
- [Brace Expansion](#brace-expansion)
- [What Is Brace Expansion?](#what-is-brace-expansion)
  - [Why Use Brace Expansion?](#why-use-brace-expansion)
  - [Rough Grammar](#rough-grammar)
  - [Other Notes](#other-notes)
  - [Status](#status)
- [Tilde Expansion](#tilde-expansion)
  - [What Is Tilde Expansion?](#what-is-tilde-expansion)
  - [Why Use Tilde Expansion?](#why-use-tilde-expansion)
  - [Rough Grammar](#rough-grammar-1)
  - [Other Notes](#other-notes-1)
  - [Status](#status-1)
- [Parameter Expansion](#parameter-expansion)
  - [What Is Parameter Expansion?](#what-is-parameter-expansion)
  - [Why Use Parameter Expansion?](#why-use-parameter-expansion)
  - [Supported Parameter Expansions](#supported-parameter-expansions)
- [Command Substitution](#command-substitution)
  - [What Is Command Substitution?](#what-is-command-substitution)
  - [Status](#status-2)
- [Arithmetic Expansion](#arithmetic-expansion)
  - [Rough Grammar](#rough-grammar-2)
  - [Status](#status-3)
- [Process Substitution](#process-substitution)
  - [Rough Grammar](#rough-grammar-3)
  - [Status](#status-4)
- [Word Splitting](#word-splitting)
  - [Rough Grammar](#rough-grammar-4)
  - [Status](#status-5)
- [Pathname Expansion](#pathname-expansion)
  - [Rough Grammar](#rough-grammar-5)
  - [Status](#status-6)
- [Escape Sequence Expansion](#escape-sequence-expansion)
  - [What Is Escape Sequence Expansion?](#what-is-escape-sequence-expansion)
  - [Rough Grammar](#rough-grammar-6)
  - [Status](#status-7)
- [Quote Removal](#quote-removal)
  - [What Is Quote Removal?](#what-is-quote-removal)
  - [Why Do Shells Perform Quote Removal?](#why-do-shells-perform-quote-removal)
  - [Status](#status-8)
- [Common Terms](#common-terms)
  - [Escaped Character](#escaped-character)
  - [Glob Pattern](#glob-pattern)
  - [Word](#word)
- [Reporting Problems](#reporting-problems)
- [Inspirations](#inspirations)

## Why Use ShellExpand?

### Who Is ShellExpand For?

We've built _ShellExpand_ for anyone who wants UNIX shell string-expansion support in their Golang projects.

We're using it ourselves in [Scriptish](https://github.com/ganbarodigital/go_scriptish), our library that helps you port your UNIX shell scripts into pure Golang.

### Why UNIX Shell String Expansion?

UNIX shell string expansion allows you to take a string of text, and replace it with the value of variables from your environment or from a backing store like [Envish](https://github.com/ganbarodigital/go_envish). You don't have to do a straight replacement either; there are lots of different ways that the value of a variable can be manipulated before it's substituted back into the original string.

Almost every non-trivial UNIX script uses string expansion. By bringing it to Golang, we hope that this will save you time when you're porting your shell scripts to Golang.

## How Does It Work?

### Getting Started

Import _ShellExpand_ into your Golang code:

```golang
import shellexpand "github.com/ganbarodigital/go_shellexpand"
```

Create a set of [expansion callbacks](#expansion-callbacks):

```golang
cb := ExpansionCallbacks{
    AssignToVar: func(...)
    LookupVar: func(...)
    LookupHomeDir: func(...)
    MatchVarNames: func(...)
}
```

Call `shellexpand.Expand()` to expand your string:

```golang
output, err := shellexpand.Expand(input, cb)
```

### How Are Errors Handled?

As a general principle, when an expansion fails, the input string is returned unmodified - and no error is returned.

That said, UNIX shells have grown organically over the last 40+ years, and there will be times where a failed expansion is simply removed from the input string instead ... and still no error returned.

Our main unit tests are in [`expand_test.go`](expand_test.go), and each supported expansion is run through a real UNIX shell too, to confirm that `ShellExpand` is as 100% compatible as possible.

Golang errors can come from two places:

* they can be returned from your [expansion callbacks](#expansion-callbacks)
* they can be caused by using invalid [glob patterns](#glob-pattern)

We return all errors back to you. When we do, the contents of the string we return is undefined.

## Expansion Callbacks

The vast majority of supported string expansions need to look things up:

* variables
* home directories for users

You need to tell _ShellExpand_ how to do this.

These callbacks are responsible for searching and updating what we call your _variable backing store_. This might be your Golang program's environment, or it might be an emulation of that like our own [Envish package](https://github.com/ganbarodigital/go_envish).

There are four callbacks that you need to define:

### ExpansionCallbacks.AssignToVar()

```golang
func AssignToVar(key, value string) error
```

`ShellExpand` will call `AssignToVar()` when it needs to set the value of a variable. Your callback must update your variable backing store.

If an error occurs, your callback should return that error back to `ShellExpand`. We'll then pass that back to the caller of `shellexpand.Expand()`.

### ExpansionCallbacks.LookupVar()

```golang
func LookupVar(key string) (string, bool)
```

`ShellExpand` will call `LookupVar()` when it needs to get the value of a variable.

* If the variable exists in your backing store, return its value and `true`
* If the variable does not exist in your backing store, return `""` (empty string) and `false`

__Please do not return `""` and `true` if the variable does not exist in your backing store.__ That behaviour may lead to undefined results from _ShellExpand_. Even if your own testing says that you can get away with it today, we do not guarantee you'll get the results you expect in a future release.

(If you're familiar with Golang's `os.LookupEnv()`, `LookupVar` does the same job.)

### ExpansionCallbacks.LookupHomeDir()

```golang
func LookupHomeDir(user string) (string, bool)
```

`ShellExpand` will call `LookupHomeDir()` when it needs to know the home directory of a user. This is needed for [tilde expansion](#tilde-expansion).

* If the user exists on your computer, return the user's home directory and `true`
* Otherwise, return `""` (empty string) and `false`

You can use Golang's `os/user` package to do this.

### ExpansionCallbacks.MatchVarNames()

```golang
func MatchVarNames(prefix string) []string
```

`ShellExpand` will call `MatchVarNames` when it needs to know which variable names start with the given prefix. This is needed for [parameter expansion](#parameter-expansion).

Your callback must return a list of all variable names that start with the given prefix. If no names match, return an empty list.

## Supported Expansions

UNIX shells perform 10 different types of string expansion. This table tracks which ones we currently support, and what we (currently) plan to do about the rest of them.

Expansion                                               | Status                    | Planned?
--------------------------------------------------------|---------------------------|---------
[Brace expansion](#brace-expansion)                     | fully supported           | n/a
[Tilde expansion](#tilde-expansion)                     | fully supported           | n/a
[Parameter expansion](#parameter-expansion)             | (almost) fully supported  | n/a
[Command substitution](#command-substitution)           | not supported             | no plans to add
[Arithmetic expansion](#arithmetic-expansion)           | not supported             | yes
[Process substitution](#process-substitution)           | not supported             | no plans to add
[Word splitting](#word-splitting)                       | not supported             | if there is a need
[Pathname expansion](#pathname-expansion)               | not supported             | if there is a need
[Quote removal](#quote-removal)                         | partial support           | depends on feedback
[Escape sequence expansion](#escape-sequence-expansion) | not supported             | no plans to

We have put more details about each of them below.

## Brace Expansion

## What Is Brace Expansion?

_Brace expansion_ is a way to generate longer strings from simple patterns or range sequences.

This example uses a `brace-pattern`:

```golang
input := "ab{c,d,e}fg"
cb := ExpansionCallbacks{}
output := shellexpand.Expand(input, cb)

// output is: abcfg abdfg abefg
```

This example uses a `brace-sequence`:

```golang
input := "ab{c..e}fg"
cb := ExpansionCallbacks{}
output := shellexpand.Expand(input, cb)

// output is: abcfe abdfg abefg
```

### Why Use Brace Expansion?

It's commonly used in shell scripts to expand a list of filenames _without_ checking that they exist (ie, without relying on [pathname expansion](#pathname-expansion)).

### Rough Grammar

Brace expansion takes the form:

`[preamble]((brace-pattern|brace-sequence)+)[postscript]`

where:

* `preamble` is optional text immediately before the `brace-pattern` or `brace-sequence`
* `brace-pattern` is `{text,text[,text...]}`
  * each `brace-pattern` is surrounded by braces
  * each `brace-pattern` contains a comma-separated list of text to substitute in
  * each `brace-pattern` must contain _at least_ two phrases to be a valid `brace-pattern`
* `brace-sequence` is `{lo..hi[..incr]}` or `{hi..lo[..incr]}`
  * `hi` and `lo` can be characters or number, as long as they're both the same type
  * `incr` is optional, and must be a number
  * `incr`'s sign is always auto-corrected to match the order you've put `hi` and `lo` in
* `postscript` is optional text immediately after the `brace-pattern` or `brace-sequence`

### Other Notes

* Brace expansions can be nested.
* Left-to-right order is preserved. The result of a brace expansion is never sorted.
* You can escape the opening brace (ie do `\\{`) to prevent a brace triggering brace expansion.

### Status

_Brace expansion_ is fully supported in v1.0.0 and later.

If you find any bugs in brace expansion, please [let us know](#reporting-problems).

## Tilde Expansion

### What Is Tilde Expansion?

_Tilde expansion_ turns `~` into the path to a user's home directory.

```golang
input := "~/.storyplayer/storyplayer.json"
cb := ExpansionCallbacks{
    // you need to provide this
    LookupHomeDir: func(...),
}
output, err := shellexpand.Expand(input)

// on my system, output would be: /home/stuart/.storyplayer/storyplayer.json
```

The `~` can be followed by a username.

```golang
input := "~stuart/.storyplayer/storyplayer.json"
cb := ExpansionCallbacks{
    // you need to provide this
    LookupHomeDir: func(...),
}
output, err := shellexpand.Expand(input)

// on my system, output would be: /home/stuart/.storyplayer/storyplayer.json
```

The `~` can be followed by a `+` or `-`. A `+` is replaced by the value of `PWD`, and a `-` is replaced by the value of `OLDPWD`.

```golang
input := "~+/storyplayer.json"
cb := ExpansionCallbacks{
    // you need to provide this
    LookupVar: func(...),
}
output, err := shellexpand.Expand(input)

// on my system, output would be: (current working directory)/storyplayer.json
```

### Why Use Tilde Expansion?

Many programs these days create/maintain/support a config file stored in the user's home directory. Tilde expansion is a very easy way

### Rough Grammar

Tilde expansion takes the form:

`~(+|-|username|<blank>)[/path/to/folder/or/file]`

where:

* `~` must be the first character of the [word](#word)
* `~+` is replaced by the value of `PWD` via a call to [LookupVar](#expansioncallbackslookupvar)
* `~-` is replaced by the value of `OLDPWD` via a call to [LookupVar](#expansioncallbackslookupvar)
* `~username` is replaced by user's home directory, via a call to [LookupHomeDir](#expansioncallbackslookuphomedir)
* `~` on its own is replaced by the value of `HOME` via a call to [LookupVar](#expansioncallbackslookupvar)

The `/path/to/folder/or/file` is optional.

### Other Notes

* Tilde expansion does not check that the expanded filepath is valid, or that whatever it points to exists.

### Status

_Tilde expansion_ is fully supported in v1.0.0 or later.

If you find any bugs in tilde expansion, please [let us know](#reporting-problems).

## Parameter Expansion

### What Is Parameter Expansion?

_Parameter expansion_ replaces variable names with their values, and allows you to change those values before they're plugged back into the result string.

Parameter expansion makes up the majority of UNIX shell string expansion (and, arguably, the majority of your average shell script's code too).

### Why Use Parameter Expansion?

In a word: convenience. Being able to assign a value to a variable, and then use that value whereever you need to ... that's one of the main reasons that shell scripts are so quick and easy to knock up.

### Supported Parameter Expansions

Here's a list of all the different types of UNIX shell parameter expansion that we know about. Most of them are fully supported by _ShellExpand_.

Syntax                        | Name                              | Status
------------------------------|-----------------------------------|-------
`$PARAM` / `${PARAM}`         | expand-to-value                   | supported
`${PARAM}`                    | expand-to-value                   | supported
`${PARAM:-word}`              | expand-with-default-value         | supported
`${PARAM:=word}`              | expand-assign-default-value       | supported
`${PARAM:?word}`              | expand-write-error                | supported
`${PARAM:+word}`              | expand-use-alternate-value        | supported
`${PARAM:offset}`             | expand-to-substring               | supported
`${PARAM:offset:length}`      | expand-to-substring-length        | supported
`${!prefix*}` / `${!prefix@}` | expand-prefix-match-names         | supported
`${!name[*]}` / `${!name[@]}` | list-of-array-keys                | not supported
`${#PARAM}`                   | expand-parameter-length           | supported
`${#*}` / `${#@}`             | expand-no-positional-params       | supported
`${PARAM#pattern}`            | expand-remove-shortest-prefix     | supported
`${PARAM##pattern}`           | expand-remove-longest-prefix      | supported
`${PARAM%pattern}`            | expand-remove-shortest-suffix     | supported
`${PARAM%%pattern}`           | expand-remove-longest-suffix      | supported
`${PARAM/old/new}`            | expand-search-replace-all-matches | supported
`${PARAM//old/new}`           | expand-search-replace-first-match | supported
`${PARAM/#old/new}`           | expand-search-replace-prefix      | supported
`${PARAM/%old/new}`           | expand-search-replace-suffix      | supported
`${PARAM^pattern}`            | expand-uppercase-first-char       | supported
`${PARAM^^pattern}`           | expand-uppercase-all-chars        | supported
`${PARAM,pattern}`            | expand-lowercase-first-char       | supported
`${PARAM,,pattern}`           | expand-lowercase-all-chars        | supported
`${PARAM@operator}`           | expand-parameter-transform        | not supported

## Command Substitution

### What Is Command Substitution?

_Command substitution_ calls an external program, and puts the output from that program in the returned string.

```bash
#!/usr/bin/env bash

CURRENT_BRACH=$(git branch --no-color | grep '^\* ' | grep -v 'no branch' | sed 's/^* //g')
```

or

```bash
#!/usr/bin/env bash

CURRENT_BRACH=`git branch --no-color | grep '^\* ' | grep -v 'no branch' | sed 's/^* //g'`
```

### Status

_Command substitution_ is __not supported__.

There are no plans to add support for command substitution at this time.

* Command substition needs [word splitting](#word-splitting) implementing first.
* If the input string has come from user input, it should be treated as _untrusted_ to avoid security problems. Calling arbitrary external programs from string expansion is asking for trouble.

If we add support in a future version, we'll make it send the command name and arguments to a callback that you provide.

## Arithmetic Expansion

### Rough Grammar

### Status

## Process Substitution

### Rough Grammar

### Status

## Word Splitting

### Rough Grammar

### Status

## Pathname Expansion

### Rough Grammar

### Status

## Escape Sequence Expansion

### What Is Escape Sequence Expansion?

_Escape sequence expansion_ is

### Rough Grammar

An escape sequence is a `\` (forward-slash) followed by one or more characters. The supported characters are in this table:

Escape Sequence | Shell Expansion
----------------|----------------
`\a`            | alert (bell)
`\b`            | backspace
`\e` and `\E`   | ANSI escape sequence
`\f`            | form feed
`\n`            | new line
`\r`            | carriage return
`\t`            | horizontal tab
`\v`            | vertical tab
`\'`            | escaped single quote
`\"`            | escaped double quote
`\?`            | escaped question mark
`\nnn`          | the 8-bit character for the octal number `nnn`
`\xNN`          | the 8-bit character for the hexadecimal `XX`
`\uHHHH`        | the Unicode character for the hexadecimal `HHHH`
`\UHHHHHHHH`    | the Unicode character for the hexadecimal `HHHHHHHH`
`\cX`           | a `control-X` character

Any other sequence starting with a `\` is treated as [an escaped character](#escaped-character).

### Status

_Escape sequence expansion_ is **not supported**.

There are no plans to add support for escape sequence expansion at the moment.

Why? Many escape sequences exist for working with interactive shells. There's no direct target to translate them to in a Golang library. Many (all?) of the rest are already supported by Golang's `fmt` package.

## Quote Removal

### What Is Quote Removal?

_Quote removal_ is the removal of:

* the `\` in front of escaped characters
* single quotes surrounding [words](#word)
* double quotes surrounding [words](#word)

### Why Do Shells Perform Quote Removal?

If a [word](#word) contains spaces - for example, a long filename on Windows or MacOS - you have to surround it with either single or double quotes. That's how a UNIX shell knows that the filename is a single word.

The quotes aren't actually part of the filename on disk. The UNIX shell has to remove them, otherwise

### Status

_Quote removal_ is partially supported in `ShellExpand` v1.0.0.

* the `\` in front of escaped characters is removed
  - this was required for our unit tests (which compares results in a real UNIX shell and our code) to work
* single and double quotes surrounding words are not removed (because [word splitting](#word-splitting) is not supported).

## Common Terms

### Escaped Character

An _escaped character_ is any character that has a `\` (forward-slash) in front of it.

They are used to tell _ShellExpand_ to treat the escaped character as a normal character (ie, just text).

For example, in [brace expansion](#brace-expansion), a `{` (brace) on its own denotes the start of a brace pattern:

```bash
# ab{c,d,e}fg becomes 'abcfg abdfg abefg'
```

If you escape the brace, _ShellExpand_ no longer treats it as the start of a brace pattern

```bash
# ab\{c,d,e}fg remains `ab{c,d,e}fg`
```

All of our expansions should correctly support escaped characters. By correctly, we mean that it should:

* always treat them as normal text
* the result should be the same that a real UNIX shell would give you

If you find any bugs related to this, please [let us know](#reporting-problems).

### Glob Pattern

### Word

## Reporting Problems

## Inspirations