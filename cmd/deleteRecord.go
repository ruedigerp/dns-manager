package cmd

import (
	"dns-manager/dnsapi"
	"fmt"

	"github.com/spf13/cobra"
)

var deleteRecordCmd = &cobra.Command{
	Use:       "del",
	Short:     "delete Cloudlfare DSN Recrods",
	Long:      `delete Cloudlfare DSN Recrods`,
	ValidArgs: []string{"--domain", "-d", "--ip", "-i", "-type", "-t"},

	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		// rtype, _ := cmd.Flags().GetString("rtype")
		// ip, _ := cmd.Flags().GetString("ip")
		// proxied, _ := cmd.Flags().GetBool("proxied")
		configPath, _ := cmd.Flags().GetString("config")

		zoneID, token = dnsapi.GetConfig(configPath)
		// fmt.Printf("ZonenID: %s, Token: %s\n", zoneID, token)
		if domain == "" {
			fmt.Printf("No domain/subdomain user -d|--domain test.domain.tld")

		} else {
			resp, recordID, _ := dnsapi.GetRecordId(zoneID, token, domain)
			if resp {
				fmt.Printf("RecordID: %s\n", recordID)
				dnsapi.DeleteRecord(zoneID, token, recordID)
				return
			} else {
				fmt.Printf("Record not exists.\n")
				return
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(deleteRecordCmd)

}
