package cli

import "github.com/pkg/errors"

type Context struct {
	Named      map[string]interface{}
	Positional []string
	flags      map[string]Flag
	optMap     map[string]string
}

func NewContext() *Context {
	return &Context{
		Named:      map[string]interface{}{},
		Positional: []string{},

		flags:  map[string]Flag{},
		optMap: map[string]string{},
	}
}

func (c *Context) clear() {
	c.Named = map[string]interface{}{}
	c.Positional = []string{}
}

// AddFlags adds flags to the Context.
// TODO set default value for BoolFlags'
func (c *Context) AddFlags(fs ...Flag) error {
	for _, f := range fs {
		if _, ok := c.flags[f.Name]; ok {
			return errors.Errorf("flag with name '%s' is already declared", f.Name)
		}
		c.flags[f.Name] = f
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
			if f.Parse != nil {
				// FIXME error if ind == len(args)
				if v, err = f.Parse(args[ind]); err != nil {
					return
				}
				n = 1
			} else {
				v, n = args[ind], 1
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

		if cur, ok := c.Named[name]; ok && f.Combine != nil {
			v = f.Combine(cur, v)
		}
		c.Named[name] = v
	}

	if ind < len(args) {
		c.Positional = args[ind:]
	}

	for name, val := range c.Named {
		if f, ok := c.flags[name]; ok && f.PostValidate != nil {
			if err = f.PostValidate(val); err != nil {
				return
			}
		}
	}

	return
}
