package cmd

import (
	"dns-manager/dnsapi"
	"os"

	"github.com/spf13/cobra"
)

var zoneID = ""
var token = ""

var rootCmd = &cobra.Command{
	Use:   "dns-manager",
	Short: "dns-manager. Manage Cloudflare DNS records.",
	Long:  `dns-manager. Manage Cloudflare DNS records.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		zoneID, token = dnsapi.GetConfig(configPath)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("domain", "d", "", "Domain/Subdomain")
	rootCmd.PersistentFlags().StringP("ip", "i", "", "IP Address")
	rootCmd.PersistentFlags().StringP("rtype", "r", "", "Record type (A,TXT,CNAME, ...)")
	rootCmd.PersistentFlags().BoolP("proxied", "p", false, "Record type (A,TXT,CNAME, ...)")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config.yaml")
}
