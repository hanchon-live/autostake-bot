package responses

type AttributeResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Index bool   `json:"index"`
}
type EventReponse struct {
	Type       string              `json:"type"`
	Attributes []AttributeResponse `json:"attributes"`
}
type LogResponse struct {
	MsgIndex int64          `json:"msg_index"`
	Log      string         `json:"log"`
	Events   []EventReponse `json:"events"`
}

type TxDetail struct {
	TypeUrl string `json:"type_url"`
	Value   string `json:"value"`
}

type TxResponse struct {
	Height    string         `json:"height"`
	TxHash    string         `json:"txhash"`
	Codespace string         `json:"codespace"`
	Code      int64          `json:"code"`
	Data      string         `json:"data"`
	RawLog    string         `json:"raw_log"`
	Logs      []LogResponse  `json:"logs"`
	Info      string         `json:"info"`
	GasWanted string         `json:"gas_wanted"`
	GasUsed   string         `json:"gas_used"`
	Tx        TxDetail       `json:"tx"`
	Timestamp string         `json:"timestamp"`
	Events    []EventReponse `json:"events"`
}

type BroadcastTxResponse struct {
	Response TxResponse `json:"tx_response"`
}
