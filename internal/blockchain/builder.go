package blockchain

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/hanchon-live/autostake-bot/internal/wallet"

	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	ethermintcodec "github.com/evmos/ethermint/crypto/codec"
)

type Sender struct {
	Sequence      uint64
	AccountNumber uint64
	PrivKey       ethsecp256k1.PrivKey
}

type Message struct {
	Msg      sdk.Msg
	Enconder codec.ProtoCodec
	Fee      sdk.Coins
	GasLimit uint64
	Memo     string
	ChainId  string
}

func CreateEnconder() *codec.ProtoCodec {
	reg := codectypes.NewInterfaceRegistry()
	bank.RegisterInterfaces(reg)
	ethermintcodec.RegisterInterfaces(reg)
	enconder := codec.NewProtoCodec(reg)
	return enconder
}

func CreatePrivateKeyFromMnemonic(mnemonic string) (ethsecp256k1.PrivKey, error) {
	w, account, err := wallet.GetWalletFromMnemonic(mnemonic)
	if err != nil {
		return ethsecp256k1.PrivKey{}, err
	}

	privBytes, err := w.PrivateKeyBytes(account)
	if err != nil {
		return ethsecp256k1.PrivKey{}, err
	}

	priv := ethsecp256k1.PrivKey{
		Key: privBytes,
	}

	return priv, nil
}

func Bech32StringToAddress(address string) (sdk.AccAddress, error) {
	sdk.GetConfig().SetBech32PrefixForAccount("evmos", "evmospub")
	return sdk.AccAddressFromBech32(address)
}

func CreateTransaction(sender Sender, message Message) ([]byte, error) {
	sdk.GetConfig().SetBech32PrefixForAccount("evmos", "evmospub")

	clientCtx := client.Context{}.
		WithHomeDir("./").
		WithViper("").
		WithCodec(&message.Enconder).
		WithChainID(message.ChainId).
		WithTxConfig(authTx.NewTxConfig(&message.Enconder, []signing.SignMode{signing.SignMode_SIGN_MODE_DIRECT}))

	txBuilder := clientCtx.TxConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(message.Msg); err != nil {
		return []byte{}, err
	}

	txBuilder.SetFeeAmount(message.Fee)
	txBuilder.SetGasLimit(message.GasLimit)
	txBuilder.SetMemo(message.Memo)

	signerData := authsigning.SignerData{
		ChainID:       message.ChainId,
		AccountNumber: sender.AccountNumber,
		Sequence:      sender.Sequence,
		PubKey:        sender.PrivKey.PubKey(),
		Address:       sdk.AccAddress(sender.PrivKey.PubKey().Address()).String(),
	}

	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   sender.PrivKey.PubKey(),
		Data:     &sigData,
		Sequence: sender.Sequence,
	}
	sigs := []signing.SignatureV2{sig}

	if err := txBuilder.SetSignatures(sigs...); err != nil {
		return []byte{}, err
	}

	bytesToSign, err := clientCtx.TxConfig.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder.GetTx())
	if err != nil {
		return []byte{}, err
	}

	// Sign those bytes
	sigBytes, err := sender.PrivKey.Sign(bytesToSign)
	if err != nil {
		return []byte{}, err
	}

	// Construct the SignatureV2 struct
	sigData = signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: sigBytes,
	}
	sig = signing.SignatureV2{
		PubKey:   sender.PrivKey.PubKey(),
		Data:     &sigData,
		Sequence: sender.Sequence,
	}
	txBuilder.SetSignatures(sig)

	txBz, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return []byte{}, err
	}
	return txBz, nil
}
