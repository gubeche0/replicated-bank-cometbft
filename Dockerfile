FROM alpine:3.19

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file

VOLUME /cometbft
WORKDIR /cometbft
EXPOSE 26656 26657 26660
CMD ["/usr/bin/bank-account", "-cmt-home", "/cometbft/node"]
STOPSIGNAL SIGTERM

COPY bank-account /usr/bin/bank-account
