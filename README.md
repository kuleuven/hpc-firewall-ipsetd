# HPC Firewall Ipset Updater

Updates ipsets based on values stored in consul. To be used with https://github.com/kuleuven/hpc-firewall

Usage:

```bash
./hpc-firewall-ipsetd --consul-addr <consul> --consul-token <token> --consul-path <path> --ipset <ipset-name>
```

Arguments:

* `--consul-addr` or `CONSUL_HTTP_ADDR`: Consul address for storage of ips
* `--consul-token` or `CONSUL_HTTP_TOKEN`: Consul token used to read from the key value store
* `--consul-path` or `CONSUL_PATH`: Root path in the consul kv store to read from
* `--ipset` (defaults to `hpcuafw`): Name of the ipset to update

The specified name of ipset is used for the IPv4 ipset, and should be of the type hash:net. For IPv6, '-6' is appended to this name.

The firewall on the server should be configured to allow connections from these ipsets.
