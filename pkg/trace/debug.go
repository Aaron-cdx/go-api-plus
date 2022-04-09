package trace

// Debug used for dialog debug info record
type Debug struct {
	Key         string      `json:"key"`          // debug key
	Value       interface{} `json:"value"`        // debug value
	CostSeconds float64     `json:"cost_seconds"` // execute cost time(unit second)
}
