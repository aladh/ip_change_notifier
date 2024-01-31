package main

import (
	"flag"
	"fmt"
	"github.com/aladh/ip_change_notifier/webhook"
	"log"
	"net"
	"os"
)

var domain string
var previousIP string
var webhookURL string
var outputFilename string

func init() {
	flag.StringVar(&domain, "d", "", "The domain to look up")
	flag.StringVar(&previousIP, "p", "", "The previous IP address")
	flag.StringVar(&webhookURL, "w", "", "A URL to send webhooks when IP changes")
	flag.StringVar(&outputFilename, "o", "", "The file to write the IP address to")
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

	if outputFilename != "" {
		err = writeToFile(outputFilename, ipString)
		if err != nil {
			log.Fatalln("Error writing to file:", err)
		}
	}
}

func writeToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
