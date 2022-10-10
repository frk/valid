#### Placeholders

- *`object_identifier`* An object identifier is a string consisting of a Go
package's import path followed by a dot (".") and the name of a type or function
that is declared in that package.  
  ```yaml
  # examples of valid object identifiers
  - "strings.HasPrefix"
  - "github.com/me/mod/pkg.MyErrorAggregator"
  ```

- *`scalar`* Scalars are used to specify literal values that will be
used by the generated code as arguments to rule functions. A scalar value
can be one of the following types: `bool`, `string`, `float`, `int`.  

  As a special case one can also use the YAML tag `!!nil` in place of a
  scalar value to indicate that the generated code should pass the predeclared
  identifier `nil` to a rule function.
  
  ```yaml
  # examples of valid scalars
  - true
  - "foo bar"
  - 0.5
  - 3
  - !!nil
  ```

---

#### config_file

```yaml
# The directory in which the tool will search for files to process.
# If not provided, the current working directory will be used by default.
#
# CLI flag: -wd
[working_directory: <string> | default = "."]

# When set to true, the tool will search the hierarchy
# of the working directory for files to process.
#
# CLI flag: -r
[recursive: <bool> | default = false]

# A list of names of specific files that the tool should process.
# The files MUST be located in the working directory's hierarchy.
#
# CLI flag: -f
file_list:
  [- <string> ...]

# A list of regular expressions to match the names of files that
# the tool should process. The regular expressions MUST match files
# that are located in the working directory's hierarchy.
#
# CLI flag: -rx
file_pattern_list:
  [- <string> ...]

# A string containing the format to be used for generating
# the name of the output files.
#
# Inside the format string the percent sign ("%") can be used as
# the placeholder to be replaced by the tool with the base name of
# the source file. If no placeholder is present in the format string
# then the base name of the source file will simply be prefixed
# to the format string.
#
# CLI flag: -o
[out_name_format: <string> | default = "%_valid.go"]

# A string containing a regular expression that will be used by the tool
# to identify the struct types for which to generate the validation code.
[validator_name_pattern: <string> | default = "(?i:validator)$"]

# Configures the code generation of the handling of validation errors.
[error_handling: <error_handling>]

# List of custom rules to be made available to the tool.
rules:
  [- <rule_config> ...]
```

---

#### error_handling

```yaml
# Configures how field keys should be constructed.
#
# Field keys are used for error reporting by the generated code. When a field
# fails validation the field's key, along with other details, will be passed as
# an argument to the client's implementation of the error handling code.
field_key:
  # If non-empty string, specifies the struct tag whose value will be
  # used for constructing the field keys. If explicitly set to an empty
  # string, the generator will default to use the fields' names for
  # constructing the field keys.
  #
  # A valid tag must begin with a letter (A-z) or an underscore (_),
  # subsequent characters in the tag can be letters, underscores,
  # and digits (0-9).
  #
  # CLI flag: -fk.tag
  [tag: <string> | default = "json"]

  # When set to true, a nested struct field's key will be constructed by
  # joining it together with all of its parent fields. When false, a nested
  # struct field's key will be constructed only from that field's tag or name.
  #
  # CLI flag: -fk.join
  [join: <bool> | default = true]

  # The separator to be used for joining fields' tags or names
  # when constructing field keys.
  #
  # The separator MUST be a single, one byte long, character.
  #
  # CLI flag: -fk.sep
  [separator: <string> | default = "."]

# The identifier of a function that the generated code should
# use for constructing custom, application-specific errors.
# The function's signature MUST be the following:
#
#     func(key string, val any, rule string, args ...any) error
#
# CLI flag: -error.constructor
[constructor: <object_identifier>]

# The identifier of a type that the generated code should use
# for constructing & aggregation of custom, application-specific
# errors. The type MUST implement the following interface:
#
#     interface {
#         Error(key string, val any, rule string, args ...any)
#         Out() error
#     }
#
# CLI flag: -error.aggregator
[aggregator: <object_identifier>]
```

---

#### rule_config

```yaml
# The function associated with the rule.
func: <object_identifier>

# The rule specification.
rule:
  name: <string>
  # The configuration for the rule's arguments.
  #
  # When args is NOT empty, then it MUST contain a number of elements
  # that is compatible with the associated function's number of parameters.
  args:
    [- <rule_arg_config> ...]

  # The minimum number of arguments that the rule MUST accept.
  [arg_min: <int>]

  # This can be used to override the upper limit of the number of arguments
  # that the rule should be allowed to accept. A negative value can be used
  # to indicate that there's no upper limit to the number of arguments.
  #
  # The value MUST be compatible with the rule's associated function's
  # signature. If not, then the tool will exit with an error.
  [arg_max: <int>]

  # Configures how the error message for the rule should be constructed.
  error:
    # The main text of the error message.
    [text: <string>]

    # When true, the generated error message will include the rule's arguments.
    [with_args: <bool>]

    # The separator that should be used to join the rule's arguments in the error message.
    [arg_sep: <string>]

    # The text to be appended after the list of arguments in the error message.
    [arg_suffix: <string>]

  # The logical operator that should be used to join multiple
  # calls to the associated function into a single expression.
  #
  # The value MUST be one of: "AND", "OR", or "NOT" (case insensitive).
  [join_op: <string>]
```

---

#### rule_arg_config

```yaml
# The rule argument's default value. If omitted, then the
# rule argument's value MUST be provided inside the struct tag.
[default: <scalar>]

# When options is empty, then ANY value can be provided for
# the argument inside the rule's struct tag.
#
# When options is NOT empty, then it is considered to represent, together
# with the default value, the COMPLETE set of valid values that can be
# provided as the argument inside the rule's struct tag.
options:
  [- <rule_arg_option> ...]
```

---

#### rule_arg_option

```yaml
# Specifies the value that the generator should supply
# as the rule's argument in the generated code.
value: <scalar>

# Alias is an alternative identifier of the argument's
# value that can be used within the rule's struct tag.
[alias: <string>]
```
