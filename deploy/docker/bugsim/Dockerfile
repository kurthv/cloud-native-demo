ARG GO_VERSION=1.17

FROM golang:${GO_VERSION} as builder

ARG VERSION="0.0.4-dev14"
ARG TIME
ARG GO_MODULE="github.com/cypherfox/cloud-native-demo"

ARG LDFLAGS="-extldflags=-static -X '${GO_MODULE}/pkg/version.BuildVersion=${VERSION}'"

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/

RUN export TIME=`date` ; \
    export FULL_LDFLAGS="${LDFLAGS} -X '${GO_MODULE}/pkg/version.BuildTime=${TIME}'"
    CGO_ENABLED=0 go build -o bin/bugsim -ldflags=${LDFLAGS} ./cmd/bugsim

FROM scratch

COPY --from=builder /workspace/bin/bugsim /bugsim

EXPOSE 80/tcp

CMD ["/bugsim"]
