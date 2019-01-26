package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/fyrchik/cli"
)

func main() {
	c := cli.NewContext(cli.Command{Name: "main"})
	err := c.Root.AddFlags(
		cli.StringFlag("name", "-n", "--name"),
		cli.BoolFlag("confirm", "--yes-i-am-really-sure"),
		cli.StringSliceFlag("multiple", "-m", "--multi"),
	)
	if err != nil {
		fatal(err)
	}
	if err = c.Parse(os.Args[1:]); err != nil {
		fatal(err)
	}
	spew.Dump(c.Named, c.Positional)
}

func fatal(err error) {
	fmt.Println(err)
	os.Exit(-1)
}
