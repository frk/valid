package command

import (
	"fmt"
	"os"
)

func printUsage() {
	fmt.Fprint(os.Stderr, usage)
}

const usage = `usage: isvalid [-wd] [-r] [-f] [-rx] [-o] [-fktag] [-fkbase] [-fksep]

isvalid generates struct field validation .... (todo: write doc)


The -wd flag specifies the directory whose files the tool will process. When used
together with the -f or -rx flags the tool will process only those files that match
the -f and -rx values. If left unespecified, the current working directory will be
used by default.


The -r flag instructs the tool to process the files in the whole hierarchy of the
working directory. When used together with the -f or -rx flags the tool will process
only those files that match the -f and -rx values.


The -f flag specifies a file to be used as input for the tool. The file must be
located in the working directory. The flag can be used more than once to specify
multiple files.


The -rx flag specifies a regular expressions to match input files that the tool should
process. The regular expressions must match files located in the working directory.
The flag can be used more than once to specify multiple regular expressions.


The -o flag specifies the format to be used for generating the name of the output files.
The format can contain one (and only one) "%s" placeholder which the tool will replace
with the input file's base name, if no placeholder is present then the input file's base
name will be prefixed to the format.
If left unspecified, the format "%s_isvalid.go" will be used by default.


The -fktag flag if set to a non-empty string, specifies the struct tag to be used
for constructing the field keys that will be used by the generator for error reporting.
A valid tag must begin with a letter (A-z) or an underscore (_), subsequent characters
in the tag can be letters, underscores, and digits (0-9). If set to "" (empty string),
the generator will default to use the field names instead of struct tags to construct
the field keys. If left unspecified, the tag "json" will be used by default.


The -fkjoin flag if set to true, specifies that a nested struct field's key will be
produced by joining it together with all of its parent fields. If set to false, such
a field's key will be produced only from that field's name/tag. If left unspecified,
the value true will be used by default.


The -fksep flag specifies the separator to be used for joining fields' tags/names
when producing the field keys. The separator can be at most one byte long.
If left unspecified, the separator "." will be used by default.

` //`
