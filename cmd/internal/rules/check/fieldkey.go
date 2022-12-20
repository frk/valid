package check

import (
	"strconv"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/types"

	"github.com/frk/tagutil"
)

// FieldKeyFunc is the type of the function called by the checker
// for each field to generate a unique key from the FieldChain.
type FieldKeyFunc func(ff types.FieldChain) (key string)

// makeFKFunc returns a function that, based on the given configuration,
// generates a unique field key for a given FieldChain.
func makeFKFunc(c *config.FieldKeyConfig) FieldKeyFunc {
	// keyset & unique are used by the returned function
	// to ensure that the generated key is unique.
	keyset := make(map[string]uint)
	unique := func(key string) string {
		if num, ok := keyset[key]; ok {
			keyset[key] = num + 1
			key += "-" + strconv.FormatUint(uint64(num), 10)
		} else {
			keyset[key] = 1
		}
		return key
	}

	if c != nil && len(c.Tag.Value) > 0 {
		if c.Join.Value {
			// Returns the joined tag values of the fields in the given slice.
			// If one of the fields does not have a tag value set, their name
			// will be used in the join as default.
			return func(ff types.FieldChain) (key string) {
				tag := c.Tag.Value
				sep := c.Separator.Value

				for _, f := range ff {
					t := tagutil.New(f.Tag)
					if t.Contains("is", "omitkey") || f.IsEmbedded {
						continue
					}

					v := t.First(tag)
					if len(v) == 0 {
						v = f.Name
					}

					// TODO(mkopriva): if the field's type an array, a slice,
					// or a map, it might be useful to add the index expression
					// delimiters with a fmt verb for the value which would be
					// supplied during error construction.
					//
					// When implementing this make sure it's used by the generator
					// only when there are values (indexes, keys) to be supplied
					// for the fmt verbs, otherwise the error message will become
					// unclear.
					key += v + sep
				}
				if len(sep) > 0 && len(key) > len(sep) {
					return unique(key[:len(key)-len(sep)])
				}
				return unique(key)
			}
		}

		// Returns the tag value of the last field, if no value was
		// set the field's name will be returned instead.
		return func(ff types.FieldChain) string {
			t := tagutil.New(ff[len(ff)-1].Tag)
			if key := t.First(c.Tag.Value); len(key) > 0 {
				return unique(key)
			}
			return unique(ff[len(ff)-1].Name)
		}
	}

	if c != nil && c.Join.Value {
		sep := c.Separator.Value

		// Returns the joined names of the fields in the given slice.
		return func(ff types.FieldChain) (key string) {
			for _, f := range ff {
				t := tagutil.New(f.Tag)
				if t.Contains("is", "omitkey") || f.IsEmbedded {
					continue
				}

				// TODO(mkopriva): if the field's type an array, a slice,
				// or a map, it might be useful to add the index expression
				// delimiters with a fmt verb for the value which would be
				// supplied during error construction.
				//
				// When implementing this make sure it's used by the generator
				// only when there are values (indexes, keys) to be supplied
				// for the fmt verbs, otherwise the error message will become
				// unclear.
				key += f.Name + sep
			}
			if len(sep) > 0 && len(key) > len(sep) {
				return unique(key[:len(key)-len(sep)])
			}
			return unique(key)
		}
	}

	// Returns the name of the last field.
	return func(ff types.FieldChain) string {
		return unique(ff[len(ff)-1].Name)
	}
}
