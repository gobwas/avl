# avl

[![GoDoc][godoc-image]][godoc-url]
[![CI][ci-badge]][ci-url]

> AVL (Adelson-Velsky and Landis) immutable tree implementation.

# Overview

This is an *immutable* implementation of the balanced binary search tree.
Its goal is to be as simple as possible in terms of API as well as correct and generic.

# Installation

```bash
go get github.com/gobwas/avl
```

# Documentation

You can read the docs at [GoDoc][godoc-url].

# Usage

```go
package main

import (
	"strings"

	"github.com/gobwas/avl"
)

func main() {
	var tree avl.Tree
	tree, _ = tree.Insert(StringItem("foo"))
	tree, _ = tree.Delete(StringItem("foo"))
	if tree.Search(StringItem("foo")) != nil {
		// whoa!
	}
}

type StringItem string

func (s StringItem) Compare(x Item) int {
	return strings.Compare(string(s), string(x.(StringItem)))
}
```


[godoc-image]: https://godoc.org/github.com/gobwas/avl?status.svg
[godoc-url]:   https://godoc.org/github.com/gobwas/avl
[ci-badge]:    https://github.com/gobwas/avl/actions/workflows/main.yml/badge.svg?branch=main
[ci-url]:      https://github.com/gobwas/avl/actions/workflows/main.yml
