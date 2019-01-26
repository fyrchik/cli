package cli

import (
	"strings"

	"github.com/pkg/errors"
)

// Context is container for flags and commands.
// It will contain all parsed results.
type Context struct {
	Root       Command
	Named      map[string]interface{}
	Positional []string
}

// NewContext returns new Context.
func NewContext(root Command) *Context {
	return &Context{
		Root:       root,
		Named:      map[string]interface{}{},
		Positional: []string{},
	}
}

// Parse parses command-line arguments. Results are stored
// in c.Named for named options and in c.Positional for positional ones (unexpected, huh?)
func (c *Context) Parse(args []string) (err error) {
	c.Named = map[string]interface{}{}
	c.Positional = []string{}

	return c.parse(&c.Root, args)
}

// parse is auxilliary function which parses (sub-)Command arguments in
// the Context c.
func (c *Context) parse(cmd *Command, args []string) (err error) {
	var (
		v      interface{}
		ind, n int
		name   string
		ok     bool
		f      Flag
	)

	for ind < len(args) {
		for i := range cmd.Subcommands {
			if cmd.Subcommands[i].Name == args[ind] {
				return c.parse(&cmd.Subcommands[i], args[ind+1:])
			}
		}

		if strings.HasPrefix(args[ind], "-") {
			if args[ind] == "--" {
				ind++
				break
			}
			if name, ok = cmd.optMap[args[ind]]; !ok {
				return errors.Errorf("unknown option '%s'", args[ind])
			}
		}

		if name, ok = cmd.optMap[args[ind]]; !ok {
			break
		}

		ind++
		f = cmd.flags[name]
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

	c.setDefaults(cmd)

	for name, f := range cmd.flags {
		if v, ok := c.Named[name]; ok && f.PostValidate != nil {
			if err = f.PostValidate(v); err != nil {
				return
			}
		}
	}
	return
}

func (c *Context) setDefaults(cmd *Command) {
	var (
		err error
		v   interface{}
	)

loop:
	for name, f := range cmd.flags {
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
