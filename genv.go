package genv

import (
	"os"
	"strconv"
)

type envVariable struct {
	key          string
	val          string
	defaultValue interface{}
	isDefined    bool
}

var envVariables = make(map[string]*envVariable)

func Key(key string) *envVariable {

	envVar, ok := envVariables[key]
	if !ok {

		val, ok := os.LookupEnv(key)
		envVariables[key] = &envVariable{key: key, val: val, isDefined: ok}

		return envVariables[key]
	}

	return envVar
}

func (e *envVariable) Default(defaultValue interface{}) *envVariable {

	e.defaultValue = defaultValue

	return e
}

func (e *envVariable) Update(value interface{}) {

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

func (e *envVariable) Bool() bool {

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

func (e *envVariable) Float() float64 {

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

func (e *envVariable) Int() int {

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

func (e *envVariable) String() string {

	var dv string
	if !e.isDefined {
		if e.defaultValue != nil {
			dv = e.defaultValue.(string)
		}

		return dv
	}

	return e.val
}
