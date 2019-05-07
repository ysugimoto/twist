package cascade_config

import (
	"os"
	"reflect"

	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env"
	"github.com/go-ini/ini"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

const (
	optionNameToml string = "toml"
	optionNameIni  string = "ini"
	optionNameYaml string = "yaml"
	optionNameJson string = "json"
	optionNameEnv  string = "env"
)

type Option struct {
	name  string
	value interface{}
}

func WithToml(tomlPath string) Option {
	return Option{
		name:  optionNameToml,
		value: tomlPath,
	}
}

func WithIni(iniPath string) Option {
	return Option{
		name:  optionNameIni,
		value: iniPath,
	}
}

func WithYaml(yamlPath string) Option {
	return Option{
		name:  optionNameYaml,
		value: yamlPath,
	}
}

func WithJson(jsonPath string) Option {
	return Option{
		name:  optionNameJson,
		value: jsonPath,
	}
}

func WithEnv() Option {
	return Option{
		name:  optionNameEnv,
		value: nil,
	}
}

func Cascade(v interface{}, opts ...Option) error {
	for _, opt := range opts {
		var file string
		switch opt.name {
		case optionNameToml:
			file = opt.value.(string)
			if err := cascadeToml(file, v); err != nil {
				return errors.Wrap(err, "Failed to cascade toml")
			}
		case optionNameIni:
			file = opt.value.(string)
			if err := cascadeIni(file, v); err != nil {
				return errors.Wrap(err, "Failed to cascade ini")
			}
		case optionNameYaml:
			file = opt.value.(string)
			if err := cascadeYaml(file, v); err != nil {
				return errors.Wrap(err, "Failed to cascade yaml")
			}
		case optionNameJson:
			file = opt.value.(string)
			if err := cascadeJson(file, &v); err != nil {
				return errors.Wrap(err, "Failed to cascade json")
			}
		case optionNameEnv:
			if err := cascadeEnv(v); err != nil {
				return errors.Wrap(err, "Failed to cascade env")
			}
		}
	}
	return nil
}

func cascadeToml(file string, v interface{}) error {
	_, err := toml.DecodeFile(file, v)
	return errors.Wrap(err, "toml decode error")
}

func cascadeIni(file string, v interface{}) error {
	src, err := ini.Load(file)
	if err != nil {
		return errors.Wrap(err, "ini load error")
	}
	if err := src.MapTo(v); err != nil {
		return errors.Wrap(err, "ini map error")
	}
	return nil
}

func cascadeYaml(file string, v interface{}) error {
	if _, err := os.Stat(file); err != nil {
		return errors.Wrap(err, "yaml file not exists")
	}
	fp, err := os.Open(file)
	if err != nil {
		return errors.Wrap(err, "yaml file open error")
	}
	defer fp.Close()
	if err := yaml.NewDecoder(fp).Decode(v); err != nil {
		return errors.Wrap(err, "yaml decode error")
	}
	return nil
}

func cascadeJson(file string, v interface{}) error {
	if _, err := os.Stat(file); err != nil {
		return errors.Wrap(err, "json file not exists")
	}
	fp, err := os.Open(file)
	if err != nil {
		return errors.Wrap(err, "json file open error")
	}
	defer fp.Close()
	if err := json.NewDecoder(fp).Decode(v); err != nil {
		return errors.Wrap(err, "json decode error")
	}
	return nil
}

func cascadeEnv(v interface{}) error {
	if err := env.Parse(v); err != nil {
		return errors.Wrap(err, "env parse error")
	}
	return nil
}
