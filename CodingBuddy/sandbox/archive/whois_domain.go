package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type WhoisInfo struct {
	Domain       string
	Registrar    string
	CreationDate string
	ExpiryDate   string
	NameServers  []string
	Authoritative bool
}

func WhoisDomain(domain string, whoisServer string) (*WhoisInfo, error) {
	var attempt int
	var maxRetries = 5
	var waitTime = 1 * time.Second

	for attempt < maxRetries {
		conn, err := net.Dial("tcp", whoisServer+":43")
		if err != nil {
			if attempt == maxRetries-1 {
				return nil, errors.New("failed to connect to WHOIS server")
			}
			time.Sleep(waitTime)
			waitTime *= 2
			attempt++
			continue
		}
		defer conn.Close()

		_, err = conn.Write([]byte(domain + "\r\n"))
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(conn)
		response := ""
		for scanner.Scan() {
			response += scanner.Text() + "\n"
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		whoisInfo, err := parseWhoisResponse(response)
		if err != nil {
			return nil, err
		}

		if whoisInfo.Authoritative {
			return whoisInfo, nil
		}

		referredServer := extractReferredServer(response)
		if referredServer == "" {
			return nil, errors.New("could not find referred WHOIS server")
		}

		return WhoisDomain(domain, referredServer)
	}

	return nil, errors.New("exceeded maximum retries")
}

func parseWhoisResponse(response string) (*WhoisInfo, error) {
	whoisInfo := &WhoisInfo{}
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Domain Name:") {
			whoisInfo.Domain = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Registrar:") {
			whoisInfo.Registrar = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Creation Date:") {
			whoisInfo.CreationDate = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Expiry Date:") {
			whoisInfo.ExpiryDate = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Name Server:") {
			whoisInfo.NameServers = append(whoisInfo.NameServers, strings.TrimSpace(strings.Split(line, ":")[1]))
		} else if strings.Contains(line, "Authoritative:") {
			whoisInfo.Authoritative = true
		}
	}
	return whoisInfo, nil
}

func extractReferredServer(response string) string {
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Refer:") {
			return strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	return ""
}

func main() {
	domain := "example.com"
	whoisServer := "whois.verisign-grs.com"
	info, err := WhoisDomain(domain, whoisServer)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Whois Info: %+v\n", info)
}
