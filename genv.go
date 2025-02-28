// Package genv is a library for Go (golang) that makes it easy to read and use
// environment variables in your projects. It also allows environment variables
// to be loaded from the .env file.
package genv

import (
	"os"
	"strconv"
)

// EnvVariable contains information about the environment variable, such as key,
// value, and default value.
type EnvVariable struct {
	Key          string
	Val          string
	DefaultValue interface{}
	IsDefined    bool
}

// EnvVariables is where environment variables are stored.
var EnvVariables = make(map[string]*EnvVariable)

// Key is used to determine the path of the environment variable to be accessed.
//
//	genv.Key("env-key").String()
func Key(key string) *EnvVariable {

	envVar, ok := EnvVariables[key]
	if !ok {

		val, ok := os.LookupEnv(key)
		EnvVariables[key] = &EnvVariable{Key: key, Val: val, IsDefined: ok}

		return EnvVariables[key]
	}

	return envVar
}

// Default is used to specify the default value for the environment
// variable to be accessed.
//
//	genv.Key("env-key").Default("defaultValue").String()
func (e *EnvVariable) Default(defaultValue interface{}) *EnvVariable {

	e.DefaultValue = defaultValue

	return e
}

// Update is used to update the value of the corresponding environment variable.
//
//	genv.Key("env-key").Update("updatedValue")
func (e *EnvVariable) Update(value interface{}) {

	switch value.(type) {
	case bool:
		e.Val = strconv.FormatBool(value.(bool))
	case float64:
		e.Val = strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case int:
		e.Val = strconv.FormatInt(int64(value.(int)), 10)
	case string:
		e.Val = value.(string)
	}

	e.IsDefined = true
	os.Setenv(e.Key, e.Val)
}

// Bool method is used for environment variables of type bool.
//
//	genv.Key("env-key").Bool()
func (e *EnvVariable) Bool() bool {

	var dv bool
	if !e.IsDefined {
		if e.DefaultValue != nil {
			dv = e.DefaultValue.(bool)
		}

		return dv
	}

	val, _ := strconv.ParseBool(e.Val)

	return val
}

// Float method is used for environment variables of type float.
//
//	genv.Key("env-key").Float()
func (e *EnvVariable) Float() float64 {

	var dv float64
	if !e.IsDefined {
		if e.DefaultValue != nil {
			dv = e.DefaultValue.(float64)
		}

		return dv
	}

	val, _ := strconv.ParseFloat(e.Val, 64)

	return val
}

// Int method is used for environment variables of type int.
//
//	genv.Key("env-key").Int()
func (e *EnvVariable) Int() int {

	var dv int
	if !e.IsDefined {
		if e.DefaultValue != nil {
			dv = e.DefaultValue.(int)
		}

		return dv
	}

	val, _ := strconv.ParseInt(e.Val, 10, 32)

	return int(val)
}

// trimQuotes removes single or double quotes from the beginning and end of a string
// if they exist at both ends.
func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// String method is used for environment variables of type string.
// It automatically trims single or double quotes from the beginning and end of the value.
//
//	genv.Key("env-key").String()
func (e *EnvVariable) String() string {

	var dv string
	if !e.IsDefined {
		if e.DefaultValue != nil {
			dv = e.DefaultValue.(string)
		}

		return trimQuotes(dv)
	}

	return trimQuotes(e.Val)
}
