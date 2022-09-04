package messages

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethermintcodec "github.com/evmos/ethermint/crypto/codec"
	"github.com/hanchon-live/autostake-bot/internal/blockchain"
)

func CreateMessageSend(sender string, receiver string, amount int64, denom string) (bank.MsgSend, codec.ProtoCodec, error) {
	// Create the proto message
	from, err := blockchain.Bech32StringToAddress(sender)
	if err != nil {
		return bank.MsgSend{}, codec.ProtoCodec{}, fmt.Errorf("Error creating from address: %q\n", err)
	}
	to, err := blockchain.Bech32StringToAddress(receiver)
	if err != nil {
		return bank.MsgSend{}, codec.ProtoCodec{}, fmt.Errorf("Error creating to address: %q\n", err)
	}

	msgSend := bank.NewMsgSend(from, to, blockchain.Int64ToCoins(amount, denom))

	// Create the Encoder
	reg := codectypes.NewInterfaceRegistry()
	bank.RegisterInterfaces(reg)
	ethermintcodec.RegisterInterfaces(reg)
	encoder := codec.NewProtoCodec(reg)

	return *msgSend, *encoder, nil

}
