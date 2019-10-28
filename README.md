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
  - [Expansion Callbacks](#expansion-callbacks)
    - [ExpansionCallbacks.AssignToVar()](#expansioncallbacksassigntovar)
  - [ExpansionCallbacks.LookupVar()](#expansioncallbackslookupvar)
  - [ExpansionCallbacks.LookupHomeDir()](#expansioncallbackslookuphomedir)
  - [ExpansionCallbacks.MatchVarNames()](#expansioncallbacksmatchvarnames)
  - [How Are Errors Handled?](#how-are-errors-handled)

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

### Expansion Callbacks

The vast majority of supported string expansions need to look things up:

* variables
* home directories for users

You need to tell _ShellExpand_ how to do this.

These callbacks are responsible for searching and updating what we call your _variable backing store_. This might be your Golang program's environment, or it might be an emulation of that like our own [Envish package](https://github.com/ganbarodigital/go_envish).

There are four callbacks that you need to define:

#### ExpansionCallbacks.AssignToVar()

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

### How Are Errors Handled?

As a general principle, when an expansion fails, the input string is returned unmodified - and no error is returned.

That said, UNIX shells have grown organically over the last 40+ years, and there will be times where a failed expansion is simply removed from the input string instead ... and still no error returned.

Our main unit tests are in [`expand_test.go`](expand_test.go), and each supported expansion is run through a real UNIX shell too, to confirm that `ShellExpand` is as 100% compatible as possible.

Golang errors can come from two places:

* they can be returned from your [expansion callbacks](#expansion-callbacks)
* they can be caused by using invalid [glob patterns](#glob-pattern)

We return all errors back to you. When we do, the contents of the string we return is undefined.

