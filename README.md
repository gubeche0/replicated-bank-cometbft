# sistema-distribuido-t2

Implement a simple replicated bank application. The bank service maintains a set of client records; a set of accounts (a client can have zero or more accounts); and operations such as deposit, withdraw, transfer, and inquiry.

## Executando o projeto

### Criando arquivos de configuração do node
```bash
go run github.com/cometbft/cometbft/cmd/cometbft@v0.38.8 init --home ./cometbft-home

```

### Compilando e executando a aplicação

```bash
go build -mod=mod
./bank-account -cmt-home ./cometbft-home
```