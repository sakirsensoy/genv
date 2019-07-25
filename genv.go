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
	key          string
	val          string
	defaultValue interface{}
	isDefined    bool
}

// EnvVariables is where environment variables are stored.
var EnvVariables = make(map[string]*EnvVariable)

// Key is used to determine the path of the environment variable to be accessed.
//
//				genv.Key("env-key").String()
//
func Key(key string) *EnvVariable {

	envVar, ok := EnvVariables[key]
	if !ok {

		val, ok := os.LookupEnv(key)
		EnvVariables[key] = &EnvVariable{key: key, val: val, isDefined: ok}

		return EnvVariables[key]
	}

	return envVar
}

// Default is used to specify the default value for the environment
// variable to be accessed.
//
//				genv.Key("env-key").Default("defaultValue").String()
//
func (e *EnvVariable) Default(defaultValue interface{}) *EnvVariable {

	e.defaultValue = defaultValue

	return e
}

// Update is used to update the value of the corresponding environment variable.
//
//				genv.Key("env-key").Update("updatedValue")
//
func (e *EnvVariable) Update(value interface{}) {

	switch value.(type) {
	case bool:
		e.val = strconv.FormatBool(value.(bool))
	case float64:
		e.val = strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case int:
		e.val = strconv.FormatInt(int64(value.(int)), 10)
	case string:
		e.val = value.(string)
	}

	e.isDefined = true
	os.Setenv(e.key, e.val)
}

// Bool method is used for environment variables of type bool.
//
//				genv.Key("env-key").Bool()
//
func (e *EnvVariable) Bool() bool {

	var dv bool
	if !e.isDefined {
		if e.defaultValue != nil {
			dv = e.defaultValue.(bool)
		}

		return dv
	}

	val, _ := strconv.ParseBool(e.val)

	return val
}

// Float method is used for environment variables of type float.
//
//				genv.Key("env-key").Float()
//
func (e *EnvVariable) Float() float64 {

	var dv float64
	if !e.isDefined {
		if e.defaultValue != nil {
			dv = e.defaultValue.(float64)
		}

		return dv
	}

	val, _ := strconv.ParseFloat(e.val, 64)

	return val
}

// Int method is used for environment variables of type int.
//
//				genv.Key("env-key").Int()
//
func (e *EnvVariable) Int() int {

	var dv int
	if !e.isDefined {
		if e.defaultValue != nil {
			dv = e.defaultValue.(int)
		}

		return dv
	}

	val, _ := strconv.ParseInt(e.val, 10, 32)

	return int(val)
}

// String method is used for environment variables of type string.
//
//				genv.Key("env-key").String()
//
func (e *EnvVariable) String() string {

	var dv string
	if !e.isDefined {
		if e.defaultValue != nil {
			dv = e.defaultValue.(string)
		}

		return dv
	}

	return e.val
}
