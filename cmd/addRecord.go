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
		isvalid := dnsapi.IsValidIP(ip)
		if !isvalid {
			fmt.Printf("IP: %s is no a valid ip address\n", ip)
			return
		}
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
				if zone == "" {
					fmt.Printf("Please provide zone. (--zone|-z example.com)\n")
					return
				}
				if ip == "" {
					fmt.Printf("Please provide ip address. (--ip|-i 123.123.123.123)\n")
					return
				}
				if rtype == "" {
					fmt.Printf("Please provide zone. (--rytpe|-r (A|MX|TXT|...))\n")
					return
				}
				dnsapi.BindInsertRecord(dnsconfig.Bind.Server, dnsconfig.Bind.Keyname, dnsconfig.Bind.Hmackey, zone, domain, ip, rtype)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(addRecordCmd)
}
