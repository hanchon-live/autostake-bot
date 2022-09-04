package messages

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hanchon-live/autostake-bot/internal/blockchain"
)

func NewMessage(
	msg sdk.Msg,
	encoder codec.ProtoCodec,
	feeAmount int64,
	feeDenom string,
	gasLimit uint64,
	memo string,
	chainId string,
) blockchain.Message {
	// Create the message
	return blockchain.Message{
		Msg:      msg,
		Encoder:  encoder,
		Fee:      blockchain.Int64ToCoins(feeAmount, feeDenom),
		GasLimit: gasLimit,
		Memo:     memo,
		ChainId:  chainId,
	}
}
