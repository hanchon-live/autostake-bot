package blockchain

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/hanchon-live/autostake-bot/internal/wallet"

	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	//	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	ethermintcodec "github.com/evmos/ethermint/crypto/codec"
)

func CreateExecMsg(grantee string) {

	// Move this to params
	sequence := 8
	accountNumber := 11

	granteeAccount := sdk.AccAddress(grantee)

	msg := []sdk.Msg{}
	_ = authz.NewMsgExec(granteeAccount, msg)

	sdk.GetConfig().SetBech32PrefixForAccount("evmos", "evmospub")

	from, err := sdk.AccAddressFromBech32("evmos10gu0eudskw7nc0ef48ce9x22sx3tft0s463el3")
	if err != nil {
		panic(err)
	}
	to, err := sdk.AccAddressFromBech32("evmos1urc5gn9x4kvl3sxu4qd9vckfdmtet7shdskm55")
	if err != nil {
		panic(err)
	}

	msgSend := bank.NewMsgSend(from, to, sdk.NewCoins(sdk.NewCoin("aevmos", sdk.NewInt(10000000))))

	reg := codectypes.NewInterfaceRegistry()
	bank.RegisterInterfaces(reg)
	ethermintcodec.RegisterInterfaces(reg)
	enconder := codec.NewProtoCodec(reg)

	chainId := "evmos_9000-1"
	clientCtx := client.Context{}.
		WithHomeDir("./").
		WithViper("").
		WithCodec(enconder).
		WithChainID(chainId).
		WithTxConfig(authTx.NewTxConfig(enconder, []signing.SignMode{signing.SignMode_SIGN_MODE_DIRECT}))

	fmt.Println(clientCtx)

	txBuilder := clientCtx.TxConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msgSend); err != nil {
		return
	}

	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("aevmos", sdk.NewInt(100000000))))
	txBuilder.SetGasLimit(uint64(150000))

	memo := "Hanchon restake"

	txBuilder.SetMemo(memo)

	w, account, err := wallet.GetWallet()
	if err != nil {
		panic(err)
	}
	privBytes, err := w.PrivateKeyBytes(account)
	if err != nil {
		panic(err)
	}

	priv := &ethsecp256k1.PrivKey{
		Key: privBytes,
	}

	signerData := authsigning.SignerData{
		ChainID:       "evmos_9000-1",
		AccountNumber: uint64(accountNumber),
		Sequence:      uint64(sequence),
		PubKey:        priv.PubKey(),
		Address:       sdk.AccAddress(priv.PubKey().Address()).String(),
	}

	fmt.Println(sdk.AccAddress(priv.PubKey().Address()).String())

	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   priv.PubKey(),
		Data:     &sigData,
		Sequence: uint64(sequence),
	}
	sigs := []signing.SignatureV2{sig}

	if err := txBuilder.SetSignatures(sigs...); err != nil {
		return
	}

	bytesToSign, err := clientCtx.TxConfig.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder.GetTx())
	if err != nil {
		return
	}

	fmt.Println("bytesToSign")
	fmt.Println(bytesToSign)

	// Sign those bytes
	sigBytes, err := priv.Sign(bytesToSign)
	if err != nil {
		return
	}
	fmt.Println("sigBytes")
	fmt.Println(sigBytes)

	// Construct the SignatureV2 struct
	sigData = signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: sigBytes,
	}

	sig = signing.SignatureV2{
		PubKey:   priv.PubKey(),
		Data:     &sigData,
		Sequence: uint64(sequence),
	}

	txBuilder.SetSignatures(sig)

	fmt.Println(txBuilder.GetTx())

	txBz, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(txBz)

	fmt.Println(string(txBz))
	/*

				sigData := signing.SingleSignatureData{
					SignMode:  signMode,
					Signature: nil,
				}
				sig := signing.SignatureV2{
					PubKey:   pubKey,
					Data:     &sigData,
					Sequence: txf.Sequence(),
				}

			// Overwrite or append signer infos.
			var sigs []signing.SignatureV2
			if overwriteSig {
				sigs = []signing.SignatureV2{sig}
			} else {
				sigs = append(prevSignatures, sig)
			}
			if err := txBuilder.SetSignatures(sigs...); err != nil {
				return err
			}


		// Generate the bytes to be signed.
		bytesToSign, err := txf.txConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
		if err != nil {
			return err
		}

		// Sign those bytes
		sigBytes, _, err := txf.keybase.Sign(name, bytesToSign)
		if err != nil {
			return err
		}

		// Construct the SignatureV2 struct
		sigData = signing.SingleSignatureData{
			SignMode:  signMode,
			Signature: sigBytes,
		}
		sig = signing.SignatureV2{
			PubKey:   pubKey,
			Data:     &sigData,
			Sequence: txf.Sequence(),
		}

		if overwriteSig {
			return txBuilder.SetSignatures(sig)
		}
	*/
	/*
		txFactory := tx.Factory{}
		txFactory = txFactory.
			WithChainID(chainId).
			WithMemo(memo).
			WithTxConfig(clientCtx.TxConfig).
			WithSequence(0).
			WithAccountNumber(0)

		tx, err := txFactory.BuildUnsignedTx()
		if err != nil {
			fmt.Println("error", err)
		} else {
			fmt.Println("tx", tx)
		}

		txFactory.PrintUnsignedTx(clientCtx, &exec)
	*/

	/*

		// WithSequence returns a copy of the Factory with an updated sequence number.
		func (f Factory) WithSequence(sequence uint64) Factory {
			f.sequence = sequence
			return f
		}

		// WithMemo returns a copy of the Factory with an updated memo.
		func (f Factory) WithMemo(memo string) Factory {
			f.memo = memo
			return f
		}

		// WithAccountNumber returns a copy of the Factory with an updated account number.
		func (f Factory) WithAccountNumber(accnum uint64) Factory {
			f.accountNumber = accnum
			return f
		}
				valTokens := sdk.TokensFromConsensusPower(100, ethermint.PowerReduction)
				createValMsg, err := stakingtypes.NewMsgCreateValidator(
					sdk.ValAddress(addr),
					valPubKeys[i],
					sdk.NewCoin(ethermint.AttoPhoton, valTokens),
					stakingtypes.NewDescription(nodeDirName, "", "", "", ""),
					stakingtypes.NewCommissionRates(sdk.OneDec(), sdk.OneDec(), sdk.OneDec()),
					sdk.OneInt(),
				)
				if err != nil {
					return err
				}

				txBuilder := clientCtx.TxConfig.NewTxBuilder()
				if err := txBuilder.SetMsgs(createValMsg); err != nil {
					return err
				}

				txBuilder.SetMemo(memo)

				txFactory := tx.Factory{}
				txFactory = txFactory.
					WithChainID(args.chainID).
					WithMemo(memo).
					WithKeybase(kb).
					WithTxConfig(clientCtx.TxConfig)

				if err := tx.Sign(txFactory, nodeDirName, txBuilder, true); err != nil {
					return err
				}

				txBz, err := clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
				if err != nil {
					return err
				}
	*/
}
