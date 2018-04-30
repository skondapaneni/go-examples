package vsl

import (
	"fmt"
	"os"
	"github.com/codegangsta/cli"
)

// AnyBool checks whether a boolean CLI flag is set either globally or on this
// individual command. This provides more flexible flag parsing behavior.
func AnyBool(c *cli.Context, key string) bool {
	return c.Bool(key) || c.GlobalBool(key)
}

// ErrCantWriteVslfile indicates that we are unable to write to the services file
var ErrCantWriteVslfile = fmt.Errorf(
	"Unable to write to %s. Maybe you need to sudo?", GetServicesPath())

// MaybeErrorln will print an error message unless -s is passed
func MaybeErrorln(c *cli.Context, message string) {
	if !AnyBool(c, "s") {
		os.Stderr.WriteString(fmt.Sprintf("%s\n", message))
	}
}

// MaybeError will print an error message unless -s is passed and then exit
func MaybeError(c *cli.Context, message string) {
	MaybeErrorln(c, message)
	os.Exit(1)
}

// MaybePrintln will print a message unless -q or -s is passed
func MaybePrintln(c *cli.Context, message string) {
	if !AnyBool(c, "q") && !AnyBool(c, "s") {
		fmt.Println(message)
	}
}

// MaybeLoadVslfile will try to load, parse, and return a Vslfile. If we
// encounter errors we will terminate, unless -f is passed.
func MaybeLoadVslfile(c *cli.Context) *Vslfile {
	vslfile, errs := LoadVslfile()
	if len(errs) > 0 && !AnyBool(c, "f") {
		for _, err := range errs {
			MaybeErrorln(c, err.Error())
		}
		MaybeError(c, "Errors while parsing vslfile. Try hostess fix")
	}
	return vslfile
}

// AlwaysLoadVslfile will load, parse, and return a Vslfile. If we encouter
// errors they will be printed to the terminal, but we'll try to continue.
func AlwaysLoadVslfile(c *cli.Context) *Vslfile {
	vslfile, errs := LoadVslfile()
	if len(errs) > 0 {
		for _, err := range errs {
			MaybeErrorln(c, err.Error())
		}
	}
	return vslfile
}

// MaybeSaveVslfile will output or write the Hostfile, or exit 1 and error.
func MaybeSaveVslfile(c *cli.Context, vslfile *Vslfile) {
	// If -n is passed, no-op and output the resultant hosts file to stdout.
	// Otherwise it's for real and we're going to write it.
	if AnyBool(c, "n") {
		fmt.Printf("%s", vslfile.Format())
	} else {
		err := vslfile.Save()
		if err != nil {
			MaybeError(c, ErrCantWriteVslfile.Error())
		}
	}
}

// Add command parses <hostname> <ip> and adds or updates a hostname in the
// hosts file. If the aff command is used the hostname will be disabled or
// added in the off state.
func Add(c *cli.Context) {
	if len(c.Args()) != 5 {
		MaybeError(c, "expected <interface> <service> <app> <role> <subnet>")
	}

	vslfile := MaybeLoadVslfile(c)
	service, err := NewVslConfig(c.Args()[0], c.Args()[1], 
			c.Args()[2],
			c.Args()[3],
			c.Args()[4])
	if err != nil {
		MaybeError(c, fmt.Sprintf("Failed to parse vsl entry: %s", err))
	}

	i, replace := vslfile.Services.Contains_i(service)

	// Note that Add() may return an error, but they are informational only. We
	// don't actually care what the error is -- we just want to add the
	// hostname and save the file. This way the behavior is idempotent.
	if (replace) {
		vslfile.Services.ReplaceIndex(service, i)
	} else {
		vslfile.Services.Add(service)
	}

	// If the user passes -n then we'll Add and show the new hosts file, but
	// not save it.
	if c.Bool("n") || AnyBool(c, "n") {
		fmt.Printf("%s", vslfile.Format())
	} else {
		MaybeSaveVslfile(c, vslfile)
		// We'll give a little bit of information about whether we added or
		// updated, but if the user wants to know they can use has or ls to
		// show the file before they run the operation. Maybe later we can add
		// a verbose flag to show more information.
		if replace {
			MaybePrintln(c, fmt.Sprintf("Updated %s", service.FormatHuman()))
		} else {
			MaybePrintln(c, fmt.Sprintf("Added %s", service.FormatHuman()))
		}
	}
}

// Del command removes any matching <intf> from the vsl  file
func Del(c *cli.Context) {
	if len(c.Args()) != 1 {
		MaybeError(c, "expected <interface>")
	}
	intf := c.Args()[0]
	service, _ := NewVslConfig(c.Args()[0], "", "", "", "")	
	vslfile := MaybeLoadVslfile(c)

	i, found := vslfile.Services.Contains_i(service)
	if found {
		vslfile.Services.RemoveIndex(i)
		if AnyBool(c, "n") {
			fmt.Printf("%s", vslfile.Format())
		} else {
			MaybeSaveVslfile(c, vslfile)
			MaybePrintln(c, fmt.Sprintf("Deleted %s", intf))
		}
	} else {
		MaybePrintln(c, fmt.Sprintf("%s not found in %s", intf, GetServicesPath()))
	}
}
