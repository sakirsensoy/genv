# genv

[![GitHub Actions](https://github.com/sakirsensoy/genv/actions/workflows/go.yml/badge.svg)](https://github.com/sakirsensoy/genv/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/sakirsensoy/genv)](https://goreportcard.com/report/github.com/sakirsensoy/genv) [![Coverage Status](https://coveralls.io/repos/github/sakirsensoy/genv/badge.svg?branch=master)](https://coveralls.io/github/sakirsensoy/genv?branch=master) [![GoDoc](https://godoc.org/github.com/sakirsensoy/genv?status.svg)](https://godoc.org/github.com/sakirsensoy/genv)

Genv is a library for Go (golang) that makes it easy to read and use environment variables in your projects. It also allows environment variables to be loaded from the .env file.

The dotenv parser is fully compatible with [motdotla/dotenv](https://github.com/motdotla/dotenv), the de facto standard dotenv implementation used by Node.js projects.

## Installation

```shell
go get github.com/sakirsensoy/genv
```

## Usage

Create a `.env` file in the root directory of your project and enter the environment variables you want to use:

```shell
# .env
APP_HOST=localhost
APP_PORT=1234
APP_DEBUG=true
```

### Supported .env Syntax

The parser supports the full [motdotla/dotenv](https://github.com/motdotla/dotenv) specification:

```shell
# Comments start with #
BASIC=basic

# Empty values
EMPTY=

# Whitespace is trimmed from unquoted values
SPACED=   some value   

# Single quotes (raw, no escape processing)
SINGLE='quoted value'

# Double quotes (expands \n and \r)
DOUBLE="line1\nline2"

# Backticks (raw, no escape processing)
BACKTICK=`raw value`

# Inline comments (for unquoted values)
WITH_COMMENT=value # this is a comment

# Export prefix is supported
export EXPORTED=value

# Colon separator is also supported
COLON: value

# JSON values work correctly
JSON={"foo": "bar", "baz": 123}
```

In the meantime, it is optional to use the `.env` file. You can also send environment variables to your project in classic ways:

```shell
APP_HOST=localhost ./myproject
```

Rather than using your environment variables directly in your project, it is better to map and match them with a struct.  Below you can see how we get our application parameters from environment variables:

```go
// config/config.go
package config

import "github.com/sakirsensoy/genv"

type appConfig struct {
	Host string
	Port int
	Debug bool
}

var  App = &appConfig{
	Host: genv.Key("APP_HOST").String(),
	Port: genv.Key("APP_PORT").Default(8080).Int(),
	Debug: genv.Key("APP_DEBUG").Default(false).Bool(),
}
```

In `main.go` we can directly use the parameters defined in `config.go`. The `.env` file is automatically loaded on the first `genv.Key()` call:

```go
// main.go
package main

import (
	"fmt"
	"myproject/config"
)

func main() {
	fmt.Println(config.App.Host) // localhost
	fmt.Println(config.App.Port) // 1234
	fmt.Println(config.App.Debug) // true
}
```

### Auto-load Configuration

By default, genv automatically loads the `.env` file on the first `Key()` call. You can customize this behavior:

```go
package main

import "github.com/sakirsensoy/genv"

func init() {
	// Disable auto-loading (must be called before any Key() calls)
	genv.DisableAutoLoad()
	
	// Or change the default path (must be called before any Key() calls)
	genv.SetAutoLoadPath(".env.local")
}
```

> **Note:** The `github.com/sakirsensoy/genv/dotenv/autoload` package is still available for backward compatibility, but it's no longer necessary since auto-loading is now enabled by default.

### Accessing Environment Variables

Genv provides an easy-to-use API for accessing environment variables.

> First we specify the key to the variable want to access

```go
var env = genv.Key("MY_VARIABLE")
```

> Define default value (optional)

```go
env = env.Default("default_value")
```

> Finally, we specify the type of the environment variable and pass its contents to another variable

```go
var myVariable = env.String()
```

### Supported Variable Types

Genv provides support for the following data types:

 - `String()`: Returns data of **String** type
 - `Int()`: Returns data of **Int32** type
 - `Float()`: Returns data of **Float64** type
 - `Bool()`: Returns data of **Bool** type

For other types, you can use type conversion:

```go
var stringValue = genv.Key("KEY").String()
var byteArrayValue = []byte(stringValue)
```

### Updating Environment Variables

You can update environment variables at runtime using the `Update()` method:

```go
genv.Key("APP_DEBUG").Update(false)
genv.Key("APP_PORT").Update(8080)
genv.Key("APP_HOST").Update("127.0.0.1")
```

The `Update()` method updates both the internal cache and the actual environment variable via `os.Setenv()`.

### Thread Safety

Genv is safe for concurrent use. All read and write operations are protected by mutexes, making it suitable for use in multi-goroutine applications.

## Contributing

Thanks in advance for your contributions :) I would appreciate it if you make sure that the API remains simple when developing.

*code changes without tests will not be accepted*

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## License

© Şakir Şensoy, 2019 ~ time.Now()

Released under the [MIT License](https://github.com/sakirsensoy/genv/blob/master/LICENSE)
