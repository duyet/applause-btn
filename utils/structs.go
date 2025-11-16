package utils

import "time"

// Item data structure to be saved to db
type Item struct {
	SourceIP  string          // Deprecated: kept for backward compatibility, use SourceIPs
	SourceIPs map[string]bool // Set of IP addresses that have clapped
	Claps     int             // Total number of claps
	Clappers  []ClapperInfo   // Information about who clapped
}

// ClapperInfo information about clapper
type ClapperInfo struct {
	Email     string
	UID       string
	CreatedAt time.Time
}

// HasClappedFrom checks if an IP has already clapped
func (i *Item) HasClappedFrom(ip string) bool {
	if i.SourceIPs == nil {
		// For backward compatibility, check legacy SourceIP field
		return i.SourceIP == ip
	}
	return i.SourceIPs[ip]
}

// AddClapFrom records a clap from an IP address
func (i *Item) AddClapFrom(ip string, claps int) {
	if i.SourceIPs == nil {
		i.SourceIPs = make(map[string]bool)
	}
	i.SourceIPs[ip] = true
	i.Claps += claps
}
