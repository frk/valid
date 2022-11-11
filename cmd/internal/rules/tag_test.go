package rules

import (
	"testing"

	"github.com/frk/compare"
)

func TestParse(t *testing.T) {
	tests := []struct {
		tag  string
		want *Tag
	}{{
		tag:  ``,
		want: nil,
	}, {
		tag:  `json:"foo,omitempty" xml:">abc"`,
		want: nil,
	}, {
		tag:  `json:"foo,omitempty" is:"r1" xml:">abc"`,
		want: &Tag{Rules: []*Rule{{Name: "r1"}}},
	}, {
		// single plain rule
		tag:  `is:"rule"`,
		want: &Tag{Rules: []*Rule{{Name: "rule"}}},
	}, {
		// single rule with argument
		tag: `is:"rule:arg"`,
		want: &Tag{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "arg", Type: ARG_STRING},
		}}}},
	}, {
		// single rule with arguments
		tag: `is:"rule:arg:123:true:0.0064"`,
		want: &Tag{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "arg", Type: ARG_STRING},
			{Value: "123", Type: ARG_INT},
			{Value: "true", Type: ARG_BOOL},
			{Value: "0.0064", Type: ARG_FLOAT},
		}}}},
	}, {
		// single rule with empty argument
		tag: `is:"rule:"`,
		want: &Tag{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "", Type: ARG_UNKNOWN},
		}}}},
	}, {
		// single rule with empty arguments
		tag: `is:"rule::::"`,
		want: &Tag{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "", Type: ARG_UNKNOWN},
			{Value: "", Type: ARG_UNKNOWN},
			{Value: "", Type: ARG_UNKNOWN},
			{Value: "", Type: ARG_UNKNOWN},
		}}}},
	}, {
		// single rule with empty & non-empty arguments
		tag: `is:"rule:arg::true:::0.0064:"`,
		want: &Tag{Rules: []*Rule{{Name: "rule", Args: []*Arg{
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
		tag: `is:"rule:\"arg\""`,
		want: &Tag{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "arg", Type: ARG_STRING},
		}}}},
	}, {
		// single rule with quoted argument
		tag: `is:"rule:\"foo \\\"bar\\\" baz\""`,
		want: &Tag{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "foo \\\"bar\\\" baz", Type: ARG_STRING},
		}}}},
	}, {
		// single rule with quoted, empty, and non-empty arguments
		tag: `is:"rule:\"foo\":bar:\"\":123::\"b \\\"a\\\" z\""`,
		want: &Tag{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "foo", Type: ARG_STRING},
			{Value: "bar", Type: ARG_STRING},
			{Value: "", Type: ARG_STRING},
			{Value: "123", Type: ARG_INT},
			{Value: "", Type: ARG_UNKNOWN},
			{Value: "b \\\"a\\\" z", Type: ARG_STRING},
		}}}},
	}, {
		// rule with abs field
		tag: `is:"rule:&F1.F2.F3"`,
		want: &Tag{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "F1.F2.F3", Type: ARG_FIELD_ABS},
		}}}},
	}, {
		// rule with rel field
		tag: `is:"rule:.f1.f2"`,
		want: &Tag{Rules: []*Rule{{Name: "rule", Args: []*Arg{
			{Value: "f1.f2", Type: ARG_FIELD_REL},
		}}}},
	}, {
		// multiple plain rules
		tag: `is:"ra,re,ri,ru,ro"`,
		want: &Tag{Rules: []*Rule{
			{Name: "ra"}, {Name: "re"}, {Name: "ri"}, {Name: "ru"}, {Name: "ro"},
		}},
	}, {
		// multiple plain rules (omit empty rules)
		tag: `is:"ra,,,re,ri,,"`,
		want: &Tag{Rules: []*Rule{
			{Name: "ra"}, {Name: "re"}, {Name: "ri"},
		}},
	}, {
		// multiple rules with arguments
		tag: `is:"ra:a:b:c,re:foo::321:,ri:1:2:3"`,
		want: &Tag{Rules: []*Rule{
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
		tag: `is:"[]ra"`,
		want: &Tag{Elem: &Tag{
			Rules: []*Rule{{Name: "ra"}},
		}},
	}, {
		// flat & nested
		tag: `is:"ra,[]re"`,
		want: &Tag{
			Rules: []*Rule{{Name: "ra"}},
			Elem: &Tag{
				Rules: []*Rule{{Name: "re"}},
			},
		},
	}, {
		// flat & nested #2
		tag: `is:"required,[]re"`,
		want: &Tag{
			Rules: []*Rule{{Name: "required"}},
			Elem: &Tag{
				Rules: []*Rule{{Name: "re"}},
			},
		},
	}, {
		// flat & nested #3
		tag: `is:"notnil,[]re"`,
		want: &Tag{
			Rules: []*Rule{{Name: "notnil"}},
			Elem: &Tag{
				Rules: []*Rule{{Name: "re"}},
			},
		},
	}, {
		// flat & nested #4
		tag: `is:"optional,[]re"`,
		want: &Tag{
			Rules: []*Rule{{Name: "optional"}},
			Elem: &Tag{
				Rules: []*Rule{{Name: "re"}},
			},
		},
	}, {
		// nested rule (elem [levels])
		tag: `is:"[][][][]ra"`,
		want: &Tag{
			Elem: &Tag{Elem: &Tag{Elem: &Tag{Elem: &Tag{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
		},
	}, {
		// nested rules (elem [levels])
		tag: `is:"[]ra,re,[]re,ri,[]ri,ru,[]ru,ro"`,
		want: &Tag{
			Elem: &Tag{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Elem: &Tag{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Elem: &Tag{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Elem: &Tag{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// nested rule (key)
		tag: `is:"[ra]"`,
		want: &Tag{Key: &Tag{
			Rules: []*Rule{{Name: "ra"}},
		}},
	}, {
		// nested rule (key [levels])
		tag: `is:"[[[[ra]]]]"`,
		want: &Tag{
			Key: &Tag{Key: &Tag{Key: &Tag{Key: &Tag{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
		},
	}, {
		// nested rules (key [levels])
		tag: `is:"[ra,re,[re,ri,[ri,ru,[ru,ro]]]]"`,
		want: &Tag{
			Key: &Tag{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Key: &Tag{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Key: &Tag{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Key: &Tag{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// nested rules (key & elem)
		tag: `is:"[ra]re"`,
		want: &Tag{Key: &Tag{
			Rules: []*Rule{{Name: "ra"}},
		}, Elem: &Tag{
			Rules: []*Rule{{Name: "re"}},
		}},
	}, {
		// nested rules (key & elem [levels])
		tag: `is:"[[[[ra]]]][][][]re"`,
		want: &Tag{
			Key: &Tag{Key: &Tag{Key: &Tag{Key: &Tag{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
			Elem: &Tag{Elem: &Tag{Elem: &Tag{Elem: &Tag{
				Rules: []*Rule{{Name: "re"}},
			}}}},
		},
	}, {
		// nested rules (key & elems [levels])
		tag: `is:"[ra,re,[re,ri,[ri,ru,[ru,ro]]]]ra,re,[]re,ri,[]ri,ru,[]ru,ro"`,
		want: &Tag{
			Key: &Tag{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Key: &Tag{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Key: &Tag{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Key: &Tag{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
			Elem: &Tag{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Elem: &Tag{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Elem: &Tag{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Elem: &Tag{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// ... with arguments and all ...
		tag: `is:"[ra,re:1:2:3,[re::\"]]\\\"[]]\":foo,ri,[ri:&MyField:::-321,ru:\"  \",[ru,ro:\"]\"]` +
			`ro:\"[\",ru]ru:foo:123::&MyOtherField:]ri,re::\"]]\\\"[]]\":foo]ra:xyz:,re:&mykey,` +
			`[la:\"]heee![\"]re,ri:,[le:a,li:b,lu:c]ri:\"foo \\\"]]]\":,ru::-abc,[c:lu,b:li,a:le]ru,ro:\"[foo]\":"`,
		want: &Tag{
			Key: &Tag{
				Rules: []*Rule{
					{Name: "ra"},
					{Name: "re", Args: []*Arg{
						{Value: "1", Type: ARG_INT},
						{Value: "2", Type: ARG_INT},
						{Value: "3", Type: ARG_INT},
					}},
				},
				Key: &Tag{
					Rules: []*Rule{
						{Name: "re", Args: []*Arg{
							{Value: "", Type: ARG_UNKNOWN},
							{Value: "]]\\\"[]]", Type: ARG_STRING},
							{Value: "foo", Type: ARG_STRING},
						}},
						{Name: "ri"},
					},
					Key: &Tag{
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
						Key: &Tag{
							Rules: []*Rule{
								{Name: "ru"},
								{Name: "ro", Args: []*Arg{
									{Value: "]", Type: ARG_STRING},
								}},
							},
						},
						Elem: &Tag{
							Rules: []*Rule{
								{Name: "ro", Args: []*Arg{
									{Value: "[", Type: ARG_STRING},
								}},
								{Name: "ru"},
							},
						},
					},
					Elem: &Tag{
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
				Elem: &Tag{
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
			Elem: &Tag{
				Rules: []*Rule{
					{Name: "ra", Args: []*Arg{
						{Value: "xyz", Type: ARG_STRING},
						{Value: "", Type: ARG_UNKNOWN},
					}},
					{Name: "re", Args: []*Arg{
						{Value: "mykey", Type: ARG_FIELD_ABS},
					}},
				},
				Key: &Tag{
					Rules: []*Rule{{Name: "la", Args: []*Arg{
						{Value: "]heee![", Type: ARG_STRING},
					}}},
				},
				Elem: &Tag{
					Rules: []*Rule{
						{Name: "re"},
						{Name: "ri", Args: []*Arg{
							{Value: "", Type: ARG_UNKNOWN},
						}},
					},
					Key: &Tag{
						Rules: []*Rule{
							{Name: "le", Args: []*Arg{{Value: "a", Type: ARG_STRING}}},
							{Name: "li", Args: []*Arg{{Value: "b", Type: ARG_STRING}}},
							{Name: "lu", Args: []*Arg{{Value: "c", Type: ARG_STRING}}},
						},
					},
					Elem: &Tag{
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
						Key: &Tag{
							Rules: []*Rule{
								{Name: "c", Args: []*Arg{{Value: "lu", Type: ARG_STRING}}},
								{Name: "b", Args: []*Arg{{Value: "li", Type: ARG_STRING}}},
								{Name: "a", Args: []*Arg{{Value: "le", Type: ARG_STRING}}},
							},
						},
						Elem: &Tag{
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
		got := parseTag(tt.tag, "is")
		if e := compare.Compare(got, tt.want); e != nil {
			t.Errorf("`%s`: %v", tt.tag, e)
		}
	}
}
