package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/aladh/ip_change_notifier/webhook"
)

const outputFilename = "ip.txt"

var domain string
var previousIP string
var webhookURL string

func init() {
	flag.StringVar(&domain, "d", "", "The domain to look up")
	flag.StringVar(&previousIP, "p", "", "The previous IP address")
	flag.StringVar(&webhookURL, "w", "", "A URL to send webhooks when IP changes")
	flag.Parse()
}

func main() {
	if domain == "" {
		log.Fatalln("Please provide a domain with -d")
	}

	log.Println("Looking up IP address for:", domain)

	ips, err := net.LookupIP(domain)
	if err != nil {
		log.Fatalln("Error looking up IP address:", err)
	}

	ipString := ips[0].String()
	if previousIP == ipString {
		log.Println("IP address has not changed")
		return
	}

	message := fmt.Sprintf("IP address changed to: %s", ipString)
	log.Println(message, ips)

	err = webhook.Send(webhookURL, message)
	if err != nil {
		log.Fatalln("Error sending webhook:", err)
	}

	file, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalln("Error creating file:", err)
	}
	defer file.Close()

	_, err = file.WriteString(ipString)
}
