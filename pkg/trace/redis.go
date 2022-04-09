package trace

// Redis used for dialog redis info record
type Redis struct {
	Timestamp   string  `json:"timestamp"`    // time format: 2006-01-02 15:04:05
	Handle      string  `json:"handle"`       // operation SET/GET...
	Key         string  `json:"key"`          // Key for redis handle
	Value       string  `json:"value"`        // Value for redis handle
	TTL         float64 `json:"ttl"`          // timeout(unit minute)
	CostSeconds float64 `json:"cost_seconds"` // execute cost time(unit second)
}
