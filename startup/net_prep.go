package startup

import (
	"log"
	"net"
	"time"

	"github.com/google/uuid"
)

// Probe struct
type probe struct {
	region string
	host   string
}

var probes = []probe{
	// North America
	{"North America", "google.com:80"},
	{"North America", "cloudflare.com:80"},

	// South America
	{"South America", "sa-east-1.amazonaws.com:80"},

	// Europe
	{"Europe", "eu-west-1.amazonaws.com:80"},
	{"Europe", "google.de:80"},

	// Russia
	{"Russia", "yandex.ru:80"},

	// Asia
	{"Asia", "ap-northeast-1.amazonaws.com:80"},
	{"Asia", "google.co.in:80"},
	{"Asia", "google.co.jp:80"},
	{"Asia", "google.com.sg:80"},

	// Middle East
	{"Middle East", "me-central-1.amazonaws.com:80"},

	// Africa
	{"Africa", "af-south-1.amazonaws.com:80"},

	// Oceania
	{"Oceania", "google.com.au:80"},
}

// Returns a string-cased UUIDv4, the servers unique ID
// (used for server listings and identification)
func GenerateServerGUID() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Failed to generate server GUID: %v", err)
	}

	return uuid.String()
}

// Ping several servers around the world and determine the best region to display
// Valid regions:
//   - North Amercica (Canada, US, Mexico)
//   - South America (includes the caribbean region)
//   - Europe (excludes Russia)
//   - Russia
//   - Asia (includes India, China, Japan, Korea, SE Asia)
//   - Middle East (includes Iran, Saudi Arabia, Qatar, etc)
//   - Africa (The entire continent for now)
//   - Oceania (Australia, New Zealand, Pacific Islands)
//
// Note: This func can be completely overridden by a cmdline arg; `--force-region=<region>`
func DetermineRegion() string {
	timeout := time.Millisecond * 500
	results := make(map[string][]time.Duration)

	// Probe each host via TCP
	for _, p := range probes {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", p.host, timeout)
		if err != nil {
			continue
		}
		conn.Close()
		results[p.region] = append(results[p.region], time.Since(start))
	}

	// Find lowest avg latency
	bestRegion := "UNKNOWN"
	bestLatency := time.Duration(1<<63 - 1)

	for region, latencies := range results {
		if len(latencies) == 0 {
			continue
		}

		var total time.Duration
		for _, latency := range latencies {
			total += latency
		}

		avg := total / time.Duration(len(latencies))
		if avg < bestLatency {
			bestLatency = avg
			bestRegion = region
		}
	}

	return bestRegion
}

// Easy ping, check if we can see IPs like 1.1.1.1 or 8.8.8.8 and ret a bool
func CheckNetConnection() (bool, error) {
	timeout := 500 * time.Millisecond
	targets := []string{
		"1.1.1.1:53", // Cloudflare DNS
		"8.8.8.8:53", // Google DNS
	}

	for _, addr := range targets {
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err == nil {
			conn.Close()
			return true, nil
		}
	}

	return false, nil
}
