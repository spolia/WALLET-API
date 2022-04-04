FROM golang:latest as BUILD
RUN apk add --update git
WORKDIR builddir
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o WALLET-API

FROM alpine
COPY --from=BUILD ./go/builddir/WALLET-API .
ENTRYPOINT ["./WALLET-API"]