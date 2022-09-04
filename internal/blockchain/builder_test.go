package blockchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	priv, err := CreatePrivateKeyFromMnemonic("flash local taste power maple fragile pool name file position drop swarm")
	if err != nil {
		t.Fatalf("Error creating priv: %q", err)
	}

	sender := Sender{
		Sequence:      uint64(8),
		AccountNumber: uint64(11),
		PrivKey:       priv,
	}

	from, err := Bech32StringToAddress("evmos10gu0eudskw7nc0ef48ce9x22sx3tft0s463el3")
	if err != nil {
		t.Fatalf("Error creating from address: %q", err)
	}
	to, err := Bech32StringToAddress("evmos1urc5gn9x4kvl3sxu4qd9vckfdmtet7shdskm55")
	if err != nil {
		t.Fatalf("Error creating to address: %q", err)
	}

	msgSend := bank.NewMsgSend(from, to, sdk.NewCoins(sdk.NewCoin("aevmos", sdk.NewInt(10000000))))

	message := Message{
		Msg:      msgSend,
		Enconder: *CreateEnconder(),
		Fee:      sdk.NewCoins(sdk.NewCoin("aevmos", sdk.NewInt(100000000))),
		GasLimit: uint64(150000),
		Memo:     "Hanchon restake",
		ChainId:  "evmos_9000-1",
	}

	tx, err := CreateTransaction(sender, message)
	if err != nil {
		t.Fatalf("Error creating transaction: %q", err)
	}

	expected := []byte{10, 164, 1, 10, 144, 1, 10, 28, 47, 99, 111, 115, 109, 111, 115, 46, 98, 97, 110, 107, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 83, 101, 110, 100, 18, 112, 10, 44, 101, 118, 109, 111, 115, 49, 48, 103, 117, 48, 101, 117, 100, 115, 107, 119, 55, 110, 99, 48, 101, 102, 52, 56, 99, 101, 57, 120, 50, 50, 115, 120, 51, 116, 102, 116, 48, 115, 52, 54, 51, 101, 108, 51, 18, 44, 101, 118, 109, 111, 115, 49, 117, 114, 99, 53, 103, 110, 57, 120, 52, 107, 118, 108, 51, 115, 120, 117, 52, 113, 100, 57, 118, 99, 107, 102, 100, 109, 116, 101, 116, 55, 115, 104, 100, 115, 107, 109, 53, 53, 26, 18, 10, 6, 97, 101, 118, 109, 111, 115, 18, 8, 49, 48, 48, 48, 48, 48, 48, 48, 18, 15, 72, 97, 110, 99, 104, 111, 110, 32, 114, 101, 115, 116, 97, 107, 101, 18, 118, 10, 89, 10, 79, 10, 40, 47, 101, 116, 104, 101, 114, 109, 105, 110, 116, 46, 99, 114, 121, 112, 116, 111, 46, 118, 49, 46, 101, 116, 104, 115, 101, 99, 112, 50, 53, 54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35, 10, 33, 3, 13, 251, 35, 125, 14, 103, 28, 123, 74, 120, 140, 166, 125, 108, 228, 103, 75, 196, 171, 1, 118, 251, 198, 71, 21, 182, 144, 67, 143, 150, 218, 123, 18, 4, 10, 2, 8, 1, 24, 8, 18, 25, 10, 19, 10, 6, 97, 101, 118, 109, 111, 115, 18, 9, 49, 48, 48, 48, 48, 48, 48, 48, 48, 16, 240, 147, 9, 26, 65, 223, 20, 117, 108, 116, 246, 33, 96, 133, 205, 124, 140, 226, 101, 91, 140, 18, 69, 252, 238, 13, 89, 251, 80, 204, 41, 125, 253, 108, 234, 232, 254, 54, 140, 21, 57, 204, 121, 143, 231, 245, 79, 206, 100, 67, 174, 141, 37, 160, 210, 246, 87, 228, 132, 178, 34, 100, 135, 34, 148, 50, 204, 243, 121, 0}
	if len(tx) != len(expected) {
		t.Fatalf("Bad tx generated")
	}
	for k, v := range tx {
		if expected[k] != v {
			t.Fatalf("Bad tx generated")
		}
	}
}
