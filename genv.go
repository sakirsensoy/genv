// Package genv is a library for Go (golang) that makes it easy to read and use
// environment variables in your projects. It also allows environment variables
// to be loaded from the .env file.
package genv

import (
	"os"
	"strconv"
	"sync"

	"github.com/sakirsensoy/genv/dotenv"
)

var (
	autoLoadOnce    sync.Once
	autoLoadEnabled = true // Enabled by default
	autoLoadPath    = ".env"

	// mu protects EnvVariables map for concurrent access
	mu sync.RWMutex
)

// DisableAutoLoad disables automatic .env file loading.
// Must be called before any genv.Key() calls.
func DisableAutoLoad() {
	autoLoadEnabled = false
}

// SetAutoLoadPath sets the path for automatic .env file loading.
// Default is ".env". Must be called before any genv.Key() calls.
func SetAutoLoadPath(path string) {
	autoLoadPath = path
}

// ensureEnvLoaded automatically loads .env file on first call.
func ensureEnvLoaded() {
	if autoLoadEnabled {
		autoLoadOnce.Do(func() {
			dotenv.Load(autoLoadPath)
		})
	}
}

// EnvVariable contains information about the environment variable, such as key,
// value, and default value.
type EnvVariable struct {
	Key          string
	Val          string
	DefaultValue interface{}
	IsDefined    bool
	mu           sync.RWMutex // protects Val and IsDefined for concurrent access
}

// EnvVariables is where environment variables are stored.
var EnvVariables = make(map[string]*EnvVariable)

// Key is used to determine the path of the environment variable to be accessed.
//
//	genv.Key("env-key").String()
func Key(key string) *EnvVariable {
	// Auto-load .env file on first Key() call
	ensureEnvLoaded()

	// Try to get from cache with read lock first
	mu.RLock()
	envVar, ok := EnvVariables[key]
	mu.RUnlock()

	if !ok {
		// Not in cache, need write lock to create new entry
		mu.Lock()
		// Double-check after acquiring write lock
		envVar, ok = EnvVariables[key]
		if !ok {
			val, defined := os.LookupEnv(key)
			envVar = &EnvVariable{Key: key, Val: val, IsDefined: defined}
			EnvVariables[key] = envVar
		}
		mu.Unlock()
	}

	return envVar
}

// Default is used to specify the default value for the environment
// variable to be accessed.
//
//	genv.Key("env-key").Default("defaultValue").String()
func (e *EnvVariable) Default(defaultValue interface{}) *EnvVariable {
	e.mu.Lock()
	e.DefaultValue = defaultValue
	e.mu.Unlock()
	return e
}

// Update is used to update the value of the corresponding environment variable.
//
//	genv.Key("env-key").Update("updatedValue")
func (e *EnvVariable) Update(value interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch v := value.(type) {
	case bool:
		e.Val = strconv.FormatBool(v)
	case float64:
		e.Val = strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		e.Val = strconv.FormatInt(int64(v), 10)
	case string:
		e.Val = v
	}

	e.IsDefined = true
	os.Setenv(e.Key, e.Val)
}

// Bool method is used for environment variables of type bool.
//
//	genv.Key("env-key").Bool()
func (e *EnvVariable) Bool() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.IsDefined {
		if e.DefaultValue != nil {
			return e.DefaultValue.(bool)
		}
		return false
	}

	val, _ := strconv.ParseBool(e.Val)
	return val
}

// Float method is used for environment variables of type float.
//
//	genv.Key("env-key").Float()
func (e *EnvVariable) Float() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.IsDefined {
		if e.DefaultValue != nil {
			return e.DefaultValue.(float64)
		}
		return 0
	}

	val, _ := strconv.ParseFloat(e.Val, 64)
	return val
}

// Int method is used for environment variables of type int.
//
//	genv.Key("env-key").Int()
func (e *EnvVariable) Int() int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.IsDefined {
		if e.DefaultValue != nil {
			return e.DefaultValue.(int)
		}
		return 0
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
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.IsDefined {
		if e.DefaultValue != nil {
			return trimQuotes(e.DefaultValue.(string))
		}
		return ""
	}

	return trimQuotes(e.Val)
}
