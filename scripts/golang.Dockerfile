FROM golang:1.17


WORKDIR /go

RUN go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest