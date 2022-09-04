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
