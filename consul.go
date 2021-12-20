package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	consul "github.com/hashicorp/consul/api"
)

// GetKVPairs returns consul kv pairs for a given prefix.
func (u *IpsetUpdater) GetKVPairs(prefix string, index uint64) ([]*consul.KVPair, uint64, error) {
	queryoptions := &consul.QueryOptions{}

	if index != 0 {
		queryoptions.WaitIndex = index
	}

	kv := u.ConsulClient.KV()

	pairs, meta, err := kv.List(prefix, queryoptions)
	if err != nil {
		return nil, 0, err
	}

	return pairs, meta.LastIndex, nil
}

// A KVIpsetRecord represents an IP in the consul kv store to be used for an ipset
type KVIpsetRecord struct {
	Since      time.Time `json:"since"`
	Expiration time.Time `json:"expiration"`
	IP         string    `json:"ip"`
}

// KVPairToIpsetEntries parses consul kv pairs for json ipset entries.
func KVPairToIpsetEntries(kp *consul.KVPair) ([]IpsetEntry, []IpsetEntry, error) {
	var (
		entries4 = []IpsetEntry{}
		entries6 = []IpsetEntry{}
		now      = time.Now()
		data     []KVIpsetRecord
	)

	err := json.Unmarshal(kp.Value, &data)
	if err != nil {
		return nil, nil, err
	}

	for _, dataEntry := range data {
		expires := dataEntry.Expiration
		starts := dataEntry.Since

		if (!expires.IsZero() && now.After(expires)) || (!starts.IsZero() && now.Before(starts)) {
			continue
		}

		var timeout uint

		if expires.IsZero() {
			timeout = ipsetMaxTimeout
		} else {
			timeout = uint(expires.Sub(now).Seconds())
			if timeout > ipsetMaxTimeout {
				timeout = ipsetMaxTimeout
			} else if timeout == 0 {
				continue
			}
		}

		addr := net.ParseIP(dataEntry.IP)

		if addr == nil {
			return nil, nil, fmt.Errorf("invalid address %s", dataEntry.IP)
		}

		ipsetEntry := IpsetEntry{
			addr:    addr.String(),
			timeout: timeout,
			comment: kp.Key,
		}

		if addr.To4() != nil {
			entries4 = append(entries4, ipsetEntry)
		} else {
			entries6 = append(entries6, ipsetEntry)
		}
	}

	return entries4, entries6, nil
}
