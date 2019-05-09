# twist

Cascade and integrate config struct from various setting files and envrironment.

Supporting configuration types:

- toml file
- yaml file
- json file
- ini file
- environment variables
- commandline arguments
- default values

## Installation

```
$ go get -u github.com/ysugimoto/twist
```

Or install via pakcage manager which you're using like dep, go mod, etc...

## Getting Started

Define your configuration strcut and call `Mix()` function with cascading configurations.

Here is toml file example:

```toml
# /path/to/setting.toml

host = "localhost"
port = 8000

[service]
name = "My Service"
description = "This is my favorite service"
```

Then this package can use as following:

```Go
package main

import (
  "log"

  "github.com/ysugimoto/twist"
)

type MyConfig struct {
  Host string `toml:"host"`
  Port int    `toml:"port"`

  Service struct{
    Name        string `toml:"name"`
    Description string `toml:"description"`
  } `toml:"service"` // <- need this
}

func main() {
  var config MyConfig
  if err := twist.Mix(&config, twist.WithToml("/path/to/setting.toml")); err != nil {
    log.Fatal(err)
  }
  log.Println(config.Host)                // => localhost
  log.Println(config.Port)                // => 8000
  log.Println(config.Service.Name)        // => My Service
  log.Println(config.Service.Description) // => This is my favorite service
}
```

## Merging configuration from various kind of files and defaults

This package can atwistept kinds of config files (eg. yaml and json and env).
To merge from some configurations, call `Cascade()` with some of `WithXXX` options and make sure put tag in your struct field which you want to assign:

```Go
package main

import (
  "log"
  "os"

  "github.com/ysugimoto/twist"
)

type MyConfig struct {
  Host string `toml:"host"` // will be used from toml file
  Port int    `yaml:"port"` // will be used from yaml file

  Service struct{
    Name        string `default:"My Service"`          // set as default
    Description string `default:"My Favorite Service"` // set as default
  } `toml:"service" yaml:"service"` // <- need if you want to assign from multiple files

  Secret string `env:"SECRET"`      // will be used from envrironment variable
}

func main() {
  var config MyConfig
  if err := twist.Mix(
    &config,
    twist.WithToml("/path/to/setting.toml"),
    twist.WithToml("/path/to/setting.yaml"),
    twist.WithEnv(),
  ); err != nil {
    log.Fatal(err)
  }
  log.Println(config.Host)                // => host in toml file
  log.Println(config.Port)                // => port in yaml file 
  log.Println(config.Service.Name)        // => My Service as default
  log.Println(config.Service.Description) // => My Favorite Service as default
  log.Println(config.Secret)              // => value of os.Getenv("SECRET")
}
```

Of course you'll confuse when cascading from different file types, so usually you can use from same file types of partial configurations and integrate it, and secret values (like accessToken, etc) should be assigned from envrironment variables:

```json
# /path/to/server.json
{
  "host": "localhost",
  "port": 8000
}

# /path/to/service.json
{
  "service": {
    "name": "My Service",
    "description": "My Favorite Service"
  }
}
```

Then cascading is:

```Go
package main

import (
  "log"

  "github.com/ysugimoto/twist"
)

type MyConfig struct {
  Host string `json:"host"`
  Port int    `json:"port"`

  Service struct{
    Name        string `json:"name"`
    Description string `json:"description"`
  } `json:"service"`

  Secret  string `env:"SECRET"`    // will be used from envrironment variable
  Verbose string `cli:"v,verbose"` // will be used from commandline arguments
}

func main() {
  var config MyConfig
  if err := twist.Mix(
    &config,
    twist.WithToml("/path/to/server.json"),  // will set only Host and Port
    twist.WithToml("/path/to/service.json"), // will set only Service.Name and Service.Description
    twist.WithEnv(),
    twist.WithCli(os.Args[1:]),
  ); err != nil {
    log.Fatal(err)
  }
  log.Println(config.Host)                // => host value in server.json
  log.Println(config.Port)                // => port value in server.json
  log.Println(config.Service.Name)        // => name value in service.json
  log.Println(config.Service.Description) // => description value in service.json
  log.Println(config.Secret)              // => value of os.Getenv("SECRET")
}
```

## Correspond Struct Tags

```Go
type Config struct {
  TomlValue    string `toml:"value"`     // for toml mapping
  YamlValue    string `yaml:"value"`     // for yaml mapping
  JsonValue    string `json:"value"`     // for json mapping
  IniValue     string `ini:"value"`      // for ini mapping
  EnvValue     string `env:"ENV_NAME"`   // for env mapping
  CliValue     string `cli:"short,long"` // for cli mapping
  DefaultValue string `default:"value"`  // set as default value
}
```

`toml`, `yaml`, `json` and `ini` are used following packages:

- `github.com/BurntSushi/toml`
- `github.com/go-yaml/yaml`
- `github.com/go-ini/ini`
- `encoding/json`

`env` and `cli` is package defined.

## Author

Yoshiaki Sugimoto

## License

MIT



PR is welcome :-)
