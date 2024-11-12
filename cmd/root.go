package cmd

import (
	"dns-manager/dnsapi"
	"os"

	"github.com/spf13/cobra"
)

var zoneID = ""
var token = ""
var defaultConfig = "./config.yaml"

var dnsconfig dnsapi.Config

var (
	rootCmd = &cobra.Command{
		Use:   "dns-manager",
		Short: "dns-manager. Manage Cloudflare DNS records.",
		Long:  `dns-manager. Manage Cloudflare DNS records.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config, _ := cmd.Flags().GetString("config")
			if config != "" {
				dnsconfig = dnsapi.GetConfig(config)
			} else {
				dnsconfig = dnsapi.GetConfig(defaultConfig)
			}
		},
	}
	config string
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("domain", "d", "", "Domain/Subdomain")
	rootCmd.PersistentFlags().StringP("zone", "z", "", "Domain zone")
	rootCmd.PersistentFlags().StringP("ip", "i", "", "IP Address")
	rootCmd.PersistentFlags().StringP("oldip", "o", "", "Old ip address")
	rootCmd.PersistentFlags().StringP("rtype", "r", "", "Record type (A,TXT,CNAME, ...)")
	rootCmd.PersistentFlags().BoolP("proxied", "p", false, "Record type (A,TXT,CNAME, ...)")
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "config.yaml")
	rootCmd.PersistentFlags().StringP("serviceprovider", "s", "", "Service Provider (Cloudflare, bind)")
}