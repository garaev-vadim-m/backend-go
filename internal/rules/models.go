package rules

type Rule struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type CreateRuleRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdateRuleRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
