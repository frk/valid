# List of Standard Library Preprocessor Rules

- [`lower`](#to-lower)
- [`upper`](#to-upper)
- [`title`](#to-title)
- [`trim`](#trim-white-space)
- [`ltrim`](#trim-left)
- [`rtrim`](#trim-right)
- [`trimprefix`](#trim-prefix)
- [`trimsuffix`](#trim-suffix)
- [`replace`](#replace)
- [`repeat`](#repeat)
- [`validutf8`](#to-valid-utf8)
- [`quote`](#quote)
- [`quoteascii`](#quote-to-ascii)
- [`quotegraphic`](#quote-to-graphic)
- [`urlqueryesc`](#url-query-escape)
- [`urlpathesc`](#url-path-escape)
- [`htmlesc`](#html-escape)
- [`htmlunesc`](#html-unescape)
- [`round`](#round)
- [`ceil`](#ceil)
- [`floor`](#floor)

## To Lower

The `lower` rule can be used to apply the [`strings.ToLower`](https://pkg.go.dev/strings@go1.19.1#ToLower)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"lower"`
	F2 *string `pre:"lower"`
}
```

</td><td>

```go
v.F1 = strings.ToLower(v.F1)
if v.F2 != nil {
	*v.F2 = strings.ToLower(*v.F2)
}
```

</td></tr>
</tbody></table>

## To Upper

The `upper` rule can be used to apply the [`strings.ToUpper`](https://pkg.go.dev/strings@go1.19.1#ToUpper)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"upper"`
	F2 *string `pre:"upper"`
}
```

</td><td>

```go
v.F1 = strings.ToUpper(v.F1)
if v.F2 != nil {
	*v.F2 = strings.ToUpper(*v.F2)
}
```

</td></tr>
</tbody></table>

## To Title

The `title` rule can be used to apply the [`strings.ToTitle`](https://pkg.go.dev/strings@go1.19.1#ToTitle)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"title"`
	F2 *string `pre:"title"`
}
```

</td><td>

```go
v.F1 = strings.ToTitle(v.F1)
if v.F2 != nil {
	*v.F2 = strings.ToTitle(*v.F2)
}
```

</td></tr>
</tbody></table>

## Trim White Space

The `trim` rule can be used to apply the [`strings.TrimSpace`](https://pkg.go.dev/strings@go1.19.1#TrimSpace)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"trim"`
	F2 *string `pre:"trim"`
}
```

</td><td>

```go
v.F1 = strings.TrimSpace(v.F1)
if v.F2 != nil {
	*v.F2 = strings.TrimSpace(*v.F2)
}
```

</td></tr>
</tbody></table>

## Trim Left

The `ltrim` rule can be used to apply the [`strings.TrimLeft`](https://pkg.go.dev/strings@go1.19.1#TrimLeft)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"ltrim:abc"`
	F2 *string `pre:"ltrim:xyz"`
}
```

</td><td>

```go
v.F1 = strings.TrimLeft(v.F1, "abc")
if v.F2 != nil {
	*v.F2 = strings.TrimLeft(*v.F2, "xyz")
}
```

</td></tr>
</tbody></table>

## Trim Right

The `rtrim` rule can be used to apply the [`strings.TrimRight`](https://pkg.go.dev/strings@go1.19.1#TrimRight)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"rtrim:abc"`
	F2 *string `pre:"rtrim:xyz"`
}
```

</td><td>

```go
v.F1 = strings.TrimRight(v.F1, "abc")
if v.F2 != nil {
	*v.F2 = strings.TrimRight(*v.F2, "xyz")
}
```

</td></tr>
</tbody></table>

## Trim Prefix

The `trimprefix` rule can be used to apply the [`strings.TrimPrefix`](https://pkg.go.dev/strings@go1.19.1#TrimPrefix)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"trimprefix:foo"`
	F2 *string `pre:"trimprefix:bar"`
}
```

</td><td>

```go
v.F1 = strings.TrimPrefix(v.F1, "foo")
if v.F2 != nil {
	*v.F2 = strings.TrimPrefix(*v.F2, "bar")
}
```

</td></tr>
</tbody></table>

## Trim Suffix

The `trimsuffix` rule can be used to apply the [`strings.TrimSuffix`](https://pkg.go.dev/strings@go1.19.1#TrimSuffix)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"trimsuffix:foo"`
	F2 *string `pre:"trimsuffix:bar"`
}
```

</td><td>

```go
v.F1 = strings.TrimSuffix(v.F1, "foo")
if v.F2 != nil {
	*v.F2 = strings.TrimSuffix(*v.F2, "bar")
}
```

</td></tr>
</tbody></table>

## Repeat

The `repeat` rule can be used to apply the [`strings.Repeat`](https://pkg.go.dev/strings@go1.19.1#Repeat)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"repeat:2"`
	F2 *string `pre:"repeat:5"`
}
```

</td><td>

```go
v.F1 = strings.Repeat(v.F1, 2)
if v.F2 != nil {
	*v.F2 = strings.Repeat(*v.F2, 5)
}
```

</td></tr>
</tbody></table>

## Replace

The `replace` rule can be used to apply the [`strings.Replace`](https://pkg.go.dev/strings@go1.19.1#Replace)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"replace:a:b:3"`
	F2 *string `pre:"replace:foo:bar"`
}
```

</td><td>

```go
v.F1 = strings.Replace(v.F1, "a", "b", 3)
if v.F2 != nil {
	*v.F2 = strings.Replace(*v.F2, "foo", "bar", -1)
}
```

</td></tr>
</tbody></table>

## To Valid UTF8

The `validutf8` rule can be used to apply the [`strings.ToValidUTF8`](https://pkg.go.dev/strings@go1.19.1#ToValidUTF8)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"validutf8"`
	F2 *string `pre:"validutf8"`
}
```

</td><td>

```go
v.F1 = strings.ToValidUTF8(v.F1)
if v.F2 != nil {
	*v.F2 = strings.ToValidUTF8(*v.F2)
}
```

</td></tr>
</tbody></table>

## Quote

The `quote` rule can be used to apply the [`strconv.Quote`](https://pkg.go.dev/strconv@go1.19.1#Quote)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"quote"`
	F2 *string `pre:"quote"`
}
```

</td><td>

```go
v.F1 = strconv.Quote(v.F1)
if v.F2 != nil {
	*v.F2 = strconv.Quote(*v.F2)
}
```

</td></tr>
</tbody></table>

## Quote To ASCII

The `quoteascii` rule can be used to apply the [`strconv.QuoteToASCII`](https://pkg.go.dev/strconv@go1.19.1#QuoteToASCII)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"quoteascii"`
	F2 *string `pre:"quoteascii"`
}
```

</td><td>

```go
v.F1 = strconv.QuoteToASCII(v.F1)
if v.F2 != nil {
	*v.F2 = strconv.QuoteToASCII(*v.F2)
}
```

</td></tr>
</tbody></table>

## Quote To Graphic

The `quotegraphic` rule can be used to apply the [`strconv.QuoteToGraphic`](https://pkg.go.dev/strconv@go1.19.1#QuoteToGraphic)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"quotegraphic"`
	F2 *string `pre:"quotegraphic"`
}
```

</td><td>

```go
v.F1 = strconv.QuoteToGraphic(v.F1)
if v.F2 != nil {
	*v.F2 = strconv.QuoteToGraphic(*v.F2)
}
```

</td></tr>
</tbody></table>

## URL Query Escape

The `urlqueryesc` rule can be used to apply the [`url.QueryEscape`](https://pkg.go.dev/net/url@go1.19.1#QueryEscape)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"urlqueryesc"`
	F2 *string `pre:"urlqueryesc"`
}
```

</td><td>

```go
v.F1 = url.QueryEscape(v.F1)
if v.F2 != nil {
	*v.F2 = url.QueryEscape(*v.F2)
}
```

</td></tr>
</tbody></table>

## URL Path Escape

The `urlpathesc` rule can be used to apply the [`url.PathEscape`](https://pkg.go.dev/net/url@go1.19.1#PathEscape)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"urlpathesc"`
	F2 *string `pre:"urlpathesc"`
}
```

</td><td>

```go
v.F1 = url.PathEscape(v.F1)
if v.F2 != nil {
	*v.F2 = url.PathEscape(*v.F2)
}
```

</td></tr>
</tbody></table>

## HTML Escape

The `htmlesc` rule can be used to apply the [`html.EscapeString`](https://pkg.go.dev/html@go1.19.1#EscapeString)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"htmlesc"`
	F2 *string `pre:"htmlesc"`
}
```

</td><td>

```go
v.F1 = html.EscapeString(v.F1)
if v.F2 != nil {
	*v.F2 = html.EscapeString(*v.F2)
}
```

</td></tr>
</tbody></table>

## HTML Unescape

The `htmlunesc` rule can be used to apply the [`html.UnescapeString`](https://pkg.go.dev/html@go1.19.1#UnescapeString)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `pre:"htmlunesc"`
	F2 *string `pre:"htmlunesc"`
}
```

</td><td>

```go
v.F1 = html.UnescapeString(v.F1)
if v.F2 != nil {
	*v.F2 = html.UnescapeString(*v.F2)
}
```

</td></tr>
</tbody></table>

## Round

The `round` rule can be used to apply the [`math.Round`](https://pkg.go.dev/math@go1.19.1#Round)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 float64  `pre:"round"`
	F2 *float64 `pre:"round"`
}
```

</td><td>

```go
v.F1 = math.Round(v.F1)
if v.F2 != nil {
	*v.F2 = math.Round(*v.F2)
}
```

</td></tr>
</tbody></table>

## Ceil

The `ceil` rule can be used to apply the [`math.Ceil`](https://pkg.go.dev/math@go1.19.1#Ceil)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 float64  `pre:"ceil"`
	F2 *float64 `pre:"ceil"`
}
```

</td><td>

```go
v.F1 = math.Ceil(v.F1)
if v.F2 != nil {
	*v.F2 = math.Ceil(*v.F2)
}
```

</td></tr>
</tbody></table>

## Floor

The `floor` rule can be used to apply the [`math.Floor`](https://pkg.go.dev/math@go1.19.1#Floor)
function to a field.

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 float64  `pre:"floor"`
	F2 *float64 `pre:"floor"`
}
```

</td><td>

```go
v.F1 = math.Floor(v.F1)
if v.F2 != nil {
	*v.F2 = math.Floor(*v.F2)
}
```

</td></tr>
</tbody></table>

