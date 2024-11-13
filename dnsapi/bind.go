package dnsapi

import (
	"fmt"
	"log"
	"time"

	"github.com/miekg/dns"
)

func BindInsertRecord(server string, keyName string, keySecret string, zone string, recordName string, ipAddress string, rtype string) {

	msg := new(dns.Msg)
	msg.SetUpdate(zone + ".")

	rr, err := dns.NewRR(fmt.Sprintf("%s. 3600 IN %s %s", recordName, rtype, ipAddress))
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des Records: %v", err)
	}

	msg.Insert([]dns.RR{rr})
	msg.SetTsig(keyName, dns.HmacSHA512, 300, time.Now().Unix())
	client := new(dns.Client)
	client.TsigSecret = map[string]string{keyName: keySecret}

	resp, _, err := client.Exchange(msg, server)
	if err != nil {
		log.Fatalf("Fehler bei DNS-Update: %v", err)
	}

	if resp.Rcode != dns.RcodeSuccess {
		log.Fatalf("DNS-Update fehlgeschlagen: %v", dns.RcodeToString[resp.Rcode])
	}

	fmt.Println("DNS-Update erfolgreich!")

}

func BindDeleteRecord(server string, keyName string, keySecret string, zone string, recordName string, ipAddress string, rtype string) {

	msg := new(dns.Msg)
	msg.SetUpdate(zone + ".")

	rr, err := dns.NewRR(fmt.Sprintf("%s. 3600 IN %s %s", recordName, rtype, ipAddress))
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des Records: %v", err)
	}

	msg.Remove([]dns.RR{rr})
	msg.SetTsig(keyName, dns.HmacSHA512, 300, time.Now().Unix())
	client := new(dns.Client)
	client.TsigSecret = map[string]string{keyName: keySecret}

	resp, _, err := client.Exchange(msg, server)
	if err != nil {
		log.Fatalf("Fehler bei DNS-Update: %v", err)
	}

	if resp.Rcode != dns.RcodeSuccess {
		log.Fatalf("DNS-Update fehlgeschlagen: %v", dns.RcodeToString[resp.Rcode])
	}

	fmt.Println("DNS-Eintrag erfolgreich gelöscht!")
}

func BindUpdateRecord(server string, keyName string, keySecret string, zone string, recordName string, newIPAddress string, oldIPAddress string, rtype string) {

	msg := new(dns.Msg)
	msg.SetUpdate(zone + ".")

	oldRR, err := dns.NewRR(fmt.Sprintf("%s. 3600 IN %s %s", recordName, rtype, oldIPAddress))
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des alten Records: %v", err)
	}

	msg.Remove([]dns.RR{oldRR})

	newRR, err := dns.NewRR(fmt.Sprintf("%s. 3600 IN %s %s", recordName, rtype, newIPAddress))
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des neuen Records: %v", err)
	}

	msg.Insert([]dns.RR{newRR})
	msg.SetTsig(keyName, dns.HmacSHA512, 300, time.Now().Unix())
	client := new(dns.Client)
	client.TsigSecret = map[string]string{keyName: keySecret}

	resp, _, err := client.Exchange(msg, server)
	if err != nil {
		log.Fatalf("Fehler bei DNS-Update: %v", err)
	}

	if resp.Rcode != dns.RcodeSuccess {
		log.Fatalf("DNS-Update fehlgeschlagen: %v", dns.RcodeToString[resp.Rcode])
	}

	fmt.Println("DNS-Update erfolgreich!")
}
