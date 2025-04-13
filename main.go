package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tangthinker/cloudflare-ddns/cloudflare"
	"github.com/tangthinker/cloudflare-ddns/network"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Cloudflare struct {
		APIToken string   `yaml:"api_token"`
		ZoneID   string   `yaml:"zone_id"`
		Domains  []string `yaml:"domains"`
	} `yaml:"cloudflare"`
	Network struct {
		Interface string `yaml:"interface"`
	} `yaml:"network"`
	Interval struct {
		Success int `yaml:"success"`
		Error   int `yaml:"error"`
	} `yaml:"interval"`
}

func loadConfig(path string) (*Config, error) {
	config := &Config{}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return config, nil
}

func updateDomains(cf *cloudflare.CloudflareClient, config *Config, ipv6 string) error {
	var lastErr error
	successCount := 0

	for _, domain := range config.Cloudflare.Domains {
		err := cf.UpdateDNSRecord(config.Cloudflare.ZoneID, domain, ipv6)
		if err != nil {
			log.Printf("Failed to update domain %s: %v", domain, err)
			lastErr = err
		} else {
			successCount++
		}
	}

	if lastErr != nil {
		if successCount > 0 {
			log.Printf("Partially successful: updated %d out of %d domains", successCount, len(config.Cloudflare.Domains))
		}
		return fmt.Errorf("error updating some domains: %v", lastErr)
	}

	return nil
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cf := cloudflare.NewCloudflareClient(config.Cloudflare.APIToken)
	net := network.NewNetworkManager()

	for {
		ipv6, err := net.GetIPv6Address(config.Network.Interface)
		if err != nil {
			log.Printf("Error getting IPv6 address: %v", err)
			time.Sleep(time.Duration(config.Interval.Error) * time.Second)
			continue
		}

		err = updateDomains(cf, config, ipv6)
		if err != nil {
			log.Printf("Error updating domains: %v", err)
			time.Sleep(time.Duration(config.Interval.Error) * time.Second)
			continue
		}

		log.Printf("Successfully updated all domains with IPv6 address: %s", ipv6)
		time.Sleep(time.Duration(config.Interval.Success) * time.Second)
	}
}
