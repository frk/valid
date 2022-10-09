# List of Standard Library Validation Rules

- [`runecount`](#has-rune-count) (has number of characters)
- [`contains`](#contains) (has substring)
- [`prefix`](#has-prefix) (has prefix)
- [`suffix`](#has-suffix) (has suffix)

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


## Has Substring

The `contains` rule can be used to check whether a string *contains* a substring or not.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"contains:foo"`
	F2 *string `is:"contains:bar"`
	F3 string  `is:"contains:foo:bar:baz"`
}
```

</td><td>

```go
if !strings.Contains(v.F1, "foo") {
	return errors.New("...")
}
if v.F2 != nil && !strings.Contains(*v.F2, "bar") {
	return errors.New("...")
}
if !strings.Contains(v.F3, "foo") && !strings.Contains(v.F3, "bar") && !strings.Contains(v.F3, "baz") {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Has Prefix

The `prefix` rule can be used to check whether a string *beings with* a substring or not.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"prefix:foo"`
	F2 *string `is:"prefix:bar"`
	F3 string  `is:"prefix:foo:bar:baz"`
}
```

</td><td>

```go
if !strings.HasPrefix(v.F1, "foo") {
	return errors.New("...")
}
if v.F2 != nil && !strings.HasPrefix(*v.F2, "bar") {
	return errors.New("...")
}
if !strings.HasPrefix(v.F3, "foo") && !strings.HasPrefix(v.F3, "bar") && !strings.HasPrefix(v.F3, "baz") {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## Has Suffix

The `suffix` rule can be used to check whether a string *ends with* a substring or not.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"suffix:foo"`
	F2 *string `is:"suffix:bar"`
	F3 string  `is:"suffix:foo:bar:baz"`
}
```

</td><td>

```go
if !strings.HasSuffix(v.F1, "foo") {
	return errors.New("...")
}
if v.F2 != nil && !strings.HasSuffix(*v.F2, "bar") {
	return errors.New("...")
}
if !strings.HasSuffix(v.F3, "foo") && !strings.HasSuffix(v.F3, "bar") && !strings.HasSuffix(v.F3, "baz") {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>
