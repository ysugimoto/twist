package twist_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	twist "github.com/ysugimoto/twist"
)

func TestMixCli(t *testing.T) {
	config := struct {
		Token  string `cli:"t,token"`
		Server struct {
			Host string `cli:"h,host"`
			Port int    `cli:"p,port"`
		}
	}{}
	// Actually provide from os.Args[1:]
	args := []string{"-h", "cli.localhost", "--port=9000", "--token", "token_from_cli"}
	err := twist.Mix(
		&config,
		twist.WithCli(args),
	)
	assert.NoError(t, err)
	assert.Equal(t, "token_from_cli", config.Token)
	assert.Equal(t, "cli.localhost", config.Server.Host)
	assert.Equal(t, 9000, config.Server.Port)
}

func TestMixIni(t *testing.T) {
	config := struct {
		Token  string `ini:"token"`
		Server struct {
			Host string `ini:"host"`
			Port int    `ini:"port"`
		} `ini:"server"`
	}{}
	err := twist.Mix(
		&config,
		twist.WithIni("./fixtures/example.ini"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "token_from_ini", config.Token)
	assert.Equal(t, "ini.localhost", config.Server.Host)
	assert.Equal(t, 7777, config.Server.Port)
}

func TestMixDefault(t *testing.T) {
	config := struct {
		Token  string `default:"default_token"`
		Server struct {
			Host string `default:"default.localhost"`
			Port int    `default:"60000"`
		}
	}{}
	err := twist.Mix(
		&config,
	)
	assert.NoError(t, err)
	assert.Equal(t, "default_token", config.Token)
	assert.Equal(t, "default.localhost", config.Server.Host)
	assert.Equal(t, 60000, config.Server.Port)
}

func TestMixToml(t *testing.T) {
	config := struct {
		TomlValue string `toml:"toml_value"`
		Token     string `toml:"token"`
		Server    struct {
			Host string `toml:"host"`
			Port int    `toml:"port"`
		} `toml:"server"`
	}{}
	err := twist.Mix(
		&config,
		twist.WithToml("./fixtures/example.toml"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "toml_value", config.TomlValue)
	assert.Equal(t, "token_from_toml", config.Token)
	assert.Equal(t, "toml.localhost", config.Server.Host)
	assert.Equal(t, 9999, config.Server.Port)
}

func TestMixJson(t *testing.T) {
	config := struct {
		JsonValue string `json:"json_value"`
		Token     string `json:"token"`
		Server    struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		} `json:"server"`
	}{}
	err := twist.Mix(
		&config,
		twist.WithJson("./fixtures/example.json"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "json_value", config.JsonValue)
	assert.Equal(t, "token_from_json", config.Token)
	assert.Equal(t, "json.localhost", config.Server.Host)
	assert.Equal(t, 6666, config.Server.Port)
}

func TestMixYaml(t *testing.T) {
	config := struct {
		YamlValue string `yaml:"yaml_value"`
		Token     string `yaml:"token"`
		Server    struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"server"`
	}{}
	err := twist.Mix(
		&config,
		twist.WithYaml("./fixtures/example.yaml"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "", config.YamlValue)
	assert.Equal(t, "token_from_yaml", config.Token)
	assert.Equal(t, "yaml.localhost", config.Server.Host)
	assert.Equal(t, 8888, config.Server.Port)
}

func TestMixEnv(t *testing.T) {
	os.Setenv("TOKEN", "token_from_env")
	os.Setenv("HOST", "env.localhost")
	os.Setenv("PORT", "3333")

	config := struct {
		Token  string `env:"TOKEN"`
		Server struct {
			Host string `env:"HOST"`
			Port int    `env:"PORT"`
		}
	}{}
	err := twist.Mix(
		&config,
		twist.WithEnv(),
	)
	assert.NoError(t, err)
	assert.Equal(t, "token_from_env", config.Token)
	assert.Equal(t, "env.localhost", config.Server.Host)
	assert.Equal(t, 3333, config.Server.Port)
}

func TestMixTomlAndJson(t *testing.T) {
	config := struct {
		TomlValue string `toml:"toml_value"`
		JsonValue string `json:"json_value"`
		Token     string `toml:"token"`
		Server    struct {
			Host string `json:"host"` // will be set JSON value
			Port int    `toml:"port"` // will be TOML value
		} `json:"server" toml:"server"`
	}{}
	err := twist.Mix(
		&config,
		twist.WithToml("./fixtures/example.toml"),
		twist.WithJson("./fixtures/example.json"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "toml_value", config.TomlValue)
	assert.Equal(t, "json_value", config.JsonValue)
	assert.Equal(t, "token_from_toml", config.Token)
	assert.Equal(t, "json.localhost", config.Server.Host)
	assert.Equal(t, 9999, config.Server.Port)
}

func TestMixJsonAndIni(t *testing.T) {
	config := struct {
		JsonValue string `json:"json_value"`
		Token     string `ini:"token"`
		Server    struct {
			Host string `json:"host"` // will be set JSON value
			Port int    `ini:"port"`  // will set INI value
		} `ini:"server" json:"server"`
	}{}
	err := twist.Mix(
		&config,
		twist.WithJson("./fixtures/example.json"),
		twist.WithIni("./fixtures/example.ini"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "json_value", config.JsonValue)
	assert.Equal(t, "token_from_ini", config.Token)
	assert.Equal(t, "json.localhost", config.Server.Host)
	assert.Equal(t, 7777, config.Server.Port)
}

func TestMixYamlAndEnv(t *testing.T) {
	config := struct {
		Token  string `yaml:"token"`
		Server struct {
			Host string `yaml:"host"` // will be set YAML value
			Port int    `env:"PORT"`  // will be set ENV value
		} `yaml:"server"`
		Other int `env:"PORT"`
	}{}
	err := twist.Mix(
		&config,
		twist.WithYaml("./fixtures/example.yaml"),
		twist.WithEnv(),
	)
	assert.NoError(t, err)
	assert.Equal(t, "token_from_yaml", config.Token)
	assert.Equal(t, "yaml.localhost", config.Server.Host)
	assert.Equal(t, 3333, config.Server.Port)
	assert.Equal(t, 3333, config.Other)
}

func TestMixMusltipleToml(t *testing.T) {
	config := struct {
		TomlValue string `toml:"toml_value"`
		Token     string `toml:"token"`
		Server    struct {
			Host string `toml:"host"`
			Port int    `toml:"port"`
		} `toml:"server"`
	}{}
	err := twist.Mix(
		&config,
		twist.WithToml("./fixtures/example.toml"),
		twist.WithToml("./fixtures/example.override.toml"),
	)
	assert.NoError(t, err)
	assert.Equal(t, "toml_value", config.TomlValue)
	assert.Equal(t, "token_from_toml", config.Token)
	assert.Equal(t, "toml.override.localhost", config.Server.Host)
	assert.Equal(t, 9999, config.Server.Port)
}

func TestMixFallbackDefault(t *testing.T) {
	var config struct {
		Token  string `yaml:"token" default:"value1"`
		Server struct {
			Host     string `yaml:"host" default:"localhost"`
			Port     int    `yaml:"port" default:"9000"`
			Protocol string `yaml:"protocol" default:"tcp"`
		} `yaml:"server"`
	}
	err := twist.Mix(&config, twist.WithYaml("./fixtures/example.yaml"))
	assert.NoError(t, err)
	assert.Equal(t, "token_from_yaml", config.Token)
	assert.Equal(t, "yaml.localhost", config.Server.Host)
	assert.Equal(t, 8888, config.Server.Port)
	assert.Equal(t, "tcp", config.Server.Protocol)
}

func TestMixCliWithBool(t *testing.T) {
	var config struct {
		Is bool `cli:"i,is" default:"false"`
	}
	err := twist.Mix(&config, twist.WithCli([]string{"-i"}))
	assert.NoError(t, err)
	assert.True(t, config.Is)
}
