run: 
	@echo "Running the bank-account service with CometBFT consensus module."
	go build -mod=mod
	./bank-account -cmt-home ./cometbft-home

clean:
	rm -rf ./cometbft-home

init: clean
	go run github.com/cometbft/cometbft/cmd/cometbft@v0.38.8 init --home ./cometbft-home

all: clean init run

docker-init:
	rm -rf .docker/config
	go run github.com/cometbft/cometbft/cmd/cometbft@v0.38.8 testnet --config .docker/config-template.toml --o .docker/config/ --starting-ip-address 192.167.10.2

docker-build:
	GOOS=linux GOARCH=amd64 go build -mod=mod
	docker build --tag cometbft/bankaccount .

.PHONY: run clean init