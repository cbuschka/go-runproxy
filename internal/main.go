package internal

func Run(commandLine []string) error {
	config, err := NewConfig(commandLine)
	if err != nil {
		return err
	}

	server, err := NewServer(config)
	if err != nil {
		return err
	}

	err = server.Run()
	return err
}
