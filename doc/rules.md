
required: [applies to pointer(s) and the base]
- the `required` rule means that the field MUST be non-nil-and-non-zero. For instance,
if you have a `F *string` field then the `required` rule would first check that the field
is not `nil` and then check if the dereferenced string value is empty.

notnil: [applies to pointer(s) and a nilable base]

optional: [applies to pointer(s) and the base]
	- the `optional` rule means that the field does NOT have to be non-nil-and-non-zero.
	- the `optional` rule is applied by default when `required` and `notnil` rules are specified.
	- if neither `required`, `notnil`, `optional`, nor `noguard` are present
	then `optional` will be applied automatically.

omitnil:
	- applies to pointers and nilable base


noguard
	- applies to pointers only

-----------------
- `F string`
	- required: string cannot be empty
	- optional: string can be empty, (any other rules are applied only if
	  non-empty, but empty does not constitute an error)
	- notnil: <does not apply>
	- omitnil: <does not apply>
	- noguard: <does not apply>
- `F *string`
	- required: pointer cannot be nil AND string cannot be empty
	- optional: pointer can be nil AND string can be empty, (any other rules are
	  applied only if non-nil-and-non-empty, but nil or empty does not constitute
	  an error)
	- notnil: pointer cannot be nil BUT string can be empty
	- omitnil: pointer can be nil AND string can be empty
	- noguard: assume pointer is not nil
- `F *any`
	- required: pointer cannot be nil AND string cannot be empty
	- optional: pointer can be nil AND string can be empty, (any other rules are
	  applied only if non-nil-and-non-empty, but nil or empty does not constitute
	  an error)
	- notnil: pointer cannot be nil BUT string can be empty
	- omitnil: pointer can be nil AND string can be empty
	- noguard: assume pointer is not nil
