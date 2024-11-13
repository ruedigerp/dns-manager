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

		if dnsapi.CheckEmpty(ip, "ip address", "-i|--ip 123.123.123.123") ||
			dnsapi.CheckEmpty(domain, "domain", "-d|--domain test.domain.tld") {
			return
		}

		if !dnsapi.IsValidIP(ip) {
			fmt.Printf("IP: %s is no a valid ip address\n", ip)
			return
		}

		if serviceprovider == "cloudflare" {

			zoneID := dnsconfig.Cloudflare.ZoneId
			token := dnsconfig.Cloudflare.Token

			resp, id, _ := dnsapi.GetRecordId(zoneID, token, domain)
			if resp {

				fmt.Printf("Record exists. (%s)\n", id)
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

			if dnsapi.CheckEmpty(zone, "zone", "-z|--zone example.com") ||
				dnsapi.CheckEmpty(rtype, "record type", "-r|--rtype (A|MX|TXT|...)") {
				return
			}

			dnsapi.BindInsertRecord(dnsconfig.Bind.Server, dnsconfig.Bind.Keyname, dnsconfig.Bind.Hmackey, zone, domain, ip, rtype)

		}

	},
}

func init() {
	rootCmd.AddCommand(addRecordCmd)
}

// switch serviceprovider {
// case "cloudflare":
// 	handleCloudflareRecord(domain, rtype, ip, proxied)
// case "bind":
// 	if dnsapi.CheckEmpty(zone, "zone", "-z|--zone example.com") {
// 		return
// 	}
// 	dnsapi.BindInsertRecord(dnsconfig.Bind.Server, dnsconfig.Bind.Keyname, dnsconfig.Bind.Hmackey, zone, domain, ip, rtype)
// default:
// 	fmt.Println("Unsupported service provider. Please specify either 'cloudflare' or 'bind'.")
// }

// Funktion zur zentralen Validierung und Abruf von Flags
// func getFlag(cmd *cobra.Command, name, usageHint string) string {
//     value, _ := cmd.Flags().GetString(name)
//     if dnsapi.CheckEmpty(value, name, usageHint) {
//         return ""
//     }
//     return value
// }

// // Separater Handler für Cloudflare-Logik
// func handleCloudflareRecord(domain, rtype, ip string, proxied bool) {
//     zoneID := dnsconfig.Cloudflare.ZoneId
//     token := dnsconfig.Cloudflare.Token

//     recordExists, _, err := dnsapi.GetRecordId(zoneID, token, domain)
//     if err != nil {
//         fmt.Printf("Error checking record: %v\n", err)
//         return
//     }
//     if recordExists {
//         fmt.Println("Record exists.")
//         return
//     }

//     dnsapi.AddRecord(zoneID, token, domain, rtype, ip, proxied)
//     _, recordID, err := dnsapi.GetRecordId(zoneID, token, domain)
//     if err != nil {
//         fmt.Printf("Error retrieving record ID: %v\n", err)
//     } else {
//         fmt.Printf("RecordID: %s\n", recordID)
//     }
// }
