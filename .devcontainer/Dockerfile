FROM alpine:latest

RUN apk add --no-cache ca-certificates

ENV PATH /usr/local/go/bin:$PATH

RUN set -eux; \
    apk add --no-cache --virtual .fetch-deps gnupg; \
    apk add --no-cache --virtual .build-deps bash gcc go musl-dev

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 1777 "$GOPATH"

RUN go install golang.org/x/tools/gopls@latest git

WORKDIR $GOPATH

ENTRYPOINT ["/bin/bash"]