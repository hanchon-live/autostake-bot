# Autostake-bot

## Configuration

Create your .env file with all your settings

## Run

```sh
go build
./autostake-bot help
```

## Grant permission using the evmosd cli

```sh
evmosd tx authz grant evmos10gu0eudskw7nc0ef48ce9x22sx3tft0s463el3 generic --msg-type /cosmos.staking.v1beta1.MsgDelegate --chain-id evmos_9000-1 --node http://localhost:26657 --from mykey --keyring-backend test --gas auto --gas-prices 25000000000.0000aevmos --gas-adjustment 1.5
```
