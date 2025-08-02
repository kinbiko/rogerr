module github.com/kinbiko/rogerr/internal/myapp

go 1.24

replace github.com/kinbiko/rogerr => ../..

replace github.com/kinbiko/rogerr/internal/mylib => ../mylib

require (
	github.com/kinbiko/rogerr v0.0.0-00010101000000-000000000000
	github.com/kinbiko/rogerr/internal/mylib v0.0.0-00010101000000-000000000000
)
