package config

func DefaultConfig() Config {
	return Config{
		WorkDir:              String{Value: "."},
		Recursive:            Bool{Value: false},
		FileList:             StringSlice{},
		FilePatternList:      StringSlice{},
		OutNameFormat:        String{Value: "%_valid.go"},
		ValidatorNamePattern: String{Value: `^(?i:\w*Validator)$`},
		ErrorHandling: ErrorHandlingConfig{
			FieldKey: FieldKeyConfig{
				Tag:       String{Value: "json"},
				Join:      Bool{Value: true},
				Separator: String{Value: "."},
			},
		},
	}
}
