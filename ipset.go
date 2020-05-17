package main

import (
	"fmt"
)

// An IpsetEntry describes an entry to be added to some ipset
type IpsetEntry struct {
	addr    string
	timeout uint
	comment string
}

// An Ipset designates a ipset fill session
type Ipset struct {
	u       *IpsetUpdater
	ipset   string
	settype string
	family  string
	timeout uint
}

// NewIpset prepares a new ipset
func (u *IpsetUpdater) NewIpset(ipset string, settype string, family string, timeout uint) (*Ipset, error) {
	cmd := fmt.Sprintf("create %s-__NEW__ %s family %s timeout %d\n", ipset, settype, family, timeout)

	if err := u.ipsetCmd(cmd); err != nil {
		return nil, err
	}

	return &Ipset{
		u:       u,
		ipset:   ipset,
		settype: settype,
		family:  family,
		timeout: timeout,
	}, nil
}

// Add a new ipset entry
func (i *Ipset) Add(entry *IpsetEntry) error {
	var cmd string

	if entry.timeout != 0 {
		cmd = fmt.Sprintf("add -exist %s %s timeout %d comment '%s'\n", i.ipset, entry.addr, entry.timeout, entry.comment)
	} else {
		cmd = fmt.Sprintf("add -exist %s %s comment '%s'\n", i.ipset, entry.addr, entry.comment)
	}

	return i.u.ipsetCmd(cmd)
}

// AddMultiple adds multiple times
func (i *Ipset) AddMultiple(entries []*IpsetEntry) error {
	for _, entry := range entries {
		if err := i.Add(entry); err != nil {
			return err
		}
	}

	return nil
}
