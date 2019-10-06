package main

import (
	"os"
	"github.com/codegangsta/cli"
        "vsl"
)

func getCommand() string {
	return os.Args[1]
}

func getArgs() []string {
	return os.Args[2:]
}

const help = `an idempotent tool for managing /etc/vsl_hosts

 * Commands will exit 0 or 1 in a sensible way to facilitate scripting.
 * vslctl operates on /etc/vsl_hosts by default. 
 * Specify the VSL_SERVICES_PATH environment variable to change this.
 * Report bugs and feedback at <REPORT_TAG>
`

func main() {
	app := cli.NewApp()
	app.Name = "vsl"
	app.Authors = []cli.Author{{Name: "srihari kondapaneni", Email: "skondapa@gmail.com"}}
	app.Usage = help
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "f",
			Usage: "operate even if there are errors or conflicts",
		},
		cli.BoolFlag{
			Name:  "n",
			Usage: "no-op. Show changes but don't write them.",
		},
		cli.BoolFlag{
			Name:  "q",
			Usage: "quiet operation -- no notices",
		},
		cli.BoolFlag{
			Name:  "s",
			Usage: "silent operation -- no errors (implies -q)",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "add",
			Usage:  "add or replace a vsl hosts entry",
			Action: vsl.Add,
			Flags:  app.Flags,
		},
		{
			Name:    "del",
			Aliases: []string{"rm"},
			Usage:   "delete a vsl hosts entry",
			Action:  vsl.Del,
			Flags:   app.Flags,
		},
/*
		{
			Name:   "has",
			Usage:  "exit 0 if entry exists, 1 if not",
			Action: vsl.Has,
			Flags:  app.Flags,
		},
		{
			Name:   "on",
			Usage:  "enable a vsl hosts entry (if if exists)",
			Action: vsl.OnOff,
			Flags:  app.Flags,
		},
		{
			Name:   "off",
			Usage:  "disable a vsl hosts entry (don't delete it)",
			Action: vsl.OnOff,
			Flags:  app.Flags,
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list entries in the vsl hosts file",
			Action:  vsl.Ls,
			Flags:   app.Flags,
		},
		{
			Name:   "dump",
			Usage:  "dump the vsl hosts file as JSON",
			Action: vsl.Dump,
			Flags:  app.Flags,
		},
		{
			Name:   "apply",
			Usage:  "add hostnames from a JSON file to the vsl hosts file",
			Action: vsl.Apply,
			Flags:  app.Flags,
		},
*/
	}

	app.Run(os.Args)
	os.Exit(0)
}
