# Genv [![Build Status](https://travis-ci.org/sakirsensoy/genv.svg?branch=master)](https://travis-ci.org/sakirsensoy/genv)

Read environment variables easily with type support and load them from .env in Go (golang).

## Installation

```shell
go get github.com/sakirsensoy/genv
```

## Example Usage

`.env`:

```shell
APP_HOST=localhost
APP_PORT=1234
APP_DEBUG=true
```

`config/config.go`:

```go
package config

import (
  "github.com/sakirsensoy/genv"
)

type appConfig struct {
  Host string
  Port int
  Debug bool
}

var App = &appConfig{
  Host: genv.Key("APP_HOST").String(),
  Port: genv.Key("APP_PORT").Default(8080).Int(),
  Debug: genv.Key("APP_DEBUG").Default(false).Bool(),
}
```

`main.go`:

```go
package main

import (
  _ "github.com/sakirsensoy/genv/dotenv/autoload"

  "fmt"
  "project/config"
)

func main() {

  fmt.Println(config.App.Host) // localhost
  fmt.Println(config.App.Port) // 1234
  fmt.Println(config.App.Debug) // true
}
```
