FROM golang:1.26 AS builder

RUN go install github.com/go-task/task/v3/cmd/task@v3.50.0

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN task server:build

# ---

FROM scratch

COPY --from=builder /app/tmp/mycorp-server /bin/mycorp-server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["/bin/mycorp-server"]
