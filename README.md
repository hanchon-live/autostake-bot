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

## Run

```sh
go build
./autostake-bot help
```

## Grant permission using the evmosd cli

Localnet example: replace denom, chain id and node for mainnet/testnet

```sh
evmosd tx authz grant evmos10gu0eudskw7nc0ef48ce9x22sx3tft0s463el3 generic --msg-type /cosmos.staking.v1beta1.MsgDelegate --chain-id evmos_9000-1 --node http://localhost:26657 --from mykey --keyring-backend test --gas auto --gas-prices 25000000000.0000aevmos --gas-adjustment 1.5
```

## TODO:

- Remove grantee from the list if grant revoked
- More tests, right now the only thing tested is the transaction generator
