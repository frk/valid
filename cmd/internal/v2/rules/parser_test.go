package rules

import (
	"testing"

	"github.com/frk/compare"
)

func TestParse(t *testing.T) {
	rr := registry{
		builtins: map[string]Rule{
			"required": {Name: "required"},
			"notnil":   {Name: "notnil"},
			"optional": {Name: "optional"},
		},
		custom: map[string]Rule{
			"rule": {Name: "rule"},
			"ra":   {Name: "ra"},
			"re":   {Name: "re"},
			"ri":   {Name: "ri"},
			"ru":   {Name: "ru"},
			"ro":   {Name: "ro"},
			"la":   {Name: "la"},
			"le":   {Name: "le"},
			"li":   {Name: "li"},
			"lu":   {Name: "lu"},
			"lo":   {Name: "lo"},
			"a":    {Name: "a"},
			"b":    {Name: "b"},
			"c":    {Name: "c"},
		},
	}

	////////////////////////////////////////////////////////////////////////

	tests := []struct {
		str  string
		want *Obj
		err  error
	}{{
		// single plain rule
		str:  "foo",
		want: nil,
		err:  &UndefinedRuleError{RuleName: "foo"},
	}, {
		// single plain rule
		str:  "rule",
		want: &Obj{Rules: []*Rule{{Name: "rule"}}},
	}, {
		// single rule with argument
		str: "rule:arg",
		want: &Obj{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "arg", Type: ARG_STRING},
		}}}},
	}, {
		// single rule with arguments
		str: "rule:arg:123:true:0.0064",
		want: &Obj{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "arg", Type: ARG_STRING},
			{Value: "123", Type: ARG_INT},
			{Value: "true", Type: ARG_BOOL},
			{Value: "0.0064", Type: ARG_FLOAT},
		}}}},
	}, {
		// single rule with empty argument
		str: "rule:",
		want: &Obj{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "", Type: ARG_UNKNOWN},
		}}}},
	}, {
		// single rule with empty arguments
		str: "rule::::",
		want: &Obj{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "", Type: ARG_UNKNOWN},
			{Value: "", Type: ARG_UNKNOWN},
			{Value: "", Type: ARG_UNKNOWN},
			{Value: "", Type: ARG_UNKNOWN},
		}}}},
	}, {
		// single rule with empty & non-empty arguments
		str: "rule:arg::true:::0.0064:",
		want: &Obj{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "arg", Type: ARG_STRING},
			{Value: "", Type: ARG_UNKNOWN},
			{Value: "true", Type: ARG_BOOL},
			{Value: "", Type: ARG_UNKNOWN},
			{Value: "", Type: ARG_UNKNOWN},
			{Value: "0.0064", Type: ARG_FLOAT},
			{Value: "", Type: ARG_UNKNOWN},
		}}}},
	}, {
		// single rule with quoted argument
		str: "rule:\"arg\"",
		want: &Obj{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "arg", Type: ARG_STRING},
		}}}},
	}, {
		// single rule with quoted argument
		str: "rule:\"foo \\\"bar\\\" baz\"",
		want: &Obj{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "foo \\\"bar\\\" baz", Type: ARG_STRING},
		}}}},
	}, {
		// single rule with quoted, empty, and non-empty arguments
		str: "rule:\"foo\":bar:\"\":123::\"b \\\"a\\\" z\"",
		want: &Obj{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "foo", Type: ARG_STRING},
			{Value: "bar", Type: ARG_STRING},
			{Value: "", Type: ARG_STRING},
			{Value: "123", Type: ARG_INT},
			{Value: "", Type: ARG_UNKNOWN},
			{Value: "b \\\"a\\\" z", Type: ARG_STRING},
		}}}},
	}, {
		// rule with abs field
		str: "rule:&F1.F2.F3",
		want: &Obj{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "F1.F2.F3", Type: ARG_FIELD_ABS},
		}}}},
	}, {
		// rule with rel field
		str: "rule:.f1.f2",
		want: &Obj{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "f1.f2", Type: ARG_FIELD_REL},
		}}}},
	}, {
		// multiple plain rules
		str: "ra,re,ri,ru,ro",
		want: &Obj{Rules: []*Rule{
			{Name: "ra"}, {Name: "re"}, {Name: "ri"}, {Name: "ru"}, {Name: "ro"},
		}},
	}, {
		// multiple plain rules (omit empty rules)
		str: "ra,,,re,ri,,",
		want: &Obj{Rules: []*Rule{
			{Name: "ra"}, {Name: "re"}, {Name: "ri"},
		}},
	}, {
		// multiple rules with arguments
		str: "ra:a:b:c,re:foo::321:,ri:1:2:3",
		want: &Obj{Rules: []*Rule{
			{Name: "ra", Args: []*Arg{
				{Value: "a", Type: ARG_STRING},
				{Value: "b", Type: ARG_STRING},
				{Value: "c", Type: ARG_STRING},
			}},
			{Name: "re", Args: []*Arg{
				{Value: "foo", Type: ARG_STRING},
				{Value: "", Type: ARG_UNKNOWN},
				{Value: "321", Type: ARG_INT},
				{Value: "", Type: ARG_UNKNOWN},
			}},
			{Name: "ri", Args: []*Arg{
				{Value: "1", Type: ARG_INT},
				{Value: "2", Type: ARG_INT},
				{Value: "3", Type: ARG_INT},
			}},
		}},
	}, {
		// nested rule (elem)
		str: "[]ra",
		want: &Obj{Elem: &Obj{
			Rules: []*Rule{{Name: "ra"}},
		}},
	}, {
		// flat & nested
		str: "ra,[]re",
		want: &Obj{
			Rules: []*Rule{{Name: "ra"}},
			Elem: &Obj{
				Rules: []*Rule{{Name: "re"}},
			},
		},
	}, {
		// flat & nested #2
		str: "required,[]re",
		want: &Obj{
			Rules: []*Rule{{Name: "required"}},
			Elem: &Obj{
				Rules: []*Rule{{Name: "re"}},
			},
		},
	}, {
		// flat & nested #3
		str: "notnil,[]re",
		want: &Obj{
			Rules: []*Rule{{Name: "notnil"}},
			Elem: &Obj{
				Rules: []*Rule{{Name: "re"}},
			},
		},
	}, {
		// flat & nested #4
		str: "optional,[]re",
		want: &Obj{
			Rules: []*Rule{{Name: "optional"}},
			Elem: &Obj{
				Rules: []*Rule{{Name: "re"}},
			},
		},
	}, {
		// nested rule (elem [levels])
		str: "[][][][]ra",
		want: &Obj{
			Elem: &Obj{Elem: &Obj{Elem: &Obj{Elem: &Obj{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
		},
	}, {
		// nested rules (elem [levels])
		str: "[]ra,re,[]re,ri,[]ri,ru,[]ru,ro",
		want: &Obj{
			Elem: &Obj{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Elem: &Obj{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Elem: &Obj{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Elem: &Obj{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// nested rule (key)
		str: "[ra]",
		want: &Obj{Key: &Obj{
			Rules: []*Rule{{Name: "ra"}},
		}},
	}, {
		// nested rule (key [levels])
		str: "[[[[ra]]]]",
		want: &Obj{
			Key: &Obj{Key: &Obj{Key: &Obj{Key: &Obj{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
		},
	}, {
		// nested rules (key [levels])
		str: "[ra,re,[re,ri,[ri,ru,[ru,ro]]]]",
		want: &Obj{
			Key: &Obj{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Key: &Obj{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Key: &Obj{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Key: &Obj{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// nested rules (key & elem)
		str: "[ra]re",
		want: &Obj{Key: &Obj{
			Rules: []*Rule{{Name: "ra"}},
		}, Elem: &Obj{
			Rules: []*Rule{{Name: "re"}},
		}},
	}, {
		// nested rules (key & elem [levels])
		str: "[[[[ra]]]][][][]re",
		want: &Obj{
			Key: &Obj{Key: &Obj{Key: &Obj{Key: &Obj{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
			Elem: &Obj{Elem: &Obj{Elem: &Obj{Elem: &Obj{
				Rules: []*Rule{{Name: "re"}},
			}}}},
		},
	}, {
		// nested rules (key & elems [levels])
		str: "[ra,re,[re,ri,[ri,ru,[ru,ro]]]]ra,re,[]re,ri,[]ri,ru,[]ru,ro",
		want: &Obj{
			Key: &Obj{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Key: &Obj{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Key: &Obj{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Key: &Obj{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
			Elem: &Obj{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Elem: &Obj{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Elem: &Obj{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Elem: &Obj{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// ... with arguments and all ...
		str: "[ra,re:1:2:3,[re::\"]]\\\"[]]\":foo,ri,[ri:&MyField:::-321,ru:\"  \",[ru,ro:\"]\"]" +
			"ro:\"[\",ru]ru:foo:123::&MyOtherField:]ri,re::\"]]\\\"[]]\":foo]ra:xyz:,re:&mykey," +
			"[la:\"]heee![\"]re,ri:,[le:a,li:b,lu:c]ri:\"foo \\\"]]]\":,ru::-abc,[c:lu,b:li,a:le]ru,ro:\"[foo]\":",
		want: &Obj{
			Key: &Obj{
				Rules: []*Rule{
					{Name: "ra"},
					{Name: "re", Args: []*Arg{
						{Value: "1", Type: ARG_INT},
						{Value: "2", Type: ARG_INT},
						{Value: "3", Type: ARG_INT},
					}},
				},
				Key: &Obj{
					Rules: []*Rule{
						{Name: "re", Args: []*Arg{
							{Value: "", Type: ARG_UNKNOWN},
							{Value: "]]\\\"[]]", Type: ARG_STRING},
							{Value: "foo", Type: ARG_STRING},
						}},
						{Name: "ri"},
					},
					Key: &Obj{
						Rules: []*Rule{
							{Name: "ri", Args: []*Arg{
								{Value: "MyField", Type: ARG_FIELD_ABS},
								{Value: "", Type: ARG_UNKNOWN},
								{Value: "", Type: ARG_UNKNOWN},
								{Value: "-321", Type: ARG_INT},
							}},
							{Name: "ru", Args: []*Arg{
								{Value: "  ", Type: ARG_STRING},
							}},
						},
						Key: &Obj{
							Rules: []*Rule{
								{Name: "ru"},
								{Name: "ro", Args: []*Arg{
									{Value: "]", Type: ARG_STRING},
								}},
							},
						},
						Elem: &Obj{
							Rules: []*Rule{
								{Name: "ro", Args: []*Arg{
									{Value: "[", Type: ARG_STRING},
								}},
								{Name: "ru"},
							},
						},
					},
					Elem: &Obj{
						Rules: []*Rule{
							{Name: "ru", Args: []*Arg{
								{Value: "foo", Type: ARG_STRING},
								{Value: "123", Type: ARG_INT},
								{Value: "", Type: ARG_UNKNOWN},
								{Value: "MyOtherField", Type: ARG_FIELD_ABS},
								{Value: "", Type: ARG_UNKNOWN},
							}},
						},
					},
				},
				Elem: &Obj{
					Rules: []*Rule{
						{Name: "ri"},
						{Name: "re", Args: []*Arg{
							{Value: "", Type: ARG_UNKNOWN},
							{Value: "]]\\\"[]]", Type: ARG_STRING},
							{Value: "foo", Type: ARG_STRING},
						}},
					},
				},
			},
			Elem: &Obj{
				Rules: []*Rule{
					{Name: "ra", Args: []*Arg{
						{Value: "xyz", Type: ARG_STRING},
						{Value: "", Type: ARG_UNKNOWN},
					}},
					{Name: "re", Args: []*Arg{
						{Value: "mykey", Type: ARG_FIELD_ABS},
					}},
				},
				Key: &Obj{
					Rules: []*Rule{{Name: "la", Args: []*Arg{
						{Value: "]heee![", Type: ARG_STRING},
					}}},
				},
				Elem: &Obj{
					Rules: []*Rule{
						{Name: "re"},
						{Name: "ri", Args: []*Arg{
							{Value: "", Type: ARG_UNKNOWN},
						}},
					},
					Key: &Obj{
						Rules: []*Rule{
							{Name: "le", Args: []*Arg{{Value: "a", Type: ARG_STRING}}},
							{Name: "li", Args: []*Arg{{Value: "b", Type: ARG_STRING}}},
							{Name: "lu", Args: []*Arg{{Value: "c", Type: ARG_STRING}}},
						},
					},
					Elem: &Obj{
						Rules: []*Rule{
							{Name: "ri", Args: []*Arg{
								{Value: "foo \\\"]]]", Type: ARG_STRING},
								{Value: "", Type: ARG_UNKNOWN},
							}},
							{Name: "ru", Args: []*Arg{
								{Value: "", Type: ARG_UNKNOWN},
								{Value: "-abc", Type: ARG_STRING},
							}},
						},
						Key: &Obj{
							Rules: []*Rule{
								{Name: "c", Args: []*Arg{{Value: "lu", Type: ARG_STRING}}},
								{Name: "b", Args: []*Arg{{Value: "li", Type: ARG_STRING}}},
								{Name: "a", Args: []*Arg{{Value: "le", Type: ARG_STRING}}},
							},
						},
						Elem: &Obj{
							Rules: []*Rule{
								{Name: "ru"},
								{Name: "ro", Args: []*Arg{
									{Value: "[foo]", Type: ARG_STRING},
									{Value: "", Type: ARG_UNKNOWN},
								}},
							},
						},
					},
				},
			},
		},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}

	for _, tt := range tests {
		got, err := Parse(tt.str, rr)
		if e := compare.Compare(err, tt.err); e != nil {
			t.Error(e)
		}
		if e := compare.Compare(got, tt.want); e != nil {
			t.Errorf("`%s`: %v", tt.str, e)
		}
	}
}
