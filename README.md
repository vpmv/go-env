Go Environment Helper
===

[![Go](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/vpmv/go-env/main/.github/badges/go.json)](https://github.com/vpmv/go-env)
[![Tests](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/vpmv/go-env/main/.github/badges/tests.json)](https://github.com/vpmv/go-env/actions)
[![Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/vpmv/go-env/main/.github/badges/coverage.json)](https://github.com/vpmv/go-env/actions)
[![Code Coverage](https://codecov.io/gh/vpmv/go-env/graph/badge.svg)](https://codecov.io/gh/vpmv/go-env)
[![Lint](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/vpmv/go-env/main/.github/badges/lint.json)](https://github.com/vpmv/go-env/actions)


This package provides a simple way to load environment variables from a variety of files - with default value support - and an atomic API to read and write environment variables. Files supported:
- .env
- .ini 
- .yaml

# Usage

## Global application environment

The package uses a global environment variable for the general application environment; default `ENV`. This can be overridden with the variable of your choosing, e.g.: `env.GlobalEnv = "APP_ENV".

The global environment can be set by calling `SetEnv`. Common environment shorthands are automatically parsed. The main environment types are:
- development
- production
- testing
- staging

The package comes with some helper functions to check the global environment, e.g. `IsDevelopment()`, `IsStaging()` et cetera.

Feel free to use a custom environment type, fitting your application's needs.

> [!Note] 
> The default environment is `development`, unless `ENV` (<= `env.GlobalVar`) is set on machine level.
> E.g.: `docker run -e ENV=production my-application`

## Manual injection

You can inject variables using the `Set()` function. This supports all stringable data types such as strings, integers, booleans, floats, and slices.

## DotEnv

Import: `github.com/vpmv/go-env/dotenv`

This library supports loading environment variables from .env files, using [joho/godotenv](https://github.com/joho/godotenv) as the file processor, making it easier to overload (custom) files into your application environment.

Files overload each other in the following order:
- .env
- .env.local
- .env.<app_env>
- .env.<app_env>.local
- <custom_file>

## INI

Import: `github.com/vpmv/go-env/ini`

This library supports loading environment variables from .ini files, using [gopkg.in/ini.v1](https://gopkg.in/ini.v1) as the file processor.

Files overload each other in the following order:
- env.ini
- env.local.ini
- env.<app_env>.ini
- env.<app_env>.local.ini
- <custom_file>

You can also map your (overloaded) files directly to a struct using `Map()`, or access the `*ini.File` object using `LoadFile()`.

## YAML

Import: `github.com/vpmv/go-env/yaml`

This library supports loading environment variables from .yaml files, using [goccy/go-yaml](https://github.com/goccy/go-yaml) as the file processor, with the support of YAML anchors.

Files overload each other in the following order:
- env.yaml
- env.local.yaml
- env.<app_env>.yaml
- env.<app_env>.local.yaml
- <custom_file>
  
You can also map your (overloaded) files directly to a struct using `Map()`, or `MapWithReferences()`.

### Foreign Anchors

> [!NOTE]
> The package [goccy/go-yaml](github.com/guccy/go-yaml) supports loading anchors from other files. Although it is designed to support `reference directories`, we explicitly only support `reference files`, because `reference directories` will evaluate all files consecutively prior to parsing. If an unknown reference is found, it'll stop execution.  

If you want to use YAML anchors defined in different files, you can supply paths to these `reference files`. This allows you to easily reuse/overwrite blocks of configuration.

Related functions are: 
- `yaml.LoadWithReferences` - loading contents into the environment
- `yaml.MapWithReferences` - mapping contents to an interface

All files are expected to be relative to the basedir.

# Examples

## Basic example
```go
package main

import (
    "github.com/vpmv/go-env"	
    "github.com/vpmv/go-env/dotenv"	
)

func main() {
    dotenv.Load(`/config/`)
    
    if env.IsDevelopment() {
        env.Set(`SEED_DB`, true)
    }
    
    database := env.MustString(`DATABASE_URL`) // will panic if unset
    databasePort := env.GetInt(`DATABASE_PORT`, 3306) // will return default value (3306) if unset
    // ...
	
	// check if variable exists
	if env.Has(`DATABASE_SEED`) {
		// ...
    }
}
```

## INI files

### Parse INI to environment

```go
package main

import (
	"fmt"

	"github.com/go-fuego/fuego"
	"github.com/vpmv/go-env"
	"github.com/vpmv/go-env/ini"
)

func main() {
	env.SetEnv(true, `app`) // set custom ENV
	
	ini.Load(`/config/`)
	host := env.GetString(`APP_HOST`, `localhost`)
	port := env.GetInt(`APP_PORT`, 8080)

	server := fuego.NewServer(
		fuego.WithAddr(fmt.Sprintf("%s:%d", host, port))
	)
	fuego.Use(server, cors.New(cors.Options{
		AllowedOrigins: env.GetStringSlice(`APP_ALLOWED_ORIGINS`, []string{`*`}),
	}))
}
````

### Map environment INI to struct
```go
package main

import (
	"github.com/vpmv/go-env"
	"github.com/vpmv/go-env/ini"
)

type Config struct {
    App struct {
        Host string `ini:"host"`
        Port int    `ini:"port"`
    } `ini:"app"`
    Meta struct {
        JWTSecret string `ini:"jwt"`
        TTL       int    `ini:"ttl"`
    } `ini:"app.meta"`
    Database struct {
        Host     string `ini:"host"`
        Port     int    `ini:"port"`
        User     string `ini:"user"`
        Password string `ini:"password"`
        Seed     bool   `ini:"bool"`
    } `ini:"database"`
}

func main() {
	env.SetEnv(true, `emergency`) // set custom ENV
	
	config := new(Config)
	_ = ini.Map(config, `/config/`)
}
```
## YAML files

### Parse YAML to environment
```go
package main

import (
	"fmt"

	"github.com/go-fuego/fuego"
	"github.com/vpmv/go-env"
	"github.com/vpmv/go-env/yaml"
)

func main() {
	err := yaml.Load(`/config`)
	if err != nil {
		// ...
    }
	host := env.GetString(`APP_HOST`, `localhost`)
	port := env.GetInt(`APP_PORT`, 8080)
	
	mysql := storage.NewClient(
		env.MustString(`DATABASE[0]_HOST`),
		env.MustInt(`DATABASE[0]_PORT`),
    )
	redis := storage.NewClient(
		env.MustString(`DATABASE[1]_HOST`),
		env.MustInt(`DATABASE[1]_PORT`),
    )
	
	server := fuego.NewServer(
		fuego.WithAddr(fmt.Sprintf("%s:%d", host, port))
	)
	fuego.Use(server, cors.New(cors.Options{
		AllowedOrigins: env.GetStringSlice(`APP_ALLOWED_ORIGINS`, []string{`*`}),
	}))
}
````

### Map environment YAML to struct
```go
package main

import (
	"github.com/vpmv/go-env"
	"github.com/vpmv/go-env/yaml"
)

type Config struct {
    App struct {
        Host     string   `yaml:"host"`
        Port     int     `yaml:"port"`
        Origins  []string `yaml:"allowed_origins"`
	} `yaml:"app"`
    Database    struct {
        Host  string  `yaml:"host"`
        Port  int     `yaml:"port"`
        User  string  `yaml:"username"`
        Pass  string  `yaml:"password"`
    } `yaml:"database"`
}

func main() {
	config := new(Config)
	err := yaml.MapWithReferences(config, `/app/config`, []string{`defaults.yaml`})
}
```


## Working with slices

```go
package main

import (
	"github.com/vpmv/go-env"
)

func main() {
	env.Set(`ALLOWED_ORIGINS`, []string{`10.0.0.0/8`, `192.168.0.0/16`})
	//  ALLOWED_ORIGINS=10.0.0.0/8;192.168.0.0/16
	
	// specify a custom delimiter for slice-types
	// NOTE: the delimiter remains in memory until changed
	env.SetDelimiter(`,`) 
	env.Set(`NUMBERS`, []int{101,202,303})
	//  NUMBERS=101,202,303
	
	
	env.SetDelimiter(`;`) // reset to default delimiter 
	origins := env.GetStringSlice(`ALLOWED_ORIGINS`, []string{`*`})
	// []string{`10.0.0.0/8`, `192.168.0.0/16`}
}
```