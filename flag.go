package cli

import (
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Combiner is a generic type for combining 2 values into 1.
type Combiner func(interface{}, interface{}) interface{}

// Validator is a generic type for checking if value satisfies
// some restrictions.
type Validator func(interface{}) error

// Parser is a generic type for transforming string into
// any other type.
type Parser func(string) (interface{}, error)

// MultiParser is a generic type for transforming sequence of
// string arguments into any other type. Returns number of strings consumed.
type MultiParser func([]string) (interface{}, int, error)

// Flag represents named command-line argument.
// TODO add required flags
type Flag struct {
	// Name under which flag is stored in Context.
	Name string

	// Options are command-line options corresponding to the flag.
	Options []string

	// Default is a default value for a flag.
	Default interface{}

	// EnvKey is the name of environment variable.
	// It has more precendence than Default and less than parse result.
	// EnvKey can specify multiple names delimited by comma ','.
	EnvKey string

	// Combine combines 2 values into 1 if multiple options
	// for the same flag are presented.
	Combine Combiner

	// Parse parses string into value of any other type.
	// Note: it returns any syntactic errors encountered (e.g. number containing letters).
	//   Use Validate for semantic errors (e.g. year must be <2050).
	Parse Parser

	// ParseMany eats many arguments at one time.
	// If takes precedence over Parse.
	// Note: be careful not to eat other flags.
	ParseMany MultiParser

	// PostValidate checks if flag's final value is correct.
	// It is executed once after all parsing was finished.
	PostValidate Validator

	// Validate checks if flag's value is correct for every occurence of
	// flag on command-line.
	Validate Validator
}

// NewFlag returns new flag.
// Note: you MUST provide Parse or ParseMany function
// via SetParse/SetParseMany.
func NewFlag(name string, options ...string) *Flag {
	return &Flag{Name: name, Options: options}
}

// SetCombine sets f.Combiner to c and returns f.
func (f *Flag) SetCombine(c Combiner) *Flag {
	f.Combine = c
	return f
}

// SetDefault sets f.Default to v and returns f.
func (f *Flag) SetDefault(v interface{}) *Flag {
	f.Default = v
	return f
}

// SetParse sets f.Parse to p and returns f.
func (f *Flag) SetParse(p Parser) *Flag {
	f.Parse = p
	return f
}

// SetParseMany sets f.ParseMany to p and returns f.
func (f *Flag) SetParseMany(p MultiParser) *Flag {
	f.ParseMany = p
	return f
}

// SetValidate sets f.Validate to v and returns f.
func (f *Flag) SetValidate(v Validator) *Flag {
	f.Validate = v
	return f
}

// SetPostValidate sets f.PostValidate to v and returns f.
func (f *Flag) SetPostValidate(v Validator) *Flag {
	f.PostValidate = v
	return f
}

// SetEnvKey sets f.EnvKey to s and returns f.
func (f *Flag) SetEnvKey(s string) *Flag {
	f.EnvKey = s
	return f
}

func (f *Flag) getEnviron() (vals []string) {
	for _, key := range strings.Split(f.EnvKey, ",") {
		if val, ok := syscall.Getenv(strings.TrimSpace(key)); ok {
			vals = append(vals, val)
			continue
		}

	}
	return
}

// StringFlag returns Flag which eats one argument and represents it as a string.
func StringFlag(name string, options ...string) *Flag {
	return &Flag{Name: name, Options: options}
}

// StringSliceFlag returns Flag which can be used multiple times
// to accumulate strings in a slice.
func StringSliceFlag(name string, options ...string) *Flag {
	return StringFlag(name, options...).
		SetCombine(func(a, b interface{}) interface{} {
			if a == nil {
				return []string{b.(string)}
			}
			return append(a.([]string), b.(string))
		})
}

// BoolFlag returns Flag which needs no argument and
// sets it's value to true if presented.
func BoolFlag(name string, options ...string) *Flag {
	return &Flag{
		Name:      name,
		Options:   options,
		Default:   false,
		ParseMany: Const(true),
	}
}

// BoolFlagT is like BoolFlag but true by default and
// can be set to false.
func BoolFlagT(name string, options ...string) *Flag {
	return &Flag{
		Name:      name,
		Options:   options,
		Default:   true,
		ParseMany: Const(false),
	}
}

// IntFlag returns Flag which parses it's argument to int64.
// Note: for other numeric types one needs to define other flags.
func IntFlag(name string, options ...string) *Flag {
	return &Flag{
		Name:    name,
		Options: options,
		Parse: func(arg string) (interface{}, error) {
			i, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return nil, err
			}
			return i, nil
		},
	}
}

// DurationFlag parses time duration intervals.
// Parsed value has type time.Duration
func DurationFlag(name string, options ...string) *Flag {
	return &Flag{
		Name:    name,
		Options: options,
		Parse: func(arg string) (interface{}, error) {
			return time.ParseDuration(arg)
		},
	}
}
