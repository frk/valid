package command

type Config struct {
	// TODO
}

var DefaultConfig = Config{
	// TODO
}

// ParseFlags unmarshals the cli flags into the receiver.
func (c *Config) ParseFlags() {
	// TODO
}

// ParseFile looks for a gosql config file in the git project's root of the receiver's
// working directory, if it finds such a file it will then unmarshal it into the receiver.
func (c *Config) ParseFile() error {
	// TODO
	return nil
}

func (c *Config) FileFilterFunc() (filter func(filePath string) bool) {
	// TODO
	return nil
}

// validate checks the config for errors and updates some of the values to a more "normalized" format.
func (c *Config) validate() (err error) {
	return nil
}
