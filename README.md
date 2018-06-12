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

## Stats

GET https://api.cndy.store/stats[?from=2018-03-03T23:05:40Z&to=2018-03-03T23:05:50Z]

If not set, `from` defaults to UNIX timestamp `0`, `to` to `now`.

```json
{
  "accounts_involved": 10,
  "amount_issued": "4000.0000000",
  "amount_transferred": "6410.0000000",
  "asset_code": "CNDY",
  "current_cursor": "35496616211263489-2",
  "effect_count": 59,
  "trustlines_created": 9
}
```

## Effects

GET https://api.cndy.store/effects[?from=2018-03-03T23:05:40Z&to=2018-03-03T23:05:50Z]

If not set, `from` defaults to UNIX timestamp `0`, `to` to `now`.

```json
{
  "effects": [
    {
      "id": "0034641642841444353-0000000002",
      "operation": "https://horizon-testnet.stellar.org/operations/34641642841444353",
      "succeeds": "https://horizon-testnet.stellar.org/effects?order=desc&cursor=34641642841444353-2",
      "precedes": "https://horizon-testnet.stellar.org/effects?order=asc&cursor=34641642841444353-2",
      "paging_token": "34641642841444353-2",
      "account": "GBET4AVL4BYLFJTFKX2UYE3Y3EHNZXOBMBO5FP7O5FFTHSEAPZ5VEHHD",
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
      "created_at": "2018-03-23T18:54:05Z",
      "amount": "10.0000000",
      "balance": "0.0000000",
      "balance_limit": "0.0000000"
    },
    {
      "id": "0034683389923561473-0000000002",
      "operation": "https://horizon-testnet.stellar.org/operations/34683389923561473",
      "succeeds": "https://horizon-testnet.stellar.org/effects?order=desc&cursor=34683389923561473-2",
      "precedes": "https://horizon-testnet.stellar.org/effects?order=asc&cursor=34683389923561473-2",
      "paging_token": "34683389923561473-2",
      "account": "GBET4AVL4BYLFJTFKX2UYE3Y3EHNZXOBMBO5FP7O5FFTHSEAPZ5VEHHD",
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
      "created_at": "2018-03-24T08:24:06Z",
      "amount": "10.0000000",
      "balance": "0.0000000",
      "balance_limit": "0.0000000"
    }
  ]
}
```

## Asset stats

GET https://api.cndy.store/history[?from=2018-03-03T23:05:40Z&to=2018-03-03T23:05:50Z]

If not set, `from` defaults to UNIX timestamp `0`, `to` to `now`.

```json
{
  "history": [
    {
      "paging_token": "34683389923561473-1",
      "asset_type": "credit_alphanum4",
      "asset_code": "CNDY",
      "asset_issuer": "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX",
      "num_accounts": 10,
      "payments": 25,
      "created_at": "2018-03-24T08:24:06Z",
      "total_amount": "4000.0000000"
    },
    {
      "paging_token": "34683389923561473-2",
      "asset_type": "credit_alphanum4",
      "asset_code": "CNDY",
      "asset_issuer": "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX",
      "num_accounts": 10,
      "payments": 28,
      "created_at": "2018-03-24T08:24:06Z",
      "total_amount": "4000.0000000"
    }
  ]
}
```
