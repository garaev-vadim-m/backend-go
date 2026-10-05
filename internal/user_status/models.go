package userstatus

type UserStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type CreateUserStatusRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdateUserStatusRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
