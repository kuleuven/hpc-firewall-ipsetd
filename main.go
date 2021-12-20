package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	config := IpsetUpdaterConfig{}

	rootCmd := &cobra.Command{
		Use:   "hpc-firewall-ipsetd",
		Short: "HPC firewall ipsetd keeps an ipset up to date",
		Long:  `HPC firewall ipsetd keeps an ipset in sync with configuration in consul.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			u, err := NewIpsetUpdater(config)
			if err != nil {
				log.Fatal(err)
			}

			err = u.Run()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.Flags().StringVar(&config.ConsulURL, "consul-addr", "", "Consul address")
	rootCmd.Flags().StringVar(&config.ConsulToken, "consul-token", "", "Consul token")
	rootCmd.Flags().StringVarP(&config.ConsulPath, "consul-path", "p", os.Getenv("CONSUL_PATH"), "Consul path")
	rootCmd.Flags().StringVarP(&config.Ipset, "ipset", "s", "hpcuafw", "Name of ipset to update")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
