# Genv [![Build Status](https://travis-ci.org/sakirsensoy/genv.svg?branch=master)](https://travis-ci.org/sakirsensoy/genv) [![GoDoc](https://godoc.org/github.com/sakirsensoy/genv?status.svg)](https://godoc.org/github.com/sakirsensoy/genv)

Genv is a library for Go (golang) that makes it easy to read and use environment variables in your projects. It also allows environment variables to be loaded from the .env file.

## Installation

```shell
go get -u github.com/sakirsensoy/genv
```

## Usage

Create a `.env` file in the root directory of your project and enter the environment variables you want to use:

```shell
# .env
APP_HOST=localhost
APP_PORT=1234
APP_DEBUG=true
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

In `main.go` we first include the package that allows you to automatically load the environment variables from the .env file. Then we can include and use the parameters defined in `config.go`:

```go
// main.go
package main

import (
	_ "github.com/sakirsensoy/genv/dotenv/autoload"

	"fmt"
	"myproject/config"
)

func  main() {
	fmt.Println(config.App.Host) // localhost
	fmt.Println(config.App.Port) // 1234
	fmt.Println(config.App.Debug) // true
}
```

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