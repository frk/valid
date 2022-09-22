[`godoc`](http://godoc.org/github.com/frk/valid)
[`pkg.go.dev`](https://pkg.go.dev/github.com/frk/valid)

# valid

The package `valid` defines a number of validation functions and other objects that can
be used, together with the [`cmd/validgen`](cmd/validgen) tool, to generate static input validation.

--------

Some (most) validators, including their tests & comments, were ported from https://github.com/validatorjs/validator.js

validgen configuration
custom validation functions
	configuring from file
	configuring through comments
specifying arguments for validation functions
passing dependencies to validation functions
error handling
- default error handling
- custom error handling

basic optional/required rules
