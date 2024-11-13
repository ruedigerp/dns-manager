package cmd

import (
	"dns-manager/dnsapi"
	"fmt"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "update DNS record",
	Long:  `update DNS record`,
	Run: func(cmd *cobra.Command, args []string) {

		zoneID := dnsconfig.Cloudflare.ZoneId
		token := dnsconfig.Cloudflare.Token

		if dnsconfig.Batch.Command == "add" {

			fmt.Printf("Add command\n")
			fmt.Printf("provider: %s\nzone: %s\nip: %s\n", dnsconfig.Batch.Provider, dnsconfig.Batch.Zone, dnsconfig.Batch.IP)

			if !dnsapi.IsValidIP(dnsconfig.Batch.IP) {
				fmt.Printf("IP: %s is no a valid ip address\n", dnsconfig.Batch.IP)
				return
			}

			for _, domain := range dnsconfig.Batch.Domains {

				fmt.Printf("Add domain: %s\n", domain)

				if dnsconfig.Batch.Provider == "cloudflare" {

					dnsapi.AddRecord(zoneID, token, domain, dnsconfig.Batch.IP, dnsconfig.Batch.IP, dnsconfig.Batch.Proxied)

				} else if dnsconfig.Batch.Provider == "bind" {

					dnsapi.BindInsertRecord(dnsconfig.Bind.Server, dnsconfig.Bind.Keyname, dnsconfig.Bind.Hmackey, dnsconfig.Batch.Zone, domain, dnsconfig.Batch.IP, dnsconfig.Batch.Rtype)

				}

			}

		} else if dnsconfig.Batch.Command == "update" {

			fmt.Printf("update command\n")
			fmt.Printf("provider: %s\nzone: %s\nip: %s\n", dnsconfig.Batch.Provider, dnsconfig.Batch.Zone, dnsconfig.Batch.IP)

			if !dnsapi.IsValidIP(dnsconfig.Batch.IP) {
				fmt.Printf("IP: %s is no a valid ip address\n", dnsconfig.Batch.IP)
				return
			}

			if !dnsapi.IsValidIP(dnsconfig.Batch.Oldip) {
				fmt.Printf("Old-IP: %s is no a valid ip address\n", dnsconfig.Batch.Oldip)
				return
			}

			for _, domain := range dnsconfig.Batch.Domains {

				fmt.Printf("Update domain: %s\n", domain)

				if dnsconfig.Batch.Provider == "cloudflare" {

					dnsapi.UpdateRecord(zoneID, token, domain, dnsconfig.Batch.IP, dnsconfig.Batch.IP, dnsconfig.Batch.Proxied)

				} else if dnsconfig.Batch.Provider == "bind" {

					dnsapi.BindUpdateRecord(dnsconfig.Bind.Server, dnsconfig.Bind.Keyname, dnsconfig.Bind.Hmackey, dnsconfig.Batch.Zone, domain, dnsconfig.Batch.IP, dnsconfig.Batch.Oldip, dnsconfig.Batch.Rtype)

				}

			}
		} else if dnsconfig.Batch.Command == "delete" {

			fmt.Printf("delete command\n")
			fmt.Printf("provider: %s\nzone: %s\nip: %s\n", dnsconfig.Batch.Provider, dnsconfig.Batch.Zone, dnsconfig.Batch.IP)

			for _, domain := range dnsconfig.Batch.Domains {

				fmt.Printf("Delete domain: %s\n", domain)

				if dnsconfig.Batch.Provider == "cloudflare" {

					_, recordID, _ := dnsapi.GetRecordId(zoneID, token, domain)
					dnsapi.DeleteRecord(zoneID, token, recordID)

				} else if dnsconfig.Batch.Provider == "bind" {

					if !dnsapi.IsValidIP(dnsconfig.Batch.IP) {
						fmt.Printf("IP: %s is no a valid ip address\n", dnsconfig.Batch.IP)
						return
					}

					dnsapi.BindDeleteRecord(dnsconfig.Bind.Server, dnsconfig.Bind.Keyname, dnsconfig.Bind.Hmackey, dnsconfig.Batch.Zone, domain, dnsconfig.Batch.IP, dnsconfig.Batch.Rtype)

				}

			}

		}

	},
}

func init() {

	rootCmd.AddCommand(batchCmd)

}
