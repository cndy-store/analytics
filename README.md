# Native

## Dependencies: PostgreSQL

Install PostgreSQL and point the following environment variables to the database:

```bash
export PGHOST=127.0.0.1
export PGUSER=cndy
export PGDATABASE=cndy
export PGSSLMODE=disable
```

Build the API natively (Dependency management is done via [Glide](https://glide.sh)):

## Build

```bash
glide install
go build
```

## Run tests

```bash
go test ./...
```

# Docker

## Run API inside Docker

Spawn PostgreSQL instance and CNDY analytics API via [Docker](https://docker.com/):

```bash
docker-compose up
```

## Run tests inside Docker

```bash
docker-compose run -e DATABASE=cndy_test api go test ./...
```

## Copy cross compiled binary for Linux from Docker container to host

```bash
docker run -v $PWD:/host --entrypoint cp cndy-store/analytics analytics /host/cndy-linux-amd64
```


# API endpoints and examples

## Latest stats

GET https://api.cndy.store/stats/latest?asset_code=CNDY&asset_issuer=GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX

```json
{
  "status": "ok",
  "latest": {
    "paging_token": "33825130903777281-1",
    "asset_type": "credit_alphanum4",
    "asset_code": "CNDY",
    "asset_issuer": "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX",
    "payments": 4,
    "accounts_with_trustline": 4,
    "accounts_with_payments": 2,
    "created_at": "2018-03-12T18:49:40Z",
    "issued": "2000.0000000",
    "transferred": "40.0000000"
  }
}
```

## Asset stats history

GET https://api.cndy.store/stats?asset_code=CNDY&asset_issuer=GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX[&from=2018-03-03T23:05:40Z&to=2018-03-03T23:05:50Z]

If not set, `from` defaults to UNIX timestamp `0`, `to` to `now`.

```json
{
  "status": "ok",
  "stats": [
    {
      "paging_token": "33864305300480001-1",
      "asset_type": "credit_alphanum4",
      "asset_code": "CNDY",
      "asset_issuer": "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX",
      "payments": 6,
      "accounts_with_trustline": 5,
      "accounts_with_payments": 2,
      "created_at": "2018-03-13T07:29:48Z",
      "issued": "2000.0000000",
      "transferred": "140.0000000"
    },
    {
      "paging_token": "33864305300480001-2",
      "asset_type": "credit_alphanum4",
      "asset_code": "CNDY",
      "asset_issuer": "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX",
      "payments": 6,
      "accounts_with_trustline": 5,
      "accounts_with_payments": 2,
      "created_at": "2018-03-13T07:29:48Z",
      "issued": "3000.0000000",
      "transferred": "140.0000000"
    }
  ]
}
```

## Current Horizon cursor

GET https://api.cndy.store/stats/cursor

```json
{
  "status": "ok",
  "current_cursor": "33877250331906049-1"
}
```

## Effects

GET https://api.cndy.store/effects?asset_code=CNDY&asset_issuer=GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX[&from=2018-03-03T23:05:40Z&to=2018-03-03T23:05:50Z]


If not set, `from` defaults to UNIX timestamp `0`, `to` to `now`.

```json
{
  "status": "ok",
  "effects": [
    {
      "id": "0033819672000335873-0000000001",
      "operation": "https://horizon-testnet.stellar.org/operations/33819672000335873",
      "succeeds": "https://horizon-testnet.stellar.org/effects?order=desc&cursor=33819672000335873-1",
      "precedes": "https://horizon-testnet.stellar.org/effects?order=asc&cursor=33819672000335873-1",
      "paging_token": "33819672000335873-1",
      "account": "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD",
      "type": "trustline_created",
      "type_i": 20,
      "starting_balance": "",
      "asset_type": "credit_alphanum4",
      "asset_code": "CNDY",
      "asset_issuer": "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX",
      "signer_public_key": "",
      "signer_weight": 0,
      "signer_key": "",
      "signer_type": "",
      "created_at": "2018-03-12T17:03:45Z",
      "amount": "0.0000000",
      "balance": "0.0000000",
      "balance_limit": "922337203685.4775807"
    },
    {
      "id": "0033820110087000065-0000000002",
      "operation": "https://horizon-testnet.stellar.org/operations/33820110087000065",
      "succeeds": "https://horizon-testnet.stellar.org/effects?order=desc&cursor=33820110087000065-2",
      "precedes": "https://horizon-testnet.stellar.org/effects?order=asc&cursor=33820110087000065-2",
      "paging_token": "33820110087000065-2",
      "account": "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX",
      "type": "account_debited",
      "type_i": 3,
      "starting_balance": "",
      "asset_type": "credit_alphanum4",
      "asset_code": "CNDY",
      "asset_issuer": "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX",
      "signer_public_key": "",
      "signer_weight": 0,
      "signer_key": "",
      "signer_type": "",
      "created_at": "2018-03-12T17:12:15Z",
      "amount": "1000.0000000",
      "balance": "0.0000000",
      "balance_limit": "0.0000000"
    }
}
```
