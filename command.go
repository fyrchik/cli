package cli

import "github.com/pkg/errors"

// Command represents command with separate flags.
// Note: it is assumed that subcommands have no flags with the same name.
type Command struct {
	Name        string
	Help        string
	Subcommands []Command

	flags  map[string]Flag
	optMap map[string]string
}

// AddFlags adds flags to the Command.
func (c *Command) AddFlags(fs ...*Flag) error {
	if c.flags == nil {
		c.flags = map[string]Flag{}
	}
	if c.optMap == nil {
		c.optMap = map[string]string{}
	}
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
