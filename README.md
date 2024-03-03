# tx-parser

Ethereum blockchain parser that will allow to query transactions for subscribed addresses.

## Run application

to run application use:

``` bash
go run ./cmd/tx-parser --starting_block <starting_block>
```

Default values:

- `starting_block=0`

## Test application

In order to test the application 3 endpoints are available:

To get current processed block do

``` bash
curl -X GET http://localhost:8080/api/current-block
```

To get transactions for the address do

``` bash
curl -X GET http://localhost:8080/api/transactions/{address}
```

To subscribe to the address do

``` bash
curl -X POST -H "Content-Type: application/json" -d '{"address":"some-address"}' http://localhost:8080/api/subscribe
```

## Implementation

Scanner scans blocks for the ethereum blockchain and stores hashes of the transactions in the database.
Hashes are stored in the storage (memory). Storage is a map which key is `address` (`From` field) or
receiver (`To` field) of the transaction.

Scanner implements function to return full transaction when hash is provided.

Storage contains a map of subscribed address. Subscribed addresses can be added with the endpoint that is provided

## Other notes

All transaction hashes are stored in the storage. I decided to store the hashes of all transactions because if new
address is provided at some point in the future, we would have to rescan the whole blockchain from the start.

When user request transactions for the address, first in the database are found hashes of the transactions and then all
data of the transactions are queried from the blockchain node. I used this approach to save the storage. Although the
storage is relatively cheap if the company maintains the ethereum node there is no need to have additional space to
store transactions. However if the node is rented from a provider like Alchemy probably it is much cheaper to store all
transactions data locally then do separate queries again when they are needed.