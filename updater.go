package main

import (
	"fmt"
	"log"

	"github.com/42wim/ipsetd/ipset"
	consul "github.com/hashicorp/consul/api"
)

var (
	ipsetPath = "/sbin/ipset"
	// Default ipset timeout
	ipsetTimeout uint = 300
	// Max ipset timeout of values from consul
	ipsetMaxTimeout uint = 86400
)

// IpsetUpdaterConfig specifies configuration for an IpsetdUpdater
type IpsetUpdaterConfig struct {
	ConsulURL   string
	ConsulToken string
	ConsulPath  string
	Ipset       string
}

// An IpsetUpdater keeps an ipset in sync with consul configuration
type IpsetUpdater struct {
	IpsetUpdaterConfig
	ConsulClient *consul.Client
	IpsetClient  *ipset.IPset
	Ipset4       *Ipset
	Ipset6       *Ipset
}

// NewIpsetUpdater creates a new IpsetUpdater
func NewIpsetUpdater(config IpsetUpdaterConfig) (*IpsetUpdater, error) {
	// Consul client
	consulConfig := consul.DefaultConfig()

	if config.ConsulURL != "" {
		consulConfig.Address = config.ConsulURL
	}

	if config.ConsulToken != "" {
		consulConfig.Token = config.ConsulToken
	}

	consulClient, err := consul.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	// Create IpsetUpdater object
	u := &IpsetUpdater{
		IpsetUpdaterConfig: config,
		ConsulClient:       consulClient,
		IpsetClient:        ipset.NewIPset(ipsetPath),
	}

	// Prepare ipsets
	u.Ipset4, err = u.NewIpset(config.Ipset, "hash:net", "inet", ipsetTimeout)
	if err != nil {
		return nil, err
	}

	u.Ipset6, err = u.NewIpset(fmt.Sprintf("%s-6", config.Ipset), "hash:net", "inet6", ipsetTimeout)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Run an IpsetUpdater
func (u *IpsetUpdater) Run() error {
	var (
		index    uint64
		pairs    []*consul.KVPair
		entries4 []*IpsetEntry
		entries6 []*IpsetEntry
		err      error
	)

	// Main loop
	for {
		pairs, index, err = u.GetKVPairs(u.ConsulPath, index)
		if err != nil {
			return err
		}

		for _, pair := range pairs {
			entries4, entries6, err = KVPairToIpsetEntries(pair)
			if err != nil {
				log.Printf("error while parsing consul data for %s: %s", pair.Key, err)
				continue
			}

			err = u.Ipset4.AddMultiple(entries4)
			if err != nil {
				return err
			}

			err = u.Ipset4.AddMultiple(entries6)
			if err != nil {
				return err
			}
		}
	}
}

// Run ipset command wrapper
func (u *IpsetUpdater) ipsetCmd(cmd string) error {
	_, err := u.IpsetClient.Cmd(cmd)

	if err != nil {
		log.Printf("error %s occurred during %s, restarting ipset client", err, cmd)

		// Recover by starting new ipset process
		u.IpsetClient = ipset.NewIPset(ipsetPath)

		return err
	}

	return nil
}
