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
- List of builtin Validation Rules
	- [`required`](#required)
	- [`notnil`](#not-nil)
	- [`optional`](#optional)
	- [`omitnil`](#omit-nil)
	- [`noguard`](#no-guard)
	- [`re`](#regular-expression) (Regular Expression)
	- [`eq`](#is-equal-to)
	- [`ne`](#is-not-equal-to)
	- [`gt`](#is-greater-than)
	- [`lt`](#is-less-than)
	- [`gte`](#is-greater-than-or-equal-to)
	- [`lte`](#is-less-than-or-equal-to)
	- [`min`](#min)
	- [`max`](#max)
	- [`rng`](#is-in-range)
	- [`len`](#has-length)
	- [`runecount`](#has-rune-count)
	- [`enum`](#enum)
	- [`isvalid`](#isvalid-interface)
	- [`-isvalid`](#omit-isvalid-interface)
	- [`prefix`](#has-prefix)
	- [`suffix`](#has-suffix)
	- [`contains`](#contains)
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

## Required

The `required` rule can be used to check that a field's *base* value is non-[zero](https://go.dev/ref/spec#The_zero_value).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string    `is:"required"`
	F2 float64   `is:"required"`
	F3 any       `is:"required"`
	F4 time.Time `is:"required"`
}
```

</td><td>

```go
if v.F1 == "" {
	return errors.New("...")
}
if v.F2 == 0.0 {
	return errors.New("...")
}
if v.F3 == nil {
	return errors.New("...")
}
if v.F4 == (time.Time{}) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

When the field's type is a pointer type, then the generated code first checks
the field against `nil` and then uses indirection to check the base value.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 *string    `is:"required"`
	F2 **float64  `is:"required"`
	F3 *any       `is:"required"`
	F4 *time.Time `is:"required"`
}
```

</td><td>

```go
if v.F1 == nil || *v.F1 == "" {
	return errors.New("...")
}
if v.F2 == nil || *v.F2 == nil || **v.F2 == 0.0 {
	return errors.New("...")
}
if v.F3 == nil || *v.F3 == nil {
	return errors.New("...")
}
if v.F4 == nil || *v.F4 == (time.Time{}) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

When the field's type is a map or a slice, the generated code will check the field's
length against `0` rather than the field's value against `nil`.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	S []string           `is:"required"`
	M *map[string]string `is:"required"`
}
```

</td><td>

```go
if len(v.S) == 0 {
	return errors.New("...")
}
if v.M == nil || len(*v.M) == 0 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Not Nil

The `notnil` rule can be used to check a field's value against `nil`. The field must be of
a type whose [zero value](https://go.dev/ref/spec#The_zero_value) is `nil` (i.e. pointer,
function, interface, slice, channel, or map), otherwise the tool will exit with an error.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 *string `is:"required"`
	F2 []int   `is:"required"`
	F3 any     `is:"required"`
}
```

</td><td>

```go
if v.F1 == nil {
	return errors.New("...")
}
if v.F2 == nil {
	return errors.New("...")
}
if v.F3 == nil {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Optional

The `optional` rule can be used to validate a field ONLY if its base value IS NOT
the [zero value](https://go.dev/ref/spec#The_zero_value).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string `is:"email,optional"`
	F2 *int64 `is:"eq:42,optional"`
}
```

</td><td>

```go
if v.F1 != "" && !valid.Email(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && *v.F2 > 0 && *v.F2 != 42 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Omit Nil

The `omitnil` rule can be used to validate a field ONLY if its value IS NOT `nil`.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 *string `is:"email,omitnil"`
	F2 []int64 `is:"len:5,omitnil"`
}
```

</td><td>

```go
if v.F1 != nil && !valid.Email(*v.F1) {
	return errors.New("...")
}
if v.F2 != nil && len(v.F2) != 5 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## No Guard

The `noguard` rule can be used to omit the by-default-generated nil-pointer checks
that guard against nil-pointer-dereference issues. Use this rule ONLY when you are
certain that the caller properly initialized the pointer field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 *string   `is:"email,noguard"`
	F2 **[]int64 `is:"len:5,noguard"`
}
```

</td><td>

```go
if !valid.Email(*v.F1) {
	return errors.New("...")
}
if len(**v.F2) != 5 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Regular Expression

The `re` rule can be used to validate a field with a regular expression.

Note; the generated code will include a top-level `init()` func that registers
(and compiles) the regular expression(s) with [`valid.RegisterRegexp`](https://pkg.go.dev/github.com/frk/valid#RegisterRegexp)
and then the validation is done using the [`valid.Match`](https://pkg.go.dev/github.com/frk/valid#Match)
function, which, under the hood, uses the regexp string as the map-key to retrieve the
registered (and compiled) regular expression and then invokes its [`MatchString`](https://pkg.go.dev/regexp@go1.19.1#Regexp.MatchString)
method to do the validation.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"re:foo"`
	F2 string  `is:"re:\"\\w+\""`
	F3 *string `is:"re:\"^[a-z]+\\[[0-9]+\\]$\""`
}
```

</td><td>

```go
func init() {
	valid.RegisterRegexp(`foo`)
	valid.RegisterRegexp(`\w+`)
	valid.RegisterRegexp(`^[a-z]+\[[0-9]+\]$`)
}

// ...

if !valid.Match(v.F1, `foo`) {
	return errors.New("...")
}
if !valid.Match(v.F2, `\w+`) {
	return errors.New("...")
}
if v.F3 != nil && !valid.Match(*v.F3, `^[a-z]+\[[0-9]+\]$`) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Is Equal To

The `eq` rule ensures that a field's value is equal to one of the rule's arguments.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string `is:"eq:foo:bar:baz"`
	F2 *int   `is:"eq:64:128"`
	F3 any    `is:"eq:foo:0.8:true"`
}
```

</td><td>

```go
if v.F1 != "foo" && v.F1 != "bar" && v.F1 != "baz" {
	return errors.New("...")
}
if v.F2 != nil && (*v.F2 != 64 && *v.F2 != 128) {
	return errors.New("...")
}
if v.F3 != "foo" && v.F3 != 0.8 && v.F3 != true {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Is Not Equal To

The `ne` rule ensures that a field's value is not equal to any of the rule's arguments.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string `is:"ne:foo:bar"`
	F2 *int   `is:"ne:64:128"`
	F3 any    `is:"ne:foo:0.8:true"`
}
```

</td><td>

```go
if v.F1 == "foo" || v.F1 == "bar" {
	return errors.New("...")
}
if v.F2 != nil && (*v.F2 == 64 || *v.F2 == 128) {
	return errors.New("...")
}
if v.F3 == "foo" || v.F3 == 0.8 || v.F3 == true {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Is Greater Than

The `gt` rule ensures that a field's numeric value is greater than the rule's argument.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 float64 `is:"gt:3.14"`
	F2 *int    `is:"gt:314"`
	F3 *uint8  `is:"gt:31,required"`
}
```

</td><td>

```go
if v.F1 <= 3.14 {
	return errors.New("...")
}
if v.F2 != nil && *v.F2 <= 314 {
	return errors.New("...")
}
if v.F3 == nil || *v.F3 == 0 {
	return errors.New("...")
} else if *v.F3 <= 31 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Is Less Than

The `lt` rule ensures that a field's numeric value is less than the rule's argument.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 float64 `is:"lt:3.14"`
	F2 *int    `is:"lt:314"`
	F3 *uint8  `is:"lt:31,required"`
}
```

</td><td>

```go
if v.F1 >= 3.14 {
	return errors.New("...")
}
if v.F2 != nil && *v.F2 >= 314 {
	return errors.New("...")
}
if v.F3 == nil || *v.F3 == 0 {
	return errors.New("...")
} else if *v.F3 >= 31 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Is Greater Than or Equal To

The `gte` rule ensures that a field's numeric value is greater than or equal to the rule's argument.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 float64 `is:"gte:3.14"`
	F2 *int    `is:"gte:314"`
	F3 *uint8  `is:"gte:31,required"`
}
```

</td><td>

```go
if v.F1 < 3.14 {
	return errors.New("...")
}
if v.F2 != nil && *v.F2 < 314 {
	return errors.New("...")
}
if v.F3 == nil || *v.F3 == 0 {
	return errors.New("...")
} else if *v.F3 < 31 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Min

The `min` rule is an alias for `gte`.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 float64 `is:"min:3.14"`
	F2 *int    `is:"min:314"`
	F3 *uint8  `is:"min:31,required"`
}
```

</td><td>

```go
if v.F1 < 3.14 {
	return errors.New("...")
}
if v.F2 != nil && *v.F2 < 314 {
	return errors.New("...")
}
if v.F3 == nil || *v.F3 == 0 {
	return errors.New("...")
} else if *v.F3 < 31 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Is Less Than or Equal To

The `lte` rule ensures that a field's numeric value is less than or equal to the rule's argument.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 float64 `is:"lte:3.14"`
	F2 *int    `is:"lte:314"`
	F3 *uint8  `is:"lte:31,required"`
}
```

</td><td>

```go
if v.F1 > 3.14 {
	return errors.New("...")
}
if v.F2 != nil && *v.F2 > 314 {
	return errors.New("...")
}
if v.F3 == nil || *v.F3 == 0 {
	return errors.New("...")
} else if *v.F3 > 31 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Max

The `max` rule is an alias for `lte`.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 float64 `is:"max:3.14"`
	F2 *int    `is:"max:314"`
	F3 *uint8  `is:"max:31,required"`
}
```

</td><td>

```go
if v.F1 > 3.14 {
	return errors.New("...")
}
if v.F2 != nil && *v.F2 > 314 {
	return errors.New("...")
}
if v.F3 == nil || *v.F3 == 0 {
	return errors.New("...")
} else if *v.F3 > 31 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Is In Range

The `rng` rule ensures that a field's numeric value is between its two arguments (inclusive).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 float64 `is:"rng:3.14:42"`
	F2 *int    `is:"rng:-8:256"`
	F3 *uint8  `is:"rng:1:2,required"`
}
```

</td><td>

```go
if v.F1 < 3.14 || v.F1 > 42 {
	return errors.New("...")
}
if v.F2 != nil && (*v.F2 < -8 || *v.F2 > 256) {
	return errors.New("...")
}
if v.F3 == nil || *v.F3 == 0 {
	return errors.New("...")
} else if *v.F3 < 1 || *v.F3 > 2 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Has Length

The `len` rule checks a field value's length. This rule takes either one integer
argument, in which case it checks the exact length of the field's value, or, it
takes two integer arguments, in which case it will check the field value's length
against min/max bounds. This rule must be used only with fields whose types have
a length. (e.g. string, slice, map, etc.)

- Length must be **equal to X**:

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F string `is:"len:8"`
}
```

</td><td>

```go
if len(v.F) != 8 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

- Length must be **between X and Y (inclusive)**:

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F *string `is:"len:8:16"`
}
```

</td><td>

```go
if v.F != nil && (len(*v.F) < 8 || len(*v.F) > 16) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

- Length must be **at least X**:

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F []byte `is:"len:8:"`
}
```

</td><td>

```go
if len(v.F) < 8 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

- Length can be **at most Y**:

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F *[]byte `is:"len::16"`
}
```

</td><td>

```go
if v.F != nil && (len(*v.F) > 16) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## Rune Count

The `runecount` rule checks a field value's `rune` count. This rule takes either
one integer argument, in which case it checks the exact count of the field's runes,
or, it takes two integer arguments, in which case it will check the field's rune
count against min/max bounds. This rule can be used only with fields of type `string`
or type `[]byte`.

- Rune count must be **equal to X**:

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F string `is:"runecount:8"`
}
```

</td><td>

```go
if utf8.RuneCountInString(v.F) != 8 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

- Rune count must be **between X and Y (inclusive)**:

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F *string `is:"runecount:8:16"`
}
```

</td><td>

```go
if v.F != nil && (utf8.RuneCountInString(*v.F) < 8 || utf8.RuneCountInString(*v.F) > 16) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

- Rune count must be **at least X**:

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F []byte `is:"runecount:8:"`
}
```

</td><td>

```go
if utf8.RuneCount(v.F) < 8 {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

- Rune count can be **at most Y**:

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F *[]byte `is:"runecount::16"`
}
```

</td><td>

```go
if v.F != nil && (utf8.RuneCount(*v.F) > 16) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Enum

The `enum` rule checks that a field's value matches one of the constants declared
with the field's type.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type mytype string

const (
	foo mytype = "foo"
	bar mytype = "bar"
	baz mytype = "baz"
)

type Validator struct {
	F mytype `is:"enum"`
}
```

</td><td>

```go
if v.F != foo && v.F != bar && v.F != baz {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

The `enum` rule can also be used with types declared in imported packages.

<table><tbody>
<tr><th>Rule Tag</th></tr>
<tr><td>

```go
import "go/ast"

type Validator struct {
	F ast.ObjKind `is:"enum"`
}
```

</td></tr>

<tr><th>Generated Output</th></tr>
<tr><td>

```go
if v.F != ast.Bad && v.F != ast.Pkg && v.F != ast.Con && v.F != ast.Typ && v.F != ast.Var && v.F != ast.Fun && v.F != ast.Lbl {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## The IsValid interface

The `isvalid` rule is applied *automatically* to any field whose type implements `interface{ IsValid() bool }`.
The `isvalid` rule validates a field by invoking the `IsValid() bool` method on the field. If a field's type
implements the interface but you wish not to use the `IsValid() bool` method for validation,
then the `-isvalid` rule can be used.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type mytype string

func (mytype) IsValid() bool {
	// ...
	return true
}

type Validator struct {
	F1 mytype
	F2 mytype `is:"isvalid"`
	F3 mytype `is:"-isvalid"`
}
```

</td><td>

```go
if !v.F1.IsValid() {
	return errors.New("...")
}
if !v.F2.IsValid() {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>
