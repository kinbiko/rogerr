/*
Package rogerr is a zero-dependency error handling support package.

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
*/
package rogerr
