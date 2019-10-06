package vsl

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

var ipv4Pattern = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
var ipv6Pattern = regexp.MustCompile(`^[(a-fA-F0-9){1-4}:]+$`)

// MatchIPv4 returns true if the IP looks like it's IPv4. This does not
// validate whether the string is a valid IP address.
func MatchIPv4(ip string) bool {
	return ipv4Pattern.MatchString(ip)
}

// MatchIPv6 returns true if the IP looks like it's IPv6. This does not
// validate whether the string is a valid IP address.
func MatchIPv6(ip string) bool {
	if !strings.Contains(ip, ":") {
		return false
	}
	return ipv6Pattern.MatchString(ip)
}

type VSLConfig struct {
	Interface string `json:"intferface"`
	Service   string `json:"service"`
	App       string `json:"app"`
	Role      string `json:"role"`
	Subnet    net.IP `json:"subnet"`
	IPv6      bool   `json:"-"`
}

// NewVslConfig creates a new VSLConfig struct and automatically sets the IPv6
// field based on the IP you pass in.
func NewVslConfig(intf string, service string, app string, role string, subnet string) (*VSLConfig, error) {
	if !MatchIPv4(subnet) && !MatchIPv6(subnet) {
		return nil, fmt.Errorf("Unable to parse IP address %q", subnet)
	}
	ip := net.ParseIP(subnet)
	return &VSLConfig{
		Interface: intf,
		Service:   service,
		App:       app,
		Role:      role,
		Subnet:    ip,
		IPv6:      MatchIPv6(subnet),
	}, nil
}

// Equal compares two VslConfig Items.
func (v *VSLConfig) Equal(n *VSLConfig) bool {
	return v.Interface == n.Interface &&
		v.Service == n.Service &&
		v.App == n.App &&
		v.Role == n.Role &&
		v.Subnet.Equal(n.Subnet)
}

// EqualIP compares an IP against this Hostname.
func (v *VSLConfig) EqualSubnet(subnet net.IP) bool {
	return v.Subnet.Equal(subnet)
}

// IsValid does a spot-check on the VSLConfig to make sure they aren't blank
func (v *VSLConfig) IsValid() bool {
	return v.Interface != "" &&
		v.Service != "" &&
		v.App != "" &&
		(v.Role == "P" || v.Role == "C") &&
		v.Subnet != nil
}

func (v *VSLConfig) FormatHuman() string {
	return fmt.Sprintf("%s %s %s %s", v.Service, v.App, v.Role, v.Subnet)
}


func (v *VSLConfig) Compare(b interface{}) bool {

    switch bv := b.(type) {
    case *VSLConfig:
        if v == bv {
            return true
        }

        if v.Interface != bv.Interface {
            return false
        }

        return true
    default:
        return false
    }
}
