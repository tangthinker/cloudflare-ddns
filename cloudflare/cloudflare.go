package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CloudflareClient struct {
	apiToken string
	client   *http.Client
}

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type DNSResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
	Result []struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
		TTL     int    `json:"ttl"`
	} `json:"result"`
}

// 用于创建记录时的响应格式
type CreateDNSResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
	Result struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
		TTL     int    `json:"ttl"`
	} `json:"result"`
}

func NewCloudflareClient(apiToken string) *CloudflareClient {
	return &CloudflareClient{
		apiToken: apiToken,
		client:   &http.Client{},
	}
}

// 获取现有的DNS记录
func (c *CloudflareClient) getDNSRecord(zoneID, domain string) (*DNSRecord, error) {
	log.Printf("Querying existing DNS record for domain %s...", domain)
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?type=AAAA&name=%s", zoneID, domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var response DNSResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if !response.Success {
		if len(response.Errors) > 0 {
			return nil, fmt.Errorf("cloudflare API error: %s", response.Errors[0].Message)
		}
		return nil, fmt.Errorf("unknown cloudflare API error")
	}

	if len(response.Result) == 0 {
		log.Printf("No existing DNS record found for domain %s", domain)
		return nil, nil // 记录不存在
	}

	log.Printf("Found existing DNS record: Type=%s, Name=%s, Content=%s, TTL=%d",
		response.Result[0].Type,
		response.Result[0].Name,
		response.Result[0].Content,
		response.Result[0].TTL,
	)

	return &DNSRecord{
		ID:      response.Result[0].ID,
		Type:    response.Result[0].Type,
		Name:    response.Result[0].Name,
		Content: response.Result[0].Content,
		TTL:     response.Result[0].TTL,
	}, nil
}

// 创建新的DNS记录
func (c *CloudflareClient) createDNSRecord(zoneID string, record DNSRecord) error {
	log.Printf("Creating new DNS record for domain %s with IPv6 address: %s", record.Name, record.Content)

	jsonData, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("error marshaling DNS record: %v", err)
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Sending create request to Cloudflare...")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var response CreateDNSResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	if !response.Success {
		if len(response.Errors) > 0 {
			return fmt.Errorf("cloudflare API error: %s", response.Errors[0].Message)
		}
		return fmt.Errorf("unknown cloudflare API error")
	}

	log.Printf("Successfully created DNS record for domain %s", record.Name)
	return nil
}

// 更新现有的DNS记录
func (c *CloudflareClient) updateDNSRecord(zoneID, recordID string, record DNSRecord) error {
	log.Printf("Updating existing DNS record for domain %s with IPv6 address: %s", record.Name, record.Content)

	jsonData, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("error marshaling DNS record: %v", err)
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Sending update request to Cloudflare...")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var response CreateDNSResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	if !response.Success {
		if len(response.Errors) > 0 {
			return fmt.Errorf("cloudflare API error: %s", response.Errors[0].Message)
		}
		return fmt.Errorf("unknown cloudflare API error")
	}

	log.Printf("Successfully updated DNS record for domain %s", record.Name)
	return nil
}

// 更新DNS记录
func (c *CloudflareClient) UpdateDNSRecord(zoneID, domain, ipv6 string) error {
	log.Printf("Starting DNS record update for domain %s, new IPv6 address: %s", domain, ipv6)

	// 获取现有记录
	existingRecord, err := c.getDNSRecord(zoneID, domain)
	if err != nil {
		return fmt.Errorf("error getting existing record: %v", err)
	}

	// 如果记录存在且内容相同，不需要更新
	if existingRecord != nil && existingRecord.Content == ipv6 {
		log.Printf("Existing IPv6 address %s for domain %s matches new address, no update needed", ipv6, domain)
		return nil // 记录已存在且内容相同，无需更新
	}

	record := DNSRecord{
		Type:    "AAAA",
		Name:    domain,
		Content: ipv6,
		TTL:     1, // Auto TTL
	}

	// 如果记录不存在，创建新记录
	if existingRecord == nil {
		return c.createDNSRecord(zoneID, record)
	}

	// 如果记录存在但内容不同，更新现有记录
	if existingRecord.Content != ipv6 {
		log.Printf("Record update required: changing from %s to %s", existingRecord.Content, ipv6)
		return c.updateDNSRecord(zoneID, existingRecord.ID, record)
	}

	return nil
}
