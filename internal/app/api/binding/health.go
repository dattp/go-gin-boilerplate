package binding

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Uptime    string `json:"uptime"`
	Timestamp int64  `json:"timestamp"`
} 