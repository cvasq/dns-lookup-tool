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

	log.Println("Resolving DNS")
	ips, err := net.LookupIP(dnsname)
	if err != nil {
		log.Println("Could not get IPs: ", err)
	}

	mx, err := net.LookupMX(dnsname)
	if err != nil {
		log.Println(err)
	}

	nss, err := net.LookupNS(dnsname)
	if err != nil {
		log.Println("NS Lookup error:", err)
	}
	if len(nss) == 0 {
		log.Println("no record")
	}
	for _, ns := range nss {
		log.Printf("%s\n", ns.Host)
	}

	ss := DNSresponse{DNSname: dnsname}

	for _, ip := range ips {
		if ip.To4() != nil {
			ss.A = append(ss.A, ip.String())
		} else {
			ss.AAAA = append(ss.AAAA, ip.String())
		}
	}

	for _, ip := range mx {
		ss.MX = append(ss.MX, fmt.Sprintf("%v\t%v", ip.Pref, ip.Host))
	}

	for _, ip := range nss {
		ss.NS = append(ss.NS, ip.Host)
	}

	return ss
}

func validateDomainName(domain string) bool {

	RegExp := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z
 ]{2,3})$`)

	return RegExp.MatchString(domain)
}
