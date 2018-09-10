package env

import (
	"os"
	"strconv"
	"strings"
)

// Var holds details for an environment variable
type Var struct {
	Key   string
	Value string
	Set   bool
}

// Get reads an environment variable from the OS
func Get(key string) Var {
	value, set := os.LookupEnv(key)
	return Var{Key: key, Value: value, Set: set}
}

// WithDefault returns the value of the environment variable if it is set.
// Otherwise it returns the provided default value.
func (e Var) WithDefault(value string) string {
	if e.Set {
		return e.Value
	}
	return value
}

// Required returns the value of the environment variable if it is set.
// Otheriwse it will call provided error func and return an empty string.
//
// Example:
//  env.Get("SOME_ENVIRONMENT_VARIABLE").Required(func(key string) { log.Fatalf("%q must be set", key) })
//
func (e Var) Required(errlog func(key string)) string {
	if !e.Set {
		errlog(e.Key)
	}
	return e.Value
}

// WithDefaultInt attempts to read an integer from the Var, returns value if
// Var is unset, and calls errlog if the Var is set to something that is not
// parsable as an integer.
func (e Var) WithDefaultInt(value int, errlog func(key string, parseErr error)) int {
	if !e.Set {
		return value
	}
	v, err := strconv.Atoi(e.Value)
	if err != nil {
		errlog(e.Key, err)
	}
	return v
}

// List returns the individual values of a comma separated list from a Var.
func (e Var) List(sep string) []string {
	return strings.Split(e.Value, sep)
}
