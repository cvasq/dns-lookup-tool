package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
)

type DNSresponse struct {
	DNSname string   `json:"dnsname"`
	A       []string `json:"A"`
	AAAA    []string `json:"AAAA"`
	MX      []string `json:"MX"`
	NS      []string `json:"NS"`
}

func resolveDNS(dnsname string) DNSresponse {

	log.Println("Resolving DNS for", dnsname)
	ips, err := net.LookupIP(dnsname)
	if err != nil {
		log.Println("Could not get IPs: ", err)
	}

	mx, err := net.LookupMX(dnsname)
	if err != nil {
		log.Println(err)
	}

	ns, err := net.LookupNS(dnsname)
	if err != nil {
		log.Println("NS Lookup error:", err)
	}

	response := DNSresponse{DNSname: dnsname}

	for _, ip := range ips {
		if ip.To4() != nil {
			response.A = append(response.A, ip.String())
		} else {
			response.AAAA = append(response.AAAA, ip.String())
		}
	}

	for _, ip := range mx {
		response.MX = append(response.MX, fmt.Sprintf("%v\t%v", ip.Pref, ip.Host))
	}

	for _, ip := range ns {
		response.NS = append(response.NS, ip.Host)
	}

	return response
}

func validateDomainName(domain string) bool {

	RegExp := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z
 ]{2,3})$`)

	return RegExp.MatchString(domain)
}
