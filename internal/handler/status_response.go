package handler

type StatusResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}
