package tables

// TODO
// - https://github.com/unicode-org/cldr/tree/release-38/common
// - http://cldr.unicode.org/index/downloads
// - http://cldr.unicode.org/

var Decimal = make(map[string]rune)

var decimal = []struct {
	char    rune
	locales []string
}{
	{char: '.', locales: []string{
		"en-US", "en-AU", "en-GB", "en-HK", "en-IN", "en-NZ", "en-ZA", "en-ZM",
	}},
	{char: '٫', locales: []string{
		"ar", "ar-AE", "ar-BH", "ar-DZ", "ar-EG", "ar-IQ", "ar-JO", "ar-KW",
		"ar-LB", "ar-LY", "ar-MA", "ar-QM", "ar-QA", "ar-SA", "ar-SD", "ar-SY",
		"ar-TN", "ar-YE",
	}},
}
