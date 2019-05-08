package twist

const (
	optionNameToml string = "toml"
	optionNameIni  string = "ini"
	optionNameYaml string = "yaml"
	optionNameJson string = "json"
	optionNameEnv  string = "env"
)

// Cascading config options
type Option struct {
	name  string
	value interface{}
}

// Will cascade from Toml config file
func WithToml(tomlPath string) Option {
	return Option{
		name:  optionNameToml,
		value: tomlPath,
	}
}

// Will cascade from ini config file
func WithIni(iniPath string) Option {
	return Option{
		name:  optionNameIni,
		value: iniPath,
	}
}

// Will cascade from ini config file
func WithYaml(yamlPath string) Option {
	return Option{
		name:  optionNameYaml,
		value: yamlPath,
	}
}

// Will cascade from JSON config file
func WithJson(jsonPath string) Option {
	return Option{
		name:  optionNameJson,
		value: jsonPath,
	}
}

// Will cascade from Environment variables
func WithEnv() Option {
	return Option{
		name:  optionNameEnv,
		value: nil,
	}
}
