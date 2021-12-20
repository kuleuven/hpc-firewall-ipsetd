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
// Important note: ipsetCmd will not return an error if an ipset command fails.
func (u *IpsetUpdater) NewIpset(ipset string, settype string, family string, timeout uint) (*Ipset, error) {
	// Try to create directly. This will probably fail silently as there is already some ipset present.
	if err := u.ipsetCmd(fmt.Sprintf("create %s %s family %s timeout %d comment\n", ipset, settype, family, timeout)); err != nil {
		return nil, err
	}

	// Destroy ipset-__new__. IpsetCmd will not complain if this fails - that is fine.
	if err := u.ipsetCmd(fmt.Sprintf("destroy %s-__new__\n", ipset)); err != nil {
		return nil, err
	}

	// Create ipset to to swap.
	if err := u.ipsetCmd(fmt.Sprintf("create %s-__new__ %s family %s timeout %d comment\n", ipset, settype, family, timeout)); err != nil {
		return nil, err
	}

	// Swap the ipsets.
	if err := u.ipsetCmd(fmt.Sprintf("swap %s %s-__new__\n", ipset, ipset)); err != nil {
		return nil, err
	}

	// Destroy ipset-__new__.
	if err := u.ipsetCmd(fmt.Sprintf("destroy %s-__new__\n", ipset)); err != nil {
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
func (i *Ipset) Add(entry IpsetEntry) error {
	var cmd string

	if entry.timeout != 0 {
		cmd = fmt.Sprintf("add -exist %s %s timeout %d comment \"%s\"\n", i.ipset, entry.addr, entry.timeout, entry.comment)
	} else {
		cmd = fmt.Sprintf("add -exist %s %s comment \"%s\"\n", i.ipset, entry.addr, entry.comment)
	}

	return i.u.ipsetCmd(cmd)
}

// AddMultiple adds multiple times
func (i *Ipset) AddMultiple(entries []IpsetEntry) error {
	for _, entry := range entries {
		if err := i.Add(entry); err != nil {
			return err
		}
	}

	return nil
}
