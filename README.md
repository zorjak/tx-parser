# tx-parser

Ethereum blockchain parser that will allow to query transactions for subscribed addresses.

## Run application

to run application use:

```bash
go run ./cmd/tx-parser --starting_block <starting_block>
```

Default values:

- `starting_block=0`

## Test application

In order to test the application 3 endpoints are available:

To get current processed block do

```bash
curl -X GET http://localhost:8080/api/current-block
```

To get transactions for the address do

```bash
curl -X GET http://localhost:8080/api/transactions/{address}
```

To subscribe to the address do

```bash
curl -X POST -H "Content-Type: application/json" -d '{"address":"some-address"}' http://localhost:8080/api/subscribe
```

## Implementation

Scanner scans blocks for the ethereum blockchain and stores hashes of the transactions in the storage.

Storage is implemented as a memory using three maps:

1. `observedAddresses`: Stores addresses subscribed to the service.
2. `transactionsOfAddress`: Stores transaction hashes for each address (sender or receiver).
3. `transactions`: Stores full transaction data for each transaction hash.

## Other notes

1. Transaction data is stored in the database to facilitate returning full transaction details. While it's possible to
   avoid storing transactions and only store hashes, doing so would require querying the EVM node later for each
   transaction, which may be costlier if the company does not maintain its own EVM node, but uses some third-party
   service (e.g. Alchemy, Infura, etc.).
2. All transactions since the starting block are stored. While this approach ensures data integrity, it may not be necessary for the application's requirements. If only transactions from the block of subscription onwards are needed, storing transactions for subscribed addresses exclusively could be considered.
