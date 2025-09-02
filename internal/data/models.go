package data

import (
	"time"
)

type Data struct {
	LastUpdate string      `json:"lastUpdate"`
	Data       []DataEntry `json:"data"`
}

type DataEntry struct {
	IP           string        `json:"ip"`
	ISP          string        `json:"isp"`
	Description  string        `json:"description"`
	StateChanges []StateChange `json:"stateChanges"`
}

type StateChange struct {
	Timestamp time.Time `json:"timestamp"`
	State     bool      `json:"state"`
}

// ParsedData represents the parsed data structure
// An IP is marked as blocked if it's blocked in at least one ISP
type ParsedDataElement struct {
	IP       string
	Provider string
	Blocked  bool
}

type ParsedData struct {
	Elements []ParsedDataElement
}

func ParseData(data *Data) ParsedData {
	var parsed ParsedData
	ipBlockMap := make(map[string]bool)
	ipProviderMap := make(map[string]string)

	for _, entry := range data.Data {
		ip := entry.IP

		// An IP could be tested from different ISPs, but the provider is the same (e.g. Cloudflare)
		// The Provider info comes from the entry.Description
		ipProviderMap[ip] = entry.Description

		// If the IP is already marked as blocked, we don't update it (conservative approach).
		if val, ok := ipBlockMap[ip]; !ok || !val {
			ipBlockMap[ip] = IsBlocked(entry.StateChanges)
		}
	}

	// Now create the ParsedData from the map
	for ip, blocked := range ipBlockMap {
		parsed.Elements = append(parsed.Elements, ParsedDataElement{
			IP:       ip,
			Provider: ipProviderMap[ip],
			Blocked:  blocked,
		})
	}

	return parsed
}

func IsBlocked(stateChanges []StateChange) bool {
	blocked := false
	var lastTimestamp *time.Time
	for _, change := range stateChanges {
		if lastTimestamp == nil || change.Timestamp.After(*lastTimestamp) {
			lastTimestamp = &change.Timestamp
			blocked = change.State
		}
	}
	return blocked
}
