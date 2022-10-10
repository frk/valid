[`pkg.go.dev`](https://pkg.go.dev/github.com/frk/valid)

# valid

The package `valid` defines a number of validation functions and other objects
that can be used, together with the [`cmd/validgen`](cmd/validgen) tool, to
generate struct field validation.

Most of validation logic in the `./valid.go` file, including the related tests and
comments, was ported over from https://github.com/validatorjs/validator.js

--------

# Table of Contents

- [An introductory example](#an-introductory-example)
- [The configuration of `validgen`](#the-configuration-of-validgen)
- [The generator rules](#the-generator-rules)
- TODO adding custom validation functions
- TODO configuring custom validation functions from file
- TODO configuring custom validation functions using comments
- TODO specifying arguments for validation functions
- TODO passing dependencies to validation functions
- TODO error handling
- TODO default error handling
- TODO custom error handling

## An Introductory Example

<table><thead><tr><th>Input</th></tr></thead><tbody>
<tr><td>

```go
type UserCreateParams struct {
	FName string `is:"len:1:300" pre:"trim"`
	LName string `is:"len:1:300,required" pre:"trim"`
	Email string `is:"email,required" pre:"lower,trim"`
	Passw string `is:"strongpass,required" pre:"trim"`
	Age   int    `is:"min:3,max:150"`
}

type UserCreateParamsValidator struct {
	UserCreateParams
}
```

</td></tr>
</tbody></table>

<table><thead><tr><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
func (v UserCreateParamsValidator) Validate() error {
	v.FName = strings.TrimSpace(v.FName)
	if len(v.FName) < 1 || len(v.FName) > 300 {
		return errors.New("FName must be of length between: 1 and 300 (inclusive)")
	}
	v.LName = strings.TrimSpace(v.LName)
	if v.LName == "" {
		return errors.New("LName is required")
	} else if len(v.LName) < 1 || len(v.LName) > 300 {
		return errors.New("LName must be of length between: 1 and 300 (inclusive)")
	}
	v.Email = strings.TrimSpace(strings.ToLower(v.Email))
	if v.Email == "" {
		return errors.New("Email is required")
	} else if !valid.Email(v.Email) {
		return errors.New("Email must be a valid email address")
	}
	v.Passw = strings.TrimSpace(v.Passw)
	if v.Passw == "" {
		return errors.New("Passw is required")
	} else if !valid.StrongPassword(v.Passw, nil) {
		return errors.New("Passw must be a strong password")
	}
	if v.Age < 3 {
		return errors.New("Age must be greater than or equal to: 3")
	} else if v.Age > 150 {
		return errors.New("Age must be less than or equal to: 150")
	}
	return nil
}
```

</td></tr>
</tbody></table>

## The Configuration of Validgen

The `validgen` tool can be configured in two ways. Either by using CLI arguments
(this is a limited approach), or with a YAML config file (this is the recommended
and complete approach). To see the documentation for the CLI arguments, you can run:

```sh
validgen --help
```

To configure the tool using a specific file you can provide the `-c file` argument, e.g.

```sh
validgen -c /path/to/config.yaml
```

When the `-c file` argument is omitted, the tool will look in the nearest git-root
directory of the current working directory for a file named `.valid.yaml`. If such
a file is found the tool will use that to configure itself. The complete documentation
of the config yaml file can be found [here](./doc/configuration.md).

## The Generator Rules

The `validgen` tool looks for particular struct tags that are then used as the instructions
for what code the tool should generate. These instructions are referred to as *rules* and
there are two distinct kinds:

- *validation rules* (denoted with `is:"..."` structs tags) are used to generate
struct field validation.
- *preprocessor rules* (denoted with `pre:"..."` struct tags) are used to generate
struct field "pre-processing".

Based on how these rules are implmenetated, they can be classified into the following categories:

1. *"builtin" validation rules*: These rules are implemented using the Go langauge's
primitive operators and builtin functionality. For a full list (with examples) of the
builtin validation rules, see: [builtin validation rules](./doc/list_of_builtin_validation_rules.md).
2. *"stdlib" validation rules*: These rules are implemented using functions of the
Go standard library. For a full list (with examples) of the stdlib validation rules,
see: [stdlib validation rules](./doc/list_of_stdlib_validation_rules.md).
3. *"included" validation rules*: These rules are implemented using functions from
the `github.com/frk/valid` package. For a full list (with examples) of the included
validation rules, see: [included validation rules](./doc/list_of_included_validation_rules.md).
4. *"custom" validation rules*: These rules are implemented with functions that are
sourced from the configuration file's `"rules"` entry.
5. *"stdlib" preprocessor rules*: These rules are implemented using functions of the
Go standard library. For a full list (with examples) of available included validation
rules, see: [stdlib preprocessor rules](./doc/list_of_stdlib_preprocessor_rules.md).
6. *"custom" preprocessor rules*: These rules are implemented with functions that are
sourced from the configuration file's `"rules"` entry.
