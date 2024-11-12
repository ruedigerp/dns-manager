package dnsapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func GetConfig(configPath string) Config {
	configfile := ""
	if configPath != "" {
		configfile = configPath
	} else {
		fmt.Printf("Using default config: ./config.yaml\n")
		configfile = "config.yaml"
	}
	configFile, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Fatalf("Fehler beim Lesen der Config-Datei: %v", err)
		os.Exit(2)
	}

	// YAML-Daten in die Config-Struktur parsen
	// var config Config
	var dnsconfig Config
	err = yaml.Unmarshal(configFile, &dnsconfig)
	if err != nil {
		log.Fatalf("Fehler beim Parsen der Config-Datei: %v", err)
	}

	// return config.Cloudflare.ZoneId, config.Cloudflare.Token
	return dnsconfig
}
