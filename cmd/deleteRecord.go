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
				dnsapi.CheckEmpty(rtype, "record type", "-r|--rtype (A|MX|TXT|...)") {
				return
			}

			dnsapi.BindDeleteRecord(dnsconfig.Bind.Server, dnsconfig.Bind.Keyname, dnsconfig.Bind.Hmackey, zone, domain, ip, rtype)

		}

	},
}

func init() {
	rootCmd.AddCommand(deleteRecordCmd)

}
