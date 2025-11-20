[![Build, Package, Release](https://github.com/ruedigerp/dns-manager/actions/workflows/main.yaml/badge.svg)](https://github.com/ruedigerp/dns-manager/actions/workflows/main.yaml)

# dns-manager 

A command-line tool to manage DNS records in Cloudflare or Bind DNS server zones. With DNS Manager, you can create, update, retrieve, and delete DNS records easily, and even perform batch operations across multiple records.

## Install

To install DNS Manager via Homebrew:

    brew tap ruedigerp/dns-manager
    brew install dns-manager

## Config

The default configuration file is config.yaml, located in the current directory. You can also use separate configuration files for different zones as needed.

### Example Configurations

config-example.com.yaml:

    cloudflare:
    zoneid: "XXXXXYourCloudZoneIDXXXXX"
    token: "XXXXXYourCloudflareAPIToken"

    bind:
    server: "ns01.example.com:53"
    keyname: "update-key-zone1."
    hmackey: "D/bOz+8YKMZEKyrtBcWythudSrSW76IPH8SnXfvqgL147IcLqPi1wUaktQZIa9AHVsMDQy49h4fublwwoT5Ruu=="

This file includes both Cloudflare and Bind DNS settings, allowing you to specify zone IDs, API tokens, and update keys for secure access.

To manage multiple zones, create additional config files like config-example.net.yaml with their specific credentials.
    

## Managing DNS Records in Cloudflare

Add a New DNS Record: 

    #> dns-manager add -d test1.example.com -s cloudflare -c config-example.com.yaml -r A -i 123.123.123.123 -p false

Output:

    RecordID: 6728448f8340dfd2da753557a70d3670

Retrieve a DNS Record: 

    #> dns-manager get -d test1.example.com -s cloudflare -c config-example.com.yaml

Output:

    RecordID: 6728448f8340dfd2da753557a70d3670
    test1.example.com A 123.123.123.123

Update a DNS Record:

    ❯ dns-manager update -d test1.example.com -s cloudflare -c config-example.com.yaml -r A -i 111.111.111.111 -p false

Output:

    RecordID: 258cd278b17f9a3d35f52a9eb9e58ab6

Verify the update:

    ❯ dns-manager get -d test1.example.com -s cloudflare -c config-example.com.yaml

Output:

    RecordID: 258cd278b17f9a3d35f52a9eb9e58ab6
    test1.example.com A 111.111.111.111

Delete a DNS Record: 

    #> dns-manager delete -d test1.example.com -s cloudflare -c config-example.com.yaml

Output:

    RecordID: 6728448f8340dfd2da753557a70d3670
    DNS-Record deleted.

## Managing DNS Records in Bind

Add a DNS Record:

    ❯ dns-manager add -d test1.zone-in-bind.net -s bind -c config-example.com.yaml -r A -i 123.123.123.123 -z zone-in-bind.net

Output:

    DNS-Update ok!

 Verify the update:

    ❯ dig +short test1.zone-in-bind.net
    123.123.123.123

Update DNS Record: 

    ❯ dns-manager update -d test1.zone-in-bind.net -s bind -c config-example.com.yaml -r A -o 123.123.123.123 -i 111.111.111.111 -z zone-in-bind.net

Output:

    DNS-Update ok!

Verify the update:

    ❯ dig +short test1.zone-in-bind.net 
    111.111.111.111

Delete DNS Record:

    ❯ dns-manager delete -d test1.zone-in-bind.net -s bind -c config-example.com.yaml -r A -i 111.111.111.111 -z zone-in-bind.net

Output:

    DNS-Eintrag deleted!

Verify the update:

    ❯ dig +short test1.zone-in-bind.net 

## Batch Mode (Currently Cloudflare Only)

Batch mode allows you to create, update, or delete multiple DNS records at once.

To enable batch mode, add the following section to your configuration file:

    batch:
    command: add
    provider: cloudflare
    zone: example.com
    ip: 123.123.123.123
    oldip: 122.122.122.122
    proxied: true
    rtype: A
    domains: 
        - api.example.com
        - user.example.com
        - ws.example.com
        - app.example.com

Batch Command Options

	•	command: Specifies the batch operation to perform. Options are:
	    •	add
	    •   update
	    •	delete
	•	zone: Only required for Bind configurations to specify the zone.
	•	proxied: Only used for Cloudflare to enable or disable proxying.


