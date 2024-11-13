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
		// domain := getFlag(cmd, "domain", "-d|--domain test.domain.tld")
		// rtype := getFlag(cmd, "rtype", "-r|--rtype (A|MX|TXT|...)")
		// ip := getFlag(cmd, "ip", "-i|--ip 123.123.123.123")
		// proxied, _ := cmd.Flags().GetBool("proxied")
		// zone := getFlag(cmd, "zone", "-z|--zone example.com")
		// serviceprovider := getFlag(cmd, "serviceprovider", "--serviceprovider (cloudflare|bind)")

		domain, _ := cmd.Flags().GetString("domain")
		rtype, _ := cmd.Flags().GetString("rtype")
		ip, _ := cmd.Flags().GetString("ip")
		oldip, _ := cmd.Flags().GetString("oldip")
		proxied, _ := cmd.Flags().GetBool("proxied")
		zone, _ := cmd.Flags().GetString("zone")
		serviceprovider, _ := cmd.Flags().GetString("serviceprovider")

		if dnsapi.CheckEmpty(ip, "ip address", "-i|--ip 123.123.123.123") ||
			dnsapi.CheckEmpty(domain, "domain", "-d|--domain test.domain.tld") {
			return
		}
		if !dnsapi.IsValidIP(ip) {
			return
		}

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

			if !dnsapi.IsValidIP(oldip) {
				return
			}

			if dnsapi.CheckEmpty(zone, "zone", "-z|--zone example.com") ||
				dnsapi.CheckEmpty(oldip, "old ip address", "-o|--oldip 11.11.11.11") ||
				dnsapi.CheckEmpty(rtype, "record type", "-r|--rtype (A|MX|TXT|...)") {
				return
			}

			dnsapi.BindUpdateRecord(dnsconfig.Bind.Server, dnsconfig.Bind.Keyname, dnsconfig.Bind.Hmackey, zone, domain, ip, oldip, rtype)

		}

	},
}

func init() {
	rootCmd.AddCommand(updateRecordCmd)
}
