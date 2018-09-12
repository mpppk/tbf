FROM golang:1 AS builder
RUN apt update && apt -y upgrade

WORKDIR /go/src/github.com/mpppk/tbf
COPY Makefile /go/src/github.com/mpppk/tbf/Makefile
COPY vendor /go/src/github.com/mpppk/tbf/vendor

# Install build tools
RUN make setup

# build fzf
RUN git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
RUN ~/.fzf/install

COPY . /go/src/github.com/mpppk/tbf
RUN make CGO_ENABLED=0 install
COPY . /go/src/github.com/mpppk/tbf

FROM alpine
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache jq

COPY --from=builder /go/bin/* /usr/local/bin/
COPY --from=builder /root/.fzf/bin/fzf /usr/local/bin/fzf

WORKDIR /go/src/github.com/mpppk/tbf
COPY ./scripts /go/src/github.com/mpppk/tbf/scripts
RUN chmod -R +x ./scripts
ENV PATH $PATH:/go/src/github.com/mpppk/tbf/scripts
CMD ["tbf"]