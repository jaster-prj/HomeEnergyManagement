FROM debian:latest

ENV GO_VERSION "1.21.6"

RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -y ca-certificates openssh-client gnupg git wget tar gcc bzip2
RUN apt-get install -y libeccodes-tools

ENV PATH /usr/local/go/bin:$PATH

RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 1777 "$GOPATH"

RUN go install golang.org/x/tools/gopls@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest

WORKDIR $GOPATH

ENTRYPOINT ["/bin/bash"]