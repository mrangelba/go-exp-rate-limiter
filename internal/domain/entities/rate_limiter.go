package entities

import "encoding/json"

type RateLimiter struct {
	Key       string `json:"key"`
	Requests  int    `json:"requests"`
	Every     int    `json:"every"`
	Remaining int    `json:"remaining"`
	Reset     int64  `json:"reset"`
}

func (r RateLimiter) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *RateLimiter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}
