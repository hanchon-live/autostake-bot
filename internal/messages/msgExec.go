package messages

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethermintcodec "github.com/evmos/ethermint/crypto/codec"
	"github.com/hanchon-live/autostake-bot/internal/blockchain"
)

type ValueToClaim struct {
	Granter     string
	Amount      sdk.Int
	Denom       string
	Validator   string
	IsValidator bool
}

func CreateMessageExec(grantee string, clameable []ValueToClaim) (authz.MsgExec, codec.ProtoCodec, error) {
	// Claim all the comission first and add the amount to the restake message
	// TODO: improve this
	fixedData := []ValueToClaim{}
	extraData := []ValueToClaim{}
	for _, toClaim := range clameable {
		if toClaim.IsValidator == true {

			found := false
			for i, toFix := range clameable {
				if toFix.IsValidator == false && toFix.Granter == toClaim.Granter && toFix.Validator == toClaim.Validator {
					clameable[i].Amount = toFix.Amount.Add(toClaim.Amount)
					found = true
				}
			}

			// If it was not found, it is because the remaining rewards are too low, so we need to create a new one
			if found == false {
				extraData = append(extraData, ValueToClaim{
					Granter:     toClaim.Granter,
					Amount:      toClaim.Amount,
					Denom:       toClaim.Denom,
					Validator:   toClaim.Validator,
					IsValidator: false,
				})
			}
		}
	}

	// Add autostake for commission if there was none to add the balances
	for _, data := range extraData {
		fixedData = append(fixedData, data)
	}

	// Execute the claim commission first and then the rewards
	for _, toClaim := range clameable {
		if toClaim.IsValidator == true {
			fixedData = append(fixedData, toClaim)
		} else {
			fixedData = append([]ValueToClaim{toClaim}, fixedData...)
		}
	}

	// Create the sender account
	granteeAccount, err := blockchain.Bech32StringToAddress(grantee)
	if err != nil {
		return authz.MsgExec{}, codec.ProtoCodec{}, fmt.Errorf("Error creating from address: %q\n", err)
	}

	var messages []sdk.Msg
	for _, toClaim := range fixedData {
		// Create the validator account
		validator, errValidator := blockchain.Bech32StringToValidatorAddress(toClaim.Validator)
		granter, errGranter := blockchain.Bech32StringToAddress(toClaim.Granter)

		if errGranter == nil && errValidator == nil {
			if toClaim.IsValidator == true {
				// Claim the validator rewards
				msg := distribution.NewMsgWithdrawValidatorCommission(validator)
				messages = append(messages, msg)
			} else {
				// Claim delegators rewards
				msg := staking.NewMsgDelegate(granter, validator, sdk.NewCoin(toClaim.Denom, toClaim.Amount))
				messages = append(messages, msg)
			}
		}

		if len(messages) > 20 {
			break
		}
	}

	msgExec := authz.NewMsgExec(granteeAccount, messages)
	// Create the Encoder
	reg := codectypes.NewInterfaceRegistry()
	staking.RegisterInterfaces(reg)
	authz.RegisterInterfaces(reg)
	distribution.RegisterInterfaces(reg)
	ethermintcodec.RegisterInterfaces(reg)
	encoder := codec.NewProtoCodec(reg)

	return msgExec, *encoder, nil
}
