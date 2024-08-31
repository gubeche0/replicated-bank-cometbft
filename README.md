# sistema-distribuido-t2

This project was created for the "Distributed Systems" course at the university.

Implement a simple replicated bank application. The bank service maintains a set of client records; a set of accounts (a client can have zero or more accounts); and operations such as deposit, withdraw, transfer, and inquiry. This project uses the CometBFT library to implement the consensus algorithm with byzantine fault tolerance.

## Execute the application

### create the genesis file and configuration
```bash
go run github.com/cometbft/cometbft/cmd/cometbft@v0.38.8 init --home ./cometbft-home

```

### compile and run the application

```bash
go build -mod=mod
./bank-account -cmt-home ./cometbft-home
```