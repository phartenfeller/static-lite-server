package main

type arguments struct {
	configFilePath string
}

func parseArguments(args []string) (arguments) {
	if len(args) < 1 {
		panic("No config file path provided")
	}

	return arguments{args[0]}
}
