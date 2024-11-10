package cmd

import (
	"dns-manager/dnsapi"
	"fmt"

	"github.com/spf13/cobra"
)

var getRecordCmd = &cobra.Command{
	Use:       "get",
	Short:     "Get Cloudflare DSN Record",
	Long:      `Get Cloudflare DSN Record`,
	ValidArgs: []string{"--domain", "-d"},
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		configPath, _ := cmd.Flags().GetString("config")

		zoneID, token = dnsapi.GetConfig(configPath)

		// fmt.Printf("ZonenID: %s, Token: %s\n", zoneID, token)
		if domain == "" {
			fmt.Printf("No domain/subdomain user -d|--domain test.domain.tld")

		} else {
			dnsapi.GetRecordId(zoneID, token, domain)
			resp, recordID, err := dnsapi.GetRecordId(zoneID, token, domain)
			if resp {
				fmt.Printf("RecordID: %s\n", recordID)
			} else {
				fmt.Printf("Fehler beim abrufen des Records: %v\n", err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(getRecordCmd)
}
