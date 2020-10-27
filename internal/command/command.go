package command

type Command struct {
	Config
}

func New(cfg Config) (*Command, error) {
	return &Command{cfg}, nil
}

func (cmd *Command) Run() error {
	return nil
}
