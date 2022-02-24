**The types in this document are here for the purposes of documentation only.**

--------------------------------------------------------------------------------

#### The `Validator` Types

When the tool runs it searches the source for types for which to
generate validation code. These types are referred to as *validator*
types. For a type to be considered a validator type, it must be a struct
type, declared at the top-level, whose name matches a specific pre-defined
pattern. By default this pattern is set to `"^(?i:\w*Validator)$"`,
i.e. any string that ends with `Validator` *(case insensitive)*. E.g.,

```go
type UserInputValidator struct {
    // ... fields to validate ...
}
```

To change the default name pattern use the
`validator_name_pattern` config entry. E.g.,

```yaml
validator_name_pattern: "^[A-z]+Input$"
```
```go
type UserInput struct {
    // ... fields to validate ...
}
```

Types that happen to have the prerequiste validator attributes, but
which should NOT have validation code generated for them, can use
the `valid:ignore` directive in their comments to instruct the tool
to ignore them. E.g,

```go
// valid:ignore
type FooValidator struct {
    // ...
}
```

---

#### The `IsValid` Type

```go
type IsValid interface {
    IsValid() bool
}
```

Types that should validate themselves can do so by implementing
the `IsValid` interface. When the tool encounters a field with
such a type it will *automatically* generate an expression that
invokes the `IsValid()` method on that field. E.g.,

```go
if !f.IsValid() {
    // ...
}
```

Fields with such a type, for which the tool should NOT generate
the `IsValid()` method call, can use the `-isvalid` rule in their
struct tags to instruct the tool to not generate that expression. E.g.,

```go
Field myType `is:"-isvalid"`
```

---

#### The `BeforeValidator` and `AfterValidator` Types

Validator types that need to execute some code right before and/or
right after validation can do so by implementing the `BeforeValidator`
and/or `AfterValidator` interfaces shown below.

```go
type BeforeValidator interface {
    BeforeValidate() error
}

type AfterValidator interface {
    AfterValidate() error
}
```

---

#### The `ErrorConstructorFunc` and `ErrorConstructor` Types

By default the tool will generate code that constructs errors using
the stdlib packages `errors` and `fmt`. To instead generate code that
returns custom, app-specific errors the `ErrorConstructorFunc` type
or the `ErrorConstructor` interface can be implemented.

```go
// The implementation should construct a new custom
// error value from the given arguments and return it.
type ErrorConstructorFunc func(key string, val any, rule string, args ...any) error

type ErrorConstructor interface {
	// The implementation should construct a new custom
	// error value from the given arguments and return it.
	Error(key string, val any, rule string, args ...any) error
}
```

1. To utilize the `ErrorConstructorFunc` type, set the config file's
`error_handling.constructor` entry to a custom function whose
signature is *identical* to that of the `ErrorConstructorFunc`.

   ```yaml
   error_handling.constructor: "example.com/me/mymod/mypkg.MyErrorConstructorFunc"
   ```

2. To utilize the `ErrorConstructor` interface, include an implementation
of that interface as a field in the validator struct:

   ```go
   type UserInputValidator struct {
       // ... fields to validate ...
       MyErr mypkg.MyErrorConstructorType
   }
   ```

---

#### The `ErrorAggregator` Type

By default the generated code will exit and return an error immediately
upon encountering the first validation fail. To instead have the generated
code validate every field and aggregate and return all the encountered
errors the `ErrorAggregator` interface can be implemented.

```go
type ErrorAggregator interface {
	// The implementation of the Error method should construct a
	// new custom error value from the given arguments and retain
	// it for until the generated validation code is done.
	Error(key string, val any, rule string, args ...any)
	// The Out method will be invoked by the generated code at
	// the end to yield the error value it returns.
	Out() error
}
```

There are two ways to utilize the `ErrorAggregator` interface:

1. In the config file set the `error_handling.aggregator` entry to the
custom type that implmenets the `ErrorAggregator` interface. Note that
the *zero value* of that custom type MUST be ready to use.

   ```yaml
   error_handling.aggregator: "example.com/me/mymod/mypkg.MyErrorAggregatorType"
   ```

2. In the validator struct's definition include a field whose type
implements the `ErrorAggregator` interface:

   ```go
   type UserInputValidator struct {
       // ... fields to validate ...
       MyErr mypkg.MyErrorAggregatorType
   }
   ```

---

#### The error parameters

The `ErrorConstructorFunc`, the `ErrorConstructor.Error` method, and
the `ErrorAggregator.Error` method, all have the same input parameters.
The following is a description of these parameters:

0. `key` (type `string`) will be the key of the field whose value failed validation.
1. `val` (type `any`) will hold the value that failed validation.
2. `rule` (type `string`) will be the name of the rule for which the value failed validation.
3. `args` (type `[]any`) will hold the list of the rule's arguments specificed in the `is` tag.

