package cmd

import (
	"dns-manager/dnsapi"
	"fmt"

	"github.com/spf13/cobra"
)

var addRecordCmd = &cobra.Command{
	Use:       "add",
	Short:     "Add Cloudflare DNS record",
	Long:      `Add Cloudflare DNS record`,
	ValidArgs: []string{"--domain", "-d", "--ip", "-i", "-rtype", "-r", "--proxied", "-p", "--zone", "-z"},
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		rtype, _ := cmd.Flags().GetString("rtype")
		ip, _ := cmd.Flags().GetString("ip")
		proxied, _ := cmd.Flags().GetBool("proxied")
		zone, _ := cmd.Flags().GetString("zone")
		serviceprovider, _ := cmd.Flags().GetString("serviceprovider")
		if domain == "" {
			fmt.Printf("No domain/subdomain user -d|--domain test.domain.tld")

		} else {
			if serviceprovider == "cloudflare" {
				zoneID := dnsconfig.Cloudflare.ZoneId
				token := dnsconfig.Cloudflare.Token
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
			} else if serviceprovider == "bind" {
				dnsapi.BindInsertRecord(dnsconfig.Bind.Server, dnsconfig.Bind.Keyname, dnsconfig.Bind.Hmackey, zone, domain, ip)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(addRecordCmd)
}
