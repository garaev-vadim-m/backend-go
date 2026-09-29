package models

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type CreateStatus struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdateStatusRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
