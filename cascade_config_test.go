package cascade_config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	cc "github.com/ysugimoto/cascade-config"
)

func TestCascadeTomlAndJson(t *testing.T) {
	config := struct {
		TomlValue string `toml:"toml_value"`
		JsonValue string `json:"json_value"`
		Token     string `toml:"token"`
		Server    struct {
			Host string `toml:"-" json:"host"` // will use JSON value
			Port int    `toml:"port" json:"-"` // will use TOML value
		}
	}{}
	err := cc.Cascade(
		&config,
		cc.WithToml("./fixtures/example.toml"),
		cc.WithJson("./fixtures/example.json"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "toml_value", config.TomlValue)
	assert.Equal(t, "json_value", config.JsonValue)
	assert.Equal(t, "token_from_json", config.Token)
	assert.Equal(t, "json.localhost", config.Server.Host)
	assert.Equal(t, 9999, config.Server.Port)
}

func TestCascadeJsonAndIni(t *testing.T) {
	config := struct {
		JsonValue string `json:"json_value"`
		Token     string `ini:"token"`
		Server    struct {
			Host string `json:"host" ini:"-"` // will use JSON value
			Port int    `ini:"port"`          // will use INI value
		} `ini:"server"`
	}{}
	err := cc.Cascade(
		&config,
		cc.WithJson("./fixtures/example.json"),
		cc.WithIni("./fixtures/example.ini"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "json_value", config.JsonValue)
	assert.Equal(t, "token_from_ini", config.Token)
	assert.Equal(t, "json.localhost", config.Server.Host)
	assert.Equal(t, 7777, config.Server.Port)
}

func TestCascadeYamlAndEnv(t *testing.T) {
	config := struct {
		Token  string `yaml:"token"`
		Server struct {
			Host string `yaml:"host"`         // will use YAML value
			Port int    `env:"PORT" yaml:"-"` // will use ENV value
		}
		Other int `env:"PORT"`
	}{}
	err := cc.Cascade(
		&config,
		cc.WithYaml("./fixtures/example.yaml"),
		cc.WithEnv(),
	)
	assert.NoError(t, err)
	assert.Equal(t, "token_from_yaml", config.Token)
	assert.Equal(t, "yaml.localhost", config.Server.Host)
	assert.Equal(t, 3333, config.Server.Port)
	assert.Equal(t, 3333, config.Other)
}
