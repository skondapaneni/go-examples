package vsl

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Vslfile represents /etc/services (or a similar file, depending on OS),
type Vslfile struct {
	Path     string
	Services *ServiceList
	data     []byte
}

// NewVslfile creates a new Vslfile object from the specified file.
func NewVslfile() *Vslfile {
	return &Vslfile{GetServicesPath(), NewServiceList(), []byte{}}
}

func GetServicesPath() string {
	path := os.Getenv("VSL_SERVICES_PATH")
	if path == "" {
		path = "/etc/vsl_hosts"
	}
	return path
}

// ParseLine parses an individual line in a services file
func ParseLine(line string) (*VSLConfig, error) {

	if len(line) == 0 {
		return nil, fmt.Errorf("line is blank")
	}

	// Parse leading # for disabled lines
	if line[0:1] == "#" {
		return nil, fmt.Errorf("line is disabled")
	}

	// Parse other #s for actual comments
	line = strings.Split(line, "#")[0]

	// Replace tabs and multispaces with single spaces throughout
	line = strings.Replace(line, "\t", " ", -1)
	for strings.Contains(line, "  ") {
		line = strings.Replace(line, "  ", " ", -1)
	}
	line = strings.TrimSpace(line)

	// Break line into words
	words := strings.Split(line, " ")
	for idx, word := range words {
		words[idx] = strings.TrimSpace(word)
	}

	// Separate the first bit (the ip) from the other bits (the domains)
	intf := words[0]
	service := words[1]
	app := words[2]
	role := words[3]
	subnet := words[4]

	vslConfig, err := NewVslConfig(intf, service, app, role, subnet)
	if err != nil {
		return nil, err
	}
	return vslConfig, nil
}

// Parse reads
func (v *Vslfile) Parse() []error {
	var errs []error
	var line = 1
	for _, s := range strings.Split(string(v.data), "\n") {
		service, err := ParseLine(s)
		if err != nil {
			fmt.Printf("Error parsing line %d: %s\n", line, err)
			errs = append(errs, err)
		} else {
			v.Services.Add(service)
		}
		line++
	}
	return errs
}

// Read the contents of the hostfile from disk
func (v *Vslfile) Read() error {
	data, err := ioutil.ReadFile(v.Path)
	if err == nil {
		v.data = data
	}
	return err
}

/* LoadVslfile creates a new Vslfile struct and tries to populate it from
 * disk. Read and/or parse errors are returned as a slice.
 */
func LoadVslfile() (vslfile *Vslfile, errs []error) {
	vslfile = NewVslfile()
	readErr := vslfile.Read()
	if readErr != nil {
		errs = []error{readErr}
		return
	}
	errs = vslfile.Parse()
	return
}

func (v *Vslfile) Format() []byte {
	return v.Services.Format()
}

/**
 * Save writes the Vslfile to disk to /etc/vsl_hosts or to the location specified
 * by the VSL_SERVICES_PATH environment variable (if set).
 */
func (v *Vslfile) Save() error {
	file, err := os.OpenFile(v.Path, os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.Write(v.Format())

	return err
}
