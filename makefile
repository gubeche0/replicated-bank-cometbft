run: 
	@echo "Running the bank-account service with CometBFT consensus module."
	go build -mod=mod
	./bank-account -cmt-home ./cometbft-home

clean:
	rm -rf ./cometbft-home

init: clean
	go run github.com/cometbft/cometbft/cmd/cometbft@v0.38.8 init --home ./cometbft-home

all: clean init run

.PHONY: run clean init