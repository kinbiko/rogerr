# rogerr

[![Build Status](https://github.com/kinbiko/rogerr/workflows/Go/badge.svg)](https://github.com/kinbiko/rogerr/actions)
[![Coverage Status](https://coveralls.io/repos/github/kinbiko/rogerr/badge.svg?branch=main)](https://coveralls.io/github/kinbiko/rogerr?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/kinbiko/rogerr)](https://goreportcard.com/report/github.com/kinbiko/rogerr)
[![Latest version](https://img.shields.io/github/tag/kinbiko/rogerr.svg?label=latest%20version&style=flat)](https://github.com/kinbiko/rogerr/releases)
[![Go Documentation](http://img.shields.io/badge/godoc-documentation-blue.svg?style=flat)](https://pkg.go.dev/github.com/kinbiko/rogerr?tab=doc)
[![License](https://img.shields.io/github/license/kinbiko/rogerr.svg?style=flat)](https://github.com/kinbiko/rogerr/blob/master/LICENSE)

A Go package for error handling with structured metadata. Zero dependencies.

[Blog post with detailed explanation](https://kinbiko.com/posts/2022-07-30-error-messages-should-be-boring/).

## Problem

Error messages that include unique data (user IDs, timestamps, etc.) break error grouping in monitoring tools like Sentry and Rollbar.

## Solution

Store unique data as structured metadata separate from the error message.
This package attaches metadata to Go's `context.Context` and preserves it when wrapping errors.

## Usage

1. Create an ErrorHandler: `handler := rogerr.NewErrorHandler()`
2. Add metadata to context: `ctx = rogerr.WithMetadatum(ctx, "userID", 123)`
3. Wrap errors with metadata: `err = handler.Wrap(ctx, err, "operation failed")`
4. Extract metadata for logging: `metadata := rogerr.Metadata(err)`

### Build Options

For cleaner stacktraces, use the `-trimpath` flag:

```bash
go build -trimpath ./cmd/myapp
```

This shows module-relative paths instead of absolute paths.
Skip this if you disable stacktraces with `rogerr.WithStacktrace(false)`.

[Full documentation](https://pkg.go.dev/github.com/kinbiko/rogerr)

