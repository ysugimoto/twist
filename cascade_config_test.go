package cascade_config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	cc "github.com/ysugimoto/cascade-config"
)

func TestCascadeIni(t *testing.T) {
	config := struct {
		Token  string `ini:"token"`
		Server struct {
			Host string `ini:"host"`
			Port int    `ini:"port"`
		} `ini:"server"`
	}{}
	err := cc.Cascade(
		&config,
		cc.WithIni("./fixtures/example.ini"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "token_from_ini", config.Token)
	assert.Equal(t, "ini.localhost", config.Server.Host)
	assert.Equal(t, 7777, config.Server.Port)
}

func TestCascadeDefault(t *testing.T) {
	config := struct {
		Token  string `default:"default_token"`
		Server struct {
			Host string `default:"default.localhost"`
			Port int    `default:"60000"`
		}
	}{}
	err := cc.Cascade(
		&config,
	)
	assert.NoError(t, err)
	assert.Equal(t, "default_token", config.Token)
	assert.Equal(t, "default.localhost", config.Server.Host)
	assert.Equal(t, 60000, config.Server.Port)
}

func TestCascadeToml(t *testing.T) {
	config := struct {
		TomlValue string `toml:"toml_value"`
		Token     string `toml:"token"`
		Server    struct {
			Host string `toml:"host"`
			Port int    `toml:"port"`
		} `toml:"server"`
	}{}
	err := cc.Cascade(
		&config,
		cc.WithToml("./fixtures/example.toml"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "toml_value", config.TomlValue)
	assert.Equal(t, "token_from_toml", config.Token)
	assert.Equal(t, "toml.localhost", config.Server.Host)
	assert.Equal(t, 9999, config.Server.Port)
}

func TestCascadeJson(t *testing.T) {
	config := struct {
		JsonValue string `json:"json_value"`
		Token     string `json:"token"`
		Server    struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		} `json:"server"`
	}{}
	err := cc.Cascade(
		&config,
		cc.WithJson("./fixtures/example.json"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "json_value", config.JsonValue)
	assert.Equal(t, "token_from_json", config.Token)
	assert.Equal(t, "json.localhost", config.Server.Host)
	assert.Equal(t, 6666, config.Server.Port)
}

func TestCascadeYaml(t *testing.T) {
	config := struct {
		YamlValue string `yaml:"yaml_value"`
		Token     string `yaml:"token"`
		Server    struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"server"`
	}{}
	err := cc.Cascade(
		&config,
		cc.WithYaml("./fixtures/example.yaml"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "", config.YamlValue)
	assert.Equal(t, "token_from_yaml", config.Token)
	assert.Equal(t, "yaml.localhost", config.Server.Host)
	assert.Equal(t, 8888, config.Server.Port)
}

func TestCascadeEnv(t *testing.T) {
	config := struct {
		Token  string `env:"TOKEN"`
		Server struct {
			Host string `env:"HOST"`
			Port int    `env:"PORT"`
		}
	}{}
	err := cc.Cascade(
		&config,
		cc.WithEnv(),
	)
	assert.NoError(t, err)
	assert.Equal(t, "token_from_env", config.Token)
	assert.Equal(t, "env.localhost", config.Server.Host)
	assert.Equal(t, 3333, config.Server.Port)
}

func TestCascadeTomlAndJson(t *testing.T) {
	config := struct {
		TomlValue string `toml:"toml_value"`
		JsonValue string `json:"json_value"`
		Token     string `toml:"token"`
		Server    struct {
			Host string `json:"host"` // will be set JSON value
			Port int    `toml:"port"` // will be TOML value
		} `json:"server" toml:"server"`
	}{}
	err := cc.Cascade(
		&config,
		cc.WithToml("./fixtures/example.toml"),
		cc.WithJson("./fixtures/example.json"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "toml_value", config.TomlValue)
	assert.Equal(t, "json_value", config.JsonValue)
	assert.Equal(t, "token_from_toml", config.Token)
	assert.Equal(t, "json.localhost", config.Server.Host)
	assert.Equal(t, 9999, config.Server.Port)
}

func TestCascadeJsonAndIni(t *testing.T) {
	config := struct {
		JsonValue string `json:"json_value"`
		Token     string `ini:"token"`
		Server    struct {
			Host string `json:"host"` // will be set JSON value
			Port int    `ini:"port"`  // will set INI value
		} `ini:"server" json:"server"`
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
			Host string `yaml:"host"` // will be set YAML value
			Port int    `env:"PORT"`  // will be set ENV value
		} `yaml:"server"`
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

func TestCascadeMusltipleToml(t *testing.T) {
	config := struct {
		TomlValue string `toml:"toml_value"`
		Token     string `toml:"token"`
		Server    struct {
			Host string `toml:"host"`
			Port int    `toml:"port"`
		} `toml:"server"`
	}{}
	err := cc.Cascade(
		&config,
		cc.WithToml("./fixtures/example.toml"),
		cc.WithToml("./fixtures/example.override.toml"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "toml_value", config.TomlValue)
	assert.Equal(t, "token_from_toml", config.Token)
	assert.Equal(t, "toml.override.localhost", config.Server.Host)
	assert.Equal(t, 9999, config.Server.Port)
}
