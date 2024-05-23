## Blockchain Project

This project is a simple implementation of a blockchain, designed to demonstrate the basic capabilities of blockchain technology, including creating and validating blocks, handling transactions, and interfacing with a blockchain via a REST API.

## Features

- Basic blockchain framework.
- REST API to interact with the blockchain.
- Transaction handling with LevelDB for data persistence.
- Block validation and generation.
- Automatic block commit and application exit after 15 seconds of inactivity.
- Calculate SHA-256 hash of blocks and all txns concurrently.

## New Feature: Automatic Block Commit on Inactivity

The blockchain now has an automatic block commit feature. If no new transactions are processed within 15 seconds, any pending transactions are committed to a new block, and the application exits. This ensures that all transactions are processed and saved before the application shuts down due to inactivity.

### How It Works

1. **Transaction Processing**: Transactions are received via the REST API and processed immediately.
2. **Inactivity Timer**: An inactivity timer is reset each time a new transaction is processed.
3. **Block Commit**: If no new transactions are received within 15 seconds, the current block is committed, and the application exits.

### Usage Example

Here is an example of how to send transactions to the blockchain via the REST API using `curl`:

```sh
curl -X POST http://localhost:8080/transactions -H "Content-Type: application/json" -d '[
    {"key": "SIM1", "value": 2, "ver": 1.0},
    {"key": "SIM2", "value": 3, "ver": 1.0},
    {"key": "SIM3", "value": 4, "ver": 2.0}
]'
```

To send transactions to the blockchain using Postman, follow these steps:
- Open Postman and create a new POST request.
- Set the request URL to http://localhost:8080/transactions.
- Set the request headers to include Content-Type: application/json.
- Copy and paste the following JSON payload into the body of the request:

```sh
[
    {"key": "SIM1", "value": 2, "ver": 1.0},
    {"key": "SIM2", "value": 3, "ver": 1.0},
    {"key": "SIM3", "value": 4, "ver": 2.0}
]
```

- Click the "Send" button to submit the request.

To check LevelDB database 
- Run tools/checkdb.go after exiting the program.
- All blocks including invalid txs are stored in the blocks.json file.
