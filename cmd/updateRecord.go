package cmd

import (
	"dns-manager/dnsapi"
	"fmt"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var updateRecordCmd = &cobra.Command{
	Use:       "update",
	Short:     "update DNS record",
	Long:      `update DNS record`,
	ValidArgs: []string{"--domain", "-d", "--ip", "-i", "--oldip", "-o", "-rtype", "-r", "--proxied", "-p", "--zone", "-z"},
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		rtype, _ := cmd.Flags().GetString("rtype")
		ip, _ := cmd.Flags().GetString("ip")
		oldip, _ := cmd.Flags().GetString("oldip")
		proxied, _ := cmd.Flags().GetBool("proxied")
		zone, _ := cmd.Flags().GetString("zone")
		serviceprovider, _ := cmd.Flags().GetString("serviceprovider")
		if dnsapi.CheckEmpty(domain, "domain", "-d|--domain test.domain.tld") {
			return
		} else {
			if serviceprovider == "cloudflare" {
				zoneID := dnsconfig.Cloudflare.ZoneId
				token := dnsconfig.Cloudflare.Token
				dnsapi.UpdateRecord(zoneID, token, domain, rtype, ip, proxied)
				resp, recordID, err := dnsapi.GetRecordId(zoneID, token, domain)
				if resp {
					fmt.Printf("RecordID: %s\n", recordID)
				} else {
					fmt.Printf("Fehler beim abrufen des Records: %v\n", err)
				}
			} else if serviceprovider == "bind" {
				if dnsapi.CheckEmpty(zone, "zone", "-z|--zone example.com") ||
					dnsapi.CheckEmpty(ip, "ip address", "-i|--ip 123.123.123.123") ||
					dnsapi.CheckEmpty(oldip, "old ip address", "-o|--oldip 11.11.11.11") ||
					dnsapi.CheckEmpty(rtype, "record type", "-r|--rtype (A|MX|TXT|...)") {
					return
				}
				dnsapi.BindUpdateRecord(dnsconfig.Bind.Server, dnsconfig.Bind.Keyname, dnsconfig.Bind.Hmackey, zone, domain, ip, oldip, rtype)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(updateRecordCmd)
}
