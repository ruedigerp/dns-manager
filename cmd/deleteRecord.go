package cmd

import (
	"dns-manager/dnsapi"
	"fmt"

	"github.com/spf13/cobra"
)

var deleteRecordCmd = &cobra.Command{
	Use:       "delete",
	Short:     "delete Cloudlfare DSN Recrods",
	Long:      `delete Cloudlfare DSN Recrods`,
	ValidArgs: []string{"--domain", "-d", "--ip", "-i", "--zone", "-z", "--rtype", "-r", "--serviceprovider", "-s"},

	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		rtype, _ := cmd.Flags().GetString("rtype")
		ip, _ := cmd.Flags().GetString("ip")
		zone, _ := cmd.Flags().GetString("zone")
		serviceprovider, _ := cmd.Flags().GetString("serviceprovider")
		isvalid := dnsapi.IsValidIP(ip)
		if !isvalid {
			fmt.Printf("IP: %s is no a valid ip address\n", ip)
			return
		}
		if dnsapi.CheckEmpty(domain, "domain", "-d|--domain test.domain.tld") {
			return
		} else {
			// if domain == "" {
			// 	fmt.Printf("No domain/subdomain user -d|--domain test.domain.tld")

			// } else {
			if serviceprovider == "cloudflare" {
				zoneID := dnsconfig.Cloudflare.ZoneId
				token := dnsconfig.Cloudflare.Token
				resp, recordID, _ := dnsapi.GetRecordId(zoneID, token, domain)
				if resp {
					fmt.Printf("RecordID: %s\n", recordID)
					dnsapi.DeleteRecord(zoneID, token, recordID)
					return
				} else {
					fmt.Printf("Record not exists.\n")
					return
				}
			} else if serviceprovider == "bind" {
				if dnsapi.CheckEmpty(zone, "zone", "-z|--zone example.com") ||
					dnsapi.CheckEmpty(ip, "ip address", "-i|--ip 123.123.123.123") ||
					dnsapi.CheckEmpty(rtype, "record type", "-r|--rtype (A|MX|TXT|...)") {
					return
				}
				// if zone == "" {
				// 	fmt.Printf("Please provide zone. (--zone|-z example.com)\n")
				// 	return
				// }
				// if ip == "" {
				// 	fmt.Printf("Please provide ip address. (--ip|-i 123.123.123.123)\n")
				// 	return
				// }
				// if rtype == "" {
				// 	fmt.Printf("Please provide zone. (--rytpe|-r (A|MX|TXT|...))\n")
				// 	return
				// }
				dnsapi.BindDeleteRecord(dnsconfig.Bind.Server, dnsconfig.Bind.Keyname, dnsconfig.Bind.Hmackey, zone, domain, ip, rtype)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(deleteRecordCmd)

}
