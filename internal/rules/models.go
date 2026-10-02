package rules

type Rules struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type CreateRulesRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdateRulesRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
