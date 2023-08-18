FROM golang:1.21 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /mumble-status-server ./cmd/server/

FROM scratch
COPY --from=builder /mumble-status-server /mumble-status-server
EXPOSE 3000

CMD ["/mumble-status-server"]
