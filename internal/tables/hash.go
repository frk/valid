package tables

var HashAlgoLen = map[string]int{
	"md5":       32,
	"md4":       32,
	"sha1":      40,
	"sha256":    64,
	"sha384":    96,
	"sha512":    128,
	"ripemd128": 32,
	"ripemd160": 40,
	"tiger128":  32,
	"tiger160":  40,
	"tiger192":  48,
	"crc32":     8,
	"crc32b":    8,
}