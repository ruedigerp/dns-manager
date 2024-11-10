package cmd

import (
	"dns-manager/dnsapi"
	"fmt"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var addRecordCmd = &cobra.Command{
	Use:       "add",
	Short:     "Add Cloudflare DNS record",
	Long:      `Add Cloudflare DNS record`,
	ValidArgs: []string{"--domain", "-d", "--ip", "-i", "-type", "-t", "--proxied", "-p"},
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		rtype, _ := cmd.Flags().GetString("rtype")
		ip, _ := cmd.Flags().GetString("ip")
		proxied, _ := cmd.Flags().GetBool("proxied")
		configPath, _ := cmd.Flags().GetString("config")
		zoneID, token = dnsapi.GetConfig(configPath)
		// fmt.Printf("ZonenID: %s, Token: %s\n", zoneID, token)
		if domain == "" {
			fmt.Printf("No domain/subdomain user -d|--domain test.domain.tld")

		} else {
			resp, _, _ := dnsapi.GetRecordId(zoneID, token, domain)
			if resp {
				fmt.Printf("Record exists.\n")
				return
			} else {
				dnsapi.AddRecord(zoneID, token, domain, rtype, ip, proxied)
				resp, recordID, err := dnsapi.GetRecordId(zoneID, token, domain)
				if resp {
					fmt.Printf("RecordID: %s\n", recordID)
				} else {
					fmt.Printf("Fehler beim abrufen des Records: %v\n", err)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(addRecordCmd)
}
