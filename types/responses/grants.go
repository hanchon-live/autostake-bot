package responses

type GrantsReponse struct {
	Grants     []Grant    `json:"grants"`
	Pagination Pagination `json:"pagination"`
}

type Grant struct {
	Granter       string        `json:"granter"`
	Grantee       string        `json:"grantee"`
	Authorization Authorization `json:"authorization"`
	Expiration    string        `json:"expiration"`
}

type Authorization struct {
	TypeUrl string `json:"@type"`
	Value   string `json:"msg"`
}

type Pagination struct {
	NextKey string `json:"next_key"`
	Total   string `json:"total"`
}
