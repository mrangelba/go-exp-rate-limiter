package entities

type RateLimiter struct {
	Key       string `json:"key"`
	Requests  int    `json:"requests"`
	Every     int    `json:"every"`
	Remaining int    `json:"remaining"`
	Reset     int64  `json:"reset"`
}
