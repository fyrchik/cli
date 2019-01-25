package cli

import (
	"github.com/pkg/errors"
)

// Context is container for flags and commands.
// It will contain all parsed results.
type Context struct {
	Named      map[string]interface{}
	Positional []string
	flags      map[string]Flag
	optMap     map[string]string
}

// NewContext returns new Context.
func NewContext() *Context {
	return &Context{
		Named:      map[string]interface{}{},
		Positional: []string{},

		flags:  map[string]Flag{},
		optMap: map[string]string{},
	}
}

// AddFlags adds flags to the Context.
// TODO set default value for BoolFlags'
func (c *Context) AddFlags(fs ...*Flag) error {
	for _, f := range fs {
		if _, ok := c.flags[f.Name]; ok {
			return errors.Errorf("flag '%s' is already declared", f.Name)
		}
		c.flags[f.Name] = *f
		for _, opt := range f.Options {
			if _, ok := c.optMap[opt]; ok {
				return errors.Errorf("flag for option '%s' is already declared", opt)
			}
			c.optMap[opt] = f.Name
		}
	}
	return nil
}

// Parse parses command-line arguments. Results are stored
// in c.Named for named options and in c.Positional for positional ones (unexpected, huh?)
func (c *Context) Parse(args []string) (err error) {
	var (
		v      interface{}
		ind, n int
		name   string
		ok     bool
		f      Flag
	)

	c.Named = map[string]interface{}{}
	c.Positional = []string{}

	for ind < len(args) {
		if name, ok = c.optMap[args[ind]]; !ok {
			break
		}

		ind++
		f = c.flags[name]
		if f.ParseMany != nil {
			if v, n, err = f.ParseMany(args[ind:]); err != nil {
				return
			}
		} else if ind < len(args) {
			v, n = args[ind], 1
			if f.Parse != nil {
				if v, err = f.Parse(args[ind]); err != nil {
					return
				}
			}
		} else {
			return errors.Errorf("expected argument for '%s'", args[ind-1])
		}

		ind += n

		if f.Validate != nil {
			if err = f.Validate(v); err != nil {
				return
			}
		}

		if f.Combine != nil {
			v = f.Combine(c.Named[name], v)
		}
		c.Named[name] = v
	}

	if ind < len(args) {
		c.Positional = args[ind:]
	}

	c.setDefaults()

	for name, v = range c.Named {
		if f, ok = c.flags[name]; ok && f.PostValidate != nil {
			if err = f.PostValidate(v); err != nil {
				return
			}
		}
	}
	return
}

func (c *Context) setDefaults() {
	var (
		err error
		v   interface{}
	)

loop:
	for name, f := range c.flags {
		if _, ok := c.Named[name]; ok {
			continue
		}

		vals := f.getEnviron()
		if len(vals) == 0 {
			if f.Default != nil {
				c.Named[name] = f.Default
			}
			continue
		}

		if f.Parse == nil {
			c.Named[name] = vals[0]
			continue loop
		}
		for _, s := range vals {
			if v, err = f.Parse(s); err != nil {
				continue
			}
			c.Named[name] = v
			continue loop
		}
	}
}
