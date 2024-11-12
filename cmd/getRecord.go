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
		if domain == "" {
			fmt.Printf("No domain/subdomain user -d|--domain test.domain.tld")

		} else {
			zoneID := dnsconfig.Cloudflare.ZoneId
			token := dnsconfig.Cloudflare.Token
			dnsapi.GetRecord(zoneID, token, domain)
			resp, recordID, msg, _ := dnsapi.GetRecord(zoneID, token, domain)
			if resp {
				fmt.Printf("RecordID: %s\n", recordID)
				fmt.Printf("%s\n", msg)
			} else {
				fmt.Printf("Record not found\n")
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(getRecordCmd)
}
