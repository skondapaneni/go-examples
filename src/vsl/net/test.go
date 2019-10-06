package main

import (
    "github.com/vishvananda/netlink"
)

func main() {
    lo, _ := netlink.LinkByName("lo0")
    addr, _ := netlink.ParseAddr("169.254.169.250/32")
    netlink.AddrAdd(lo, addr)
}
