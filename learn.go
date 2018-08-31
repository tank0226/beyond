package main

import (
	"flag"
	"net"
	"net/http/httputil"
	"net/url"
	"strings"
)

var (
	learnNexthops = flag.Bool("learn-nexthops", true, "set false to require explicit whitelisting")

	learnHTTPSPorts = flag.String("learn-https-ports", "443,4443,8443,9443", "try learning these backend HTTPS ports (csv)")
	learnHTTPPorts  = flag.String("learn-http-ports", "80,8080,6000,6060,7000,8000,9000,9200,15672", "after HTTPS, try these HTTP ports (csv)")
)

func learn(host string) *httputil.ReverseProxy {
	newbase := learnHostScheme(host)
	if newbase == "" {
		return nil
	}
	u, err := url.Parse(newbase)
	if err != nil {
		return nil
	}
	return httputil.NewSingleHostReverseProxy(u)
}

func learnHostScheme(hostname string) string {
	for _, httpsPort := range strings.Split(*learnHTTPSPorts, ",") {
		c, err := net.Dial("tcp", hostname+":"+httpsPort)
		if err == nil {
			c.Close()
			if httpsPort == "443" {
				return "https://" + hostname
			}
			return "https://" + hostname + ":" + httpsPort
		}
	}
	for _, httpPort := range strings.Split(*learnHTTPPorts, ",") {
		c, err := net.Dial("tcp", hostname+":"+httpPort)
		if err == nil {
			c.Close()
			if httpPort == "80" {
				return "http://" + hostname
			}
			return "http://" + hostname + ":" + httpPort
		}
	}
	return ""
}
