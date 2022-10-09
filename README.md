[`godoc`](http://godoc.org/github.com/frk/valid)
[`pkg.go.dev`](https://pkg.go.dev/github.com/frk/valid)

# valid

The package `valid` defines a number of validation functions and other objects
that can be used, together with the [`cmd/validgen`](cmd/validgen) tool, to
generate static input validation.

Most of validation logic in the `./valid.go` file, including the related tests and
comments, was ported over from https://github.com/validatorjs/validator.js

--------

## Table of Contents

- TODO `validgen` configuration
- TODO adding custom validation functions
- TODO configuring custom validation functions from file
- TODO configuring custom validation functions using comments
- TODO specifying arguments for validation functions
- TODO passing dependencies to validation functions
- TODO error handling
- TODO default error handling
- TODO custom error handling
- List of builtin Preprocessor Rules
	- [`repeat`](#repeat)
	- [`replace`](#replace)
	- [`lower`](#to-lower)
	- [`upper`](#to-upper)
	- [`title`](#to-title)
	- [`validutf8`](#to-valid-utf8)
	- [`trim`](#trim-space)
	- [`ltrim`](#trim-left)
	- [`rtrim`](#trim-right)
	- [`trimprefix`](#trim-prefix)
	- [`trimsuffix`](#trim-suffix)
	- [`quote`](#quote)
	- [`quoteascii`](#quote-ascii)
	- [`quotegraphic`](#quote-to-graphic)
	- [`urlqueryesc`](#url-query-escape)
	- [`urlpathesc`](#url-path-escape)
	- [`htmlesc`](#html-escape)
	- [`htmlunesc`](#html-unescape)
	- [`round`](#round)
