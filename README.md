Go Environment Helper
---

[![Code Coverage](https://codecov.io/gh/vpmv/goenv/graph/badge.svg)](https://codecov.io/gh/vpmv/goenv)

This package provides a simple way to load environment variables from .env and .ini files into your Go application, and an atomic API to read and write environment variables.

It supports basic value types such as strings, integers, booleans, floats, and slices.

# About


## Application Environment
The application uses a global environment variable for the general application environment: `ENV`.
The value can be overwritten by calling `SetEnv`. Common environment shorthands are automatically parsed. The main environment types are:
- development
- production
- testing
- staging

The package comes with some helper functions to check the global environment, e.g. `IsDevelopment()`.
Feel free to use a custom environment type, fitting your application's needs.

Note: The default environment is `development`, unless `ENV` is set on machine level.

## DotEnv

This package relies on [joho/godotenv](https://github.com/joho/godotenv) to load environment variables from .env files, adding some syntactic sugar to make it easier to overload (custom) files into your application environment.

Files overload each other in the following order:
- .env
- .env.local
- .env.<app_env>
- .env.<app_env>.local
- <custom_file>

## INI

This package also supports loading environment variables from .ini files, using [gopkg.in/ini.v1](https://gopkg.in/ini.v1) as the file processor.

Files overload each other in the following order:
- env.ini
- env.ini.local
- env.<app_env>.ini
- env.<app_env>.local.ini
- <custom_file>

You can also map your (overloaded) files directly to a struct using `MapIni()`, or access the `*ini.File` object using `LoadIniFile()`.

## Manual injection

You can inject variables using the `Set()` function. This supports all aforementioned basic data types.

# Examples

## Basic example
```go
package main

import (
    env "github.com/vpmv/goenv"	
)

func main() {
    env.LoadDotEnv(`/config/`)
    
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
	env "github.com/vpmv/goenv"
)

func main() {
	env.SetEnv(true, `app`) // set custom ENV
	
	env.LoadIni(`/config/`)
	host := env.GetString(`APP_HOST`, `localhost`)
	port := env.GetInt(`APP_PORT`, 8080)

	fuego.NewServer(
		fuego.WithAddr(fmt.Sprintf("%s:%d", host, port))
	)
}
````

### Map environment INI to struct
```go
package main

import (
	env "github.com/vpmv/goenv"
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
		Seed bool       `ini:"bool"`
	} `ini:"database"`
}

func main() {
	env.SetEnv(true, `app`) // set custom ENV
	
	config := new(Config)
	_ = env.MapIni(config, `/config/`)
}
```


## Working with slices

```go
package main

import (
	env "github.com/vpmv/goenv"
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
}
```