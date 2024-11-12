package dnsapi

type Response struct {
	Result []Record `json:"result"`
}

type Record struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Name    string `json:"name"`
	Proxied bool   `json:"proxied"`
}

type DNSRecord struct {
	Type    string `json:"type"`    // Typ des DNS-Records (z.B. "A" oder "CNAME")
	Name    string `json:"name"`    // Domainname des Records
	Content string `json:"content"` // IP-Adresse oder Zielinhalt des Records
	TTL     int    `json:"ttl"`     // Time-To-Live in Sekunden (z.B. 1 für "Auto")
	Proxied bool   `json:"proxied"` // Ob der Record über Cloudflare geleitet werden soll
}

type Config struct {
	Cloudflare struct {
		ZoneId string `yaml:"zoneid"`
		Token  string `yaml:"token"`
	} `yaml:"cloudflare"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Bind struct {
		Server  string `yaml:"server"`
		Keyname string `yaml:"keyname"`
		Hmackey string `yaml:"hmackey"`
	} `yaml:"bind"`
	Batch struct {
		Command  string   `yaml:"command"`
		Provider string   `yaml:"provider"`
		Zone     string   `yaml:"zone"`
		IP       string   `yaml:"ip"`
		Oldip    string   `yaml:"oldip"`
		Proxied  bool     `yaml:"proxied"`
		Rtype    string   `yaml:"rtype"`
		Domains  []string `yaml:"domains"`
	}
}
