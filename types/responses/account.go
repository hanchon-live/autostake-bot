package responses

type PubKeyResponse struct {
	Pubkeytype string `json:"@type"`
	Key        string `json:"key"`
}

type BaseAccount struct {
	Address       string         `json:"address"`
	Pubkey        PubKeyResponse `json:"pub_key"`
	AccountNumber string         `json:"account_number"`
	Sequence      string         `json:"sequence"`
}

type AccountResponse struct {
	AccountType string      `json:"@type"`
	BaseAccount BaseAccount `json:"base_account"`
	CodeHash    string      `json:"code_hash"`
}

type AuthAddressResponse struct {
	Account AccountResponse
}
