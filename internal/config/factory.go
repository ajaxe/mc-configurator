package config

import "log"

func NewConfigGenerator(opts GeneratorOptions) *ConfigGenerator {
	if opts.Source == "" {
		log.Fatalf("source file path is required")
	}
	return &ConfigGenerator{options: opts}
}

func NewGeneratorOptions(args []string) GeneratorOptions {
	src := args[1]
	if src == "" {
		log.Fatalf("error: source server.properties file path is required")
	}

	dest := "./server.properties"
	if len(args) > 2 {
		dest = args[2]
	}
	return GeneratorOptions{
		Source: args[0],
		Destination: dest,
	}
}
