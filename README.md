# rogerr

[![Build Status](https://github.com/kinbiko/rogerr/workflows/Go/badge.svg)](https://github.com/kinbiko/rogerr/actions)
[![Coverage Status](https://coveralls.io/repos/github/kinbiko/rogerr/badge.svg?branch=main)](https://coveralls.io/github/kinbiko/rogerr?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/kinbiko/rogerr)](https://goreportcard.com/report/github.com/kinbiko/rogerr)
[![Latest version](https://img.shields.io/github/tag/kinbiko/rogerr.svg?label=latest%20version&style=flat)](https://github.com/kinbiko/rogerr/releases)
[![Go Documentation](http://img.shields.io/badge/godoc-documentation-blue.svg?style=flat)](https://pkg.go.dev/github.com/kinbiko/rogerr?tab=doc)
[![License](https://img.shields.io/github/license/kinbiko/rogerr.svg?style=flat)](https://github.com/kinbiko/rogerr/blob/master/LICENSE)

Consistent and greppable errors makes your logger and error reporting tools happy. This zero-dependency error handling support package for Go helps you achieve just that.

[Blog post explaining the problem and the solution in detail](https://kinbiko.com/posts/2022-07-30-error-messages-should-be-boring/).

## Usage

When creating errors, **do not include goroutine-specific or request-specific information as part of the error message itself**.
Error messages with these specific bits of information often break filtering/grouping algorithms, e.g. as used by error reporting tools like Sentry/Rollbar/etc. (If you use Bugsnag, I recommend [kinbiko/bugsnag](https://github.com/kinbiko/bugsnag) for an **even better** experience than this package).

Instead this information should be treated as structured data, akin to structured logging solutions like Logrus and Zap.
In Go, it's conventional to attach this kind of request specific 'diagnostic' metadata to a `context.Context` type, and that's what this package enables too.

At a high level:

1. Create an ErrorHandler with `handler := rogerr.NewErrorHandler()`.
1. Attach metadata to your context with `rogerr.WithMetadata` or `rogerr.WithMetadatum`.
1. When you come across an error, use `err = handler.Wrap(ctx, err, msg)` to attach the metadata accumulated so far to the wrapped error.
1. Return the error as you would normally, and at the time of logging/reporting, extract the metadata with `md := rogerr.Metadata(err)`.
1. Record the _structured_ metadata alongside the error message.

For more details, see [the official docs](https://pkg.go.dev/github.com/kinbiko/rogerr).

### Build Recommendations

For cleaner stacktrace file paths, build with the `-trimpath` flag:

```bash
go build -trimpath ./cmd/myapp
```

This removes local build path prefixes, showing module-relative paths instead of absolute machine-specific paths.
Not needed if stacktraces are disabled with `WithStacktrace(false)`, e.g. for performance reasons.
