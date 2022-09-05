# Autostake-bot

Simple cli application to restake your evmos:

- Add your wallet to the `.env` file
- Call `./autostake-bot init` to create the database
- Call `./autostake-bot updateGranters` to get all the granters assinged to a wallet
- Call `./autostake-bot start` to restake all the rewards from your granters

Other options:

- Call `./autostake-bot listGranters` to read all the granters stored in the database
- Call `./autostake-bot clean` to remove the database file from your system

## Configuration

Copy the `.env.example` file to `./env` and replace the values with your settings

- rest: endpoint used to connect to the Evmos chain
- mnemonic: your wallet mnemonic (the bot wallet)
- derivationpath: derivation path used by the bot wallet
- granteewallet: bot wallet in bech32 evmos format (evmos1...)
- fee: amount to pay for each transaction
- feedenom: denom for the fee
- gaslimit: max gas to be used by the transaction
- memo: transaction memo
- chainid: chain id in cosmos format
- validator: validator account where the rewards will be restaked (evmosvaloper1...)
- minreward: min coin amount needed to be included in the transaction, (i.e, claim if the user has more 3000000000aevmos as rewards)

## Run

```sh
go build
./autostake-bot help
```

## Grant permission using the evmosd cli

Localnet example:

```sh
evmosd tx authz grant evmos10gu0eudskw7nc0ef48ce9x22sx3tft0s463el3 generic --msg-type /cosmos.staking.v1beta1.MsgDelegate --chain-id evmos_9000-1 --node http://localhost:26657 --from mykey --keyring-backend test --gas auto --gas-prices 25000000000.0000aevmos --gas-adjustment 1.5
evmosd tx authz grant evmos10gu0eudskw7nc0ef48ce9x22sx3tft0s463el3 generic --msg-type /cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission --chain-id evmos_9000-1 --node http://localhost:26657 --from mykey --keyring-backend test --gas auto --gas-prices 25000000000.0000aevmos --gas-adjustment 1.5
```

Testnet example:

```sh
evmosd tx authz grant evmos1h2n8tfp9z75xvck580ny3gv7hn74fe2vqtxpd0 generic --msg-type /cosmos.staking.v1beta1.MsgDelegate --chain-id evmos_9000-4 --node https://tendermint.bd.evmos.dev:26657/ --from testnet --keyring-backend file --fees 2000000000000atevmos -b block
```

## TODO:

- Remove grantee from the list if grant revoked
- More tests, right now the only thing tested is the transaction generator
- Allow more than 21 granters
- Validate if the grant is limited to a validator or by amount
