# List of Included Validation Rules

- [`re`](#match-regular-expression): match regular expression
- [`ascii`](#is-ascii-string): is ASCII string
- [`alpha`](#is-alphabetic-string): is alphabetic string
- [`alnum`](#is-alphanumeric-string): is alphanumeric string
- [`bic`](#is-bank-identification-code): is bank identification code
- [`btc`](#is-bitcoin-address): is bitcoin address
- [`base32`](#is-base-32-string): is base32 string
- [`base58`](#is-base-58-string): is base58 string
- [`base64`](#is-base-64-string): is base64 string
- [`binary`](#is-binary-integer-string): is binary integer string
- [`bool`](#is-boolean-value): is boolean value
- [`cidr`](#is-classless-inter-domain-routing-notation): is classless inter-domain routing notation
- [`cvv`](#is-card-verification-value): is card verification value
- [`ccy`](#is-currency-amount): is currency amount
- [`datauri`](#is-data-uri): is data URI
- [`decimal`](#is-decimal-number): is decimal number
- [`digits`](#is-string-of-digits): is string of digits
- [`ean`](#is-european-article-number): is european article number
- [`ein`](#is-employer-identification-number): is employer identification number
- [`eth`](#is-ethereum-address): is ethereum address
- [`email`](#is-email-address): is email address
- [`fqdn`](#is-fully-qualified-domain-name): is fully qualified domain name
- [`float`](#is-floating-point-number): is floating point number
- [`hsl`](#is-hsl-color): is HSL color
- [`hash`](#is-hash-of-algorithm): is hash of algorithm
- [`hex`](#is-hexadecimal-string): is hexadecimal string
- [`hexcolor`](#is-hexadecimal-color-code): is hexadecimal color code
- [`iban`](#is-international-bank-account-number): is international bank account number
- `ic [TODO]`: is identification card
- [`imei`](#is-international-mobile-equipment-identity-number): is international mobile equipment identity number
- [`ip`](#is-internet-protocol-address): is internet protocol address
- [`iprange`](#is-internet-protocol-address-range): is internet protocol address range
- [`isbn`](#is-international-standard-book-number): is international standard book number
- [`isin`](#is-international-securities-identification-number): is international securities identification number
- [`iso639`](#is-iso-639-string): is ISO 639 string
- [`iso31661a`](#is-iso-3166-1a-string): is ISO 3166-1A string
- [`iso4217`](#is-iso-4217-string): is ISO 4217 string
- [`isrc`](#is-international-standard-recording-code): is international standard recording code
- [`issn`](#is-international-standard-serial-number): is international standard serial number
- [`in`](#is-in): is in list / is one of
- [`int`](#is-integer-number): is integer number
- [`json`](#is-json-value): is JSON value
- [`jwt`](#is-json-web-token): is JSON web token
- [`latlong`](#is-latitude-longitude-string): is latitude longitude string
- [`locale`](#is-locale-code): is locale code
- [`lower`](#is-lower-case-string): is lower case string
- [`mac`](#is-mac-address): is MAC address
- [`md5`](#is-md5-hash): is MD5 hash
- [`mime`](#is-mime-type): is MIME type
- [`magneturi`](#is-magnet-uri): is magnet URI
- [`mongoid`](#is-mongo-id): is mongo ID
- [`numeric`](#is-numeric-string): is numeric string
- [`octal`](#is-octal-number): is octal number
- [`pan`](#is-primary-account-number): is primary account number
- `passport [TODO]`: is passport number
- [`phone`](#is-phone-number): is phone number
- [`port`](#is-port-number): is port number
- [`rgb`](#is-rgb-color): is RGB color
- [`ssn`](#is-social-security-number): is social security number
- [`semver`](#is-semantic-version-number): is semantic version number
- [`slug`](#is-slug): is slug
- [`strongpass`](#is-strong-password): is strong password
- `url [TODO]`: is uniform resource location
- [`uuid`](#is-universally-unique-identification-number): is universally unique identification number
- [`uint`](#is-unsigned-integer-number): is unsigned integer number
- [`upper`](#is-upper-case-string): is upper case string
- [`vat`](#is-value-added-tax-number): is value added tax number
- [`zip`](#is-zip-code): is zip code / is postal code

## match regular expression

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

## is ASCII string

The `ascii` rule can be used to check if a field's value is a valid ASCII string.

The validation is implemented by [`valid.ASCII`](https://pkg.go.dev/github.com/frk/valid#ASCII).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"ascii"`
	F2 *string `is:"ascii"`
}
```

</td><td>

```go
if !valid.ASCII(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.ASCII(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is alphabetic string

The `alpha[:lang]` rule can be used to check if a field's value is a valid alphabetic string.

The optional `lang` argument can be used to specify the alphabet's language. When not specified
the `lang` argument will default to `"en"`.

The validation is implemented by [`valid.Alpha`](https://pkg.go.dev/github.com/frk/valid#Alpha).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"alpha"`
	F2 *string `is:"alpha:ja"`
}
```

</td><td>

```go
if !valid.Alpha(v.F1, "en") {
	return errors.New("...")
}
if v.F2 != nil && !valid.Alpha(*v.F2, "ja") {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is alphanumeric string

The `alnum[:lang]` rule can be used to check if a field's value is a valid alphanumeric string.

The optional `lang` argument can be used to specify the alphabet's language. When not specified
the `lang` argument will default to `"en"`.

The validation is implemented by [`valid.Alnum`](https://pkg.go.dev/github.com/frk/valid#Alnum).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"alnum"`
	F2 *string `is:"alnum:bg"`
}
```

</td><td>

```go
if !valid.Alnum(v.F1, "en") {
	return errors.New("...")
}
if v.F2 != nil && !valid.Alnum(*v.F2, "bg") {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is bank identification code

The `bic` rule can be used to check if a field's value is a valid Bank Identification Code (or SWIFT code).

The validation is implemented by [`valid.BIC`](https://pkg.go.dev/github.com/frk/valid#BIC).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"bic"`
	F2 *string `is:"bic"`
}
```

</td><td>

```go
if !valid.BIC(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.BIC(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is bitcoin address

The `btc` rule can be used to check if a field's value is a valid bitcoin address.

The validation is implemented by [`valid.BTC`](https://pkg.go.dev/github.com/frk/valid#BTC).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"btc"`
	F2 *string `is:"btc"`
}
```

</td><td>

```go
if !valid.BTC(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.BTC(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is base-32 string

The `base32` rule can be used to check if a field's value is a valid base-32 string.

The validation is implemented by [`valid.Base32`](https://pkg.go.dev/github.com/frk/valid#Base32).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"base32"`
	F2 *string `is:"base32"`
}
```

</td><td>

```go
if !valid.Base32(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Base32(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is base-58 string

The `base58` rule can be used to check if a field's value is a valid base-58 string.

The validation is implemented by [`valid.Base58`](https://pkg.go.dev/github.com/frk/valid#Base58).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"base58"`
	F2 *string `is:"base58"`
}
```

</td><td>

```go
if !valid.Base58(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Base58(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is base-64 string

The `base64[:urlsafe]` rule can be used to check if a field's value is a valid base-64 string.

The optional `urlsafe` boolean argument can be used to specify whether the base-64
encoding is expected to be standard or url-safe. For readability the word `url` can
be used as an alias for `true`. When not provided, the `urlsafe` argument will default
to `false`.

The validation is implemented by [`valid.Base64`](https://pkg.go.dev/github.com/frk/valid#Base64).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"base64"`
	F2 *string `is:"base64"`
	F3 string  `is:"base64:url"`
	F4 *string `is:"base64:true"`
}
```

</td><td>

```go
if !valid.Base64(v.F1, false) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Base64(*v.F2, false) {
	return errors.New("...")
}
if !valid.Base64(v.F3, true) {
	return errors.New("...")
}
if v.F4 != nil && !valid.Base64(*v.F4, true) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is binary integer string

The `binary` rule can be used to check if a field's value is a valid binary integer string.

The validation is implemented by [`valid.Binary`](https://pkg.go.dev/github.com/frk/valid#Binary).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"binary"`
	F2 *string `is:"binary"`
}
```

</td><td>

```go
if !valid.Binary(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Binary(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is boolean value

The `bool` rule can be used to check if a field's string value represents a valid boolean.
The following are considered valid boolean strings: `"true"`, `"false"`, `"TRUE"`, `"FALSE"`.

The validation is implemented by [`valid.Bool`](https://pkg.go.dev/github.com/frk/valid#Bool).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"bool"`
	F2 *string `is:"bool"`
}
```

</td><td>

```go
if !valid.Bool(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Bool(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is classless inter-domain routing notation

The `cidr` rule can be used to check if a field's value is a valid Classless Inter-Domain Routing (CIDR) notation.

The validation is implemented by [`valid.CIDR`](https://pkg.go.dev/github.com/frk/valid#CIDR).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"cidr"`
	F2 *string `is:"cidr"`
}
```

</td><td>

```go
if !valid.CIDR(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.CIDR(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is card verification value

The `cvv` rule can be used to check if a field's value is a valid Card Verification Value (CVV).

The validation is implemented by [`valid.CVV`](https://pkg.go.dev/github.com/frk/valid#CVV).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"cvv"`
	F2 *string `is:"cvv"`
}
```

</td><td>

```go
if !valid.CVV(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.CVV(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is currency amount

The `ccy[:code][:opts]` rule can be used to check if a field's value is a valid currency amount.

The optional `code` argument can be used to specify the currency's [ISO-4217 code](https://en.wikipedia.org/wiki/ISO_4217).
When not specified, the `code` argument will default to `"usd"`.

The optional `opts` argument, which must be of type [`*valid.CurrencyOpts`](https://pkg.go.dev/github.com/frk/valid#CurrencyOpts),
can be used to provide additional options to the validation function. When not specified, the `opts` argument will default to `nil`,
which, in turn, will cause the implementation to use the [`valid.CurrencyOptsDefault`](https://pkg.go.dev/github.com/frk/valid#CurrencyOptsDefault) value.

The validation is implemented by [`valid.Currency`](https://pkg.go.dev/github.com/frk/valid#Currency).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"ccy"`
	F2 *string `is:"ccy"`
	F3 string  `is:"ccy:gbp"`
	F4 *string `is:"ccy:eur:&ccyOpts"`
	F5 string  `is:"ccy::&ccyOpts"`

	ccyOpts *valid.CurrencyOpts
}
```

</td><td>

```go
if !valid.Currency(v.F1, "usd", nil) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Currency(*v.F2, "usd", nil) {
	return errors.New("...")
}
if !valid.Currency(v.F3, "gbp", nil) {
	return errors.New("...")
}
if v.F4 != nil && !valid.Currency(*v.F4, "eur", v.ccyOpts) {
	return errors.New("...")
}
if !valid.Currency(v.F5, "usd", v.ccyOpts) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is data URI

The `datauri` rule can be used to check if a field's value is a valid [data URI](https://en.wikipedia.org/wiki/Data_URI_scheme).

The validation is implemented by [`valid.DataURI`](https://pkg.go.dev/github.com/frk/valid#DataURI).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"datauri"`
	F2 *string `is:"datauri"`
}
```

</td><td>

```go
if !valid.DataURI(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.DataURI(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is decimal number

The `decimal[:locale]` rule can be used to check if a field's value is a valid decimal number.

The optional `locale` argument can be used to specify the decimal number's locale.
When not specified, the `locale` argument will default to `"en"`.

The validation is implemented by [`valid.Decimal`](https://pkg.go.dev/github.com/frk/valid#Decimal).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"decimal"`
	F2 *string `is:"decimal"`
	F3 string  `is:"decimal:bg"`
	F4 *string `is:"decimal:hu"`
}
```

</td><td>

```go
if !valid.Decimal(v.F1, "en") {
	return errors.New("...")
}
if v.F2 != nil && !valid.Decimal(*v.F2, "en") {
	return errors.New("...")
}
if !valid.Decimal(v.F3, "bg") {
	return errors.New("...")
}
if v.F4 != nil && !valid.Decimal(*v.F4, "hu") {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is string of digits

The `digits` rule can be used to check if a field's value is a string of digits.

The validation is implemented by [`valid.Digits`](https://pkg.go.dev/github.com/frk/valid#Digits).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"digits"`
	F2 *string `is:"digits"`
}
```

</td><td>

```go
if !valid.Digits(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Digits(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is european article number

The `ean` rule can be used to check if a field's value is a valid European Article Number (EAN).

The validation is implemented by [`valid.EAN`](https://pkg.go.dev/github.com/frk/valid#EAN).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"ean"`
	F2 *string `is:"ean"`
}
```

</td><td>

```go
if !valid.EAN(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.EAN(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is employer identification number

TODO

## is ethereum address

The `eth` rule can be used to check if a field's value is a valid Ethereum address.

The validation is implemented by [`valid.ETH`](https://pkg.go.dev/github.com/frk/valid#ETH).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"eth"`
	F2 *string `is:"eth"`
}
```

</td><td>

```go
if !valid.ETH(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.ETH(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is email address

The `email` rule can be used to check if a field's value is a valid email address.

The validation is implemented by [`valid.Email`](https://pkg.go.dev/github.com/frk/valid#Email).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"email"`
	F2 *string `is:"email"`
}
```

</td><td>

```go
if !valid.Email(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Email(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is fully qualified domain name

The `fqdn` rule can be used to check if a field's value is a valid Fully Qualified Domain Name (FQDN).

The validation is implemented by [`valid.FQDN`](https://pkg.go.dev/github.com/frk/valid#FQDN).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"fqdn"`
	F2 *string `is:"fqdn"`
}
```

</td><td>

```go
if !valid.FQDN(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.FQDN(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is floating point number

The `float` rule can be used to check if a field's value is a valid floating point number.

The validation is implemented by [`valid.Float`](https://pkg.go.dev/github.com/frk/valid#Float).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"float"`
	F2 *string `is:"float"`
}
```

</td><td>

```go
if !valid.Float(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Float(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is HSL color

The `hsl` rule can be used to check if a field's value is a valid HSL (hue, saturation, lightness) color value.

The validation is implemented by [`valid.HSL`](https://pkg.go.dev/github.com/frk/valid#HSL).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"hsl"`
	F2 *string `is:"hsl"`
}
```

</td><td>

```go
if !valid.HSL(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.HSL(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is hash of algorithm

The `hash{:algo}` rule can be used to check if a field's value is a valid hash of the specified algorithm.

The required `algo` argument specifies the algorithm of the hash.

The validation is implemented by [`valid.Hash`](https://pkg.go.dev/github.com/frk/valid#Hash).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"hash:md5"`
	F2 *string `is:"hash:sha512"`
}
```

</td><td>

```go
if !valid.Hash(v.F1, "md5") {
	return errors.New("...")
}
if v.F2 != nil && !valid.Hash(*v.F2, "sha512") {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is hexadecimal string

The `hex` rule can be used to check if a field's value is a valid hexadecimal string.

The validation is implemented by [`valid.Hex`](https://pkg.go.dev/github.com/frk/valid#Hex).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"hex"`
	F2 *string `is:"hex"`
}
```

</td><td>

```go
if !valid.Hex(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Hex(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is hexadecimal color code

The `hexcolor` rule can be used to check if a field's value is a valid hexadecimal color code.

The validation is implemented by [`valid.HexColor`](https://pkg.go.dev/github.com/frk/valid#HexColor).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"hexcolor"`
	F2 *string `is:"hexcolor"`
}
```

</td><td>

```go
if !valid.HexColor(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.HexColor(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is international back account number

The `iban` rule can be used to check if a field's value is a valid International Bank Account Number (IBAN).

The validation is implemented by [`valid.IBAN`](https://pkg.go.dev/github.com/frk/valid#IBAN).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"iban"`
	F2 *string `is:"iban"`
}
```

</td><td>

```go
if !valid.IBAN(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.IBAN(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is identity card number

TODO

## is international mobile equipment identity number

The `imei` rule can be used to check if a field's value is a valid International Mobile-Equipment Identity (IMEI) number.

The validation is implemented by [`valid.IMEI`](https://pkg.go.dev/github.com/frk/valid#IMEI).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"imei"`
	F2 *string `is:"imei"`
}
```

</td><td>

```go
if !valid.IMEI(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.IMEI(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is internet protocol adress

The `ip[:ver]` rule can be used to check if a field's value is a valid Internet Protocol (IP) address.

The optional `ver` argument can be used to specify the IP version against which to check the field's value.
When not provided, the `ver` argument will default to `0`, which, in turn, causes the implementation to check
the value against both, version `4` and version `6`.

The validation is implemented by [`valid.IP`](https://pkg.go.dev/github.com/frk/valid#IP).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"ip"`
	F2 *string `is:"ip"`
	F3 string  `is:"ip:4"`
	F4 *string `is:"ip:6"`
}
```

</td><td>

```go
if !valid.IP(v.F1, 0) {
	return errors.New("...")
}
if v.F2 != nil && !valid.IP(*v.F2, 0) {
	return errors.New("...")
}
if !valid.IP(v.F3, 4) {
	return errors.New("...")
}
if v.F4 != nil && !valid.IP(*v.F4, 6) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is internet protocol adress range

The `iprange` rule can be used to check if a field's value is a valid Internet Protocol (IP) version 4 address range.

The validation is implemented by [`valid.IPRange`](https://pkg.go.dev/github.com/frk/valid#IPRange).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"iprange"`
	F2 *string `is:"iprange"`
}
```

</td><td>

```go
if !valid.IPRange(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.IPRange(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is international standard book number

The `isbn[:ver]` rule can be used to check if a field's value is a valid International Standard Book Number (ISBN).

The optional `ver` argument can be used to specify the ISBN version against which to check the field's value.
When not provided, the `ver` argument will default to `0`, which, in turn, causes the implementation to check
the value against both, version `10` and version `13`.

The validation is implemented by [`valid.ISBN`](https://pkg.go.dev/github.com/frk/valid#ISBN).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"isbn"`
	F2 *string `is:"isbn"`
	F3 string  `is:"isbn:10"`
	F4 *string `is:"isbn:13"`
}
```

</td><td>

```go
if !valid.ISBN(v.F1, 0) {
	return errors.New("...")
}
if v.F2 != nil && !valid.ISBN(*v.F2, 0) {
	return errors.New("...")
}
if !valid.ISBN(v.F3, 10) {
	return errors.New("...")
}
if v.F4 != nil && !valid.ISBN(*v.F4, 13) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is international securities identification number

The `isin` rule can be used to check if a field's value is a valid International Securities Identification Number (ISIN).

The validation is implemented by [`valid.ISIN`](https://pkg.go.dev/github.com/frk/valid#ISIN).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"isin"`
	F2 *string `is:"isin"`
}
```

</td><td>

```go
if !valid.ISIN(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.ISIN(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is ISO 639 string

The `iso639[:num]` rule can be used to check if a field's value is a valid
[ISO-639](https://en.wikipedia.org/wiki/ISO_639) language code.

The optional `num` argument can be used to specify the particular part of the
standard against which to check the field's value. When not provided, the `num`
argument will default to `0`, which, in turn, causes the implementation to check
the value against part `639-1` and part `639-2` of the standard.

The validation is implemented by [`valid.ISO639`](https://pkg.go.dev/github.com/frk/valid#ISO639).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"iso639"`
	F2 *string `is:"iso639"`
	F3 string  `is:"iso639:1"`
	F4 *string `is:"iso639:2"`
}
```

</td><td>

```go
if !valid.ISO639(v.F1, 0) {
	return errors.New("...")
}
if v.F2 != nil && !valid.ISO639(*v.F2, 0) {
	return errors.New("...")
}
if !valid.ISO639(v.F3, 1) {
	return errors.New("...")
}
if v.F4 != nil && !valid.ISO639(*v.F4, 2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is ISO 3166-1A string

The `iso31661a[:num]` rule can be used to check if a field's value is a valid
[ISO 3166-1](https://en.wikipedia.org/wiki/ISO_3166-1) country code (alpha only).

The optional `num` argument can be used to specify the particular set of the standard
against which to check the field's value. When not provided, the `num` argument will
default to `0`, which, in turn, causes the implementation to check the value against
both, the `Alpha-2` set and the `Alpha-3` set, of the standard.

The validation is implemented by [`valid.ISO31661A`](https://pkg.go.dev/github.com/frk/valid#ISO31661A).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"iso31661a"`
	F2 *string `is:"iso31661a"`
	F3 string  `is:"iso31661a:2"`
	F4 *string `is:"iso31661a:3"`
}
```

</td><td>

```go
if !valid.ISO31661A(v.F1, 0) {
	return errors.New("...")
}
if v.F2 != nil && !valid.ISO31661A(*v.F2, 0) {
	return errors.New("...")
}
if !valid.ISO31661A(v.F3, 2) {
	return errors.New("...")
}
if v.F4 != nil && !valid.ISO31661A(*v.F4, 3) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is ISO 4217 string

The `iso4217` rule can be used to check if a field's value is a valid
[ISO 4217](https://en.wikipedia.org/wiki/ISO_4217) currency code.

The validation is implemented by [`valid.ISO4217`](https://pkg.go.dev/github.com/frk/valid#ISO4217).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"iso4217"`
	F2 *string `is:"iso4217"`
}
```

</td><td>

```go
if !valid.ISO4217(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.ISO4217(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is international standard recording code

The `isrc` rule can be used to check if a field's value is a valid Internation Standard Recording Code (ISRC).

The validation is implemented by [`valid.ISRC`](https://pkg.go.dev/github.com/frk/valid#ISRC).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"isrc"`
	F2 *string `is:"isrc"`
}
```

</td><td>

```go
if !valid.ISRC(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.ISRC(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is international standard serial number

The `issn{:requireHyphen:caseSensitive}` rule can be used to check if a field's
value is a valid Internation Standard Serial Number (ISSN).

The *required* `requireHyphen` boolean argument specifies whether or not the validation
function should require the number to contain the "separator" hyphen.

The *required* `caseSensitive` boolean argument specifies whether or not the validation
function should disallow the number to contain the lower case `x`.

The validation is implemented by [`valid.ISSN`](https://pkg.go.dev/github.com/frk/valid#ISSN).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"issn:true:true"`
	F2 *string `is:"issn:false:false"`
}
```

</td><td>

```go
if !valid.ISSN(v.F1, true, true) {
	return errors.New("...")
}
if v.F2 != nil && !valid.ISSN(*v.F2, false, false) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is in list / is one of

The `in{:arg1...[:argN]}` rule can be used to check if a field's value is present
in the specified list of arguments.

The validation is implemented by [`valid.In`](https://pkg.go.dev/github.com/frk/valid#In).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string `is:"in:foo:bar:baz"`
	F2 *int   `is:"in:10:20:30:40"`
}
```

</td><td>

```go
if !valid.In(v.F1, "foo", "bar", "baz") {
	return errors.New("...")
}
if v.F2 != nil && !valid.In(*v.F2, 10, 20, 30, 40) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is integer number

The `int` rule can be used to check if a field's value is a valid integer number.

The validation is implemented by [`valid.Int`](https://pkg.go.dev/github.com/frk/valid#Int).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"int"`
	F2 *string `is:"int"`
}
```

</td><td>

```go
if !valid.Int(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Int(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is JSON value

The `json` rule can be used to check if a field's value is a valid JSON value.

The validation is implemented by [`valid.JSON`](https://pkg.go.dev/github.com/frk/valid#JSON).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"json"`
	F2 *[]byte `is:"json"`
}
```

</td><td>

```go
if !valid.JSON([]byte(v.F1)) {
	return errors.New("...")
}
if v.F2 != nil && !valid.JSON(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is JSON web token

The `jwt` rule can be used to check if a field's value is a valid JSON Web Token (JWT).

The validation is implemented by [`valid.JWT`](https://pkg.go.dev/github.com/frk/valid#JWT).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"jwt"`
	F2 *string `is:"jwt"`
}
```

</td><td>

```go
if !valid.JWT(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.JWT(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is latitude longitude string

The `latlong[:dms]` rule can be used to check if a field's value is a valid latitude-longitude coordinate string.

The optional `dms` boolean argument can be used to specify whether the coordinate
string should be in the DMS (degrees, minutes, and seconds) notaion or not. For
readability the word `dms` can be used as an alias for `true`. When not provided,
the `dms` argument will default to `false`.

The validation is implemented by [`valid.LatLong`](https://pkg.go.dev/github.com/frk/valid#LatLong).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"latlong"`
	F2 *string `is:"latlong:dms"`
}
```

</td><td>

```go
if !valid.LatLong(v.F1, false) {
	return errors.New("...")
}
if v.F2 != nil && !valid.LatLong(*v.F2, true) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is locale code

The `locale` rule can be used to check if a field's value is a valid locale code.

The validation is implemented by [`valid.Locale`](https://pkg.go.dev/github.com/frk/valid#Locale).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"locale"`
	F2 *string `is:"locale"`
}
```

</td><td>

```go
if !valid.Locale(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Locale(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is lower case string

The `lower` rule can be used to check if a field's string value contains only lower case.

The validation is implemented by [`valid.LowerCase`](https://pkg.go.dev/github.com/frk/valid#LowerCase).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"lower"`
	F2 *string `is:"lower"`
}
```

</td><td>

```go
if !valid.LowerCase(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.LowerCase(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is MAC address

The `mac[:space]` rule can be used to check if a field's value is a valid Media Access Control (MAC) address.

The optional `space` argument can be used to specify the MAC's numbering space against which to check the field's value.
When not provided, the `space` argument will default to `0`, which, in turn, causes the implementation to check
the value against both, space `6` (EUI-48) and space `8` (EUI-64).

The validation is implemented by [`valid.MAC`](https://pkg.go.dev/github.com/frk/valid#MAC).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"mac"`
	F2 *string `is:"mac"`
	F3 string  `is:"mac:6"`
	F4 *string `is:"mac:8"`
}
```

</td><td>

```go
if !valid.MAC(v.F1, 0) {
	return errors.New("...")
}
if v.F2 != nil && !valid.MAC(*v.F2, 0) {
	return errors.New("...")
}
if !valid.MAC(v.F3, 6) {
	return errors.New("...")
}
if v.F4 != nil && !valid.MAC(*v.F4, 8) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is MD5 hash

The `md5` rule can be used to check if a field's value is a valid MD5 hash.

The validation is implemented by [`valid.MD5`](https://pkg.go.dev/github.com/frk/valid#MD5).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"md5"`
	F2 *string `is:"md5"`
}
```

</td><td>

```go
if !valid.MD5(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.MD5(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is MIME type

The `mime` rule can be used to check if a field's value is a valid MIME type.

The validation is implemented by [`valid.MIME`](https://pkg.go.dev/github.com/frk/valid#MIME).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"mime"`
	F2 *string `is:"mime"`
}
```

</td><td>

```go
if !valid.MIME(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.MIME(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is magnet URI

The `magneturi` rule can be used to check if a field's value is a valid Magnet URI.

The validation is implemented by [`valid.MagnetURI`](https://pkg.go.dev/github.com/frk/valid#MagnetURI).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"magneturi"`
	F2 *string `is:"magneturi"`
}
```

</td><td>

```go
if !valid.MagnetURI(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.MagnetURI(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is mongo ID

The `mongoid` rule can be used to check if a field's value is a valid MongoDB ID.

The validation is implemented by [`valid.MongoId`](https://pkg.go.dev/github.com/frk/valid#MongoId).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"mongoid"`
	F2 *string `is:"mongoid"`
}
```

</td><td>

```go
if !valid.MongoId(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.MongoId(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is numeric string

The `numeric` rule can be used to check if a field's value is a valid numeric string.

The validation is implemented by [`valid.Numeric`](https://pkg.go.dev/github.com/frk/valid#Numeric).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"numeric"`
	F2 *string `is:"numeric"`
}
```

</td><td>

```go
if !valid.Numeric(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Numeric(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is octal number

The `octal` rule can be used to check if a field's value is a valid octal number.

The validation is implemented by [`valid.Octal`](https://pkg.go.dev/github.com/frk/valid#Octal).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"octal"`
	F2 *string `is:"octal"`
}
```

</td><td>

```go
if !valid.Octal(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Octal(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is primary account number

The `pan` rule can be used to check if a field's value is a valid Primary Account Number (or Credit Card number).

The validation is implemented by [`valid.PAN`](https://pkg.go.dev/github.com/frk/valid#PAN).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"pan"`
	F2 *string `is:"pan"`
}
```

</td><td>

```go
if !valid.PAN(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.PAN(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is phone number

The `phone[:cc]` rule can be used to check if a field's value is a valid phone number.

The optional `cc` argument can be used to specify the [ISO-3166-1 country code](https://en.wikipedia.org/wiki/ISO_3166-1).
When not specified, the `cc` argument will default to `"us"`.

The validation is implemented by [`valid.Phone`](https://pkg.go.dev/github.com/frk/valid#Phone).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"phone"`
	F2 *string `is:"phone"`
	F3 string  `is:"phone:gb"`
	F4 *string `is:"phone:jp"`
}
```

</td><td>

```go
if !valid.Phone(v.F1, "us") {
	return errors.New("...")
}
if v.F2 != nil && !valid.Phone(*v.F2, "us") {
	return errors.New("...")
}
if !valid.Phone(v.F3, "gb") {
	return errors.New("...")
}
if v.F4 != nil && !valid.Phone(*v.F4, "jp") {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is port number

The `port` rule can be used to check if a field's value is a valid port number.

The validation is implemented by [`valid.Port`](https://pkg.go.dev/github.com/frk/valid#Port).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"port"`
	F2 *string `is:"port"`
}
```

</td><td>

```go
if !valid.Port(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Port(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is RGB color

The `rgb` rule can be used to check if a field's value is a valid RGB color value.

The validation is implemented by [`valid.RGB`](https://pkg.go.dev/github.com/frk/valid#RGB).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"rgb"`
	F2 *string `is:"rgb"`
}
```

</td><td>

```go
if !valid.RGB(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.RGB(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is social security number

The `ssn` rule can be used to check if a field's value is a valid Social Security Number.

The validation is implemented by [`valid.SSN`](https://pkg.go.dev/github.com/frk/valid#SSN).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"ssn"`
	F2 *string `is:"ssn"`
}
```

</td><td>

```go
if !valid.SSN(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.SSN(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is semantic version number

The `semver` rule can be used to check if a field's value is a valid [Semantic Version](https://semver.org/) number.

The validation is implemented by [`valid.SemVer`](https://pkg.go.dev/github.com/frk/valid#SemVer).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"semver"`
	F2 *string `is:"semver"`
}
```

</td><td>

```go
if !valid.SemVer(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.SemVer(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

## is slug

The `slug` rule can be used to check if a field's value is a valid slug string.

The validation is implemented by [`valid.Slug`](https://pkg.go.dev/github.com/frk/valid#Slug).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"slug"`
	F2 *string `is:"slug"`
}
```

</td><td>

```go
if !valid.Slug(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Slug(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is strong password

The `strongpass[:opts]` rule can be used to check if a field's value is a strong password.

The optional `opts` argument, which must be of type [`*valid.StrongPasswordOpts`](https://pkg.go.dev/github.com/frk/valid#StrongPasswordOpts),
can be used to provide additional options to the validation function. When not specified, the `opts` argument will default to `nil`,
which, in turn, will cause the implementation to use the [`valid.StrongPasswordOptsDefault`](https://pkg.go.dev/github.com/frk/valid#StrongPasswordOptsDefault) value.

The validation is implemented by [`valid.StrongPassword`](https://pkg.go.dev/github.com/frk/valid#StrongPassword).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"strongpass"`
	F2 *string `is:"strongpass"`
	F3 string  `is:"strongpass:&pwOpts"`

	pwOpts *valid.StrongPasswordOpts
}
```

</td><td>

```go
if !valid.StrongPassword(v.F1, nil) {
	return errors.New("...")
}
if v.F2 != nil && !valid.StrongPassword(*v.F2, nil) {
	return errors.New("...")
}
if !valid.StrongPassword(v.F3, v.pwOpts) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is universally unique identification number

The `uuid[:ver]` rule can be used to check if a field's value is a valid Universally Unique Identification (UUID) number.

The optional `ver` argument can be used to specify the UUID version against which
to check the field's value. When not provided, the `ver` argument will default to `4`.
Currently the supported versions are `3`, `4`, and `5` with aliases `v3`, `v4`,
and `v5` respectively.

The validation is implemented by [`valid.UUID`](https://pkg.go.dev/github.com/frk/valid#UUID).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"uuid"`
	F2 *string `is:"uuid"`
	F3 string  `is:"uuid:3"`
	F4 string  `is:"uuid:v4"`
	F5 string  `is:"uuid:v5"`
}
```

</td><td>

```go
if !valid.UUID(v.F1, 4) {
	return errors.New("...")
}
if v.F2 != nil && !valid.UUID(*v.F2, 4) {
	return errors.New("...")
}
if !valid.UUID(v.F3, 3) {
	return errors.New("...")
}
if !valid.UUID(v.F4, 4) {
	return errors.New("...")
}
if !valid.UUID(v.F5, 5) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is unsigned integer number

The `uint` rule can be used to check if a field's value is a valid unsigned integer number.

The validation is implemented by [`valid.Uint`](https://pkg.go.dev/github.com/frk/valid#Uint).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"uint"`
	F2 *string `is:"uint"`
}
```

</td><td>

```go
if !valid.Uint(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.Uint(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is upper case string

The `upper` rule can be used to check if a field's string value contains only upper case.

The validation is implemented by [`valid.UpperCase`](https://pkg.go.dev/github.com/frk/valid#UpperCase).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"upper"`
	F2 *string `is:"upper"`
}
```

</td><td>

```go
if !valid.UpperCase(v.F1) {
	return errors.New("...")
}
if v.F2 != nil && !valid.UpperCase(*v.F2) {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is value added tax number

The `vat[:cc]` rule can be used to check if a field's value is a valid Value Added Tax (VAT) number.

The optional `cc` argument can be used to specify the [ISO-3166-1 country code](https://en.wikipedia.org/wiki/ISO_3166-1).
When not specified, the `cc` argument will default to `"us"`.

The validation is implemented by [`valid.VAT`](https://pkg.go.dev/github.com/frk/valid#VAT).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"vat"`
	F2 *string `is:"vat"`
	F3 string  `is:"vat:gb"`
	F4 *string `is:"vat:ru"`
}
```

</td><td>

```go
if !valid.VAT(v.F1, "us") {
	return errors.New("...")
}
if v.F2 != nil && !valid.VAT(*v.F2, "us") {
	return errors.New("...")
}
if !valid.VAT(v.F3, "gb") {
	return errors.New("...")
}
if v.F4 != nil && !valid.VAT(*v.F4, "ru") {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>


## is zip code

The `zip[:cc]` rule can be used to check if a field's value is a valid zip / postal code.

The optional `cc` argument can be used to specify the [ISO-3166-1 country code](https://en.wikipedia.org/wiki/ISO_3166-1).
When not specified, the `cc` argument will default to `"us"`.

The validation is implemented by [`valid.Zip`](https://pkg.go.dev/github.com/frk/valid#Zip).

<table><thead><tr><th>Rule Tag</th><th>Generated Output</th></tr></thead><tbody>
<tr><td>

```go
type Validator struct {
	F1 string  `is:"zip"`
	F2 *string `is:"zip"`
	F3 string  `is:"zip:gb"`
	F4 *string `is:"zip:ru"`
}
```

</td><td>

```go
if !valid.Zip(v.F1, "us") {
	return errors.New("...")
}
if v.F2 != nil && !valid.Zip(*v.F2, "us") {
	return errors.New("...")
}
if !valid.Zip(v.F3, "gb") {
	return errors.New("...")
}
if v.F4 != nil && !valid.Zip(*v.F4, "ru") {
	return errors.New("...")
}
```

</td></tr>
</tbody></table>

