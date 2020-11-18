package analysis

import (
	"testing"

	"github.com/frk/compare"
)

func TestParseRuleTag(t *testing.T) {
	tests := []struct {
		tag  string
		err  error
		want *RuleTag
	}{{
		tag:  ``,
		want: &RuleTag{},
	}, {
		tag:  `json:"foo,omitempty" xml:">abc"`,
		want: &RuleTag{},
	}, {
		tag:  `json:"foo,omitempty" is:"r1" xml:">abc"`,
		want: &RuleTag{Rules: []*Rule{{Name: "r1"}}},
	}, {
		tag:  `is:"-"`,
		want: &RuleTag{},
	}, {
		// single plain rule
		tag:  `is:"rule"`,
		want: &RuleTag{Rules: []*Rule{{Name: "rule"}}},
	}, {
		// single rule with arg
		tag: `is:"rule:arg"`,
		want: &RuleTag{Rules: []*Rule{{Name: "rule", Args: []*RuleArg{
			{Value: "arg", Type: ArgTypeString},
		}}}},
	}, {
		// single rule with args
		tag: `is:"rule:arg:123:true:0.0064"`,
		want: &RuleTag{Rules: []*Rule{{Name: "rule", Args: []*RuleArg{
			{Value: "arg", Type: ArgTypeString},
			{Value: "123", Type: ArgTypeInt},
			{Value: "true", Type: ArgTypeBool},
			{Value: "0.0064", Type: ArgTypeFloat},
		}}}},
	}, {
		// single rule with empty arg
		tag: `is:"rule:"`,
		want: &RuleTag{Rules: []*Rule{{Name: "rule", Args: []*RuleArg{
			{Value: "", Type: ArgTypeUnknown},
		}}}},
	}, {
		// single rule with empty args
		tag: `is:"rule::::"`,
		want: &RuleTag{Rules: []*Rule{{Name: "rule", Args: []*RuleArg{
			{Value: "", Type: ArgTypeUnknown},
			{Value: "", Type: ArgTypeUnknown},
			{Value: "", Type: ArgTypeUnknown},
			{Value: "", Type: ArgTypeUnknown},
		}}}},
	}, {
		// single rule with empty & non-empty args
		tag: `is:"rule:arg::true:::0.0064:"`,
		want: &RuleTag{Rules: []*Rule{{Name: "rule", Args: []*RuleArg{
			{Value: "arg", Type: ArgTypeString},
			{Value: "", Type: ArgTypeUnknown},
			{Value: "true", Type: ArgTypeBool},
			{Value: "", Type: ArgTypeUnknown},
			{Value: "", Type: ArgTypeUnknown},
			{Value: "0.0064", Type: ArgTypeFloat},
			{Value: "", Type: ArgTypeUnknown},
		}}}},
	}, {
		// single rule with quoted arg
		tag: `is:"rule:\"arg\""`,
		want: &RuleTag{Rules: []*Rule{{Name: "rule", Args: []*RuleArg{
			{Value: "arg", Type: ArgTypeString},
		}}}},
	}, {
		// single rule with quoted arg
		tag: `is:"rule:\"foo \\\"bar\\\" baz\""`,
		want: &RuleTag{Rules: []*Rule{{Name: "rule", Args: []*RuleArg{
			{Value: "foo \\\"bar\\\" baz", Type: ArgTypeString},
		}}}},
	}, {
		// single rule with quoted, empty, and non-empty args
		tag: `is:"rule:\"foo\":bar:\"\":123::\"b \\\"a\\\" z\""`,
		want: &RuleTag{Rules: []*Rule{{Name: "rule", Args: []*RuleArg{
			{Value: "foo", Type: ArgTypeString},
			{Value: "bar", Type: ArgTypeString},
			{Value: "", Type: ArgTypeString},
			{Value: "123", Type: ArgTypeInt},
			{Value: "", Type: ArgTypeUnknown},
			{Value: "b \\\"a\\\" z", Type: ArgTypeString},
		}}}},
	}, {
		// multiple plain rules
		tag: `is:"ra,re,ri,ru,ro"`,
		want: &RuleTag{Rules: []*Rule{
			{Name: "ra"}, {Name: "re"}, {Name: "ri"}, {Name: "ru"}, {Name: "ro"},
		}},
	}, {
		// multiple plain rules (omit empty rules)
		tag: `is:"ra,,,re,ri,,"`,
		want: &RuleTag{Rules: []*Rule{
			{Name: "ra"}, {Name: "re"}, {Name: "ri"},
		}},
	}, {
		// multiple rules with args
		tag: `is:"ra:a:b:c,re:foo::321:,ri:1:2:3"`,
		want: &RuleTag{Rules: []*Rule{
			{Name: "ra", Args: []*RuleArg{
				{ArgTypeString, "a"},
				{ArgTypeString, "b"},
				{ArgTypeString, "c"},
			}},
			{Name: "re", Args: []*RuleArg{
				{ArgTypeString, "foo"},
				{ArgTypeUnknown, ""},
				{ArgTypeInt, "321"},
				{ArgTypeUnknown, ""},
			}},
			{Name: "ri", Args: []*RuleArg{
				{ArgTypeInt, "1"},
				{ArgTypeInt, "2"},
				{ArgTypeInt, "3"},
			}},
		}},
	}, {
		// nested rule (elem)
		tag: `is:"[]ra"`,
		want: &RuleTag{Elem: &RuleTag{
			Rules: []*Rule{{Name: "ra"}},
		}},
	}, {
		// nested rule (elem [levels])
		tag: `is:"[][][][]ra"`,
		want: &RuleTag{
			Elem: &RuleTag{Elem: &RuleTag{Elem: &RuleTag{Elem: &RuleTag{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
		},
	}, {
		// nested rules (elem [levels])
		tag: `is:"[]ra,re,[]re,ri,[]ri,ru,[]ru,ro"`,
		want: &RuleTag{
			Elem: &RuleTag{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Elem: &RuleTag{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Elem: &RuleTag{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Elem: &RuleTag{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// nested rule (key)
		tag: `is:"[ra]"`,
		want: &RuleTag{Key: &RuleTag{
			Rules: []*Rule{{Name: "ra"}},
		}},
	}, {
		// nested rule (key [levels])
		tag: `is:"[[[[ra]]]]"`,
		want: &RuleTag{
			Key: &RuleTag{Key: &RuleTag{Key: &RuleTag{Key: &RuleTag{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
		},
	}, {
		// nested rules (key [levels])
		tag: `is:"[ra,re,[re,ri,[ri,ru,[ru,ro]]]]"`,
		want: &RuleTag{
			Key: &RuleTag{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Key: &RuleTag{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Key: &RuleTag{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Key: &RuleTag{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// nested rules (key & elem)
		tag: `is:"[ra]re"`,
		want: &RuleTag{Key: &RuleTag{
			Rules: []*Rule{{Name: "ra"}},
		}, Elem: &RuleTag{
			Rules: []*Rule{{Name: "re"}},
		}},
	}, {
		// nested rules (key & elem [levels])
		tag: `is:"[[[[ra]]]][][][]re"`,
		want: &RuleTag{
			Key: &RuleTag{Key: &RuleTag{Key: &RuleTag{Key: &RuleTag{
				Rules: []*Rule{{Name: "ra"}},
			}}}},
			Elem: &RuleTag{Elem: &RuleTag{Elem: &RuleTag{Elem: &RuleTag{
				Rules: []*Rule{{Name: "re"}},
			}}}},
		},
	}, {
		// nested rules (key & elems [levels])
		tag: `is:"[ra,re,[re,ri,[ri,ru,[ru,ro]]]]ra,re,[]re,ri,[]ri,ru,[]ru,ro"`,
		want: &RuleTag{
			Key: &RuleTag{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Key: &RuleTag{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Key: &RuleTag{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Key: &RuleTag{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
			Elem: &RuleTag{
				Rules: []*Rule{{Name: "ra"}, {Name: "re"}},
				Elem: &RuleTag{
					Rules: []*Rule{{Name: "re"}, {Name: "ri"}},
					Elem: &RuleTag{
						Rules: []*Rule{{Name: "ri"}, {Name: "ru"}},
						Elem: &RuleTag{
							Rules: []*Rule{{Name: "ru"}, {Name: "ro"}},
						},
					},
				},
			},
		},
	}, {
		// ... with arguments and all ...
		tag: `is:"[ra,re:1:2:3,[re::\"]]\\\"[]]\":foo,ri:@my_ctx,[ri:&MyField:::-321,ru:\"  \",[ru,ro:\"]\"]` +
			`ro:\"[\",ru]ru:foo:123::&MyOtherField:]ri:@my_ctx,re::\"]]\\\"[]]\":foo]ra:xyz:,re:&mykey:@MyCtx,` +
			`[la:\"]heee![\"]re,ri:,[le:a,li:b,lu:c]ri:\"foo \\\"]]]\":,ru::-abc,[c:lu,b:li,a:le]ru,ro:\"[foo]\":"`,
		want: &RuleTag{
			Key: &RuleTag{
				Rules: []*Rule{
					{Name: "ra"},
					{Name: "re", Args: []*RuleArg{
						{Value: "1", Type: ArgTypeInt},
						{Value: "2", Type: ArgTypeInt},
						{Value: "3", Type: ArgTypeInt},
					}},
				},
				Key: &RuleTag{
					Rules: []*Rule{
						{Name: "re", Args: []*RuleArg{
							{Value: "", Type: ArgTypeUnknown},
							{Value: "]]\\\"[]]", Type: ArgTypeString},
							{Value: "foo", Type: ArgTypeString},
						}},
						{Name: "ri", Context: "my_ctx"},
					},
					Key: &RuleTag{
						Rules: []*Rule{
							{Name: "ri", Args: []*RuleArg{
								{Value: "MyField", Type: ArgTypeField},
								{Value: "", Type: ArgTypeUnknown},
								{Value: "", Type: ArgTypeUnknown},
								{Value: "-321", Type: ArgTypeInt},
							}},
							{Name: "ru", Args: []*RuleArg{
								{Value: "  ", Type: ArgTypeString},
							}},
						},
						Key: &RuleTag{
							Rules: []*Rule{
								{Name: "ru"},
								{Name: "ro", Args: []*RuleArg{
									{Value: "]", Type: ArgTypeString},
								}},
							},
						},
						Elem: &RuleTag{
							Rules: []*Rule{
								{Name: "ro", Args: []*RuleArg{
									{Value: "[", Type: ArgTypeString},
								}},
								{Name: "ru"},
							},
						},
					},
					Elem: &RuleTag{
						Rules: []*Rule{
							{Name: "ru", Args: []*RuleArg{
								{Value: "foo", Type: ArgTypeString},
								{Value: "123", Type: ArgTypeInt},
								{Value: "", Type: ArgTypeUnknown},
								{Value: "MyOtherField", Type: ArgTypeField},
								{Value: "", Type: ArgTypeUnknown},
							}},
						},
					},
				},
				Elem: &RuleTag{
					Rules: []*Rule{
						{Name: "ri", Context: "my_ctx"},
						{Name: "re", Args: []*RuleArg{
							{Value: "", Type: ArgTypeUnknown},
							{Value: "]]\\\"[]]", Type: ArgTypeString},
							{Value: "foo", Type: ArgTypeString},
						}},
					},
				},
			},
			Elem: &RuleTag{
				Rules: []*Rule{
					{Name: "ra", Args: []*RuleArg{
						{Value: "xyz", Type: ArgTypeString},
						{Value: "", Type: ArgTypeUnknown},
					}},
					{Name: "re", Context: "MyCtx", Args: []*RuleArg{
						{Value: "mykey", Type: ArgTypeField},
					}},
				},
				Key: &RuleTag{
					Rules: []*Rule{{Name: "la", Args: []*RuleArg{
						{Value: "]heee![", Type: ArgTypeString},
					}}},
				},
				Elem: &RuleTag{
					Rules: []*Rule{
						{Name: "re"},
						{Name: "ri", Args: []*RuleArg{
							{Value: "", Type: ArgTypeUnknown},
						}},
					},
					Key: &RuleTag{
						Rules: []*Rule{
							{Name: "le", Args: []*RuleArg{{Value: "a", Type: ArgTypeString}}},
							{Name: "li", Args: []*RuleArg{{Value: "b", Type: ArgTypeString}}},
							{Name: "lu", Args: []*RuleArg{{Value: "c", Type: ArgTypeString}}},
						},
					},
					Elem: &RuleTag{
						Rules: []*Rule{
							{Name: "ri", Args: []*RuleArg{
								{Value: "foo \\\"]]]", Type: ArgTypeString},
								{Value: "", Type: ArgTypeUnknown},
							}},
							{Name: "ru", Args: []*RuleArg{
								{Value: "", Type: ArgTypeUnknown},
								{Value: "-abc", Type: ArgTypeString},
							}},
						},
						Key: &RuleTag{
							Rules: []*Rule{
								{Name: "c", Args: []*RuleArg{{Value: "lu", Type: ArgTypeString}}},
								{Name: "b", Args: []*RuleArg{{Value: "li", Type: ArgTypeString}}},
								{Name: "a", Args: []*RuleArg{{Value: "le", Type: ArgTypeString}}},
							},
						},
						Elem: &RuleTag{
							Rules: []*Rule{
								{Name: "ru"},
								{Name: "ro", Args: []*RuleArg{
									{Value: "[foo]", Type: ArgTypeString},
									{Value: "", Type: ArgTypeUnknown},
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
