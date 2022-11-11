package testdata

type Validator struct {
	User struct {
		Email   string `is:"email"`
		Address *struct {
			Zip    string `is:"zip"`
			Coords struct {
				Lat  float64 `is:"rng:-90:90"`
				Long float64 `is:"rng:-180:180"`
			}
		}
	}
}
