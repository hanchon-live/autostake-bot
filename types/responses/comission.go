package responses

type CommissionRewards struct {
	Commission []Reward `json:"commission"`
}
type CommissionResponse struct {
	Commission CommissionRewards `json:"commission"`
}
