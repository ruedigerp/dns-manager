package dnsapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func AddRecord(zoneID string, token string, domain string, rtype string, ip string, proxied bool) (string, error) {

	var msg = ""
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)

	dnsRecord := DNSRecord{
		Type:    rtype,
		Name:    domain,
		Content: ip,
		TTL:     1,
		Proxied: false,
	}

	jsonData, err := json.Marshal(dnsRecord)
	if err != nil {
		fmt.Println("Fehler beim Erstellen der JSON-Daten:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Fehler beim Erstellen des POST-Requests:", err)
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Fehler beim Senden des POST-Requests:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		code := "200"
		msg = fmt.Sprintf("DNS-Record erfolgreich angelegt. %s\n", code)
	} else {
		fmt.Printf("Fehler beim Anlegen des DNS-Records. Status Code: %d\n", resp.StatusCode)
	}
	return msg, nil
}

func UpdateRecord(zoneID string, token string, domain string, rtype string, ip string, proxied bool) (string, error) {

	exists, recordID, _ := GetRecordId(zoneID, token, domain)

	var msg = ""

	if exists {
		url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)

		dnsRecord := DNSRecord{
			Type:    rtype,
			Name:    domain,
			Content: ip,
			TTL:     1,
			Proxied: proxied,
		}

		jsonData, err := json.Marshal(dnsRecord)
		if err != nil {
			fmt.Println("Fehler beim Erstellen der JSON-Daten:", err)
			return "", err
		}

		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Fehler beim Erstellen des PUT-Requests:", err)
			return "", err
		}

		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Fehler beim Senden des PUT-Requests:", err)
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			msg = "DNS-Record erfolgreich aktualisiert.\n"
		} else {
			fmt.Printf("Fehler beim Aktualisieren des DNS-Records. Status Code: %d\n", resp.StatusCode)
		}

	}

	return msg, nil

}

func GetRecord(zoneID string, token string, domain string) (bool, string, string, error) {

	url := "https://api.cloudflare.com/client/v4/zones/" + zoneID + "/dns_records?type=A&name=" + domain

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Fehler beim Erstellen des Requests:", err)
		return false, "", "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, "Fehler beim Senden des Requests:", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "Fehler beim Lesen der Antwort:", "", err
	}

	var response Response
	var id string
	var msg string

	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, "Fehler beim Unmarshalling der JSON-Antwort", "", err
	}

	if len(response.Result) > 0 {
		id = response.Result[0].ID
		msg = fmt.Sprintf("%s %s %s\n", response.Result[0].Name, response.Result[0].Type, response.Result[0].Content)
	} else {
		return false, "Kein Record gefunden", "", err
	}

	return true, string(id), msg, nil

}

func GetRecordId(zoneID string, token string, domain string) (bool, string, error) {

	url := "https://api.cloudflare.com/client/v4/zones/" + zoneID + "/dns_records?type=A&name=" + domain

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Fehler beim Erstellen des Requests:", err)
		return false, "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, "Fehler beim Senden des Requests:", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "Fehler beim Lesen der Antwort:", err
	}

	var response Response
	var id string

	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, "Fehler beim Unmarshalling der JSON-Antwort", err
	}

	if len(response.Result) > 0 {
		id = response.Result[0].ID
	} else {
		return false, "Kein Record gefunden", err
	}

	return true, string(id), nil

}

func DeleteRecord(zoneID string, token string, recordID string) (bool, error) {

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Fehler beim Erstellen des DELETE-Requests:", err)
		return false, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Fehler beim Senden des DELETE-Requests:", err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("DNS-Record erfolgreich gelöscht.")
	} else {
		fmt.Printf("Fehler beim Löschen des DNS-Records. Status Code: %d\n", resp.StatusCode)
		return false, err
	}

	return true, nil

}
