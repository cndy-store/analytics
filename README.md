# Build

Dependency management is done via [Glide](https://glide.sh).

Build natively:

```bash
glide install
go build
```

Cross compile for Linux with [Docker](https://docker.com/):

```bash
docker-compose up
```

Run tests

```bash
docker-compose run cndy-analytics /bin/bash -c 'go test ./...'
```


# API

## Stats

GET http://api.cndy.store:3144/stats[?from=2018-03-03T23:05:40Z&to=2018-03-03T23:05:50Z]

If not set, `from` defaults to UNIX timestamp `0`, `to` to `now`.

```json
{
  "accounts_involved": 3,
  "amount_issued": 1108,
  "amount_transferred": 2119,
  "asset_code": "CNDY",
  "current_cursor": "33592975036518402-2",
  "trustlines_created": 1,
  "effect_count": 25
}
```

## Effects

GET http://api.cndy.store:3144/effects[?from=2018-03-03T23:05:40Z&to=2018-03-03T23:05:50Z]

If not set, `from` defaults to UNIX timestamp `0`, `to` to `now`.

```json
{
  "effects": [
    {
      "id": "0033170775456358401-0000000001",
      "operation": "https://horizon-testnet.stellar.org/operations/33170775456358401",
      "paging_token": "33170775456358401-1",
      "account": "GDQLK4Y26S5H6X3WAAMKFIT2575N54B667TCJA3KTS5XBYL5KUJWMFRM",
      "amount": "12.0000000",
      "type": "account_credited",
      "starting_balance": "",
      "balance": "",
      "balance_limit": "",
      "asset_type": "credit_alphanum4",
      "asset_code": "CNDY",
      "asset_issuer": "GD7YB3R3TKUU3OHTE3DO5BIVBLQVFKYRHPW5Y6NHVSQVNNEOQ5I2RKLU",
      "created_at": "2018-03-03T23:05:43Z"
    },
    {
      "id": "0033170775456358401-0000000002",
      "operation": "https://horizon-testnet.stellar.org/operations/33170775456358401",
      "paging_token": "33170775456358401-2",
      "account": "GD7YB3R3TKUU3OHTE3DO5BIVBLQVFKYRHPW5Y6NHVSQVNNEOQ5I2RKLU",
      "amount": "12.0000000",
      "type": "account_debited",
      "starting_balance": "",
      "balance": "",
      "balance_limit": "",
      "asset_type": "credit_alphanum4",
      "asset_code": "CNDY",
      "asset_issuer": "GD7YB3R3TKUU3OHTE3DO5BIVBLQVFKYRHPW5Y6NHVSQVNNEOQ5I2RKLU",
      "created_at": "2018-03-03T23:05:43Z"
    }
  ]
}
```

## Asset stats

GET http://api.cndy.store:3144/history[?from=2018-03-03T23:05:40Z&to=2018-03-03T23:05:50Z]

If not set, `from` defaults to UNIX timestamp `0`, `to` to `now`.

```json
{
  "history": [
    {
      "asset_type": "credit_alphanum4",
      "asset_code": "GD7YB3R3TKUU3OHTE3DO5BIVBLQVFKYRHPW5Y6NHVSQVNNEOQ5I2RKLU",
      "asset_issuer": "CNDY",
      "num_accounts": 1,
      "num_effects": 1,
      "created_at": "2018-03-03T23:05:43Z"
    },
    {
      "asset_type": "credit_alphanum4",
      "asset_code": "GD7YB3R3TKUU3OHTE3DO5BIVBLQVFKYRHPW5Y6NHVSQVNNEOQ5I2RKLU",
      "asset_issuer": "CNDY",
      "total_amount": "12.0",
      "num_accounts": 2,
      "num_effects": 2,
      "created_at": "2018-03-03T23:05:43Z"
    }
  ]
}
```
