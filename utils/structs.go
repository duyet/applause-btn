package utils

import "time"

// Item data structure to be saved to db
type Item struct {
	SourceIP string
	Claps    int
	Clappers []ClapperInfo
}

// ClapperInfo information about clapper
type ClapperInfo struct {
	Email     string
	UID       string
	CreatedAt time.Time
}
