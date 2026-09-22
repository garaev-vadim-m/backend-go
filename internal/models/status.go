package models

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateStatus struct {
	Name string `json:"name"`
}

type UpdateStatusRequest struct {
	Name string `json:"name"`
}
