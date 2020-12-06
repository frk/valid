package analysis

import (
	"testing"

	"github.com/frk/compare"
)

func TestParseRuleTag(t *testing.T) {
	tests := []struct {
		tag  string
		err  error
		want *TagNode
	}{{
		tag:  ``,
		want: &TagNode{},
	}, {
		tag:  `json:"foo,omitempty" xml:">abc"`,
		want: &TagNode{},
	}, {
		tag:  `json:"foo,omitempty" is:"r1" xml:">abc"`,
		want: &TagNode{Rules: []*Rule{{Name: "r1"}}},
	}, {
		tag:  `is:"-"`,
		want: &TagNode{},
	}, {
		// single plain rule
		tag:  `is:"rule"`,
		want: &TagNode{Rules: []*Rule{{Name: "rule"}}},
	}, {
		// single rule with option
		tag: `is:"rule:opt"`,
		want: &TagNode{Rules: []*Rule{{Name: "rule", Options: []*RuleOption{
			{Value: "opt", Type: OptionTypeString},
		}}}},
	}, {
		// single rule with options
		tag: `is:"rule:opt:123:true:0.0064"`,
		want: &TagNode{Rules: []*Rule{{Name: "rule", Options: []*RuleOption{
			{Value: "opt", Type: OptionTypeString},
			{Value: "123", Type: OptionTypeInt},
			{Value: "true", Type: OptionTypeBool},
			{Value: "0.0064", Type: OptionTypeFloat},
		}}}},
	}, {
		// single rule with empty option
		tag: `is:"rule:"`,
		want: &TagNode{Rules: []*Rule{{Name: "rule", Options: []*RuleOption{
			{Value: "", Type: OptionTypeUnknown},
		}}}},
	}, {
		// single rule with empty options
		tag: `is:"rule::::"`,
		want: &TagNode{Rules: []*Rule{{Name: "rule", Options: []*RuleOption{
			{Value: "", Type: OptionTypeUnknown},
			{Value: "", Type: OptionTypeUnknown},
			{Value: "", Type: OptionTypeUnknown},
			{Value: "", Type: OptionTypeUnknown},
		}}}},
	}, {
		// single rule with empty & non-empty options
		tag: `is:"rule:opt::true:::0.0064:"`,
		want: &TagNode{Rules: []*Rule{{Name: "rule", Options: []*RuleOption{
			{Value: "opt", Type: OptionTypeString},
			{Value: "", Type: OptionTypeUnknown},
			{Value: "true", Type: OptionTypeBool},
			{Value: "", Type: OptionTypeUnknown},
			{Value: "", Type: OptionTypeUnknown},
			{Value: "0.0064", Type: OptionTypeFloat},
			{Value: "", Type: OptionTypeUnknown},
		}}}},
	}, {
		// single rule with quoted option
		tag: `is:"rule:\"opt\""`,
		want: &TagNode{Rules: []*Rule{{Name: "rule", Options: []*RuleOption{
			{Value: "opt", Type: OptionTypeString},
		}}}},
	}, {
		// single rule with quoted option
		tag: `is:"rule:\"foo \\\"bar\\\" baz\""`,
		want: &TagNode{Rules: []*Rule{{Name: "rule", Options: []*RuleOption{
			{Value: "foo \\\"bar\\\" baz", Type: OptionTypeString},
		}}}},
	}, {
		// single rule with quoted, empty, and non-empty options
		tag: `is:"rule:\"foo\":bar:\"\":123::\"b \\\"a\\\" z\""`,
		want: &TagNode{Rules: []*Rule{{Name: "rule", Options: []*RuleOption{
			{Value: "foo", Type: OptionTypeString},
			{Value: "bar", Type: OptionTypeString},
			{Value: "", Type: OptionTypeString},
			{Value: "123", Type: OptionTypeInt},
			{Value: "", Type: OptionTypeUnknown},
			{Value: "b \\\"a\\\" z", Type: OptionTypeString},
		}}}},
	}, {
		// multiple plain rules
		tag: `is:"ra,re,ri,ru,ro"`,
		want: &TagNode{Rules: []*Rule{
			{Name: "ra"}, {Name: "re"}, {Name: "ri"}, {Name: "ru"}, {Name: "ro"},
		}},
	}, {
		// multiple plain rules (omit empty rules)
		tag: `is:"ra,,,re,ri,,"`,
		want: &TagNode{Rules: []*Rule{
			{Name: "ra"}, {Name: "re"}, {Name: "ri"},
		}},
	}, {
		// multiple rules with options
		tag: `is:"ra:a:b:c,re:foo::321:,ri:1:2:3"`,
		want: &TagNode{Rules: []*Rule{
			{Name: "ra", Options: []*RuleOption{
				{Value: "a", Type: OptionTypeString},
				{Value: "b", Type: OptionTypeString},
				{Value: "c", Type: OptionTypeString},
			}},
			{Name: "re", Options: []*RuleOption{
				{Value: "foo", Type: OptionTypeString},
				{Value: "", Type: OptionTypeUnknown},
				{Value: "321", Type: OptionTypeInt},
				{Value: "", Type: OptionTypeUnknown},
			}},
			{Name: "ri", Options: []*RuleOption{
				{Value: "1", Type: OptionTypeInt},
				{Value: "2", Type: OptionTypeInt},
				{Value: "3", Type: OptionTypeInt},
			}},
		}},
	}, {
		// nested rule (elem)
		tag: `is:"[]ra"`,
		want: &TagNode{Elem: &TagNode{
			Rules: []*Rule{{Name: "ra"}},
		}},
	}, {
		// nested rule (elem [levels])
		tag: `is:"[][][][]ra"`,
		want: &TagNode{
			Elem: &TagNode{Elem: &TagNode{Elem: &TagNode{Elem: &TagNode{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
		},
	}, {
		// nested rules (elem [levels])
		tag: `is:"[]ra,re,[]re,ri,[]ri,ru,[]ru,ro"`,
		want: &TagNode{
			Elem: &TagNode{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Elem: &TagNode{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Elem: &TagNode{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Elem: &TagNode{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// nested rule (key)
		tag: `is:"[ra]"`,
		want: &TagNode{Key: &TagNode{
			Rules: []*Rule{{Name: "ra"}},
		}},
	}, {
		// nested rule (key [levels])
		tag: `is:"[[[[ra]]]]"`,
		want: &TagNode{
			Key: &TagNode{Key: &TagNode{Key: &TagNode{Key: &TagNode{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
		},
	}, {
		// nested rules (key [levels])
		tag: `is:"[ra,re,[re,ri,[ri,ru,[ru,ro]]]]"`,
		want: &TagNode{
			Key: &TagNode{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Key: &TagNode{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Key: &TagNode{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Key: &TagNode{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// nested rules (key & elem)
		tag: `is:"[ra]re"`,
		want: &TagNode{Key: &TagNode{
			Rules: []*Rule{{Name: "ra"}},
		}, Elem: &TagNode{
			Rules: []*Rule{{Name: "re"}},
		}},
	}, {
		// nested rules (key & elem [levels])
		tag: `is:"[[[[ra]]]][][][]re"`,
		want: &TagNode{
			Key: &TagNode{Key: &TagNode{Key: &TagNode{Key: &TagNode{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
			Elem: &TagNode{Elem: &TagNode{Elem: &TagNode{Elem: &TagNode{
				Rules: []*Rule{{Name: "re"}},
			}}}},
		},
	}, {
		// nested rules (key & elems [levels])
		tag: `is:"[ra,re,[re,ri,[ri,ru,[ru,ro]]]]ra,re,[]re,ri,[]ri,ru,[]ru,ro"`,
		want: &TagNode{
			Key: &TagNode{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Key: &TagNode{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Key: &TagNode{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Key: &TagNode{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
			Elem: &TagNode{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Elem: &TagNode{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Elem: &TagNode{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Elem: &TagNode{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// ... with options and all ...
		tag: `is:"[ra,re:1:2:3,[re::\"]]\\\"[]]\":foo,ri:@my_ctx,[ri:&MyField:::-321,ru:\"  \",[ru,ro:\"]\"]` +
			`ro:\"[\",ru]ru:foo:123::&MyOtherField:]ri:@my_ctx,re::\"]]\\\"[]]\":foo]ra:xyz:,re:&mykey:@MyCtx,` +
			`[la:\"]heee![\"]re,ri:,[le:a,li:b,lu:c]ri:\"foo \\\"]]]\":,ru::-abc,[c:lu,b:li,a:le]ru,ro:\"[foo]\":"`,
		want: &TagNode{
			Key: &TagNode{
				Rules: []*Rule{
					{Name: "ra"},
					{Name: "re", Options: []*RuleOption{
						{Value: "1", Type: OptionTypeInt},
						{Value: "2", Type: OptionTypeInt},
						{Value: "3", Type: OptionTypeInt},
					}},
				},
				Key: &TagNode{
					Rules: []*Rule{
						{Name: "re", Options: []*RuleOption{
							{Value: "", Type: OptionTypeUnknown},
							{Value: "]]\\\"[]]", Type: OptionTypeString},
							{Value: "foo", Type: OptionTypeString},
						}},
						{Name: "ri", Context: "my_ctx"},
					},
					Key: &TagNode{
						Rules: []*Rule{
							{Name: "ri", Options: []*RuleOption{
								{Value: "MyField", Type: OptionTypeField},
								{Value: "", Type: OptionTypeUnknown},
								{Value: "", Type: OptionTypeUnknown},
								{Value: "-321", Type: OptionTypeInt},
							}},
							{Name: "ru", Options: []*RuleOption{
								{Value: "  ", Type: OptionTypeString},
							}},
						},
						Key: &TagNode{
							Rules: []*Rule{
								{Name: "ru"},
								{Name: "ro", Options: []*RuleOption{
									{Value: "]", Type: OptionTypeString},
								}},
							},
						},
						Elem: &TagNode{
							Rules: []*Rule{
								{Name: "ro", Options: []*RuleOption{
									{Value: "[", Type: OptionTypeString},
								}},
								{Name: "ru"},
							},
						},
					},
					Elem: &TagNode{
						Rules: []*Rule{
							{Name: "ru", Options: []*RuleOption{
								{Value: "foo", Type: OptionTypeString},
								{Value: "123", Type: OptionTypeInt},
								{Value: "", Type: OptionTypeUnknown},
								{Value: "MyOtherField", Type: OptionTypeField},
								{Value: "", Type: OptionTypeUnknown},
							}},
						},
					},
				},
				Elem: &TagNode{
					Rules: []*Rule{
						{Name: "ri", Context: "my_ctx"},
						{Name: "re", Options: []*RuleOption{
							{Value: "", Type: OptionTypeUnknown},
							{Value: "]]\\\"[]]", Type: OptionTypeString},
							{Value: "foo", Type: OptionTypeString},
						}},
					},
				},
			},
			Elem: &TagNode{
				Rules: []*Rule{
					{Name: "ra", Options: []*RuleOption{
						{Value: "xyz", Type: OptionTypeString},
						{Value: "", Type: OptionTypeUnknown},
					}},
					{Name: "re", Context: "MyCtx", Options: []*RuleOption{
						{Value: "mykey", Type: OptionTypeField},
					}},
				},
				Key: &TagNode{
					Rules: []*Rule{{Name: "la", Options: []*RuleOption{
						{Value: "]heee![", Type: OptionTypeString},
					}}},
				},
				Elem: &TagNode{
					Rules: []*Rule{
						{Name: "re"},
						{Name: "ri", Options: []*RuleOption{
							{Value: "", Type: OptionTypeUnknown},
						}},
					},
					Key: &TagNode{
						Rules: []*Rule{
							{Name: "le", Options: []*RuleOption{{Value: "a", Type: OptionTypeString}}},
							{Name: "li", Options: []*RuleOption{{Value: "b", Type: OptionTypeString}}},
							{Name: "lu", Options: []*RuleOption{{Value: "c", Type: OptionTypeString}}},
						},
					},
					Elem: &TagNode{
						Rules: []*Rule{
							{Name: "ri", Options: []*RuleOption{
								{Value: "foo \\\"]]]", Type: OptionTypeString},
								{Value: "", Type: OptionTypeUnknown},
							}},
							{Name: "ru", Options: []*RuleOption{
								{Value: "", Type: OptionTypeUnknown},
								{Value: "-abc", Type: OptionTypeString},
							}},
						},
						Key: &TagNode{
							Rules: []*Rule{
								{Name: "c", Options: []*RuleOption{{Value: "lu", Type: OptionTypeString}}},
								{Name: "b", Options: []*RuleOption{{Value: "li", Type: OptionTypeString}}},
								{Name: "a", Options: []*RuleOption{{Value: "le", Type: OptionTypeString}}},
							},
						},
						Elem: &TagNode{
							Rules: []*Rule{
								{Name: "ru"},
								{Name: "ro", Options: []*RuleOption{
									{Value: "[foo]", Type: OptionTypeString},
									{Value: "", Type: OptionTypeUnknown},
								}},
							},
						},
					},
				},
			},
		},
	}}

	for _, tt := range tests {
		got, err := parseRuleTag(tt.tag)
		if e := compare.Compare(err, tt.err); e != nil {
			t.Errorf("Error: %v", e)
		}
		if e := compare.Compare(got, tt.want); e != nil {
			t.Errorf("`%s`: %v", tt.tag, e)
		}
	}
}
