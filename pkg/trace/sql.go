package trace

// SQL used for sql dialog info record
type SQL struct {
	Timestamp   string  `json:"timestamp"`     // time, format: 2006-01-02 15:04:05
	Stack       string  `json:"stack"`         // file address and line number
	SQL         string  `json:"sql"`           // SQL statement
	Rows        int64   `json:"rows_affected"` // sql affected rows
	CostSeconds float64 `json:"cost_seconds"`  // execute cost time(unit second)
}
