package cmd

import (
	"dns-manager/dnsapi"
	"fmt"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var batchCmd = &cobra.Command{
	Use:       "batch",
	Short:     "update DNS record",
	Long:      `update DNS record`,
	ValidArgs: []string{"--domain", "-d", "--ip", "-i", "--oldip", "-o", "-rtype", "-r", "--proxied", "-p", "--zone", "-z"},
	Run: func(cmd *cobra.Command, args []string) {
		zoneID := dnsconfig.Cloudflare.ZoneId
		token := dnsconfig.Cloudflare.Token
		if dnsconfig.Batch.Command == "add" {
			fmt.Printf("Add command\n")
			fmt.Printf("provider: %s\nzone: %s\nip: %s\n", dnsconfig.Batch.Provider, dnsconfig.Batch.Zone, dnsconfig.Batch.IP)

			for _, domain := range dnsconfig.Batch.Domains {
				fmt.Printf("Add domain: %s\n", domain)
				dnsapi.AddRecord(zoneID, token, domain, "A", dnsconfig.Batch.IP, true)
			}
		} else if dnsconfig.Batch.Command == "update" {
			fmt.Printf("update command\n")
			fmt.Printf("provider: %s\nzone: %s\nip: %s\n", dnsconfig.Batch.Provider, dnsconfig.Batch.Zone, dnsconfig.Batch.IP)

			for _, domain := range dnsconfig.Batch.Domains {
				fmt.Printf("Update domain: %s\n", domain)
				dnsapi.UpdateRecord(zoneID, token, domain, "A", dnsconfig.Batch.IP, true)
			}
		} else if dnsconfig.Batch.Command == "delete" {
			fmt.Printf("delete command\n")
			fmt.Printf("provider: %s\nzone: %s\nip: %s\n", dnsconfig.Batch.Provider, dnsconfig.Batch.Zone, dnsconfig.Batch.IP)

			for _, domain := range dnsconfig.Batch.Domains {
				fmt.Printf("Delete domain: %s\n", domain)
				_, recordID, _ := dnsapi.GetRecordId(zoneID, token, domain)
				dnsapi.DeleteRecord(zoneID, token, recordID)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(batchCmd)
}
