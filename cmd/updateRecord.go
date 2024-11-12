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
		if domain == "" {
			fmt.Printf("No domain/subdomain user -d|--domain test.domain.tld")

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
				if zone == "" {
					fmt.Printf("Please provide zone. (--zone|-z example.com)\n")
					return
				}
				if ip == "" {
					fmt.Printf("Please provide ip address. (--ip|-i 123.123.123.123)\n")
					return
				}
				if oldip == "" {
					fmt.Printf("Please provide old ip address. (--oldip|-o 11.11.11.11)\n")
					return
				}
				if rtype == "" {
					fmt.Printf("Please provide zone. (--rytpe|-r (A|MX|TXT|...))\n")
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
