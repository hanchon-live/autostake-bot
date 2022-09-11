package blockchain

import (
	"encoding/json"
	"fmt"

	"github.com/hanchon-live/autostake-bot/internal/requester"
	"github.com/hanchon-live/autostake-bot/types/responses"
)

func GetAccountFromBlockchain(address string) (responses.AuthAddressResponse, error) {
	url := "/cosmos/auth/v1beta1/accounts/" + address
	if resp, err := requester.MakeGetRequest("rest", url); err != nil {
		return responses.AuthAddressResponse{}, fmt.Errorf("Failed to get the address: %q\n", err)
	} else {
		m := &responses.AuthAddressResponse{}
		err = json.Unmarshal([]byte(resp), m)
		if err != nil {
			return responses.AuthAddressResponse{}, fmt.Errorf("Error decoding response: %q", err)

		}
		return *m, nil
	}
}

func GetDistributionRewards(address string) (responses.RewardsResponse, error) {
	url := "/cosmos/distribution/v1beta1/delegators/" + address + "/rewards"
	if resp, err := requester.MakeGetRequest("rest", url); err != nil {
		return responses.RewardsResponse{}, fmt.Errorf("Failed to get the rewards: %q\n", err)
	} else {
		m := &responses.RewardsResponse{}
		err = json.Unmarshal([]byte(resp), m)
		if err != nil {
			return responses.RewardsResponse{}, fmt.Errorf("Error decoding rewards response: %q", err)
		}
		return *m, nil
	}
}

func GetCommission(address string) (responses.CommissionResponse, error) {
	url := "/cosmos/distribution/v1beta1/validators/" + address + "/commission"
	if resp, err := requester.MakeGetRequest("rest", url); err != nil {
		return responses.CommissionResponse{}, fmt.Errorf("Failed to get the comission: %q\n", err)
	} else {
		m := &responses.CommissionResponse{}
		err = json.Unmarshal([]byte(resp), m)
		if err != nil {
			return responses.CommissionResponse{}, fmt.Errorf("Error decoding comission response: %q", err)
		}
		return *m, nil
	}
}

func GetValidator(address string) (string, error) {
	if address == "evmos1dgpv4leszpeg2jusx2xgyfnhdzghf3rfzw06t3" {
		return "evmosvaloper1dgpv4leszpeg2jusx2xgyfnhdzghf3rf0qq22v", nil
	} else if address == "evmos197ahcv2x9jj0nmvnen4sqqfffhygjga7wc7qkp" {
		return "evmosvaloper197ahcv2x9jj0nmvnen4sqqfffhygjga7rk3shu", nil
	}
	return "", fmt.Errorf("Validator address not registered")
}
func Broadcast(tx []byte) (string, error) {
	body := `{"tx_bytes":` + ByteArrayToStringArray(tx) + `,"mode":"BROADCAST_MODE_BLOCK"}`
	val, err := requester.MakePostRequest("rest", "/cosmos/tx/v1beta1/txs", []byte(body))
	if err != nil {
		return "", fmt.Errorf("Error sending transaction: %q\n", err)
	}

	m := &responses.BroadcastTxResponse{}
	err = json.Unmarshal([]byte(val), m)
	if err != nil {
		return "", fmt.Errorf("Error reading the transaction response: %q", err)
	}

	if m.Response.Code == 0 {
		return m.Response.TxHash, nil
	}

	fmt.Println("Error sending transaction: log...")
	fmt.Println(m)

	return "", fmt.Errorf("Error sending the transaction error code: %d", m.Response.Code)
}
