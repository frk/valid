package analysis

type ContainsRulesTest1Validator struct {
	yes string `is:"required"`
	no  string
}

type ContainsRulesTest2Validator struct {
	yes struct {
		f1 string
		f2 string `is:"required"`
	}
	no struct {
		f1 string
		f2 string
	}
}

type ContainsRulesTest3Validator struct {
	yes struct {
		f1 string
		f2 *struct {
			f1 []*struct {
				f1 string
				f2 string
				f3 string
				f4 string `is:"required"`
			}
		}
	}
	no struct {
		f1 string
		f2 *struct {
			f1 []*struct {
				f1 string
				f2 string
				f3 string
				f4 string
			}
		}
	}
}

type ContainsRulesTest4Validator struct {
	yes struct {
		f1 string
		f2 *struct {
			f1 map[string][5][]*struct {
				f1 string
				f2 string
				f3 string
				f4 string `is:"required"`
			}
		}
	}
	no struct {
		f1 string
		f2 *struct {
			f1 map[string][5][]*struct {
				f1 string
				f2 string
				f3 string
				f4 string
			}
		}
	}
}

type ContainsRulesTest5Validator struct {
	yes struct {
		f1 string
		f2 *struct {
			f1 []map[struct {
				f1 string
				f2 string
				f3 string `is:"required"`
			}][5][]struct {
				f1 string
				f2 string
			}
		}
	}
	no struct {
		f1 string
		f2 *struct {
			f1 []map[struct {
				f1 string
				f2 string
				f3 string
			}][5][]struct {
				f1 string
				f2 string
			}
		}
	}
}
