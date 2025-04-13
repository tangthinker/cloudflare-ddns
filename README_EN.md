# Cloudflare DDNS

A Go-based Dynamic DNS client for Cloudflare that automatically updates IPv6 DNS records based on your local network interface.

## Features

- Automatically detects IPv6 address from specified network interface
- Updates multiple Cloudflare DNS records using their API
- Configurable update intervals
- Error handling with retry mechanism
- YAML-based configuration
- Support for multiple domains
- Smart record management (create/update based on existence)

## Configuration

Create a `config.yaml` file with the following structure:

```yaml
cloudflare:
  api_token: "your-api-token-here"  # Get from Cloudflare dashboard
  zone_id: "your-zone-id-here"      # Get from Cloudflare dashboard
  domains:                          # List of domains to update
    - "sub1.example.com"
    - "sub2.example.com"
    - "sub3.example.com"

network:
  interface: "en0"  # Network interface name to monitor

interval:
  success: 600  # Wait time after successful update (seconds)
  error: 10     # Wait time after error (seconds)
```

## Prerequisites

1. A Cloudflare account with API access
2. A domain managed by Cloudflare
3. IPv6 connectivity
4. Go 1.24 or higher

## Getting Started

1. Get your Cloudflare API token:
   - Log in to Cloudflare dashboard
   - Go to "My Profile" > "API Tokens"
   - Create a new token with DNS edit permissions
   - Copy the token

2. Get your Zone ID:
   - Select your domain in Cloudflare dashboard
   - Find "Zone ID" in the right sidebar
   - Copy the ID

3. Update the `config.yaml` file:
   - Paste your API token
   - Paste your Zone ID
   - Add your domains to the list
   - Set your network interface name
   - Adjust intervals if needed

4. Run the program:
```bash
go run main.go
```

Or build and run:
```bash
go build
./cloudflare-ddns
```

## How It Works

1. The program reads the configuration from `config.yaml`
2. It checks the specified network interface for IPv6 addresses every 10 minutes (configurable)
3. For each domain in the configuration:
   - Checks if an AAAA record exists
   - Creates a new record if none exists
   - Updates the record if the IPv6 address has changed
4. If any errors occur, it will retry after 10 seconds (configurable)
5. Logs all successful updates and errors

## Error Handling

- If getting IPv6 address fails, retries after error interval
- If updating a domain fails, continues with other domains
- Reports partial success if some domains are updated successfully
- Detailed error messages in logs

## Logging

The program provides detailed logs for:
- IPv6 address detection
- DNS record queries
- Record creation/updates
- Success/failure status
- Error details

## Contributing

Feel free to submit issues and enhancement requests! 