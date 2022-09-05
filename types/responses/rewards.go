package responses

type Reward struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Rewards struct {
	ValidatorAddress string   `json:"validator_address"`
	Reward           []Reward `json:"reward"`
}

type RewardsResponse struct {
	Rewards []Rewards `json:"rewards"`
	Total   []Reward  `json:"total"`
}
