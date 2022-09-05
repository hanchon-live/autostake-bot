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
	Amount      int64
	Denom       string
	Validator   string
	IsValidator bool
}

func CreateMessageExec(grantee string, clameable []ValueToClaim) (authz.MsgExec, codec.ProtoCodec, error) {
	// Create the sender account
	granteeAccount, err := blockchain.Bech32StringToAddress(grantee)
	if err != nil {
		return authz.MsgExec{}, codec.ProtoCodec{}, fmt.Errorf("Error creating from address: %q\n", err)
	}

	var messages []sdk.Msg
	for _, toClaim := range clameable {
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
				msg := staking.NewMsgDelegate(granter, validator, sdk.NewCoin(toClaim.Denom, sdk.NewInt(toClaim.Amount)))
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
